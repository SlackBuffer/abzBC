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