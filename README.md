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


## More Links
[Searcher Testing Guide](docs/searcher-testing.md) - This guide walks you through the setup process for connecting searcher emulators to your boost for e2e testing and monitoring.

[Builder Modifications](https://hackmd.io/wmmCgKJdTom9WXht2PcdLA) - A detailed walkthrough with example commits, to modify a builder to add the primev relay and enable boost.

