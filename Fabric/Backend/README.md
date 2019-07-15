- > https://godoc.org/github.com/hyperledger/fabric-sdk-go
- > https://github.com/hyperledger/fabric-sdk-go
- > https://github.com/chainHero/heroes-service/tree/v1.0.5
- Config
    - https://github.com/hyperledger/fabric-sdk-go/tree/master/test/fixtures/config
    - https://github.com/hyperledger/fabric-sdk-go/blob/master/pkg/core/config/testdata/template/config.yaml
- Some components (peers) also deploys docker containers to [ ] separate data (channel)
# `fabricsdk`
- Package fabricsdk enables Go developers to build solutions that interact with Hyperledger Fabric
- Basic workflow
    1. Instantiate a `fabsdk` instance using a configuration
        - `fabsdk` maintains caches so you should minimize instances of `fabsdk` itself
    2. Create a context based on a user and organization, using your `fabsdk` **instance**
        - Note: A channel context additionally requires the channel ID
    3. Create a client instance using its `New` func, passing the context
        - Note: you **create a new client instance for each context you need**
    4. Use the funcs provided by each client to create your solution
    5. Call `fabsdk.Close()` to release resources and caches
## `pkg/fabsdk`
- The main package of the Fabric SDK
- This package enables creation of contexts based on configuration
- Package `fabsdk` enables client usage of a Hyperledger Fabric network
## `pkg/client/channel`
- Provides channel transaction capabilities
- Package `channel` enables access to a channel on a Fabric network
- A channel client instance provides a handler to interact with peers on specified channel
- Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel
- An application that requires interaction with multiple channels should create a separate instance of the channel client for each channel
- Basic flow
    1. Prepare channel client context
    2. Create channel client
    3. Execute chaincode
    4. Query chaincode
## `pkg/client/event`
- Provides channel event capabilities
- Package `event` enables access to a channel events on a Fabric network
- Event client receives events such as block, filtered block, chaincode, and transaction status events
- Basic flow
    1. Prepare channel client context
    2. Create event client
    3. Register for events
    4. Process events (or timeout)
    5. Unregister
## `pkg/client/ledger`
- Enables queries to a channel's underlying ledger
- Package `ledger` enables ledger queries on specified channel on a Fabric network
- An application that requires ledger queries from multiple channels should create a separate instance of the ledger client for each channel
- Ledger client supports the following queries: `QueryInfo`, `QueryBlock`, `QueryBlockByHash`, `QueryBlockByTxID`, `QueryTransaction` and `QueryConfig`
- Basic flow
    1. Prepare channel context
    2. Create ledger client
    3. Query ledger
## `pdk/client/resmgmt`
- Provides resource management capabilities such as installing chaincode
- Package `resmgmt` enables creation and update of resources on a Fabric network
    - It allows administrators to create and/or update channels, and for peers to join channels
    - Administrators can also perform chaincode related operations on a peer, such as installing, instantiating, and upgrading chaincode
- Basic flow
    1. Prepare client context
    2. Create resource management client
    3. create new channel
    4. Peer(s) join channel
    5. Install chaincode onto peer(s) filesystem
    6. Instantiate chaincode on channel
    7. Query peer for channels, installed/instantiated chaincodes etc.
## `pkg/client/msp`
- Enables identity management capability
- Package `msp` enables creation and update of users on a Fabric network
- Msp client supports the following actions: `Enroll`, `Reenroll`, `Register`, `Revoke` and `GetSigningIdentity`
- Basic flow
    1. Prepare client context
    2. Create msp client
    3. Register user
    4. Enroll user