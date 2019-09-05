- 账本不可纂改地记录着由智能合约生成的交易
- Fabric 的账本包括世界状态和交易日志两部分
    - 世界状态描述某个时间点时账本的状态，是账本的数据库
    - 交易日志记录了形成当前世界状态的所有交易，是世界状态的更新历史
- 存储世界状态的数据库可配置，默认是 LevelDB (key-value database)
- 交易日志不需要可插拔
<!-- - Hyperledger Fabric has a ledger subsystem comprising two components: the **world state** and the **transaction log**
    - The world state component describes **the state of the ledger at a given point in time**. It’s the database of the ledger
    - The transaction log component records all transactions which have resulted in the current value of the world state; it’s **the update history** for the world state -->
<!-- - The ledger is comprised of a blockchain (‘chain’) to store the immutable, sequenced record in blocks, as well as a state database to maintain current fabric state -->
<!-- - The ledger has a replaceable data store for the world state
    - By default, this is a LevelDB key-value store database -->
<!-- - The transaction log does not need to be pluggable
    - It simply **records the before and after values** of the ledger database being used by the blockchain network -->
<!-- - The immutable, shared ledger encodes the entire transaction history for each channel, and includes SQL-like query capability for efficient auditing and dispute resolution -->
- Fabric 用键值对描述资产
    - 资产可以用二进制、JSON 的形式表示
    - 可以用 Hyperledger Composer 工具定义、使用资产
- 交易是调用合约产生的结果，使得一系列键值对提交到账本
- 每个 channel 有各自的账本
- 账本包含一个配置块，配置块定义了策略、访问控制等
<!-- - Assets are represented in Hyperledger Fabric as a collection of key-value pairs, with state changes recorded as transactions on a channel ledger
- Assets can be represented in binary and/or JSON form
    - Assets in Hyperledger Fabric applications can be easily defined and used using the Hyperledger Composer tool -->
<!-- - State transitions are a result of **chaincode invocations** (‘transactions’) submitted by participating parties. Each transaction results in a set of asset key-value pairs that are committed to the ledger as creates, updates, or deletes -->
<!-- - There is **one ledger per channel** -->
<!-- - Each peer maintains a copy of the ledger for each channel of which they are a member -->
<!-- - Features of a Fabric ledger
    - Query and update ledger using key-based lookups, range queries, and composite key queries
    - Read-only queries using a rich query language (if using CouchDB as state database)
    - Read-only history queries — Query ledger history for a key, enabling data provenance scenarios
    - Transactions consist of the versions of keys/values that were read in chaincode (read set) and keys/values that were written in chaincode (write set)
    - Transactions contain **signatures of every endorsing peer** and are submitted to ordering service
    - Transactions are ordered into blocks and are “delivered” from an ordering service to peers
    - Peers **validate** transactions against endorsement policies and **enforce** the policies
    - Prior to appending a block, a **versioning check** is performed to **ensure that states for assets that were read have not changed since chaincode execution time**
    - There is immutability once a transaction is validated and committed
    - A channel’s ledger contains **a configuration block** defining policies, access control lists, and other pertinent information -->
- 一个账本可以有多个可以访问账本的智能合约
- 账本不保存业务对象，保存的是这些对象的真相（facts）- 当前状态和通向当前状态的交易历史
- 账本包括 2 部分
    1. 世界状态
        - 账本的当前值
        - 用键值对表示
        - 每个 chaincode 都有它自己的独立于其它 chaincode 的世界状态
            - World states are in a namespace so that only **smart contracts within the same chaincode** can access a given namespace
    2. 区块链
        - 交易日志，记录了所有达到点前世界状态的所有变更
        - [ ] A blockchain is not namespaced. It contains transactions from many different smart contract namespaces
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/ledger.diagram.3.png)
- 世界状态用数据库实现
- 版本号在每次状态改变后都会递增，用于检查当前状态和背书交易时是否一致
- 账本刚创建时，世界状态为空
- 由于区块链记录了所有世界状态的变更，可以很方便地从区块生成世界状态
- 每个区块的区块头包含区块内交易的 hash 和前一个区块头里的 hash
- 区块链用文件实现
- 创世块不包含用户交易，包含一条配置交易（configuration transaction）
- 区块
    1. 区块头
        1. 区块编号（起始为 0）
        2. 当前区块 hash - 当前区块内所有交易的 hash
        3. 父区块 hash - 父区块所有交易 hash 的副本
        ![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/ledger.diagram.4.png)
    2. 区块数据
        - 排好序的交易
    4. 区块 metadata
        - 包括区块的写入时间，证书，区块写入者的公钥和签名
        - 随后区块提交者会对每一条交易添加有效/无效的标识；该信息包括在 hash 里，hash 在区块创建时就已生成
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/ledger.diagram.5.png) 
- 交易
    - Header
        - 包含交易的 metadata 信息，如 chaincode 名称和版本号
    - Signature
        - 由客户端的私钥生成的密码学签名，用于校验交易详情未被篡改
    - Proposal
        - 应用向智能合约提供的创建了 proposed leger update 的入参
        - 入参和当前世界状态共同决定新的世界状态
    - Response
        - 以 read write set (RW-set) 的形式记录变更前后世界状态
        - 是智能合约的输出（output）
    - Endorsements
        - 一组背书策略要求的组织签过名的 transaction response
- DB
  - A LevelDB database is closely co-located with a network node – it is embedded within the same operating system process
  - CouchDB is a particularly appropriate choice when ledger states are structured as JSON documents because CouchDB supports the rich queries and update of richer data types often found in business transactions
      - Implementation-wise, CouchDB runs in a separate operating system process, but there is still a 1:1 relation between a peer node and a CouchDB instance