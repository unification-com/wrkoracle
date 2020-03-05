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

WRKOracle can then be initialised with default values by running:

```bash
wrkoracle init
```

This will create a skeleton configuration file in `$HOME/.und_wrkoracle/config/config.toml` as
follows:

```toml
broadcast-mode = "block"
chain-id = ""
frequency = "60"
from = ""
hash1 = false
hash2 = false
hash3 = false
indent = true
keyring-backend = "os"
mainchain-rest = ""
node = ""
output = "json"
parent-hash = false
trust-node = false
wrkchain-id = ""
wrkchain-rpc = ""
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
- `hash1`, `hash2`, `hash3`: boolean values whether or not to submit any of the three optional
hashes to Mainchain. See section **Hash mapping** below. **Required**
- `parent-hash`: whether or not to optionally submit the WRKChain block header parent hash. **Required**
- `mainchain-rest`: The REST server for Mainchain, e.g. https://rest-testnet.unification.io. **Required**
- `node`: Mainchain node to broadcast Txs to, e.g. `tcp://localhost:26656` if you are running you
own local full Mainchain node. **Required**
- `trust-node`: Trust connected full node (don't verify proofs for responses). **Required**
- `wrkchain-id`: The integer ID of your WRKChain, as given when the WRKChain was registered on Mainchain. **Required**
- `wrkchain-rpc`: The RPC node where WRKOracle can query your WRKChain, e.g. http://127.0.0.1:7545. **Required**

## Running in automated mode

Once your WRKOracle has been configured with the options outlined above, it can run automatically
and poll your WRKChain according to the defined frequency, submitting the latest WRKChain block header
hashes to Mainchain:

```bash
wrkoracle run
```

WRKOracle will output its status as follows:

```bash
starting 2020-03-04 14:43:41.580232302 +0000 GMT
polling WRKChain for latest block
Get block for WRKChain 'wrkchain1', type 'geth' at http://127.0.0.1:7545
Got WRKChain block
WRKChain Height: 61
WRKChain Block Hash: 0x81a23fe7c73711260bebae8b027a6d0897f305407fdb13e6ed0a4effdd2d6e74
WRKChain Parent Hash: 0xe3cdf99657d23461ee590477eb1aec1873b52824910e5369964039757defffdc
WRKChain Hash1: 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
WRKChain Hash2: 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
WRKChain Hash3: 0xe534373922c7a6e23ba0a96bcdf53ade135ad7117b27cc9b0706c3e02360e653
recording latest WRKChain block
Generate msg
Broadcasting Tx and waiting for response...
WRKChain header hash recording fee: 1000000000nund
gas estimate: 140275
Tx Hash: D1F138682C7CB49A67E777CF40B9E242D5C89F754781BF6EDAE78918B3080996
Success! Recorded in Mainchain Block #2737
Gas used: 92202
Done. Next poll due at 2020-03-04 14:44:41.580299426 +0000 GMT
-----------------------------------
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
getting WRKChain 'wrkchain1' block 2424 and recording
Get block for WRKChain 'wrkchain1', type 'geth' at http://127.0.0.1:7545
Got WRKChain block
WRKChain Height: 2424
WRKChain Block Hash: 0x4fee5d6dd69b21b37c0923d1c1ded45ace4c94af3d1f18a423ea2e25052c25d6
WRKChain Parent Hash: 0x8ae2443997e24ec247116efe275af3cae7bbe1a62071cf52c43cd0e233fac551
WRKChain Hash1: 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
WRKChain Hash2: 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
WRKChain Hash3: 0x7cab995324db9a8556416274bd3367ec385bd3c8643e95052eab60a8e9537681
Generate msg
Broadcasting Tx and waiting for response...
WRKChain header hash recording fee: 1000000000nund
gas estimate: 140275
Tx Hash: 4EAFE4B59198AB4F34A8167FA83042AA1D501A1F2B3A192B8F5496F2DE92E0A3
Success! Recorded in Mainchain Block #2789
Gas used: 92202
```

## Hash mapping

The `Hash1`, `Hash2` and `Hash3` values that can be submitted to Mainchain are optional, and are
automatically mapped by WRKOracle to different hash values, depending on the WRKChain type:

### `geth` based chains

`Hash1`: Merkle root hash for the Receipts: `Header.ReceiptHash`  
`Hash2`: Merkle root hash for the Tx: `Header.TxHash`  
`Hash3`: Merkle root hash for Root: `Header.Root`  
