# Introduction
- 整体上看，Fabric 包括以下几个模块化组件
    1. 可插拔共识服务，对交易的顺序建立共识，并将区块广播给 peer
    2. 可插拔 MSP，负责将网络中的实体和密码学身份建立关联
    3. 可选的 peer-to-peer gossip 服务，将排序服务产生的区块在 peer 节点之间传播
    4. 智能合约，运行在容器环境以实现隔离
    5. 账本，支持多种数据库
    6. 可插拔的共识和校验策略
<!-- - At a high level, Fabric is comprised of the following modular components
    1. A pluggable **ordering service** establishes consensus on the **order of transactions** and then **broadcasts blocks to peers**
    2. A pluggable membership service provider is responsible for **associating entities in the network with cryptographic identities**
    3. An optional peer-to-peer gossip service disseminates the blocks output by [ ] ordering service to other peers
    4. Smart contracts (“chaincode”) run within a container environment (e.g. Docker) for **isolation**
        - They can be written in standard programming languages but **do not have direct access to the ledger state**
    5. The ledger can be configured to support a variety of DBMSs
    6. A pluggable [ ] endorsement and validation policy enforcement that can be independently configured per application -->

<!-- - Most existing smart-contract capable blockchain platforms follow an order-execute architecture in which the consensus protocol:
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
    - This first phase also eliminates any non-determinism, as inconsistent results can be **filtered out before ordering** -->
- 无权限控制的采用 PoW 作为共识的公链，交易执行每个节点
    - 所有交易、合约对于网络上的所有节点都可见，不存在保密性（用私密性换 PoW 带来的 byzantine fault tolerant consensus）
        - 加密数据是一种提供私密性的机制，但若有足够的时间和算力，依然可以破解
        - ZKP，缺点是计算 ZKP 需要很多事件和计算资源，造成性能损耗（性能换私密性）
<!-- - In a public, permissionless blockchain network that leverages PoW for its consensus model, transactions are executed on every node
    - This means that neither can there be confidentiality of the contracts themselves, nor of the transaction data that they process
    - Every transaction, and the code that implements it, is visible to every node in the network. In this case, we have traded confidentiality of contract and data for byzantine fault tolerant consensus delivered by PoW
    - Encrypting data is one approach to providing confidentiality; however, in a permissionless network leveraging PoW for its consensus, the encrypted data is sitting on every node. Given enough time and computational resource, the encryption could be broken
    - [ ] Zero knowledge proofs (ZKP) are another area of research being explored to address this problem, the trade-off here being that, presently, computing a ZKP requires considerable time and computational resources. Hence, the trade-off in this case is **performance for confidentiality** -->
<!-- - Fabric has added support for private data and is working on zero knowledge proofs (ZKP) available in the future -->
- DLT: distributed ledger technology
- > [A peer reviewed paper that evaluated the architecture and performance of Hyperledger Fabric](https://arxiv.org/abs/1801.10228v1)
<!-- # Blockchain
- The information recorded to a blockchain is **append-only**, using cryptographic techniques that guarantee that once a transaction has been added to the ledger it cannot be modified
- Unlike today’s systems, where a participant’s private programs are used to update their private ledgers, a blockchain system has shared programs to update shared ledgers -->
<!-- ## Consensus
- The process of **keeping the ledger transactions synchronized across the network** — to ensure that **ledgers update only when transactions are approved by the appropriate participants**, and that when ledgers do update, they **update with the same transactions in the same order** — is called consensus -->
<!-- ## Fabric
- The members of a Hyperledger Fabric network enroll through a trusted Membership Service Provider (MSP)
- Hyperledger Fabric supports networks where privacy (using channels) is a key operational requirement as well as networks that are comparatively open
    - A channel allows a group of participants to create a **separate ledger** of transactions -->

- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/ledger.html

- 大多数支持智能合约的区块链采用排序-执行的架构（Ethereum, Tendermint, Chain, Quorum）
    1. 共识协议校验交易并将交易排好序后传播给所有 peer 节点
    2. 每个 peer 节点按顺序执行交易
    - 采用此种架构的区块链的智能合约必须具有确定性，否则共识可能永远不能达成
        - 为了解决非确定的问题，许多平台要求智能合约要用非标准的语言或是特定领域语言来编写，来消除非确定性的操作
    - 所有交易在所有 peer 上顺序执行带来性能和扩展性的问题
    - 智能合约代码在系统的所有节点上执行，为保证系统安全性而所需的要求很高
- Fabric 采用执行-排序-校验架构
    1. 执行一个交易，校验其正确性，为交易背书
    2. 用共识合约为交易排序
    3. 按照每个应用对应的背书策略策略校验交易，通过后提交到账本
        - 一个应用的背书策略指定哪些 peer 节点、多少 peer 节点需要为为一个智能合约的执行结果背书
    - Fabric 在排序前先执行一次交易
        - 消除了不确定性，不一致的结果在排序之前就已被剔除
        - 使得每个 peer 可以同时执行多个交易，同时加快交易向排序服务递送的速度
    - 每个交易只需要被背书策略要求的一部分 peer 执行（即背书）
<!-- - Most existing smart-contract capable blockchain platforms follow an order-execute architecture in which the consensus protocol:
    1. validates and orders transactions then propagates them to all peer nodes
    2. **each peer** then executes the transactions **sequentially**
    - > Ethereum, Tendermint, Chain, Quorum
    - Smart contracts executing in a blockchain that operates with the order-execute architecture must be **deterministic**; otherwise, **consensus might never be reached**
    - To address the non-determinism issue, many platforms require that the smart contracts be written in a non-standard, or domain-specific language (such as Solidity) so that non-deterministic operations can be eliminated
    - Further, since all transactions are executed sequentially by all nodes, performance and scale is limited
    - The fact that the smart contract code **executes on every node** in the system demands that complex measures be taken to protect the overall system from potentially malicious contracts in order to ensure resiliency of the overall system -->

<!-- - Fabric follows an execute-order-validate architecture
    1. **execute** a transaction and check its correctness, thereby **endorsing** it
    2. **order** transactions via a (pluggable) consensus protocol
    3. **validate transactions against** an application-specific **endorsement policy before committing them to the ledger**
    - **Fabric executes transactions before reaching final agreement on their order**
    - In Fabric, an application-specific endorsement policy **specifies** which peer nodes, or how many of them, need to **vouch for the correct execution of a given smart contract**
        - Thus, each transaction need only be executed (endorsed) by the subset of the peer nodes necessary to satisfy the transaction’s endorsement policy
    - This first phase also eliminates any non-determinism, as inconsistent results can be **filtered out before ordering** -->
<!-- - The ordering of transactions is **delegated** to a modular component for consensus that is logically decoupled from the peers that execute transactions and maintain the ledger
    - Since consensus is modular, its implementation can be tailored to the trust assumption of a particular deployment or solution -->

<!-- - Transaction **execution** is separated from transaction **ordering** and **commitment** -->
<!-- - Executing transactions prior to ordering them enables each peer node to process multiple transactions simultaneously. This **concurrent execution** increases processing efficiency on each peer and accelerates delivery of transactions to the ordering service -->
<!-- - The ***division*** of labor unburdens ordering nodes from the demands of transaction execution and ledger maintenance, while peer nodes are freed from ordering (consensus) workloads
    - Orderer 只做排序，不执行交易、维护账本
    - Peer 只执行交易、维护账本，不做排序 -->
- This bifurcation of roles also limits the processing required for authorization and authentication; all peer nodes do not have to trust all ordering nodes, and vice versa, so [ ] processes on one can run independently of verification by the other
<!-- - Transaction flow: proposal, endorsement, ordering, validation, commitment -->