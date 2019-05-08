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
- 账本分叉（ledger "fork"）
    - 网络参与方对交易的顺序无法达成共识
- 无权限区块链网络依靠基于概率的共识算法来最大程度保证账本一致性
    - Rely on probabilistic consensus algorithms which eventually guarantee ledger consistency to a high degree of probability
- Fabric 依靠的是确定性公式算法（deterministic consensus algorithms）
- orderer 除了排序的角色之外，还维护允许创建 channel 的组织列表（称作集团，consortium）
    - 列表保存在 “orderer system channel” (also known as the “ordering system channel”) 的配置中
    - 默认情况下，只有 orderer 的管理员可以编辑此列表和保存该列表的 channel
    - 一个 orderer 可以维护多个组织列表
- Orderers also enforce basic access control for channels, restricting who can read and write data to them, and who can configure them
- Remember that who is authorized to modify a configuration element in a channel is subject to the policies that the relevant administrators set when they created the consortium or the channel
- Configuration transactions are processed by the orderer, as it needs to know the current set of policies to execute its basic form of access control
    - In this case, the orderer processes the configuration update to make sure that the requestor has the proper administrative rights
    - If so, the orderer validates the update request against the existing configuration, generates a new configuration transaction, and packages it into a block that is relayed to all peers on the channel
    - The peers then process the configuration transactions in order to verify that the modifications approved by the orderer do indeed satisfy the policies defined in the channel
- orderer 节点从属于一个组织
## 共识实现
### Solo
- 单个 orderer 节点；not fault tolerant
- 适用于测试应用和合约；创建 proofs of concept (PoC) 网络
- 不能用于生产环境
    - For networks which wish to start with only a single ordering node but might wish to grow in the future, a **single node Raft cluster** is a better option
### Raft
- v1.4.1 引入
- 是 crash fault tolerant (CFT) 的排序服务
- 基于 `etcd` 中的 Raft 协议来实现
- 遵循 "leader and follower" 模型
    - A leader node is elected (per channel) and its decisions are replicated by the followers
- 相较基于 Kafka 的排序服务更易搭建和管理
- 运行不同的组织贡献 orderer 节点
- 可以容忍损失节点，只要大多数（quorum，法定人数）节点还存在，则 Raft 依然是 CFT
- 与 Kafka 的区别
    - 易部署
    - Kafka 和 Zookeeper 并非为大型网络设计，实际使用中需让一个组织运行 Kafka 集群，因此 Raft 的去中心化更好
        - Having ordering nodes run by different organizations when using Kafka (which Fabric supports) doesn’t give you much in terms of decentralization because the nodes will all go to the same Kafka cluster which is under the control of a single organization
    - Kafka 与 Fabric 兼容，但需额外的 Kafka 和 Zookeeper 镜像并要会使用；Raft 则是 Fabric 原生支持的
    - Where Kafka uses a pool of servers (called “Kafka brokers”) and the admin of the orderer organization specifies how many nodes they want to use on a particular channel, Raft allows the users to specify which ordering nodes will be deployed to which channel. In this way, peer organizations can make sure that, if they also own an orderer, this node will be made a part of a ordering service of that channel, rather than trusting and depending on a central admin to manage the Kafka nodes
    - Raft is the first step toward Fabric’s development of a byzantine fault tolerant (BFT) ordering service
- > Similar to Solo and Kafka, a Raft ordering service can lose transactions after acknowledgement of receipt has been sent to a client. For example, if the leader crashes at approximately the same time as a follower provides acknowledgement of receipt. Therefore, application clients should listen on peers for transaction commit events regardless (to check for transaction validity), but extra care should be taken to ensure that the client also gracefully tolerates a timeout in which the transaction does not get committed in a configured timeframe. Depending on the application, it may be desirable to resubmit the transaction or collect a new set of endorsements upon such a timeout
- 每个 channel 运行一个单独的 Raft 实例，每个实例有各自的 leader
- 所有 Raft 节点必须是系统 channel 的一部分，但不一定要是所有应用 channel 的一部分
    - Channel creators (and channel admins) have the ability to pick a subset of the available orderers and to add or remove ordering nodes as needed (as long as only a single node is added or removed at a time)
- Raft 中，交易（proposal 或是配置更新）会自动由接收到该交易的 orderer 转发给 leader，因此 peer 和应用不需要知道哪个是 leader 节点
- Raft 节点可能有 3 个状态：follower, candidate, leader
- 所有节点初始都是 follower，此时可以接收 leader 发来的 log entries，或者为 leader 投票
- 若一定时间没有收到 log entries 或心跳，follower 自荐为 candidate，并要求其它节点来投票；若收到法定人数的投票，则当选 leader
- 新 leader 接收新的 log entries 并复制给 followers
- http://thesecretlivesofdata.com/raft/
- 用户通过 snapshotting 定义 log 里存储多少字节的数据，这一部分数据会和一定数目的区块匹配
    - 该数目由区块大小决定
    - 只有完整的区块才会被保存在 snapshot 里
    - > For example, let’s say lagging replica `R1` was just reconnected to the network. Its latest block is `100`. Leader `L` is at block `196`, and is configured to snapshot at amount of data that in this case represents 20 blocks. `R1` would therefore receive block `180` from `L` and then make a Deliver request for blocks `101` to `180` ([ ] 4 snapshots?). Blocks `180` to `196` would then be replicated to `R1` through the normal Raft protocol
#### Log entry
- The primary unit of work in a Raft ordering service is a “log entry”, with the full sequence of such entries known as the “log”
- We consider the log consistent if a majority (a quorum, in other words，需大于 50%) of members agree on the entries and their order, making the logs on the various orderers replicated
#### Consenter set (赞成者集)
- The ordering nodes actively participating in the consensus mechanism for a given channel and receiving replicated logs for the channel
- This can be all of the nodes available (either in a single cluster or in multiple clusters contributing to the system channel), or a subset of those nodes
#### Finite-State Machine (FSM)
- Every ordering node in Raft has an FSM and collectively they’re used to ensure that the sequence of logs in the various ordering nodes is deterministic (written in the same sequence)
#### Quorum
- Describes the minimum number of consenters that need to affirm a proposal so that transactions can be ordered
- For every consenter set, this is a **majority** of nodes
    - In a cluster with five nodes, three must be available for there to be a quorum
- If a quorum of nodes is unavailable for any reason, the ordering service cluster becomes unavailable for both read and write operations on the channel, and no new logs can be committed
#### Leader
- At any given time, a channel’s consenter set elects a **single** node to be the leader.
- The leader is responsible for ingesting new log entries, replicating them to follower ordering nodes, and managing when an entry is considered committed
- This is not a special type of orderer. It is only a **role** that an orderer may have at certain times, and then not others, as circumstances determine
#### Follower
- Followers receive the logs from the leader and replicate them deterministically, ensuring that logs remain consistent
- The followers also receive **“heartbeat”** messages from the leader
    - In the event that the leader stops sending those message for a configurable amount of time, the followers will initiate a leader election and one of them will be elected the new leader
### Kafka
- 是 crash fault tolerant (CFT) 的排序服务
- 遵循 "leader and follower" 模型
- Kafka 集群的管理，包括任务调度，cluster membership，访问控制，controller election 等，由 Zookeeper 管理
- https://kafka.apache.org/quickstart
- [sample configuration file](https://github.com/hyperledger/fabric/blob/release-1.1/bddtests/dc-orderer-kafka.yml)