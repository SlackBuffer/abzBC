# Quick tour
- Generate network artifacts

    ```bash
    # fabric-samples/first-network
    ./byfn.sh generate
    ``` 

    - Generates all of the **certificates** and **keys** for various network entities, the **genesis block** used to bootstrap the ordering service, and a collection of **configuration transactions** required to configure a channel
- Bring up the network

    ```bash
    # fabric-samples/first-network
    ./byfn.sh up
    # ./byfn.sh up -l node -o etcdraft
    ```

    - Compile Golang chaincode images and spin up the corresponding containers
    - Only one language can be tried unless you bring down and recreate the network
- Bring down the network

    ```bash
    # fabric-samples/first-network
    ./byfn.sh down
    # ./byfn.sh up -l node -o etcdraft
    ```
# Manual steps
- The manual steps assume that the `FABRIC_LOGGING_SPEC` in the `cli` container is set to `DEBUG`
    - Set this by modifying the `docker-compose-cli.yaml` file in the `first-network` directory
## Crypto generator
- `cryptogen` tool will be used to generate the cryptographic material (**x.509 certs** and **signing keys**) for our various network entities
    - These certificates are representative of identities, and they allow for sign/verify authentication to take place as entities communicate and transact
- `cryptogen` consumes `crypto-config.yaml` that contains **network topology** and allows us to generate a set of certificates and keys for both the organizations and the components that belong to those organizations
- Each organization is provisioned a unique root certificate (`ca-cert`) that binds specific components (peers and orderers) to that organization
    - By assigning each organization a unique CA certificate, we are mimicking a typical network where a participating member would use its own certificate authority
- Transactions and communications within Hyperledger Fabric are signed by an entity’s private key (`keystore`), and then verified by means of a public key (`signcerts`)
- `count` variable in `crypto-config.yaml` is used to specify the number of peers per organization
- The generated certificates and keys will be saved to a folder titled `crypto-config`
- The `crypto-config.yaml` file lists five orderers as being tied to the orderer organization. While the `cryptogen` tool will create certificates for all five of these orderers, unless the `Raft` or `Kafka` ordering services are being used, only one of these orderers will be used in a `Solo` ordering service implementation and be used to create the system channel and `mychannel`
## Configuration transaction generator
- The `configtxgen` tool is used to create 4 configuration artifacts
    - orderer genesis block,
    - channel configuration transaction,
    - two anchor peer transactions - one for each peer org
- The orderer block is the genesis block for the ordering service, and the channel configuration transaction file is broadcast to the orderer at channel creation time. The anchor peer transactions specify each org’s anchor peer on this channel
- `Configtxgen` consumes `configtx.yaml` that contains the definitions for the sample network
- There are 3 members - one orderer org (`OrdererOrg`) and 2 peer orgs (`Org1`, `Org2`) each managing and maintaining two peer nodes
- This file also specifies a consortium - `SampleConsortium` - consisting of the 2 peer orgs
    - `SampleConsortium` is defined in the system-level profile and then referenced by the channel-level profile
    - Channels exist within the purview of a consortium, and all consortia must be defined in the scope of the network at large
- `Profiles` section
    - `TwoOrgsOrdererGenesis`: generates the genesis block for a Solo ordering service.
    - `SampleMultiNodeEtcdRaft`: generates the genesis block for a Raft ordering service. Only used if you issue the `-o` flag and specify `etcdraft`
    - `SampleDevModeKafka`: generates the genesis block for a Kafka ordering service. Only used if you issue the `-o` flag and specify `kafka`
    - `TwoOrgsChannel`: generates the genesis block for our channel, `mychannel`
## Run the tools
- Generate the certificates that will be used for the network configuration

    ```bash
    # fabric-samples/first-network
    ../bin/cryptogen generate --config=./crypto-config.yaml
    # output dir: `fabric-samples/first-network/crypto-config`
    ```

    - The certs and keys (i.e. the MSP material) will be output into a directory - `crypto-config` - at the root of the first-network directory
- Tell the `configtxgen` where to look for the `configtx.yaml` file that it needs to ingest

    ```bash
    export FABRIC_CFG_PATH=$PWD
    ```

- Create the orderer genesis block

    ```bash
    # fabric-samples/first-network
    # `channelID` is the name of the system channe
    ../bin/configtxgen -profile TwoOrgsOrdererGenesis -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block
    # output dir: `fabric-samples/first-network/channel-artifacts`

    # output a genesis block for a Raft ordering service
    # ../bin/configtxgen -profile SampleMultiNodeEtcdRaft -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block

    # output a genesis block for a Kafka ordering service
    # ../bin/configtxgen -profile SampleDevModeKafka -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block
    ```

- Create the channel transaction artifact

    ```bash
    # The channel.tx artifact contains the definitions for our sample channel
    # fabric-samples/first-network
    export CHANNEL_NAME=mychannel  && ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
    ```

    - You don’t have to issue a special command for the channel if you are using a Raft or Kafka ordering service. The `TwoOrgsChannel` profile will use the ordering service configuration you specified when creating the genesis block for the network
- Define the anchor peer for `Org1` and `Org2` on the channel

    ```bash
    # fabric-samples/first-network
    ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
    ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
    ```

