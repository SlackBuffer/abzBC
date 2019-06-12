# Chaincode
- Chaincode runs in a secured Docker container isolated from the endorsing peer process
- Chaincode initializes and manages ledger state through transactions submitted by applications
- A chaincode typically handles business logic agreed to by members of the network, so it may be considered as a “smart contract”
- State created by a chaincode is scoped exclusively to that chaincode and can’t be accessed directly by another chaincode
- Within the same network, given the appropriate permission a chaincode may invoke another chaincode (e.g., cross channels) to access its state
    - If the called chaincode is on a different channel from the calling chaincode, only read query is allowed. That is, the called chaincode on a different channel is only a `Query`, which does not participate in state validation checks in subsequent commit phase
## [For developers](https://hyperledger-fabric.readthedocs.io/en/release-1.4/chaincode4ade.html)
- Every chaincode program must implement the `Chaincode` interface whose methods are called in response to received transactions
    - The `Init` method is called when a chaincode receives an `instantiate` or `upgrade` transaction so that the chaincode may perform any necessary initialization, including initialization of application state
        - Chaincode upgrade also calls this function to reset or to migrate data, so be careful to avoid a scenario where you inadvertently clobber ledger's data
        - When writing a chaincode that will upgrade an existing one, make sure to modify the `Init` function appropriately. In particular, provide an empty “Init” method if there’s no “migration” or nothing to be initialized as part of the upgrade
    - The `Invoke` method is called in response to receiving an `invoke` transaction to process transaction proposals
- `ChaincodeStubInterface` interface is used to access and modify the ledger, and to make invocations between chaincodes
### Simple asset chaincode

```bash
mkdir -p $GOPATH/src/sacc && cd $GOPATH/src/sacc
touch sacc.go

go get -u github.com/hyperledger/fabric/core/chaincode/shim
go build
```

### Testing using dev mode
- Normally chaincodes are started and maintained by peer
- In “dev mode”, chaincode is built and started by the user
    - This mode is useful during chaincode development phase for rapid code/build/run/debug cycle turnaround
- Start “dev mode” by leveraging pre-generated orderer and channel artifacts for a sample dev network
    - As such, the user can immediately jump into the process of compiling chaincode and driving calls
### Install hyperledger fabric samples

```bash
# fabric-samples
cd chaincode-docker-devmode
# open 3 terminals
```

### Terminal 1 - start the network
- `docker-compose -f docker-compose-simple.yaml up`
### Terminal 2 - build & start the chaincode

```bash
docker exec -it chaincode bash
cd sacc
go build
# run chaincode
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./sacc
```

### Ternimal 3 - use the chaincode
- Even though you are in `--peer-chaincodedev` mode, you still have to install the chaincode so the life-cycle system chaincode can go through its checks normally
    - This requirement may be removed in future when in `--peer-chaincodedev` mode

```bash
docker exec -it cli bash

peer chaincode install -p chaincodedev/chaincode/sacc -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":["a","10"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["set", "a", "20"]}' -C myc
peer chaincode query -n mycc -c '{"Args":["query","a"]}' -C myc
```

- Can easily test different chaincodes by adding them to the chaincode subdirectory and relaunching your network
    - At this point they will be accessible in the `chaincode` container
### Chaincode access control
- Chaincode can utilize the client (submitter) certificate for access control decisions by calling the `GetCreator()` function
- Additionally the Go shim provides extension APIs that extract client identity from the submitter’s certificate that can be used for access control decisions, whether that is based on client identity itself, or the org identity, or on a client identity attribute
- https://github.com/hyperledger/fabric/blob/master/core/chaincode/shim/ext/cid/README.md
### Chaincode encryption
- In certain scenarios, it may be useful to encrypt values associated with a key in their entirety or simply in part
    - > For example, if a person’s social security number or address was being written to the ledger, then you likely would not want this data to appear in plaintext
- Chaincode encryption is achieved by leveraging the [entities extension](https://github.com/hyperledger/fabric/tree/master/core/chaincode/shim/ext/entities) which is a BCCSP (Blockchain Cryptographic Service Provider) wrapper with commodity factories and functions to perform cryptographic operations such as encryption and elliptic curve digital signatures
    - For example, to encrypt, the invoker of a chaincode passes in a cryptographic key via the transient field. The same key may then be used for subsequent query operations, allowing for proper decryption of the encrypted state values
- `fabric/examples/chaincode/go/enccc_example`
    - `utils.go` helper program loads the chaincode shim APIs and Entities extension and builds a new class of functions (e.g. `encryptAndPutState` & `getStateAndDecrypt`) that the sample encryption chaincode then leverages. As such, the chaincode can now marry the basic shim APIs of `Get` and `Put` with the added functionality of `Encrypt` and `Decrypt`
### Managing external dependencies for chaincode written in Go
- https://github.com/golang/go/wiki/PackageManagementTools
- If your chaincode requires packages not provided by the Go standard library, you will need to include those packages with your chaincode
- It is also a good practice to add the shim and any extension libraries to your chaincode as a dependency
- `govendor`

    ```bash
    govendor init
    govendor add +external  // Add all external package, or
    govendor add github.com/external/pkg // Add specific external package
    ```

    - This imports the external dependencies into a local `vendor` directory
    - If you are vendoring the Fabric shim or shim extensions, clone the Fabric repository to your `$GOPATH/src/github.com/hyperledger` directory, before executing the `govendor` commands
- Once dependencies are vendored in your chaincode directory, `peer chaincode package` and `peer chaincode install` operations will then include code associated with the dependencies into the chaincode package
## [For operators](https://hyperledger-fabric.readthedocs.io/en/release-1.4/chaincode4noah.html)
### Chaincode lifecycle
- The Hyperledger Fabric API enables interaction with the various nodes in a blockchain network - the peers, orderers and MSPs
- It also allows one to package, install, instantiate and upgrade chaincode on the endorsing peer nodes
- The Hyperledger Fabric language-specific SDKs abstract the specifics of the Hyperledger Fabric API to facilitate application development, though it can be used to manage a chaincode’s lifecycle
- The Hyperledger Fabric API can be accessed directly via the CLI
### Packaging
- The chaincode package consists of 3 parts
    - the chaincode, as defined by `ChaincodeDeploymentSpec` or CDS
        - The CDS defines the chaincode package in terms of the code and other properties such as name and version
    - an optional instantiation policy which can be syntactically described by the same policy used for endorsement and described in Endorsement policies, and
    - a set of signatures by the entities that “own” the chaincode
- The signatures serve the following purposes:
    - to establish an ownership of the chaincode
    - to allow verification of the contents of the package
    - to allow detection of package tampering
- The creator of the instantiation transaction of the chaincode on a channel is validated against the instantiation policy of the chaincode
#### Creating the package
- There are 2 approaches to packaging chaincode
    - One for when you want to have multiple owners of a chaincode, and hence need to have the chaincode package signed by multiple identities
        - This workflow requires that we initially create a signed chaincode package (a `SignedCDS`) which is subsequently passed serially to each of the other owners for signing
    - The simpler workflow is for when you are deploying a `SignedCDS` that has only the signature of the identity of the node that is issuing the install transaction
- Create a signed chaincode package
	
    ```bash
    peer chaincode package -n mycc -p github.com/hyperledger/fabric/examples/chaincode/go/example02/cmd -v 0 -s -S -i "AND('OrgA.admin')" ccpack.out
    ```

    - The `-s` option creates a package that can be signed by multiple owners as opposed to simply creating a raw CDS
        - When `-s` is specified, the `-S` option must also be specified if other owners are going to need to sign. Otherwise, the process will create a `SignedCDS` that includes only the instantiation policy in addition to the CDS
    - The `-S` option directs the process to sign the package using the MSP identified by the value of the `localMspid` property in `core.yaml`
        - The `-S` option is optional
        - However if a package is created without a signature, it cannot be signed by any other owner using the `signpackage` command
    - The optional `-i` option allows one to specify an instantiation policy for the chaincode
        - The instantiation policy has the same format as an endorsement policy and specifies which identities can instantiate the chaincode
        - If no policy is provided, the default policy is used, which only allows the `admin` identity of the peer’s MSP to instantiate chaincode
#### Packaging signing
- A chaincode package that was signed at creation can be handed over to other owners for inspection and signing
- The workflow supports out-of-band signing of chaincode package
- The <u>[`ChaincodeDeploymentSpec`](https://github.com/hyperledger/fabric/blob/master/protos/peer/chaincode.proto#L78)</u> may be optionally be signed by the collective owners to create a <u>[`SignedChaincodeDeploymentSpec`](https://github.com/hyperledger/fabric/blob/master/protos/peer/signed_cc_dep_spec.proto#L26)</u> (or SignedCDS)
- The SignedCDS contains 3 elements
    - The CDS contains the source code, the name, and version of the chaincode
    - An instantiation policy of the chaincode, expressed as endorsement policies
        - This endorsement policy is determined out-of-band to provide proper MSP principals when the chaincode is instantiated on some channels
        - If the instantiation policy is not specified, the default policy is any MSP administrator of the channel
    - The list of chaincode owners, defined by means of Endorsement
- Each owner endorses the `ChaincodeDeploymentSpec` by combining it with that owner’s identity (e.g. certificate) and signing the combined result
- A chaincode owner can sign a previously created signed package using the following command
	
    ```bash
    peer chaincode signpackage ccpack.out signedccpack.out
    ```

    - `signedccpack.out` contains an additional signature over the package signed using the Local MSP
#### Installing chaincode
- The `install` transaction packages a chaincode’s source code into a prescribed format called a `ChaincodeDeploymentSpec` (or CDS) and installs it on a peer node that will run that chaincode
    - You must install the chaincode on each endorsing peer node of a channel that will run your chaincode
- When the `install` API is given simply a `ChaincodeDeploymentSpec`, it will default the instantiation policy and include an empty owner list
- Chaincode should only be installed on endorsing peer nodes of the owning members of the chaincode to protect the confidentiality of the chaincode logic from other members on the network
    - Those members without the chaincode, can’t be the endorsers of the chaincode’s transactions; that is, they can’t execute the chaincode. However, they can still validate and commit the transactions to the ledger
- To install a chaincode, send a `SignedProposal` to the lifecycle system chaincode (LSCC) 
    - For example, to install the `sacc` sample chaincode described in section Simple Asset Chaincode using the CLI, the command would look like the following `peer chaincode install -n asset_mgmt -v 1.0 -p sacc`
    - The CLI internally creates the `SignedChaincodeDeploymentSpec` for `sacc` and sends it to the local peer, which calls the `Install` method on the LSCC
    - The argument to the `-p` option specifies the path to the chaincode, which must be located within the source tree of the user’s `GOPATH`, e.g. `$GOPATH/src/sacc`
    - Note if using `-l node` or `-l java` for node chaincode or java chaincode, use `-p` with the absolute path of the chaincode location
- In order to install on a peer, the signature of the `SignedProposal` must be from 1 of the peer’s local MSP administrators
#### Instantiate
- The `instantiate` transaction invokes the lifecycle System Chaincode (LSCC) to create and initialize a chaincode on a channel
- This is a chaincode-channel binding process: a chaincode may be bound to any number of channels and operate on each channel individually and independently
    - Regardless of how many other channels on which a chaincode might be installed and instantiated, state is kept isolated to the channel to which a transaction is submitted
- The creator of an `instantiate` transaction must satisfy the instantiation policy of the chaincode included in `SignedCDS` and must also be a writer on the channel, which is configured as part of the channel creation
    - This is important for the security of the channel to prevent rogue entities from deploying chaincodes or tricking members to execute chaincodes on an unbound channel
- The instantiate transaction also sets up the endorsement policy for that chaincode on the channel. The endorsement policy describes the attestation requirements for the transaction result to be accepted by members of the channel
- After being successfully instantiated, the chaincode enters the active state on the channel and is ready to process any transaction proposals of type `ENDORSER_TRANSACTION`
#### Upgrade
- A chaincode may be upgraded any time by changing its version, which is part of the `SignedCDS`
    - Other parts, such as owners and instantiation policy are optional
- However, the chaincode name must be the same; otherwise it would be considered as a totally different chaincode
- Prior to upgrade, the new version of the chaincode must be installed on the required endorsers
- Upgrade is a transaction similar to the instantiate transaction, which binds the new version of the chaincode to the channel
    - Other channels bound to the old version of the chaincode still run with the old version. In other words, the `upgrade` transaction only affects one channel at a time, the channel to which the transaction is submitted
- Multiple versions of a chaincode may be active simultaneously. The upgrade process doesn’t automatically remove the old versions, so user must manage this for the time being
- There’s one subtle difference with the `instantiate` transaction: the `upgrade` transaction is checked against the current chaincode instantiation policy, not the new policy (if specified)
    - This is to ensure that only existing members specified in the current instantiation policy may upgrade the chaincode
- During upgrade, the chaincode `Init` function is called to perform any data related updates or re-initialize it, so care must be taken to avoid resetting states when upgrading chaincode
#### Stop and start
- `stop` and `start` lifecycle transactions have not yet been implemented
- However, you may stop a chaincode manually by removing the chaincode container and the `SignedCDS` package from each of the endorsers
    - This is done by deleting the chaincode’s container on each of the hosts or virtual machines on which the endorsing peer nodes are running, and then deleting the SignedCDS from each of the endorsing peer nodes
    	
        ```bash
        docker rm -f <container id>
        rm /var/hyperledger/production/chaincodes/<ccname>:<ccversion>
        ```
    
        - In order to delete the CDS from the peer node, you would need to enter the peer node’s container first
- Stop would be useful in the workflow for doing upgrade in controlled manner, where a chaincode can be stopped on a channel on all peers before issuing an upgrade
### System chaincode
- Unlike a user chaincode, a system chaincode is not installed and instantiated using proposals from SDKs or CLI. It is registered and deployed by the peer at start-up
- System chaincode has the same programming model except that it runs within the peer process rather than in an isolated container like normal chaincode
- Therefore, system chaincode is built into the peer executable and doesn’t follow the same lifecycle described above
    - In particular, `install`, `instantiate` and `upgrade` do not apply to system chaincodes
- The purpose of system chaincode is to shortcut gRPC communication cost between peer and chaincode, and tradeoff the flexibility in management
    - For example, a system chaincode can only be upgraded with the peer binary. It must also register with a fixed set of parameters compiled in and doesn’t have endorsement policies or endorsement policy functionality
- System chaincode is used in Hyperledger Fabric to implement a number of system behaviors so that they can be replaced or modified as appropriate by a system integrator
- The current list of system chaincodes:
    1. LSCC Lifecycle system chaincode handles lifecycle requests
    2. CSCC Configuration system chaincode handles channel configuration on the peer side
    3. QSCC Query system chaincode provides ledger query APIs such as getting blocks and transactions
- The former system chaincodes for endorsement and validation have been replaced by the pluggable endorsement and validation function
- System chaincodes can be linked to a peer in two ways: statically, and dynamically using Go plugins
#### System chaincode plugins
- A system chaincode is a program written in Go and loaded using the Go plugin package
- A plugin includes a main package with exported symbols and is built with the command `go build -buildmode=plugin`
- Every system chaincode must implement the `Chaincode` Interface and export a constructor method that matches the signature `func New() shim.Chaincode` in the `main` package
    - An example can be found in the repository at `examples/plugin/scc`
- Existing chaincodes such as the QSCC can also serve as templates for certain features, such as access control, that are typically implemented through system chaincodes
- The existing system chaincodes also serve as a reference for best-practices on things like logging and testing
- The Go standard library requires that a plugin must include the same version of imported packages as the host application (Fabric, in this case)
- Plugins are configured in the `chaincode.systemPlugin` section in `core.yaml`
	
    ```yaml
    chaincode:
        systemPlugins:
            - enabled: true
              name: mysyscc
              path: /opt/lib/syscc.so
              invokableExternal: true
              invokableCC2CC: true
    ```

- A system chaincode must also be whitelisted in the `chaincode.system` section in `core.yaml`
	
    ```yaml
    chaincode:
        system:
            mysyscc: enable
    ```