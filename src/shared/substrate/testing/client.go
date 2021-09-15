// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package subtest

import (
	"bytes"
	"math/big"
	"testing"

	utils "github.com/selendra/selendra-bridge/ChainBridge/shared/substrate"
	"github.com/selendra/selendra-bridge/chainbridge-utils/msg"
	"github.com/ChainSafe/log15"
	"github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3/signature"
	"github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3/types"
)

func CreateClient(t *testing.T, key *signature.KeyringPair, endpoint string) *utils.Client {
	client, err := utils.CreateClient(key, endpoint)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func AddRelayer(t *testing.T, client *utils.Client, relayer types.AccountID) {
	err := client.AddRelayer(relayer)
	if err != nil {
		t.Fatal(err)
	}
}

func WhitelistChain(t *testing.T, client *utils.Client, id msg.ChainId) {
	err := client.WhitelistChain(id)
	if err != nil {
		t.Fatal(err)
	}
}

func InitiateNativeTransfer(t *testing.T, client *utils.Client, amount types.U128, recipient []byte, destId msg.ChainId) {
	err := client.InitiateNativeTransfer(amount, recipient, destId)
	if err != nil {
		t.Fatal(err)
	}
}

func RegisterResource(t *testing.T, client *utils.Client, id msg.ResourceId, method string) {
	err := client.RegisterResource(id, method)
	if err != nil {
		t.Fatal(err)
	}
}

func BalanceOf(t *testing.T, client *utils.Client, publicKey []byte) *big.Int {
	balance, err := utils.BalanceOf(client, publicKey)
	if err != nil {
		t.Fatal(err)
	}
	return balance
}

func AssertBalanceOf(t *testing.T, client *utils.Client, publicKey []byte, expected *big.Int) {
	current, err := utils.BalanceOf(client, publicKey)
	if err != nil {
		t.Fatal(err)
	}

	if expected.Cmp(current) != 0 {
		t.Fatalf("Balance does not match expected. Expected: %s Got: %s (change %s) \n", expected.String(), current.String(), big.NewInt(0).Sub(current, expected).String())
	}
}

func GetDepositNonce(t *testing.T, client *utils.Client, chain msg.ChainId) uint64 {
	count, err := client.GetDepositNonce(chain)
	if err != nil {
		t.Fatal(err)
	}
	return count
}

func NewNativeTransferCall(t *testing.T, client *utils.Client, amount types.U128, recipient []byte, destId msg.ChainId) types.Call {
	call, err := client.NewNativeTransferCall(amount, recipient, destId)
	if err != nil {
		t.Fatal(err)
	}
	return call
}
