// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3/types"

	utils "github.com/selendra/selendra-bridge/ChainBridge/shared/substrate"
	subtest "github.com/selendra/selendra-bridge/ChainBridge/shared/substrate/testing"
	message "github.com/selendra/selendra-bridge/chainbridge-utils/msg"
)

func assertProposalState(t *testing.T, conn *Connection, prop *proposal, votes *voteState, hasValue bool) {
	var voteRes voteState
	srcId, err := types.EncodeToBytes(prop.sourceId)
	if err != nil {
		t.Fatal(err)
	}
	propBz, err := prop.encode()
	if err != nil {
		t.Fatal(err)
	}
	ok, err := conn.queryStorage(utils.BridgeStoragePrefix, "Votes", srcId, propBz, &voteRes)
	if err != nil {
		t.Fatalf("failed to query votes: %s", err)
	}
	if hasValue {
		if !reflect.DeepEqual(&voteRes, votes) {
			t.Fatalf("Vote state incorrect.\n\tExpected: %#v\n\tGot: %#v", votes, &voteRes)
		}
	}

	if !ok && hasValue {
		t.Fatalf("expected vote to exists but is None")
	}
}

func Test_ContainsVote(t *testing.T) {
	votes := []types.AccountID{types.NewAccountID(AliceKey.PublicKey)}
	if !containsVote(votes, types.NewAccountID(AliceKey.PublicKey)) {
		t.Error("Voter has votes")
	}

	if containsVote(votes, types.NewAccountID(BobKey.PublicKey)) {
		t.Error("Voter has not voted")
	}
}

func TestWriter_ResolveMessage_FungibleProposal(t *testing.T) {
	// Assert Bob's starting balances
	var startingBalance types.U128
	getFreeBalance(context.writerBob.conn, &startingBalance)

	// Setup message and params
	var rId [32]byte
	subtest.QueryConst(t, context.client, "Example", "NativeTokenId", &rId)
	// Construct the message to initiate a vote
	amount := big.NewInt(10000000)
	m := message.NewFungibleTransfer(ForeignChain, ThisChain, 0, amount, rId, context.writerBob.conn.key.PublicKey)
	// Create a proposal to help us check results
	prop, err := context.writerAlice.createFungibleProposal(m)
	if err != nil {
		t.Fatal(err)
	}

	// First, ensure the proposal doesn't already exist
	assertProposalState(t, context.writerAlice.conn, prop, nil, false)

	// Submit the message for processing
	ok := context.writerAlice.ResolveMessage(m)
	if !ok {
		t.Fatal("Alice failed to resolve the message")
	}

	// Now check if the assetTxProposal exists on chain
	singleVoteState := &voteState{
		VotesFor: []types.AccountID{types.NewAccountID(context.writerAlice.conn.key.PublicKey)},
		Status:   voteStatus{IsActive: true},
	}
	assertProposalState(t, context.writerAlice.conn, prop, singleVoteState, true)

	// Submit a second vote from Bob this time
	ok = context.writerBob.ResolveMessage(m)
	if !ok {
		t.Fatalf("Bob failed to resolve the message")
	}

	// Check the vote was added
	finalVoteState := &voteState{
		VotesFor: []types.AccountID{
			types.NewAccountID(context.writerAlice.conn.key.PublicKey),
			types.NewAccountID(context.writerBob.conn.key.PublicKey),
		},
		Status: voteStatus{IsApproved: true},
	}
	assertProposalState(t, context.writerAlice.conn, prop, finalVoteState, true)

	// Assert balance has changed
	var bBal types.U128
	getFreeBalance(context.writerBob.conn, &bBal)
	if bBal == startingBalance {
		t.Fatalf("Internal transaction failed to update Bobs balance")
	} else {
		t.Logf("Bob's new balance: %s (amount: %s)", bBal.String(), big.NewInt(0).Sub(bBal.Int, startingBalance.Int).String())
	}

	select {
	case err = <-context.wSysErr:
		t.Fatal(err)
	default:
	}
}

func TestWriter_ResolveMessage_Duplicate(t *testing.T) {
	// Setup message and params
	var rId [32]byte
	subtest.QueryConst(t, context.client, "Example", "NativeTokenId", &rId)
	// Construct the message to initiate a vote
	amount := big.NewInt(10000000)
	context.latestInNonce++
	m := message.NewFungibleTransfer(ForeignChain, ThisChain, context.latestInNonce, amount, rId, context.writerBob.conn.key.PublicKey)
	// Create a proposal to help us check results
	prop, err := context.writerAlice.createFungibleProposal(m)
	if err != nil {
		t.Fatal(err)
	}

	// First, ensure the proposal doesn't already exist
	assertProposalState(t, context.writerAlice.conn, prop, nil, false)

	// Submit the message for processing
	ok := context.writerAlice.ResolveMessage(m)
	if !ok {
		t.Fatal("Alice failed to resolve the message")
	}

	// Now check if the proposal exists on chain
	singleVoteState := &voteState{
		VotesFor: []types.AccountID{types.NewAccountID(context.writerAlice.conn.key.PublicKey)},
		Status:   voteStatus{IsActive: true},
	}
	assertProposalState(t, context.writerAlice.conn, prop, singleVoteState, true)

	// Submit a second vote from Bob this time
	ok = context.writerAlice.ResolveMessage(m)
	if !ok {
		t.Fatalf("Bob failed to resolve the message")
	}

	assertProposalState(t, context.writerAlice.conn, prop, singleVoteState, true)

	select {
	case err = <-context.wSysErr:
		t.Fatal(err)
	default:
	}

}
