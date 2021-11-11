module github.com/selendra/selendra-bridge/chainbridge-substrate-events

go 1.13

replace github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3 => ../go-substrate-rpc-client

require (
	github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3 v3.0.0-00010101000000-000000000000
	github.com/vedhavyas/go-subkey v1.0.2
)
