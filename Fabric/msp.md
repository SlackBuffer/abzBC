- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/msp.html

- In Fabric, an application-specific endorsement policy **specifies** which peer nodes, or how many of them, need to **vouch for the correct execution of a given smart contract**
        - Thus, each transaction need only be executed (endorsed) by the subset of the peer nodes necessary to satisfy the transaction’s endorsement policy
- Access control lists can be used to provide additional layers of permission through authorization of specific network operations
    - For example, a specific user ID could be permitted to invoke a chaincode application, but be blocked from deploying new chaincode

- Public key infrastructure is used to generate cryptographic certificates which are tied to organizations, network components, and end users or client applications
    - As a result, data access control can be manipulated and governed on the broader network and on channel levels