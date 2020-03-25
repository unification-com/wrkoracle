![Unification](https://raw.githubusercontent.com/unification-com/wrkoracle/master/unification_logoblack.png "Unification")

[![Go Report Card](https://goreportcard.com/badge/github.com/unification-com/wrkoracle)](https://goreportcard.com/report/github.com/unification-com/wrkoracle) [![Join the chat at https://gitter.im/unification-com/wrkoracle](https://badges.gitter.im/unification-com/wrkoracle.svg)](https://gitter.im/unification-com/wrkoracle?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

# WRKOracle

The official Unification WRKChain Oracle software for recording WRKChain block header hashes 
to Mainchain

## Build and installation

### Prerequisites

Go **1.13+** is required to install the WRKOracle binary

Install go by following the [official docs](https://golang.org/doc/install).
Once Go is installed, set your `$PATH` environment variable:

```bash
$ mkdir -p $HOME/go/bin
$ echo "export PATH=$PATH:$(go env GOPATH)/bin">>$HOME/.bash_profile
$ source $HOME/.bash_profile
```

### Build

The `build` Make target can be used to build. The binary will be output to `./build/wrkoracle`:

```bash
make build
```

### Install

Use:

```bash
make install
```

to install the `wrkoracle` binary into your `$GOPATH/bin`

Run:

```bash
wrkoracle version --long
```

to verify it has installed correctly

## Initialisation

First, you will need to import a valid Mainchain account key into the `WrkOracle` keyring. This
can be either a new key:

```bash
wrkoracle keys add my_wrkoracle_acc
```

or by importing a key from an existing Mnemonic:

```bash
wrkoracle keys add my_wrkoracle_acc --recover
```

In either case, the account will need sufficient UND to run the Oracle and submit hashes to
Mainchain.

WRKOracle must then be initialised with default values by running:

```bash
wrkoracle init [wrkchain_type]
```

E.g.

```bash
wrkoracle init geth
```

This will create a skeleton configuration file in `$HOME/.und_wrkoracle/config/config.toml` as
follows:

```toml
broadcast-mode = "block"
chain-id = ""
frequency = "60"
from = ""
hash1 = "ReceiptsRoot"
hash2 = "TxRoot"
hash3 = "StateRoot"
indent = true
keyring-backend = "os"
mainchain-rest = ""
node = ""
output = "json"
parent-hash = true
trust-node = false
wrkchain-id = ""
wrkchain-rpc = ""
wrkchain-type = "geth"
```

## Configuration options

The configuration values can be set in `$HOME/.und_wrkoracle/config/config.toml`, or passed to
the binary at runtime as `--flags` (e.g. `--chain-id`).

- `broadcast-mode`: should remain as `block`, so that `wrkoracle` waits for the Tx to be processed in 
a Mainchain block. **Required**
- `chain-id`: The chain ID of Mainchain hashes are being submitted to, e.g. `UND-Mainchain-DevNet`, 
`UND-Mainchain-TestNet`, or `UND-Mainchain-MainNet` **Required**
- `frequency`: frequency in seconds that the WRKOracle should poll your WRKChain for the latest
block header and submit the hashes to Mainchain. **Required**
- `from`: default account that should be used by WRKOracle to sign the transactions, as named when
importing the account above, e.g. `my_wrkoracle_acc`. **Required**
- `hash1`, `hash2`, `hash3`: optional values mapped to various header hashes, depending on the WRKChain type
hashes to Mainchain. See section **Hash mapping** below. If left empty, no value will be submitted.
- `parent-hash`: whether or not to optionally submit the WRKChain block header parent hash. **Required**
- `mainchain-rest`: The REST server for Mainchain, e.g. https://rest-testnet.unification.io. **Required**
- `node`: Mainchain node to broadcast Txs to, e.g. `tcp://localhost:26656` if you are running you
own local full Mainchain node. **Required**
- `trust-node`: Trust connected full node (don't verify proofs for responses). **Required**
- `wrkchain-id`: The integer ID of your WRKChain, as given when the WRKChain was registered on Mainchain. **Required**
- `wrkchain-rpc`: The RPC node where WRKOracle can query your WRKChain, e.g. `http://127.0.0.1:7545`, 
`http://172.25.0.3:26661` etc.. **Required**

## Running in automated mode

Once your WRKOracle has been configured with the options outlined above, it can run automatically
and poll your WRKChain according to the defined frequency, submitting the latest WRKChain block header
hashes to Mainchain:

```bash
wrkoracle run
```

WRKOracle will output its status as follows:

```bash
I[2020-03-11|12:55:36.872] Check WRKChain metadata                      pkg=mainchain
I[2020-03-11|12:55:36.880] Start running WRKOracle                      pkg=oracle
I[2020-03-11|12:55:36.880] start poll                                   pkg=oracle time=2020-03-11T12:55:36.880866909Z
I[2020-03-11|12:55:36.880] polling WRKChain for latest block            pkg=oracle
I[2020-03-11|12:55:36.881] Get block for WRKChain                       pkg=wrkchains moniker=wrkchain1 type=geth rpc=http://127.0.0.1:7545
I[2020-03-11|12:55:36.893] Got WRKChain block                           pkg=wrkchains
I[2020-03-11|12:55:36.893] WRKChain Height                              pkg=wrkchains height=13289
I[2020-03-11|12:55:36.893] WRKChain Block Hash                          pkg=wrkchains blockhash=0x349723088c3fa5a8871c31e256cd2a8ff5e1c19d75c5e76d48b85e28c1038f0d
I[2020-03-11|12:55:36.893] WRKChain Parent Hash                         pkg=wrkchains parenthash=0xb9673407d6ee07ccbb3d8f7808666b879ec79a1e423d6cdae05c486223a4fc00
I[2020-03-11|12:55:36.893] WRKChain Hash1                               pkg=wrkchains ref=ReceiptsRoot value=0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
I[2020-03-11|12:55:36.893] WRKChain Hash2                               pkg=wrkchains ref=TxRoot value=0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
I[2020-03-11|12:55:36.893] WRKChain Hash3                               pkg=wrkchains ref=StateRoot value=0xcf4a880479b49b6439dd833b27a387d58aff0d97df68d1c31e40d9bc2f814b7d
I[2020-03-11|12:55:36.897] recording latest WRKChain block              pkg=oracle
I[2020-03-11|12:55:36.897] Generate msg                                 pkg=mainchain
I[2020-03-11|12:55:36.897] Broadcasting Tx and waiting for response...  pkg=mainchain
I[2020-03-11|12:55:36.897] WRKChain header hash recording fee           pkg=mainchain fee=1000000000nund
I[2020-03-11|12:55:36.911] gas estimate: 147430                         pkg=mainchain
I[2020-03-11|12:55:41.880] Tx broadcast                                 pkg=mainchain hash=3C0AFAFC23E061677391E31495658F5388F1A6815A99FF2A680C682398FFABDB
I[2020-03-11|12:55:41.880] Success! Recorded in Mainchain Block         pkg=mainchain height=2531
I[2020-03-11|12:55:41.881] Gas used:                                    pkg=mainchain gas=96972
I[2020-03-11|12:55:41.881] end poll. Next poll due:                     pkg=oracle due=2020-03-11T12:55:46.880958102Z
I[2020-03-11|12:55:41.881] -----------------------------------          pkg=oracle
```

## Submitting single block headers

Individual WRKChain block headers can be submitted manually. This is useful if you wish to submit
historical data to Mainchain:

```bash
wrkoracle record [height]
```

E.g.

```bash
wrkoracle record 2424
```

The result will be output:

```bash
I[2020-03-11|12:56:32.424] Check WRKChain metadata                      pkg=mainchain
I[2020-03-11|12:56:32.427] getting requested WRKChain block header and recording pkg=oracle moniker=wrkchain1 height=12222
I[2020-03-11|12:56:32.427] Get block for WRKChain                       pkg=wrkchains moniker=wrkchain1 type=geth rpc=http://127.0.0.1:7545
I[2020-03-11|12:56:32.444] Got WRKChain block                           pkg=wrkchains
I[2020-03-11|12:56:32.444] WRKChain Height                              pkg=wrkchains height=12222
I[2020-03-11|12:56:32.444] WRKChain Block Hash                          pkg=wrkchains blockhash=0x30e0ddcb301abc5e4312a3f84d4b8dc184d47b88d21d75e2247b2bee4affb824
I[2020-03-11|12:56:32.444] WRKChain Parent Hash                         pkg=wrkchains parenthash=0x2cb59e0daeaac6372f7205e8dd32a6e9bf64ccadd7a1dc85823e9cc6d9a04a08
I[2020-03-11|12:56:32.444] WRKChain Hash1                               pkg=wrkchains ref=ReceiptsRoot value=0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
I[2020-03-11|12:56:32.444] WRKChain Hash2                               pkg=wrkchains ref=TxRoot value=0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
I[2020-03-11|12:56:32.444] WRKChain Hash3                               pkg=wrkchains ref=StateRoot value=0xea0d716b09464a350b67f13001e67ab04518a58ba7f745f67078693698a47e34
I[2020-03-11|12:56:32.446] Generate msg                                 pkg=mainchain
I[2020-03-11|12:56:32.446] Broadcasting Tx and waiting for response...  pkg=mainchain
I[2020-03-11|12:56:32.446] WRKChain header hash recording fee           pkg=mainchain fee=1000000000nund
I[2020-03-11|12:56:32.451] gas estimate: 140335                         pkg=mainchain
I[2020-03-11|12:56:35.492] Tx broadcast                                 pkg=mainchain hash=1F0A5DD144D7B6F616552E5461F81769674700B9707D22EF380A9447E0D63881
I[2020-03-11|12:56:35.492] Success! Recorded in Mainchain Block         pkg=mainchain height=2541
I[2020-03-11|12:56:35.492] Gas used:                                    pkg=mainchain gas=92242
```

## Hash mapping

The `Hash1`, `Hash2` and `Hash3` are optional values that can be submitted to Mainchain, and are
initially mapped by WRKOracle during initialisation to some default values, depending on 
the WRKChain type.

The mapping can be configured in `$HOME/.und_wrkoracle/config/config.toml` by setting the
corresponding entries for `hash1`, `hash2` and `hash3`. Leaving the entries empty will result
in the hashes being omitted from the WRKChain hash submission. The initialised defaults for 
each chain type are listed below.

**Note**: _Neither Mainchain or WRKOracle currently keep any internal records regarding what
is mapped onto the `hash1`, `hash2` and `hash3` values at submission time. It is up to the 
WRKChain and WRKOracle operators to keep track of this information externally so that any 
validation process can compare the correct hash values. This is especially important if the 
hash mapping is changed at any point during the life of the WRKChain. For example, from 
WRKChain block height 100,000, the operator may wish to change `hash1` submissions from 
`ReceiptsRoot` to `UncleHash`, in which case the operator should keep a record of this 
and update any validation processes accordingly._

### `geth` based chains

For `geth` based WRKChains, WRKOracle supports the following 5 optional [block header hashes](https://github.com/ethereum/go-ethereum/blob/master/core/types/block.go#L70) 
to be submitted:

1. `ReceiptsRoot` - Merkle root hash for the Receipts (`Header.ReceiptHash`)
2. `TxRoot` - Merkle root hash for the Tx (`Header.TxHash`)
3. `StateRoot` - Merkle root hash for State Root (`Header.Root`)
4. `UncleHash` - Uncle Hash (`Header.UncleHash`)
5. `MixHash` - Mix Digest hash (`Header.MixDigest`)

By default during initialisation, WRKOracle maps them as follows:

`hash1` = `ReceiptsRoot`  
`hash2` = `TxRoot`  
`hash3` = `StateRoot`

### `tendermint` / `cosmos` based chains

For `tendermint` and `cosmos` based WRKChains, WRKOracle supports the following 8 optional
[block header hashes](https://github.com/tendermint/tendermint/blob/master/types/block.go#L323) to be submitted:

1. `Block.Header.DataHash` - MerkleRoot of transaction hashes in this block
2. `Block.Header.AppHash` - state after txs from the previous block
3. `Block.Header.ValidatorsHash` - validators for the current block
4. `Block.Header.LastResultsHash` - root hash of all results from the txs from the previous block
5. `Block.Header.LastCommitHash` - commit from validators from the last block
6. `Block.Header.ConsensusHash` - consensus params for current block
7. `Block.Header.NextValidatorsHash` - validators for the next block
8. `Block.Header.EvidenceHash` - evidence included in the block

By default during initialisation, WRKOracle maps the following hashes:

`hash1` = `DataHash`  
`hash2` = `AppHash`  
`hash3` = `ValidatorsHash`

### `neo` based chains

For `neo` based WRKChains, WRKOracle supports the following 6 optional additional data:

1. `MerkleRoot`
2. `NextConsensus`
3. `NextBlockHash`
4. `Nonce`
5. `Script.Invocation`
6. `Script.Verification`

By default during initialisation, WRKOracle maps the following hashes:

`hash1` = `MerkleRoot`  
`hash2` = `NextConsensus`  
`hash3` = `ScriptVerification`

### `pseudochain`

`pseudochain` is a fake chain for development purposes - it will generate random hashes for each block
