- The process of **keeping the ledger transactions synchronized** across the network — to ensure that ledgers update only when transactions are **approved by the appropriate participants**, and that when ledgers do update, they update with the **same transactions** in the **same order** — is called consensus

- Transactions must be written to the ledger **in the order in which they occur**, even though they might be between different sets of participants within the network
- For this to happen, the **order of transactions must be established** and a **method for rejecting bad transactions** that have been inserted into the ledger in error (or maliciously) must be put into place
    - PBFT (Practical Byzantine Fault Tolerance) can provide a mechanism for file replicas to communicate with each other to keep each copy consistent, even in the event of corruption
- In order to mitigate this absence of trust, permissionless blockchains typically employ a “mined” native cryptocurrency or transaction fees to provide economic incentive to **offset the extraordinary costs of participating in a form of byzantine fault tolerant consensus** based on “proof of work” (PoW)
    - In Bitcoin, ordering happens through a process called mining where competing computers race to solve a cryptographic puzzle which [ ] **defines the order** that all processes subsequently build upon
    - Requiring protocols like “proof of work” to **validate transactions and secure the network** [ ] how so and how to?
    - [ ] 公链的共识；PoW 为何能解决共识问题

- [ ] Fabric can leverage consensus protocols that do not require a native cryptocurrency to incent costly mining or to fuel smart contract execution
- Pluggable consensus protocols enable the platform to be more effectively customized to fit particular use cases and trust models
    - [ ] crash fault-tolerant (CFT) consensus protocol
    - [ ] byzantine fault tolerant (BFT) consensus protocol
- Fabric currently offers 2 CFT ordering service implementations. The first is based on the [etcd](https://coreos.com/etcd/) library of the [Raft protocol](https://raft.github.io/raft.pdf). The other is [Kafka](https://kafka.apache.org/) (which uses [Zookeeper](https://zookeeper.apache.org/) internally)
    - These are not mutually exclusive. **A Fabric network can have multiple ordering services** supporting different applications or application requirements