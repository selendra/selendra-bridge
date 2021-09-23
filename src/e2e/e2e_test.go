// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package e2e

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"

	ethChain "github.com/selendra/selendra-bridge/ChainBridge/chains/ethereum"
	subChain "github.com/selendra/selendra-bridge/ChainBridge/chains/substrate"
	eth "github.com/selendra/selendra-bridge/ChainBridge/e2e/ethereum"
	sub "github.com/selendra/selendra-bridge/ChainBridge/e2e/substrate"
	"github.com/selendra/selendra-bridge/ChainBridge/shared"
	ethutils "github.com/selendra/selendra-bridge/ChainBridge/shared/ethereum"
	ethtest "github.com/selendra/selendra-bridge/ChainBridge/shared/ethereum/testing"
	subutils "github.com/selendra/selendra-bridge/ChainBridge/shared/substrate"
	subtest "github.com/selendra/selendra-bridge/ChainBridge/shared/substrate/testing"
	"github.com/selendra/selendra-bridge/chainbridge-utils/core"
	"github.com/selendra/selendra-bridge/chainbridge-utils/msg"
	log "github.com/ChainSafe/log15"
	"github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3/types"
	"github.com/ethereum/go-ethereum/common"
)

const EthAChainId = msg.ChainId(0)
const SubChainId = msg.ChainId(1)
const EthBChainId = msg.ChainId(2)

var logFiles = []string{}

type test struct {
	name string
	fn   func(*testing.T, *testContext)
}

var tests = []test{
	{"Erc20ToSubstrate", testErc20ToSubstrate},
	{"SubstrateToErc20", testSubstrateToErc20},
	{"Erc20toErc20", testErc20ToErc20},
	{"Erc20 to Substrate Round Trip", testErc20SubstrateRoundTrip},

	{"Three chain with parallel submission", testThreeChainsParallel},
}

type testContext struct {
	ethA      *eth.TestContext
	ethB      *eth.TestContext
	subClient *subutils.Client

	EthSubErc20ResourceId  msg.ResourceId
	EthEthErc20ResourceId  msg.ResourceId
}

func createAndStartBridge(t *testing.T, name string, contractsA, contractsB *ethutils.DeployedContracts) (*core.Core, log.Logger) {
	// Create logger to write to a file, and store the log file name in global var
	logger := log.Root().New()
	sysErr := make(chan error)
	ethACfg := eth.CreateConfig(name, EthAChainId, contractsA, eth.EthAEndpoint)
	ethA, err := ethChain.InitializeChain(ethACfg, logger.New("relayer", name, "chain", "ethA"), sysErr, nil)
	if err != nil {
		t.Fatal(err)
	}

	subCfg := sub.CreateConfig(name, SubChainId)
	subA, err := subChain.InitializeChain(subCfg, logger.New("relayer", name, "chain", "sub"), sysErr, nil)
	if err != nil {
		t.Fatal(err)
	}

	ethBCfg := eth.CreateConfig(name, EthBChainId, contractsB, eth.EthBEndpoint)
	ethB, err := ethChain.InitializeChain(ethBCfg, logger.New("relayer", name, "chain", "ethB"), sysErr, nil)
	if err != nil {
		t.Fatal(err)
	}

	bridge := core.NewCore(sysErr)
	bridge.AddChain(ethA)
	bridge.AddChain(subA)
	bridge.AddChain(ethB)

	err = ethA.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = subA.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = ethB.Start()
	if err != nil {
		t.Fatal(err)
	}

	return bridge, logger
}

func attemptToPrintLogs() {
	for _, fileName := range logFiles {
		dat, err := ioutil.ReadFile(fileName)
		if err != nil {
			continue
		}
		name := strings.Split(fileName, ".")[0]
		fmt.Printf("\n\tOutput from %s:\n\n", name)
		fmt.Print(string(dat))
		os.Remove(fileName)
	}
}

func assertChanError(t *testing.T, errs <-chan error) {
	select {
	case err := <-errs:
		t.Fatalf("BridgeA Fatal Error: %s", err)
	default:
		// Do nothing
		fmt.Println("No errors here!")
		return
	}
}

func setupFungibleTests(t *testing.T, ctx *testContext) {
	mintAmount := big.NewInt(1000)
	approveAmount := big.NewInt(500)

	// Deploy Sub<>Eth erc20 on ethA, register resource ID, add handler as minter
	erc20ContractASub := ethtest.Erc20DeployMint(t, ctx.ethA.Client, mintAmount)
	log.Info("Deployed erc20 contract", "address", erc20ContractASub)
	ethtest.Erc20Approve(t, ctx.ethA.Client, erc20ContractASub, ctx.ethA.BaseContracts.ERC20HandlerAddress, approveAmount)
	ethtest.RegisterResource(t, ctx.ethA.Client, ctx.ethA.BaseContracts.BridgeAddress, ctx.ethA.BaseContracts.ERC20HandlerAddress, ctx.EthSubErc20ResourceId, erc20ContractASub)
	ethtest.SetBurnable(t, ctx.ethA.Client, ctx.ethA.BaseContracts.BridgeAddress, ctx.ethA.BaseContracts.ERC20HandlerAddress, erc20ContractASub)

	// Deploy Eth<>Eth erc20 on ethA, register resource ID, add handler as minter
	erc20ContractAEth := ethtest.Erc20DeployMint(t, ctx.ethA.Client, mintAmount)
	log.Info("Deployed erc20 contract", "address", erc20ContractAEth)
	ethtest.Erc20Approve(t, ctx.ethA.Client, erc20ContractAEth, ctx.ethA.BaseContracts.ERC20HandlerAddress, approveAmount)
	ethErc20ResourceId := msg.ResourceIdFromSlice(append(common.LeftPadBytes(erc20ContractAEth.Bytes(), 31), 0))
	ethtest.RegisterResource(t, ctx.ethA.Client, ctx.ethA.BaseContracts.BridgeAddress, ctx.ethA.BaseContracts.ERC20HandlerAddress, ethErc20ResourceId, erc20ContractAEth)
	ethtest.SetBurnable(t, ctx.ethA.Client, ctx.ethA.BaseContracts.BridgeAddress, ctx.ethA.BaseContracts.ERC20HandlerAddress, erc20ContractAEth)

	// Deploy Eth<>Eth erc20 on ethB, add handler as minter
	erc20ContractBEth := ethtest.Erc20DeployMint(t, ctx.ethB.Client, mintAmount)
	log.Info("Deployed erc20 contract", "address", erc20ContractBEth)
	ethtest.Erc20Approve(t, ctx.ethB.Client, erc20ContractBEth, ctx.ethB.BaseContracts.ERC20HandlerAddress, approveAmount)
	ethtest.RegisterResource(t, ctx.ethB.Client, ctx.ethB.BaseContracts.BridgeAddress, ctx.ethB.BaseContracts.ERC20HandlerAddress, ethErc20ResourceId, erc20ContractBEth)
	ethtest.SetBurnable(t, ctx.ethB.Client, ctx.ethB.BaseContracts.BridgeAddress, ctx.ethB.BaseContracts.ERC20HandlerAddress, erc20ContractBEth)

	ctx.ethA.TestContracts.Erc20Sub = erc20ContractASub
	ctx.ethA.TestContracts.Erc20Eth = erc20ContractAEth
	ctx.ethB.TestContracts.Erc20Eth = erc20ContractBEth
	ctx.EthEthErc20ResourceId = ethErc20ResourceId
}

// This tests three relayers connected to three chains (2 ethereum, 1 substrate).
//
// EthA:
//  - Native erc20 token
// Eth B:
//  - Synthetic erc20 token
// Substrate:
//  - Synthetic token (native to chain)
//
func Test_ThreeRelayers(t *testing.T) {
	shared.SetLogger(log.LvlTrace)
	threshold := 3

	// Setup test client connections for each chain
	ethClientA := ethtest.NewClient(t, eth.EthAEndpoint, eth.AliceKp)
	ethClientB := ethtest.NewClient(t, eth.EthBEndpoint, eth.AliceKp)
	subClient := subtest.CreateClient(t, sub.AliceKp.AsKeyringPair(), sub.TestSubEndpoint)

	// First lookup the substrate resource IDs
	var rawRId types.Bytes32
	subtest.QueryConst(t, subClient, "Example", "NativeTokenId", &rawRId)
	subErc20ResourceId := msg.ResourceIdFromSlice(rawRId[:])

	// Base setup for ethA
	contractsA := eth.DeployTestContracts(t, ethClientA, eth.EthAEndpoint, EthAChainId, big.NewInt(int64(threshold)))
	// Base setup for ethB ERC20 - handler can mint
	contractsB := eth.DeployTestContracts(t, ethClientB, eth.EthBEndpoint, EthBChainId, big.NewInt(int64(threshold)))

	// Create the initial context, added to in setup functions
	ctx := &testContext{
		ethA: &eth.TestContext{
			BaseContracts: contractsA,
			TestContracts: eth.TestContracts{},
			Client:        ethClientA,
		},
		ethB: &eth.TestContext{
			BaseContracts: contractsB,
			TestContracts: eth.TestContracts{},
			Client:        ethClientB,
		},
		subClient:              subClient,
		EthSubErc20ResourceId:  subErc20ResourceId,
	}

	setupFungibleTests(t, ctx)

	// Setup substrate client, register resource, add relayers
	resources := map[msg.ResourceId]subutils.Method{
		subErc20ResourceId:    subutils.ExampleTransferMethod,
	}
	subtest.EnsureInitializedChain(t, subClient, sub.RelayerSet, []msg.ChainId{EthAChainId}, resources, uint32(threshold))

	// Create and start three bridges with both chains
	bridgeA, bobLog := createAndStartBridge(t, "bob", contractsA, contractsB)
	bridgeB, charlieLog := createAndStartBridge(t, "charlie", contractsA, contractsB)
	bridgeC, aliceLog := createAndStartBridge(t, "dave", contractsA, contractsB)

	assertChanError(t, bridgeA.Errors())
	assertChanError(t, bridgeB.Errors())
	assertChanError(t, bridgeC.Errors())

	for _, tt := range tests {
		tt := tt

		// Swap handler
		fileName := "bob" + ".output"
		logFiles = append(logFiles, fileName)
		bobLog.SetHandler(log.Must.FileHandler(fileName, log.TerminalFormat()))

		// Swap handler
		fileName = "charlie" + ".output"
		logFiles = append(logFiles, fileName)
		charlieLog.SetHandler(log.Must.FileHandler(fileName, log.TerminalFormat()))

		// Swap handler
		fileName = "dave" + ".output"
		logFiles = append(logFiles, fileName)
		aliceLog.SetHandler(log.Must.FileHandler(fileName, log.TerminalFormat()))

		t.Run(tt.name, func(t *testing.T) {
			tt.fn(t, ctx)

			assertChanError(t, bridgeA.Errors())
			assertChanError(t, bridgeB.Errors())
			assertChanError(t, bridgeC.Errors())
		})
		// Flush logs
		attemptToPrintLogs()
	}
}
