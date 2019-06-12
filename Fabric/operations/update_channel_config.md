- Channel configurations contain all of the information relevant to the administration of a channel
    - The channel configuration specifies which organizations are members of channel
    - It also includes other channel-wide configuration information such as channel access policies and block batch sizes
- This configuration is stored on the ledger in a block, and is therefore known as a configuration (config) block
- Configuration blocks contain a single configuration
    - The first of these blocks is known as the “genesis block” and contains the initial configuration required to bootstrap a channel
    - Each time the configuration of a channel changes it is done through a new configuration block, with the latest configuration block representing the current channel configuration
    - Orderers and peers keep the current channel configuration in memory to facilitate all channel operations such as cutting a new block and validating block transactions
- Because configurations are stored in blocks, updating a config happens through a process called a “configuration transaction” (even though the process is a little different from a normal transaction)
- Updating a config is a process of pulling the config, translating into a format that humans can read, modifying it and then submitting it for approval
- For a more in-depth look at the process for pulling a config and translating it into JSON, check out Adding an Org to a Channel
## Edit a config
- [Config file](../tutorials/extend_network_add_org_2_channel/config.json)
- Beyond the definitions of the policies – defining who can do certain things at the channel level, and who has the permission to change, who can change the config – channels also have other kinds of features that can be modified using a config update
- Besides adding an org to a channel, some other things are possible to change with a config update
    - `BatchSize`
    	
        ```json
        {
            "absolute_max_bytes": 102760448,
            "max_message_count": 10,
            "preferred_max_bytes": 524288
        }
        ```
    
        - No block will appear larger than `absolute_max_bytes` large or with more than `max_message_count` transactions inside the block
        - If it is possible to construct a block under `preferred_max_bytes`, then a block will be cut prematurely, and transactions larger than this size will appear in their own block
    - `BatchTimeout`
    	
        ```json
        { "timeout": "2s" }
        ```
    
        - The amount of time to wait after the first transaction arrives for additional transactions before cutting a block
        - Decreasing this value will improve latency, but decreasing it too much may decrease throughput by not allowing the block to fill to its maximum capacity
    - `ChannelRestrictions`
    	
        ```json
        {
            "max_count": 1000
        }
        ```
    
        - The total number of channels the orderer is willing to allocate may be specified as `max_count`
        - This is primarily useful in pre-production environments with weak consortium `ChannelCreation` policies
    - Channel creation policy
    	
        ```json
        {
            "type": 3,
            "value": {
                "rule": "ANY",
                "sub_policy": "Admins"
            }
        }
        ```
    
        - Defines the policy value which will be set as the `mod_policy` for the Application group of new channels for the consortium it is defined in. The signature set attached to the channel creation request will be checked against the instantiation of this policy in the new channel to ensure that the channel creation is authorized
        - Note that this config value is only set in the orderer system channel
    - Kafka brokers
    	
        ```json
        {
            "brokers": [
                "kafka0:9092",
                "kafka1:9092",
                "kafka2:9092",
                "kafka3:9092"
            ]
        }
        ```
    
        - When `ConsensusType` is set to `kafka`, the `brokers` list enumerates some subset (or preferably all) of the Kafka brokers for the orderer to initially connect to at startup
        - Note that it is not possible to change your consensus type after it has been established (during the bootstrapping of the genesis block
    - `anchor_peers`
    	
        ```json
        {
            "host": "peer0.org2.example.com",
            "port": 9051
        }
        ```
    
        - Defines the location of the anchor peers for each Org
    - `BlockDataHashingStructure`
    	
        ```json
        { "width": 4294967295 }
        ```
    
        - The hash of the block data is computed as a Merkle tree. This value specifies the width of that Merkle tree
        - For the time being, this value is fixed to `4294967295` which corresponds to a simple flat hash of the concatenation of the block data bytes
    - `HashingAlgorithm`
    	
        ```json
        { "name": "SHA256" }
        ```
    
        - The algorithm used for computing the hash values encoded into the blocks of the blockchain. In particular, this affects the data hash, and the previous block hash fields of the block
        - Note, this field currently only has one valid value (`SHA256`) and should not be changed
    - `BlockValidation`
    	
        ```json
        {
            "type": 3,
            "value": {
                "rule": "ANY",
                "sub_policy": "Writers"
            }
        }
        ```
    
        - This policy specifies the signature requirements for a block to be considered valid
        - By default, it requires a signature from some member of the ordering org
    - `OrdererAddresses`
    	
        ```json
        {
            "addresses": [
                "orderer.example.com:7050"
            ]
        }
        ```
    
        - A list of addresses where clients may invoke the orderer `Broadcast` and `Deliver` functions
        - The peer randomly chooses among these addresses and fails over between them for retrieving blocks
- Just as we add an Org by adding their artifacts and MSP information, you can remove them by reversing the process
- Note that once the consensus type has been defined and the network has been bootstrapped, it is not possible to change it through a configuration update
- There is another important channel configuration (especially for v1.1) known as [Capability Requirements](https://hyperledger-fabric.readthedocs.io/en/release-1.4/capability_requirements.html)
### Edit block batch size

```bash
export MAXBATCHSIZEPATH=".channel_group.groups.Orderer.values.BatchSize.value.max_message_count"
# display the value of that property
jq "$MAXBATCHSIZEPATH" config.json # returns 10

# set new batch size
jq “$MAXBATCHSIZEPATH = 20” config.json > modified_config.json
# display new value
jq “$MAXBATCHSIZEPATH” modified_config.json
# convert and submit...
```

## Get necessary signature
- Once you’ve successfully generated the protobuf file, it’s time to get it signed
- To do this, you need to know the relevant policy for whatever it is you’re trying to change
- By default, editing the configuration of
    - a particular org (for example, changing anchor peers) requires only the admin signature of that org
    - the application (like who the member orgs are) requires a majority of the application organizations’ admins to sign
    - the orderer requires a majority of the ordering organizations’ admins (of which there are by default only 1)
    - The top level channel group requires both the agreement of a majority of application organization admins and orderer organization admins
- If you have made changes to the default policies in the channel, you’ll need to compute your signature requirements accordingly
- The actual process of getting these signatures will depend on how you’ve set up your system, but there are 2 main implementations
    - Currently, the Fabric command line defaults to a “pass it along” system
        - That is, the Admin of the Org proposing a config update sends the update to someone else (another Admin, typically) who needs to sign it. This Admin signs it (or doesn’t) and passes it along to the next Admin, and so on, until there are enough signatures for the config to be submitted
        - This has the virtue of simplicity – when there are enough signatures, the last Admin can simply submit the config transaction (in Fabric, the peer channel update command includes a signature by default)
        - However, this process will only be practical in smaller channels, since the “pass it along” method can be time consuming
    - The other option is to submit the update to every Admin on a channel and wait for enough signatures to come back
        - These signatures can then be stitched together and submitted
        - This makes life a bit more difficult for the Admin who created the config update (forcing them to deal with a file per signer) but is the recommended workflow for users that are developing Fabric management applications
- Once the config has been added to the ledger, it will be a best practice to pull it and convert it to JSON to check to make sure everything was added correctly
    - This will also serve as a useful copy of the latest config