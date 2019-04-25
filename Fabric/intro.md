# Introduction
- At a high level, Fabric is comprised of the following modular components
    1. A pluggable **ordering service** establishes consensus on the **order of transactions** and then **broadcasts blocks to peers**
    2. A pluggable membership service provider is responsible for **associating entities in the network with cryptographic identities**
    3. An optional peer-to-peer gossip service disseminates the blocks output by [ ] ordering service to other peers
    4. Smart contracts (“chaincode”) run within a container environment (e.g. Docker) for **isolation**
        - They can be written in standard programming languages but **do not have direct access to the ledger state**
    5. The ledger can be configured to support a variety of DBMSs
    6. A pluggable [ ] endorsement and validation policy enforcement that can be independently configured per application

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

- In a public, permissionless blockchain network that leverages PoW for its consensus model, transactions are executed on every node
    - This means that neither can there be confidentiality of the contracts themselves, nor of the transaction data that they process
    - Every transaction, and the code that implements it, is visible to every node in the network. In this case, we have traded confidentiality of contract and data for byzantine fault tolerant consensus delivered by PoW
    - Encrypting data is one approach to providing confidentiality; however, in a permissionless network leveraging PoW for its consensus, the encrypted data is sitting on every node. Given enough time and computational resource, the encryption could be broken
    - [ ] Zero knowledge proofs (ZKP) are another area of research being explored to address this problem, the trade-off here being that, presently, computing a ZKP requires considerable time and computational resources. Hence, the trade-off in this case is **performance for confidentiality**
- Fabric has added support for private data and is working on zero knowledge proofs (ZKP) available in the future
- DLT: distributed ledger technology
- [A peer reviewed paper that evaluated the architecture and performance of Hyperledger Fabric](https://arxiv.org/abs/1801.10228v1)
# Blockchain
- The information recorded to a blockchain is **append-only**, using cryptographic techniques that guarantee that once a transaction has been added to the ledger it cannot be modified
- Unlike today’s systems, where a participant’s private programs are used to update their private ledgers, a blockchain system has shared programs to update shared ledgers
## Consensus
- The process of **keeping the ledger transactions synchronized across the network** — to ensure that **ledgers update only when transactions are approved by the appropriate participants**, and that when ledgers do update, they **update with the same transactions in the same order** — is called consensus
## Fabric
- The members of a Hyperledger Fabric network enroll through a trusted Membership Service Provider (MSP)
- Hyperledger Fabric supports networks where privacy (using channels) is a key operational requirement as well as networks that are comparatively ope
    - A channel allows a group of participants to create a **separate ledger** of transactions