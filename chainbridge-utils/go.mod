module github.com/selendra/selendra-bridge/chainbridge-utils

go 1.13

replace github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3 => ../go-substrate-rpc-client

require (
	github.com/ChainSafe/log15 v1.0.0
	github.com/ethereum/go-ethereum v1.10.11
	github.com/prometheus/client_golang v1.4.1
	github.com/prometheus/procfs v0.0.10 // indirect
	github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3 v3.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)
