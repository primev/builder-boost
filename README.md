# Builder Boost (in Beta)
![tests](https://github.com/primevprotocol/builder-boost/actions/workflows/tests.yml/badge.svg?branch=main)

Builder Boost is sidecar software that facilitates a blockbuilder's participation in the Primev network. The diagram below outlines the module's position within a builder's local build environment.

### Builder Boost Diagram

<img src="./diagrams/bb-vhighlevel.png" width="50%"/>

We recommend you run builder boost as an independent cloud instance if you're operations use a cloud provider.

### Pre Requistes
- golang version 1.20.4+
- updated builder with primev relay enabled
- git

# Setting Up Builder Boost

## Step 1. Ensure Builder is able to send payloads to builder boost
Begin by reviewing the Builder Reference Implementation instructions provided [here](https://hackmd.io/wmmCgKJdTom9WXht2PcdLA). After performing the necessary modifications, you can proceed to set up Builder Boost to receive block templates from your builder.


## Step 2. Download Boost via Github
For now simply do a git clone:
``` shell
$ git@github.com:primevprotocol/builder-boost.git
```

## Step 3. Build the Boost instance
``` shell
 $ cd builder-boost
 $ make all
```
You should now see two bianaries built in the root folder:
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

Note, the `BUILDER_AUTH_TOKEN` should be set to the same value set for `--builder-auth-key` in your modified builder implementation. This token is used to authenticate your builders to boost, to ensure no threat actor can send false payloads.

## Running builder and searcher in docker-compose

1. Copy `.env.example` to `.env` and set required variables
2. Run `docker-compose up --build` to run builder and searcher connected to builder

# Registering Builder in Primev Contract

## Open Primev Contract in Etherscan

Head to [Update Builder Method](https://sepolia.etherscan.io/address/0x6e100446995f4456773Cd3e96FA201266c44d4B8#writeContract#F4) in Etherscan. This method allows you to register new builder or update settings for existing builder.

## Connect Web3 Wallet

On top left corner of **Write Contract** section press **Connect to Web3**. Make sure you use builder address (the one you use for running builder geth instance). This address should be funded with some Sepolia ETH to cover transaction costs.

## Specify Builder Parameters and Register


In `updateBuilder` method there is 2 fields to specify:
- `_minimalStake`: minimal amount for searcher to deposit to this builder in order to start receiving builder hints
- `_minimalSubscriptionPeriod`: minimal subscription period given to searcher for depositing minimal stake. In case searcher deposits more than minimal stake, subscription period will be extended linearly.

After all required fields specified, press **Write** button and confirm transaction using your Web3 provider. Once transaction is confirmed, searchers could start depositing to your builder. Make sure you run builder boost instance before registering new builder.

# Depositing to Builder in Primev Contract as a Searcher

## Open Primev Contract in Etherscan

Head to [Deposit Method](https://sepolia.etherscan.io/address/0x6e100446995f4456773Cd3e96FA201266c44d4B8#writeContract#F1) in Etherscan. This method allows you to deposit stake to some particular builder on behalf of searcher.

## Obtain Commitment Hash from Builder

Making deposit requires specifying commitment hash which is unique to every searcher-builder pair. To obtain this hash you need to send a HTTP request to builder endpoint with specified `token` parameter. Token parameter is logged when running **Searcher** instance, make sure you run Searcher with correct `SEARCHER_KEY` and `BOOST_ADDR` variables.

Obtain Commitment Hash, assuming `http://localhost:18550` is a builder boost endpoint and `ABCD` is a searcher token for this builder:
```bash
$ curl http://localhost:18550/commitment?token=ABCD
```

Output will look like this:
```json
{"commitment":"0x688d0031ba0ce02c2049786ca6fd70d04869688dd84d6310b7fdb052d199612f"}
```

## Specify Deposit Parameters and Send Transaction

In `deposit` method there is 3 fields to specify:
- `payableAmount`: the amount of stake you wish to deposit for particular builder
- `_builder`: the builder address
- `_commitment`: the commitment hash obtained on previous step

After all required fields specified, press **Write** button and confirm transaction using your Web3 provider. Once transaction is confirmed, searcher will be able to receive builder hints.

## More Links
[Searcher Testing Guide](docs/searcher-testing.md) - This guide walks you through the setup process for connecting searcher emulators to your boost for e2e testing and monitoring.

[Builder Modifications](https://hackmd.io/wmmCgKJdTom9WXht2PcdLA) - A detailed walkthrough with example commits, to modify a builder to add the primev relay and enable boost.

