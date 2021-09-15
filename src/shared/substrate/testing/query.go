// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package subtest

import (
	"bytes"
	"math/big"
	"testing"

	utils "github.com/selendra/selendra-bridge/ChainBridge/shared/substrate"
	"github.com/ChainSafe/log15"
	"github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3/types"
)

func QueryStorage(t *testing.T, client *utils.Client, prefix, method string, arg1, arg2 []byte, result interface{}) bool {
	exists, err := utils.QueryStorage(client, prefix, method, arg1, arg2, result)
	if err != nil {
		t.Fatal(err)
	}
	return exists
}

func QueryConst(t *testing.T, client *utils.Client, prefix, name string, res interface{}) {
	err := utils.QueryConst(client, prefix, name, res)
	if err != nil {
		t.Fatal(err)
	}
}
