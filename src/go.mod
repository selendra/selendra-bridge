module github.com/selendra/selendra-bridge/ChainBridge

go 1.15

replace github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3 => ../go-substrate-rpc-client

require (
	github.com/ChainSafe/go-schnorrkel v0.0.0-20210527232834-58622d036665 // indirect
	github.com/ChainSafe/log15 v1.0.0
	github.com/ethereum/go-ethereum v1.10.12
	github.com/mimoo/StrobeGo v0.0.0-20210601165009-122bf33a46e0 // indirect
	github.com/prometheus/client_golang v1.4.1
	github.com/selendra/selendra-bridge/chainbridge-substrate-events v0.0.0-20210923141139-605f5971ed1d
	github.com/selendra/selendra-bridge/chainbridge-utils v0.0.0-20210923141139-605f5971ed1d
	github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3 v3.0.0-20210915054013-020acc2ab7e5
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
)
