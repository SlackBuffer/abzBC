- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/ledger.html

- Most existing smart-contract capable blockchain platforms follow an order-execute architecture in which the consensus protocol:
    1. validates and orders transactions then propagates them to all peer nodes
    2. **each peer** then executes the transactions **sequentially**
    - > Ethereum, Tendermint, Chain, Quorum
    - Smart contracts executing in a blockchain that operates with the order-execute architecture must be **deterministic**; otherwise, **consensus might never be reached**
    - To address the non-determinism issue, many platforms require that the smart contracts be written in a non-standard, or domain-specific language (such as Solidity) so that non-deterministic operations can be eliminated
    - Further, since all transactions are executed sequentially by all nodes, performance and scale is limited
    - The fact that the smart contract code **executes on every node** in the system demands that complex measures be taken to protect the overall system from potentially malicious contracts in order to ensure resiliency of the overall system

- Fabric follows an execute-order-validate architecture
    1. **execute** a transaction and check its correctness, thereby **endorsing** it
    2. **order** transactions via a (pluggable) consensus protocol
    3. **validate transactions against** an application-specific **endorsement policy before committing them to the ledger**
    - **Fabric executes transactions before reaching final agreement on their order**
    - In Fabric, an application-specific endorsement policy **specifies** which peer nodes, or how many of them, need to **vouch for the correct execution of a given smart contract**
        - Thus, each transaction need only be executed (endorsed) by the subset of the peer nodes necessary to satisfy the transaction’s endorsement policy
    - This first phase also eliminates any non-determinism, as inconsistent results can be **filtered out before ordering**
- The ordering of transactions is **delegated** to a modular component for consensus that is logically decoupled from the peers that execute transactions and maintain the ledger
    - Since consensus is modular, its implementation can be tailored to the trust assumption of a particular deployment or solution

- Transaction **execution** is separated from transaction **ordering** and **commitment**
- Executing transactions prior to ordering them enables each peer node to process multiple transactions simultaneously. This **concurrent execution** increases processing efficiency on each peer and accelerates delivery of transactions to the ordering service
- The ***division*** of labor unburdens ordering nodes from the demands of transaction execution and ledger maintenance, while peer nodes are freed from ordering (consensus) workloads
    - Orderer 只做排序，不执行交易、维护账本
    - Peer 只执行交易、维护账本，不做排序
- This bifurcation of roles also limits the processing required for authorization and authentication; all peer nodes do not have to trust all ordering nodes, and vice versa, so [ ] processes on one can run independently of verification by the other
- Transaction flow: proposal, endorsement, ordering, validation, commitment
- Consensus is achieved ultimately when the order and results of a block's transactions have met the explicit policy criteria checks
    - These checks and balances take place during the lifecycle of a transaction, and include the usage of endorsement policies to dictate which specific members must endorse a certain transaction class, as well as system chaincodes to ensure that these policies are enforced and upheld
    - Prior to commitment, peers will employ these chaincodes to make sure that enough endorsement are present, and they were derived from the appropriate entities
    - Moreover, a versioning check check will take place during which the current state of the ledger is agreed or consented upon, before any blocks containing transactions are appended to the ledger
        - This final check provides protection against double spend operations and other threats that might compromise data integrity, and allows for functions to be executed against non-static variables