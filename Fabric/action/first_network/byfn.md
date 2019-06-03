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
- The manual steps assume that the <u>`FABRIC_LOGGING_SPEC`</u> in the `cli` container is set to `DEBUG` (set this by modifying the `docker-compose-cli.yaml` file in the `first-network` directory)
## Crypto generator
- `cryptogen` tool will be used to generate the cryptographic material (**x.509 certs** and **signing keys**) for our various network entities
    - These certificates are representative of identities, and they allow for sign/verify authentication to take place as entities communicate and transact
- `cryptogen` consumes `crypto-config.yaml` that contains **network topology** and allows us to generate a set of certificates and keys for both the organizations and the components that belong to those organizations
- Each organization is provisioned a unique root certificate (`ca-cert`) that binds specific components (peers and orderers) to that organization
    - > By assigning each organization a unique CA certificate, we are mimicking a typical network where a participating member would use its own certificate authority
- Transactions and communications within Hyperledger Fabric are signed by an entity’s private key (`keystore`), and then verified by means of a public key (<u>`signcerts`</u>)
- <u>`count`</u> variable in `crypto-config.yaml` is used to specify the number of peers per organization
- The generated certificates and keys will be saved to a folder titled `crypto-config`
- The `crypto-config.yaml` file lists 5 orderers as being tied to the orderer organization. While the `cryptogen` tool will create certificates for all 5 of these orderers, unless the `Raft` or `Kafka` ordering services are being used, only one of these orderers will be used in a `Solo` ordering service implementation and be used to create the system channel and `mychannel`
## Configuration transaction generator
- The `configtxgen` tool is used to create 4 configuration artifacts
    - Orderer genesis block
        - The orderer block is the genesis block for the ordering service
    - Channel configuration transaction
        - The channel configuration transaction file is broadcast to the orderer at channel creation time
    - 2 anchor peer transactions, one for each peer org
        - The anchor peer transactions specify each org’s anchor peer on this channel
- `Configtxgen` consumes `configtx.yaml` that contains the definitions for the sample network
    - There are 3 members - 1 orderer org (`OrdererOrg`) and 2 peer orgs (`Org1`, `Org2`) each managing and maintaining 2 peer nodes
    - This file also specifies a consortium - `SampleConsortium` - consisting of the 2 peer orgs
        - `SampleConsortium` is defined in the system-level profile and then referenced by the channel-level profile
        - Channels exist within the purview of a consortium, and all consortia must be defined in the scope of the network at large
    - <u>`Profiles`</u> section
        - `TwoOrgsOrdererGenesis`: generates the genesis block for a Solo ordering service
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

    - The certs and keys (i.e. the MSP material) will be output into `crypto-config` at the root of the first-network directory
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

- Create the channel transaction artifact (`channel.tx`)

    ```bash
    # The channel.tx artifact contains the definitions for our sample channel
    # fabric-samples/first-network
    export CHANNEL_NAME=mychannel  && ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
    ```

    - > You don’t have to issue a special command for the channel if you are using a Raft or Kafka ordering service. The `TwoOrgsChannel` profile will use the ordering service configuration you specified when creating the genesis block for the network
- Define the anchor peer for `Org1` and `Org2` on the channel

    ```bash
    # fabric-samples/first-network
    # export CHANNEL_NAME=mychannel
    ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

    ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
    ```

## Start the network
- Start the network
	
    ```bash
    docker-compose -f docker-compose-cli.yaml up -d
    ```

### Create & join channel
- > Create additional channel configuration transactions using the same or different profiles in the `configtx.yaml` that's passed to the configtxgen tool. Then you can repeat the process defined in this section to establish those other channels in your network
- Enter the CLI container
	
    ```bash
    docker exec -it cli bash
    ```

- These variables for `peer0.org1.example.com` are baked into the CLI container, therefore we can operate without passing them
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container
    # printenv | less

    # Environment variables for PEER0
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    CORE_PEER_ADDRESS=peer0.org1.example.com:7051
    CORE_PEER_LOCALMSPID="Org1MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    ```

    - If you want to send calls to other peers or the orderer, **override** the environment variables
- **Channel names** must be **all lower case**, less than 250 characters long and match the regular expression `[a-z][a-z0-9.-]*`
- Create channel
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    export CHANNEL_NAME=mychannel

    # the `channel.tx` file is mounted in the channel-artifacts directory within your CLI container
    # as a result, we pass the full path for the file
    # we also pass the path for the orderer ca-cert in order to verify the TLS handshake
    # be sure to export or replace the `$CHANNEL_NAME` variable appropriately

    peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    ```

    - This command returns a genesis block `<CHANNEL_NAME.block>` (binary) which we will use to join the channel
        - It contains the configuration information specified in `channel.tx`
        - If you have not made any modifications to the default channel name, then the command will return you a proto titled `mychannel.block`
    - > `--cafile` is the local path to the orderer’s root cert, allowing us to verify the TLS handshake
- Join `peer0.org1.example.com` to the channel
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    # By default, this joins `peer0.org1.example.com` only
    # the `<CHANNEL_NAME.block>` was returned by the previous command
    # if you have not modified the channel name, you will join with `mychannel.block`
    # if you have created a different channel name, then pass in the appropriately named block

    peer channel join -b mychannel.block
    ```

    - > Make other peers join the channel as necessary by making appropriate changes in the 4 environment variables
- Join `peer0.org2.example.com` to the channel
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container
    
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:9051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel join -b mychannel.block
    ```

    - > Prior to v1.4.1 all peers within the docker network used port `7051`. If using a version of fabric-samples prior to v1.4.1, modify all occurrences of CORE_PEER_ADDRESS in this tutorial to use port `7051`
### Update the anchor peers
- The following commands are channel updates and they will propagate to the definition of the channel
    - In essence, we adding additional configuration information on top of the channel’s genesis block
    - Note that we are not modifying the genesis block, but simply adding deltas into the chain that will define the anchor peers
- Update the channel definition to define the anchor peer for `Org1` as `peer0.org1.example.com`
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/Org1MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    ```

- Update the channel definition to define the anchor peer for `Org2` as `peer0.org2.example.com`
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    # preface this call with appropriate environment variables
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:9051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/Org2MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    ```

### Install & instantiate chaincode
- Applications interact with the blockchain ledger through chaincode. As such we need to **install** the chaincode **on every peer** that will execute and endorse our transactions, and then **instantiate** the chaincode **on the channel**
- You can only install one version of the source code per chaincode name and version
    - The source code exists on the peer’s file system in the context of chaincode name and version; it is language agnostic
        - Similarly the instantiated chaincode container will be reflective of whichever language has been installed on the peer
- Install the sample chaincode onto the `peer0` in `Org1`
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    # this installs the Go chaincode. For go chaincode `-p` takes the **relative path** from `$GOPATH/src`
    peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/chaincode_example02/go/

    # this installs the Node.js chaincode
    # make note of the `-l` flag to indicate "node" chaincode
    # for node chaincode `-p` takes the **absolute path** to the node.js chaincode
    # peer chaincode install -n mycc -v 1.0 -l node -p /opt/gopath/src/github.com/chaincode/chaincode_example02/node/

    # make note of the `-l` flag to indicate "java" chaincode
    # for java chaincode `-p` takes the **absolute path** to the java chaincode
    peer chaincode install -n mycc -v 1.0 -l java -p /opt/gopath/src/github.com/chaincode/chaincode_example02/java/
    ```

- When we instantiate the chaincode on the channel, the endorsement policy will be set to require endorsements from a peer in both `Org1` and `Org2`. [ ] Therefore, we also need to install the chaincode on a peer in `Org2`
- Modify the following 4 environment variables to issue the install command against `peer0` in `Org2`
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    # Environment variables for PEER0 in Org2

    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    CORE_PEER_ADDRESS=peer0.org2.example.com:9051
    CORE_PEER_LOCALMSPID="Org2MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    ```

- Install the sample chaincode onto a `peer0` in `Org2`
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/chaincode_example02/go/
    ```

- Instantiate the chaincode on the channel
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside CLI container

    # be sure to replace the $CHANNEL_NAME environment variable if you have not exported it
    # if you did not install your chaincode with a name of `mycc`, then modify that argument as well

    # “endorsement” from a peer belonging to `Org1` AND `Org2` (i.e. 2 endorsement)
    peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"

    # The instantiation of the Node.js chaincode will take roughly a minute. The command is not hanging; rather it is installing the fabric-shim layer as the image is being compiled
    # if you did not install your chaincode with a name of `mycc`, then modify that argument as well
    # notice that we must pass the `-l` flag after the chaincode name to identify the language
    # peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -l node -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"

    # Java chaincode instantiation might take time as it compiles chaincode and downloads docker container with java environment
    # peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -l java -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"

    ```

    - This will **initialize the chaincode on the channel**, **set the endorsement policy for the chaincode**, and **launch a chaincode container for the targeted peer**
    - `-P` argument specifies policy where we specify the required level of endorsement for a transaction against this chaincode to be validated
- If you want additional peers to interact with ledger, then you will need to join them to the channel, and install the same name, version and language of the chaincode source onto the appropriate peer’s filesystem
    - A chaincode container will be launched for each peer as soon as they try to interact with that specific chaincode
- Once the chaincode has been instantiated on the channel, we can forgo the `l` flag. We need only pass in the channel identifier and name of the chaincode
### Query
- Query the value of `a`

    ```bash
    # be sure to set the -C and -n flags appropriately
    peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
    ```

### Invoke
- Move 10 from `a` to `b`
	
    ```bash
    peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["invoke","a","b","10"]}'
    ```

- [x] Error: error getting endorser client for invoke: endorser client failed to connect to peer0.org1.example.com:9051: failed to create new connection: connection error: desc = "transport: error while dialing: dial tcp 172.19.0.5:9051: connect: connection refused"
    - `peer0.org1.example.com:7051`
### Behind the scenes
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/build_network.html#what-s-happening-behind-the-scenes
- Chaincode MUST be installed on a peer in order for it to successfully perform read/write operations against the ledger
- A chaincode container is not started for a peer until an init or traditional transaction - read/write - is performed against that chaincode (e.g. query for the value of “a”)
    - The transaction causes the container to start
- All peers in a channel maintain an exact copy of the ledger which comprises the blockchain to store the immutable, sequenced record in blocks, as well as a state database to maintain a snapshot of the current state
    - This includes those peers that do not have chaincode installed on them
- The chaincode is accessible after it is installed because it has already been instantiated
- Check the logs for the CLI Docker container to see transactions
	
    ```bash
    docker logs -f cli
    ```

- Inspect the individual chaincode containers to see the separate transactions executed against each container
	
    ```bash
    docker logs dev-peer0.org2.example.com-mycc-1.0
    ```

### The Docker Compose topology
- `docker-compose-e2e.yaml` is constructed to run end-to-end tests using the Node.js 
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/build_network.html#understanding-the-docker-compose-topology
### Using CouchDB
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/build_network.html#using-couchdb
- A kind of NoSQL solution