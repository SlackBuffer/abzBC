- To support the **consistent update of information**, and to enable **a whole host of ledger functions** (transacting, querying, etc) — a blockchain network uses smart contracts to provide **controlled access to the ledger**
- Smart contracts are not only a key mechanism for encapsulating information and keeping it simple across the network, they can also be written to allow participants to execute certain aspects of transactions automatically
    - A smart contract can, for example, be written to stipulate the cost of shipping an item where the shipping charge changes depending on how quickly the item arrives
    -  With the ***terms agreed to by both parties and written to the ledger***, the appropriate funds change hands automatically when the item is received
- A chaincode functions as a **trusted distributed application** that gains its security/trust from the blockchain and the underlying consensus among the peers. It is the **business logic** of a blockchain application
    - [ ] 部署哪里；部署几个；谁能部署

- Smart contracts are written in chaincode and are invoked by **an application external to the blockchain** when that application needs to interact with the ledger
    - In most cases, chaincode interacts only with the database component of the ledger, the world state (querying it, for example), and not the transaction log
- Chaincodes do not have **direct** access to the ledger state
- Chaincode applications **encode logic** that is invoked by specific types of transactions on the channel
    - > Chaincode that defines **parameters** for a change of asset ownership ensures that all transactions that transfer ownership are subject to the same rules and requirements
    - **System chaincode** is distinguished as chaincode that defines operating parameters for **the entire channel**
        1. Lifecycle and configuration system chaincode defines the rules for the channel
        2. Endorsement and validation system chaincode defines the requirements for endorsing and validating transactions