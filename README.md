# Builder Boost
![tests](https://github.com/primevprotocol/builder-boost/actions/workflows/tests.yml/badge.svg?branch=main)

Builder Boost is sidecar software that facilitates a blockbuilder or sequencer's participation in the Primev network. The diagram below outlines the module's position within a builder's local build environment.

### Builder Boost Diagram

<img src="./diagrams/bb-vhighlevel.png" width="50%"/>

We recommend you run builder boost as an independent cloud instance if you're operations use a cloud provider.

### Pre Requistes
- golang version 1.20.4+
- updated builder or execution layer with primev relay enabled
- git

# Setting Up Builder Boost

## Step 1. Ensure Builder is able to send payloads to builder boost
Begin by reviewing the Builder Reference Implementation instructions provided [here](https://hackmd.io/wmmCgKJdTom9WXht2PcdLA) for basic geth modifications required (~5 minutes). You can proceed to set up Builder Boost to receive block templates from your builder or execution layer after the modifications are in place.


## Step 2. Download Boost via Github
Do a git clone:
``` shell
$ git@github.com:primevprotocol/builder-boost.git
```

## Step 3. Build the Boost instance
``` shell
 $ cd builder-boost
 $ make all
```
You should now see two binaries built in the root folder:
``` 
.
├── boost
├── ...
├── searcher
...
```
## Step 4. Run boost
We highly recommend you use systemctl to run boost, the following is a recommended systemctl file that can be used and placed in your `/home/systemd/system/<filename>` folder.
``` bash
[Unit]
Description=boost service
After=network.target

[Service]
User=ec2-user
ExecStart=/home/ec2-user/builder-boost/boost
restart=always
Environment="ENV=sepolia"
Environment="ROLLUP_KEY=<builder-key>"
Environment="BUILDER_AUTH_TOKEN=<builder-x-boost-token>"

[Install]
WantedBy=multi-user.target
```

Note the `BUILDER_AUTH_TOKEN` should be set to the same value set for `--builder-auth-key` in your modified builder implementation. This token is used to authenticate your builders to boost, to ensure actors can't send fraudulent payloads.

## [optional] Run builder and searcher in docker-compose

1. Copy `.env.example` to `.env` and set required variables
2. Run `docker-compose up --build` to run builder and searcher connected to builder

# Registering your Builder or Execution Layer in the Primev Contract (on Sepolia Testnet during Pilot phase)

## Open Primev Contract

Head to [Update Builder Method](https://sepolia.etherscan.io/address/0x6e100446995f4456773Cd3e96FA201266c44d4B8#writeContract#F4) on Etherscan. This method allows you to register your execution layer or update settings for your existing registration.

## Connect your Wallet

On the top left corner of the **Write Contract** section press **Connect to Web3**. Make sure you use your builder or sequencer address (the one you use for running your execution layer / builder geth instance). This address should be funded with some ETH to cover transaction costs.

## Specify Builder Parameters and Register

Use the `updateBuilder` method to specify 2 fields to register:
- `_minimalStake`: minimal amount for a data consumer to deposit to in order to authorize a connection. We suggest 0.1 ETH as a minimum.
- `_minimalSubscriptionPeriod`: minimal subscription period in blocks given to a data consumer for depositing minimal stake. We suggest 216000 blocks, which is about 30 days worth. If the depositor provides more than the minimum amount, the subscription period will be extended in a linear fashion.

After the required fields are specified, press the **Write** button and confirm the transaction using your Web3 provider. Once the transaction is confirmed, data consumers could start connecting to your builder-boost instance. Make sure to run your builder boost instance before registering a new builder or execution client.

# Depositing to a builder using the Primev Contract as a Data consumer

## Obtain Commitment Hash from the registered Builder

Making a deposit requires specifying a commitment hash which is unique to every searcher-builder pair. To obtain this hash you need to send a HTTP request to the builder endpoint with the specified `token` parameter. The token parameter is logged when running a **Searcher** instance, make sure you run Searcher with correct `SEARCHER_KEY` and `BOOST_URL` variables.

Obtain a Commitment Hash, assuming `http://localhost:18550` is a builder boost endpoint and `ABCD` is a searcher token for this builder:
```bash
$ curl -H "X-Primev-Signature: ABCD" http://localhost:18550/commitment
```

Output will look like this:
```json
{"commitment":"0x688d0031ba0ce02c2049786ca6fd70d04869688dd84d6310b7fdb052d199612f"}
```

## Open the Primev Contract in Etherscan

Head to [Deposit Method](https://sepolia.etherscan.io/address/0x6e100446995f4456773Cd3e96FA201266c44d4B8#writeContract#F1) on Etherscan. This method allows you to deposit funds to a builder to authorize your connection.

## Connect your Wallet

On the top left corner of **Write Contract** section press **Connect to Web3**. Make sure you use the address you will transact with the builder with. The contract preserves your privacy as there's a signature mechanism in place to make sure only the builder you're registering with can see your address.

## Specify your deposit parameters and send the transaction

In the `deposit` method there are 3 fields you need to specify:
- `payableAmount`: the amount you want to deposit for the builder you're registering with
- `_builder`: the builder address you want to authorize a connection for
- `_commitment`: the commitment hash obtained from the previous step

After these fields are specified, press the **Write** button and confirm your transaction using your Web3 provider. Once the transaction is confirmed, you will be able to receive builder hints in real time.

## More Links
[Searcher Testing Guide](docs/searcher-testing.md) - This guide walks you through the setup process for connecting searcher emulators to your boost for e2e testing and monitoring.

[Builder Modifications](https://hackmd.io/wmmCgKJdTom9WXht2PcdLA) - A detailed walkthrough with example commits, to modify a builder to add the primev relay and enable boost.

