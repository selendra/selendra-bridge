# Selendra Bridge

# Contents

- [Getting Started](#getting-started)
- [Installation](#installation)
- [Configuration](#configuration)
- [Chain Implementations](#chain-implementations)

# Getting Started
- Check out our [documentation](https://chainbridge.chainsafe.io).

# Installation

## Dependencies

- [Subkey](https://substrate.dev/docs/en/knowledgebase/integrate/subkey): 
Used for substrate key management. Only required if connecting to a substrate chain.

## Building

`make build`: Builds `Selendra Bridge` in `./build`.

**or**

`make install`: Uses `go install` to add `Selendra Bridge` to your GOBIN.

# Configuration

> Note: TOML configs have been deprecated in favour of JSON

A chain configurations take this form:

```
{
    "name": "eth",                      // Human-readable name
    "type": "ethereum",                 // Chain type (eg. "ethereum" or "substrate")
    "id": "0",                          // Chain ID
    "endpoint": "ws://<host>:<port>",   // Node endpoint
    "from": "0xff93...",                // On-chain address of relayer
    "opts": {},                         // Chain-specific configuration options (see below)
}
```

See `config.json.example` for an example configuration. 

### Ethereum Options

Ethereum chains support the following additional options:

```
{
    "bridge": "0x12345...",          // Address of the bridge contract (required)
    "erc20Handler": "0x1234...",     // Address of erc20 handler (required)
    "maxGasPrice": "0x1234",         // Gas price for transactions (default: 20000000000)
    "gasLimit": "0x1234",            // Gas limit for transactions (default: 6721975)
    "gasMultiplier": "1.25",         // Multiplies the gas price by the supplied value (default: 1)
    "http": "true",                  // Whether the chain connection is ws or http (default: false)
    "startBlock": "1234",            // The block to start processing events from (default: 0)
    "blockConfirmations": "10"       // Number of blocks to wait before processing a block
    "useExtendedCall": "true"        // Extend extrinsic calls to substrate with ResourceID. Used for backward compatibility with example pallet. *Default: false*
}
```

### Substrate Options

Substrate supports the following additonal options:

```
{
    "startBlock": "1234" // The block to start processing events from (default: 0)
}
```

## Blockstore

The blockstore is used to record the last block the relayer processed, so it can pick up where it left off. 

If a `startBlock` option is provided (see [Configuration](#configuration)), then the greater of `startBlock` and the latest block in the blockstore is used at startup.

To disable loading from the blockstore specify the `--fresh` flag. A custom path for the blockstore can be provided with `--blockstore <path>`. For development, the `--latest` flag can be used to start from the current block and override any other configuration.

## Keystore

Selendra Bridge requires keys to sign and submit transactions, and to identify each bridge node on chain.

To use secure keys, see `selendrabridge accounts --help`. The keystore password can be supplied with the `KEYSTORE_PASSWORD` environment variable.

To import external ethereum keys, such as those generated with geth, use `Selendra Bridge accounts import --ethereum /path/to/key`.

To import private keys as keystores, use `selendrabridge account import --privateKey key`.

For testing purposes, Selendra Bridge provides 5 test keys. The can be used with `--testkey <name>`, where `name` is one of `Alice`, `Bob`, `Charlie`, `Dave`, or `Eve`. 

# Chain Implementations

- EVM: [Selendra Bridge-Cli](https://github.com/selendra/bridge-cli) 

    The Solidity contracts required for Selendra Bridge. Includes deployment and interaction CLI.
    
    The bindings for the contracts live in `bindings/`. To update the bindings modify `scripts/setup-contracts.sh` and then run `make clean && make setup-contracts`

- Selendra: [Selendra Chain](https://github.com/selendra/selendra-chain)

    A substrate pallet that can be integrated into a chain.


## License

Selendra-Bridge is implement from [ChainBridge](https://github.com/ChainSafe/ChainBridge) under license [GPL 3.0 licensed](LICENSE-GPL3).