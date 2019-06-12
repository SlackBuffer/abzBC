# Adding an org to a channel
## Setup the environment

```bash
# fabric-samples/first-network

# tidy up
./byfn.sh down

# generate the default BYFN artifacts
./byfn.sh generate

# launch the network
./byfn.sh up
```

## Bring org3 into the channel with the script

```bash
./eyfn.sh up
```

- `eyfn.sh` can be used with the same Node.js chaincode and database options as `byfn.sh` by issuing the following (instead of `./byfn.sh up`)
	
    ```bash
    ./byfn.sh up -c testchannel -s couchdb -l node
    ./eyfn.sh up -c testchannel -s couchdb -l node
    ```

## Manually
- Bring network down first (if necessary)

    ```bash
    # ./byfn.sh down   
    ./eyfn.sh down
    ```
- The manual steps assume that the <u>`FABRIC_LOGGING_SPEC`</u> in the `cli` and `Org3cli` containers is set to `DEBUG`
    - Modify the `docker-compose-cli.yaml` file in the `first-network` directory
    - Modify the `docker-compose-org3.yaml` file in the `first-network` directory
- Bring network back up
	
    ```bash
    ./byfn.sh generate
    ./byfn.sh up    
    ```

### Generate the org3 crypto material
- Leverages `cryptogen` to generate the keys and certificates for an Org3 CA as well as two peers bound to this new Org
	
    ```bash
    # fabric-samples/first-network
    cd org3-artifacts
    # put into a newly generated `crypto-config` folder within the present working directory
    ../../bin/cryptogen generate --config=./org3-crypto.yaml
    ```

- Use the `configtxgen` utility to print out the Org3-specific configuration material in JSON
	
    ```bash
    # fabric-samples/first-network/org3-artifacts
    # preface the command by telling the tool to look in the current directory for the `configtx.yaml` file that it needs to ingest
    export FABRIC_CFG_PATH=$PWD && ../../bin/configtxgen -printOrg Org3MSP > ../channel-artifacts/org3.json
    ```

    - This file contains the policy definitions for Org3, as well as 3 important certificates presented in base 64 format: the admin user certificate (which will be needed to act as the admin of Org3 later on), a CA root cert, and a TLS root cert
- Port the Orderer Org’s MSP material into the Org3 `crypto-config` directory. In particular, we are concerned with the Orderer’s TLS root cert, which will allow for secure communication between Org3 entities and the network’s ordering node
	
    ```bash
    cd ../ && cp -r crypto-config/ordererOrganizations org3-artifacts/crypto-config/
    ```

### Prepare the cli environment
- The update process makes use of the configuration translator tool – `configtxlator`
    - This tool provides a stateless REST API independent of the SDK
    - Additionally it provides a CLI, to simplify configuration tasks in Fabric networks
    - The tool allows for the easy conversion between different equivalent data representations/formats (in this case, between protobufs and JSON)
    - Additionally, the tool can compute a configuration update transaction based on the differences between two channel configurations
- The `cli` container has been mounted with the BYFN crypto-config library, giving us access to the MSP material for the two original peer organizations and the Orderer Org
	
    ```bash
    docker exec -it cli bash
    ```

- The bootstrapped identity is the Org1 admin user, meaning that any steps where we want to act as Org2 will require the export of MSP-specific environment variables
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside `cli` container
    # Export the ORDERER_CA and CHANNEL_NAME variables
    export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  && export CHANNEL_NAME=mychannel

    # Check to make sure the variables have been properly set
    echo $ORDERER_CA && echo $CHANNEL_NAME
    ```

### Fetch the configuration
- Fetch the most recent config block for the channel – `mychannel`
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside `cli` container
    peer channel fetch config config_block.pb -o orderer.example.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
    # 2019-06-03 06:46:24.023 UTC [cli.common] readBlock -> INFO 048 Received block: 2
    ```

    - This command saves the binary protobuf channel configuration block to `config_block.pb` (the choice of name and file extension is arbitrary)
        - > Following a convention which identifies both the type of object being represented and its encoding (protobuf or JSON) is recommended
    - After the command is issued, the most recent configuration block for mychannel is actually block **2**, NOT the genesis block
        - block 0: genesis block
        - block 1: Org1 anchor peer update
        - block 2: Org2 anchor peer update
        - > The BYFN script defined anchor peers for the 2 organizations – `Org1` and `Org2` – in 2 separate channel update transactions
    - By default, the `peer channel fetch config` command returns the most recent configuration block for the targeted channel
- The reason why we have to pull the latest version of the config is because channel config elements are versioned
    - Versioning prevents config changes from being repeated or replayed (for instance, reverting to a channel config with old CRLs would represent a security risk)
    - Also it helps ensure concurrency (if you want to remove an Org from your channel, for example, after a new Org has been added, versioning will help prevent you from removing both Orgs, instead of just the Org you want to remove)
### Convert the configuration to JSON and trim it down
- Make use of the `configtxlator` tool to decode this channel configuration block into JSON format (which can be read and modified by humans); strip away all of the headers, metadata, creator signatures, and so on that are irrelevant to the change we want to make
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside `cli` container   
    configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json
    ```

### Add the or3 crypto material
- > The steps that have been taken up to this point will be nearly identical no matter what kind of config update you’re trying to make
- Use the `jq` tool once more to append the Org3 configuration definition – `org3.json` – to the channel’s application `groups` field, and name the output – `modified_config.json`
	
    ```bash
    # `/opt/gopath/src/github.com/hyperledger/fabric/peer` inside `cli` container   
    jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"Org3MSP":.[1]}}}}}' config.json ./channel-artifacts/org3.json > modified_config.json
    ```

    - The initial json file contains only Org1 and Org2 material, whereas “modified” json file contains all 3 Orgs
- Translate `config.json` back into a protobuf called `config.pb`
	
    ```bash
    configtxlator proto_encode --input config.json --type common.Config --output config.pb
    ```

- Encode `modified_config.json` to `modified_config.pb`
	
    ```bash
    configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
    ```

- Use `configtxlator` to calculate the **delta** between these two config protobufs
	
    ```bash
    configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output org3_update.pb
    ```

    - Output a new protobuf binary named `org3_update.pb`
        - `org3_update.pb` contains the Org3 definitions and high level pointers to the Org1 and Org2 material
    - We are able to forgo the extensive MSP material and modification policy information for Org1 and Org2 because this data is already present within the channel’s genesis block. As such, we only need the delta between the two configurations
- Decode `org3_update.pb` into editable JSON format and call it `org3_update.json`
	
    ```bash
    configtxlator proto_decode --input org3_update.pb --type common.ConfigUpdate | jq . > org3_update.json
    ```

- Get back the header field stripped away earlier and name the output file `org3_update_in_envelope.json`
	
    ```bash
    echo '{"payload":{"header":{"channel_header":{"channel_id":"mychannel", "type":2}},"data":{"config_update":'$(cat org3_update.json)'}}}' | jq . > org3_update_in_envelope.json
    ```

- Leverage the `configtxlator` tool one last time and convert it into the fully fledged protobuf format that Fabric requires
	
    ```bash
    configtxlator proto_encode --input org3_update_in_envelope.json --type common.Envelope --output org3_update_in_envelope.pb
    ```

### Sign and submit the config update
- Signatures from the requisite Admin users are needed before the config can be written to the ledger
- The modification policy (mod_policy) for our channel Application group is set to the default of “MAJORITY”, which means that we need a majority of existing org admins to sign it
    - Because we have only two orgs – Org1 and Org2 – and the majority of two is two, we need both of them to sign
    - Without both signatures, the ordering service will reject the transaction for failing to fulfill the policy
- Sign this update proto as the Org1 Admin
	
    ```bash
    peer channel signconfigtx -f org3_update_in_envelope.pb
    ```

    - The CLI container is bootstrapped with the Org1 MSP material, so we simply need to issue the peer channel signconfigtx command
- The final step is to switch the CLI container’s identity to reflect the Org2 Admin user
    - Do this by exporting 4 environment variables specific to the Org2 MSP
    	
        ```bash
        export CORE_PEER_LOCALMSPID="Org2MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
        export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
        export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
        ```
    
        - > Switching between organizations to sign a config transaction (or to do anything else) is not reflective of a real-world Fabric operation. A single container would never be mounted with an entire network’s crypto material. Rather, the config update would need to be securely passed out-of-band to an Org2 Admin for inspection and approval
- Issue the `peer channel update` command
	
    ```bash
    # The upcoming update call to the ordering service will undergo a series of systematic signature and policy checks
    # another shell
    docker logs -f orderer.example.com

    peer channel update -f org3_update_in_envelope.pb -c $CHANNEL_NAME -o orderer.example.com:7050 --tls --cafile $ORDERER_CA
    # 2019-06-03 07:35:20.327 UTC [channelCmd] update -> INFO 04d Successfully submitted channel update
    ```

    - The Org2 Admin signature will be attached to this call so there is no need to manually sign the protobuf a second time
- The successful channel update call returns a new block – block 5 – to all of the peers on the channel
  - Blocks 3 and 4 are the instantiation and invocation of the `mycc` chaincode
  - Block 5 serves as the most recent channel configuration with Org3 now defined on the channel
  	
      ```bash
      docker logs -f peer0.org1.example.com
      ```
  
### Configuring leader election
- This section is included as a general reference for understanding the leader election settings when adding organizations to a network after the initial channel configuration has completed
- This sample defaults to **dynamic** leader election, which is set for all peers in the network in `fabric-samples/first-network/base/peer-base.yaml`
- Newly joining peers are bootstrapped with the genesis block, which does not contain information about the organization that is being added in the channel configuration update
    - Therefore new peers are not able to utilize gossip as they cannot verify blocks forwarded by other peers [ ] from their own (newly addeded?) organization until they get the configuration transaction which added the organization to the channel
- Newly added peers must therefore have one of the following configurations so that they receive blocks from the ordering service
    1. To utilize static leader mode, configure the peer to be an organization leader
    	
        ```yaml
        CORE_PEER_GOSSIP_USELEADERELECTION=false
        CORE_PEER_GOSSIP_ORGLEADER=true
        ```
    
        - This configuration must be the same for all new peers added to the channel
    2. To utilize dynamic leader election, configure the peer to use leader election
    	
        ```yaml
        CORE_PEER_GOSSIP_USELEADERELECTION=true
        CORE_PEER_GOSSIP_ORGLEADER=false
        ```
    
        - Because peers of the newly added organization won’t be able to form membership view, this option will be similar to the static configuration, as each peer will start proclaiming itself to be a leader
        - However, once they get updated with the configuration transaction that adds the organization to the channel, there will be only 1 active leader for the organization. Therefore, it is recommended to leverage this option if you eventually want the organization’s peers to utilize leader election
### Join org3 to the channel
- At this point, the channel configuration has been updated to include our new organization – `Org3` – meaning that peers attached to it can now join `mychannel`
- Launch the containers for the Org3 peers and an Org3-specific CLI
	
    ```bash
    docker-compose -f docker-compose-org3.yaml up -d
    ```

    - This new compose file has been configured to bridge across our initial network, so the two peers and the CLI container will be able to resolve with the existing peers and ordering node
- Exec into the Org3-specific cli container
	
    ```bash
    docker exec -it Org3cli bash
    ```

- Export the two key environment variables: `ORDERER_CA` and `CHANNEL_NAME`
	
    ```bash
    export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem && export CHANNEL_NAME=mychannel

    echo $ORDERER_CA && echo $CHANNEL_NAME
    ```

- Send a call to the ordering service asking for the genesis block of `mychannel`
    - The ordering service is able to verify the Org3 signature attached to this call as a result of our successful channel update
    - If Org3 has not been successfully appended to the channel config, the ordering service should reject this request
	
    ```bash
    peer channel fetch 0 mychannel.block -o orderer.example.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
    ```

    - We are passing a `0` to indicate that we want the first block on the channel’s ledger (i.e. the genesis block)
        - If we simply passed the `peer channel fetch config` command, then we would have received block 5 – the updated config with Org3 defined
        - However, we can’t begin our ledger with a downstream block – we must start with block 0
- Issue the peer channel join command and pass in the genesis block – `mychannel.block`
	
    ```bash
    peer channel join -b mychannel.block
    ```

    - If you want to join the second peer for Org3, export the `TLS` and `ADDRESS` variables and reissue the `peer channel join` command
    	
        ```bash
        export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer1.org3.example.com/tls/ca.crt && export CORE_PEER_ADDRESS=peer1.org3.example.com:12051

        peer channel join -b mychannel.block
        ```
    
### Upgrade and invoke chaincode
- The final piece of the puzzle is to increment the chaincode version and update the endorsement policy to include Org3
    - Since we know that an upgrade is coming, we can forgo the futile exercise of installing version 1 of the chaincode
    - We are solely concerned with the new version where Org3 will be part of the endorsement policy, therefore we’ll jump directly to version 2 of the chaincode
- From the Org3 cli
	
    ```bash
    peer chaincode install -n mycc -v 2.0 -p github.com/chaincode/chaincode_example02/go/
    ```

    - Modify the environment variables accordingly and reissue the command if you want to install the chaincode on the second peer of Org3
        - A second installation is not mandated, as you only need to install chaincode on peers that are going to serve as endorsers or otherwise interface with the ledger (i.e. query only)
        - Peers will still run the validation logic and serve as committers without a running chaincode container
- Jump back to the **original CLI** container and install the new version on the Org1 and Org2 peers
	
    ```bash
    # Org1
    peer chaincode install -n mycc -v 2.0 -p github.com/chaincode/chaincode_example02/go/

    # Org2
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
    peer chaincode install -n mycc -v 2.0 -p github.com/chaincode/chaincode_example02/go/
    ```

- Upgrade the chaincode
	
    ```bash
    # original cli container
    peer chaincode upgrade -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -v 2.0 -c '{"Args":["init","a","90","b","210"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')"
    ```

    - There have been no modifications to the underlying source code, we are simply adding Org3 to the endorsement policy for a chaincode – `mycc` – on `mychannel`
    - A chaincode upgrade requires usage of the `init` method (`c` flag)
        - If your chaincode requires arguments be passed to the `init` method, then you will need to do so here
    - Any identity satisfying the chaincode’s instantiation policy can issue the upgrade call. By default, these identities are the channel Admins
- The upgrade call adds a new block – block 6 – to the channel’s ledger and allows for the Org3 peers to execute transactions during the endorsement phase
- Operate from Org3 cli
	
    ```bash
    peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -c '{"Args":["invoke","a","b","10"]}'
    ```

## Conclusion
- There is a logical method to the various steps - the endgame is to form a **delta transaction object** represented **in protobuf binary format** and then acquire the requisite number of **admin signatures** such that the channel configuration update transaction **fulfills the channel’s modification policy**
- The `configtxlator` and `jq` tools, along with the ever-growing `peer channel` commands, provide us with the functionality to accomplish this task