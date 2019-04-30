- 系统合约和普通合约的不同之处在于系统合约定义了整个 channel 的操作参数（operating parameters；入参）
    - 生命周期和配置的系统合约定义了 channel 的规则
    - 背书和校验的系统合约定义了参与校验、背书交易的要求（defines the requirements for endorsing and validating transactions）
- 合约通过发起交易 proposal 对账本的世界状态（数据库）进行操作
- 合约的执行结果是一组可以提交到网络并应用到所有 peer 节点的账本上的 key-value writes (write set)
- 合约可以只安装在需要访问账本的资源状态的 peer 上
<!-- - Chaincodes can be installed only on specific peers that need to access the asset states of a ledger -->
<!-- - To support the **consistent update of information**, and to enable **a whole host of ledger functions** (transacting, querying, etc) — a blockchain network uses smart contracts to provide **controlled access to the ledger** -->
<!-- - Smart contracts are not only a key mechanism for encapsulating information and keeping it simple across the network, they can also be written to allow participants to execute certain aspects of transactions automatically
    - A smart contract can, for example, be written to stipulate the cost of shipping an item where the shipping charge changes depending on how quickly the item arrives
    -  With the ***terms agreed to by both parties and written to the ledger***, the appropriate funds change hands automatically when the item is received -->
<!-- - A chaincode functions as a **trusted distributed application** that gains its security/trust from the blockchain and the underlying consensus among the peers
    - Chaincode is software defining an asset or assets, and the transaction instructions for modifying the asset(s)
    - It is the **business logic** of a blockchain application -->
<!-- - Smart contracts are written in chaincode and are invoked by **an application external to the blockchain** when that application needs to interact with the ledger
    - In most cases, chaincode interacts only with the database component of the ledger, the world state (querying it, for example), and not the transaction log -->
<!-- - Chaincode applications **encode logic** that is invoked by specific types of transactions on the channel
    - > Chaincode that defines **parameters** for a change of asset ownership ensures that all transactions that transfer ownership are subject to the same rules and requirements
    - **System chaincode** is distinguished as chaincode that defines operating parameters for **the entire channel**
        1. Lifecycle and configuration system chaincode defines the rules for the channel
        2. Endorsement and validation system chaincode defines the requirements for endorsing and validating transactions -->

<!-- - Chaincode functions execute against the ledger’s current state database and are ***initiated through a transaction proposal*** -->
<!-- - Chaincodes do not have **direct** access to the ledger state -->
<!-- - Chaincode execution results in a set of key-value writes (write set) that can be submitted to the network and applied to the ledger on all peers -->