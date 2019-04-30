- 共识是保证账本的交易在网络上保持同步的过程
- 共识保证了账本只在特定的参与方认可交易后，才将相同的交易以相同的顺序更新到账本
- CFT, BFT, PBFT
- Fabric 使用的共识，无需数字货币激励挖矿或奖励合约的执行
- Fabric 目前提供 2 种 CFT 的排序实现
    1. [etcd](https://coreos.com/etcd/) library of the [Raft protocol](https://raft.github.io/raft.pdf)
    2. [Kafka](https://kafka.apache.org/) (which uses [Zookeeper](https://zookeeper.apache.org/) internally)
- 一个 Fabric 可以同时有多个不同的排序服务
- 共识在交易顺序和结果通过特定策略的检查后最终达成
- Consensus is achieved ultimately when the order and results of a block's transactions have met the explicit policy criteria checks
    <!-- - These checks and balances take place during the lifecycle of a transaction, and include the usage of endorsement policies to dictate which specific members must endorse a certain transaction class, as well as system chaincodes to ensure that these policies are enforced and upheld -->
    <!-- - Prior to commitment, peers will employ these chaincodes to make sure that enough endorsement are present, and they were derived from the appropriate entities -->
    - Moreover, a **versioning check** will take place during which the current state of the ledger is agreed or consented upon, before any blocks containing transactions are appended to the ledger
        - This final check provides protection against double spend operations and other threats that might compromise data integrity, and allows for functions to be executed against non-static variables
<!-- - The process of **keeping the ledger transactions synchronized** across the network — to ensure that ledgers update only when transactions are **approved by the appropriate participants**, and that when ledgers do update, they update with the **same transactions** in the **same order** — is called consensus -->

<!-- - Transactions must be written to the ledger **in the order in which they occur**, even though they might be between different sets of participants within the network -->
<!-- - For this to happen, the **order of transactions must be established** and a **method for rejecting bad transactions** that have been inserted into the ledger in error (or maliciously) must be put into place
    - PBFT (Practical Byzantine Fault Tolerance) can provide a mechanism for file replicas to communicate with each other to keep each copy consistent, even in the event of corruption -->
<!-- - In order to mitigate this absence of trust, permissionless blockchains typically employ a “mined” native cryptocurrency or transaction fees to provide economic incentive to **offset the extraordinary costs of participating in a form of byzantine fault tolerant consensus** based on “proof of work” (PoW)
    - In Bitcoin, ordering happens through a process called mining where competing computers race to solve a cryptographic puzzle which [ ] **defines the order** that all processes subsequently build upon
    - Requiring protocols like “proof of work” to **validate transactions and secure the network** [ ] how so and how to?
    - [ ] 公链的共识；PoW 为何能解决共识问题 -->

<!-- - [ ] Fabric can leverage consensus protocols that do not require a native cryptocurrency to incent costly mining or to fuel smart contract execution -->
<!-- - Pluggable consensus protocols enable the platform to be more effectively customized to fit particular use cases and trust models
    - [ ] crash fault-tolerant (CFT) consensus protocol
    - [ ] byzantine fault tolerant (BFT) consensus protocol -->
<!-- - Fabric currently offers 2 CFT ordering service implementations. The first is based on the [etcd](https://coreos.com/etcd/) library of the [Raft protocol](https://raft.github.io/raft.pdf). The other is [Kafka](https://kafka.apache.org/) (which uses [Zookeeper](https://zookeeper.apache.org/) internally)
    - These are not mutually exclusive. **A Fabric network can have multiple ordering services** supporting different applications or application requirements -->
<!-- - Consensus is defined as the full-circle verification of the correctness of a set of transactions comprising a block -->
- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/txflow.html
- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/arch-deep-dive.html#swimlane
