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
```