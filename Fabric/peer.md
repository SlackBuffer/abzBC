- Peer 上 host 着账本和智能合约的实例
    - 智能合约和账本分别封装了网络上共享的处理逻辑和信息
    <!-- - Smart contracts and ledgers are used to encapsulate the shared processes and shared information in a network, respectively -->
- 一个 peer 上可以有很多个账本实例和 chaincode 实例
- 应用和管理员必须通过 peer 才能与账本、合约等资源交互
- 初始创建一个 peer 时，peer 上没有账本或者 chaincode
- Peer 上通常至少会安装一个 chaincode 来与账本交互
- Peer 上总是有系统合约
## 应用和 peer 交互
- 应用要连接到 peer 才能与账本和 chaincode 交互
- 应用通过 Fabric SDK 的 API 连接到 peer，调用 chaincode 生成交易，提交交易到网络，接收事件
- 应用和 peer 建立连接后，应用可以执行 chaincode 来查询或更新账本
    - 查询交易的操作会立刻返回结果  
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/peers.diagram.6.png)
- 对于应用的查询操作，应用 A 接收到 proposal response 后就已经完成了
- 对于更新账本的操作
    - A builds a transaction from all of the responses, which it sends it to O1 for ordering
    - O1 将网络中的搜集到的交易打包成块，分发到所有 peer，包括 P1
    - P1 校验交易后将区块提交到 L1
    - L1 更新完成后，P1 生成一个事件，A 收到该事件表示操作完成
- peer 上保存了账本实例，所以接收到应用的查询请求会直接返回结果，而不会去查询其它 peer
- 应用可以向一个或多个 peer 发起查询请求
    - For example, to corroborate a result between multiple peers, or retrieve a more up-to-date result from a different peer if there’s a suspicion that information might be out of date
- 更新账本的操作无法通过一个 peer 单独完成，其它 peer 必须先同意该改动才行，即达成共识
    - peer 先返回给应用一个 update proposal response
    - 应用向 orderer 发起交易排序请求
    - orderer 排序交易、打包成块后将块广播给网络上的 peer
    - peer 校验区块后将区块写入各自的账本
    - 应用收到异步的事件通知
- 组织为网络提供 peer 节点等资源  
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/peers.diagram.9.png)
- 每个 peer 都有拥有它的组织为其颁发的数字证书
- peer 连接到 channel 时，它的数字证书通过 channel MSP 识别出拥有它的组织
- 一个 peer 只能被一个组织拥有，因此也只与一个 MSP 关联
- peer 的数字证书的颁发者（CA）决定谁是它的拥有者，而不是 peer 被部署的物理位置
## peer 和 orderer 交互
### 1. Proposal
- 应用向要求的背书 peer 发送交易 proposal
- 每个背书 peer 独自用交易 proposal 执行 chaincode，生成 endorsed transaction proposal response 返回应用
- 应用收到足数的 transaction proposal response 后，proposal 阶段即完成  
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/peers.diagram.10.png)
- 应用 A1 生成交易 T1 的 proposal P，发送给 channel C 上的 peer P1 和 P2
- P1 用 T1 的 proposal P 执行 chaincode S1，生成交易 T1 由 经过 E1 背书的返回值 R1；P2 同理
- A1 收到交易 T1 的两个经过背书的返回信息，分别是 E1 和 E2
- 最初，应用选出一组 peer 为更新账本的交易 proposal 背书
    - peer 的选择由 chaincode 背书策略决定，背书策略规定了一个更改账本的提议必须经过一组选定的组织背书后才会被网络接受
- peer 背书 proposal response
    - 添加自己的数字签名后，将整个 payload 用私钥签名
- peer 的背书随后还可用作某个组织的 peer 生成了某个 response 的证明
    - In our example, if peer P1 is owned by organization Org1, endorsement E1 corresponds to a digital proof that “Transaction T1 response R1 on ledger L1 has been provided by Org1’s peer P1!”
- 返回的 transaction response 不一致
    - 节点的账本的状态不同步
        - The result was generated at different times on different peers with ledgers at different states
    - Chaincode is non-deterministic
        - 单个 peer 不知道会否存在不一致，必须对比多个 transaction response 才能检测出不一致
- 对于不一致的 transaction response，应用可以选择丢弃，终止后续的流程
### 2. Ordering and packaging transactions into blocks
- Orderer 收到来自多个应用的包含经过背书的 transaction proposal response 的交易，将交易排序、打包成块
- 排序服务节点同时接收到多个应用的经过背书的 propose transaction response
- 一个区块的交易数有 channel 配置参数（`BatchSize`, `BatchTimeout`）控制
- 区块保存到 orderer 的账本并分发到 channel 中的 peer
    - 若 peer 不在线，则会在重新连接到 orderer 后接收到区块，或是通过与其它 peer 的 gossip 获得
### 3. Validation and commit
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/peers.diagram.12.png)
- Orderer 将区块分发到所有与之相连的 peer
- 每个 peer 独立处理区块
- 不需要将所有 peer 都连到 orderer，peer 直接可以通过 gossip 协议同步区块
- peer 按照块中交易的顺序组个处理交易
    - 对于每个交易，peer 会按照生成交易的 chaincode 的背书策略的要求校验交易是否已经过必要的组织的背书
    <!-- - For every transaction, each peer will verify that the transaction has been endorsed by the required organizations according to the endorsement policy of the chaincode which generated the transaction -->
    - 这一过程可以检验相关的组织是否产生了相同的结果
- 若交易已被正确地背书，peer 会尝试将其添加到账本
- 添加之前，peer 会检查账本的一致性
    - Verify that the current state of the ledger is compatible with the state of the ledger when the proposed update was generated
- Failed transactions are not applied to the ledger, but they are retained for audit purposes, as are successful transactions
- 此阶段不需要运行 chaincode
- chaincode 只需要安装在背书节点
    - This is often helpful as it keeps the logic of the chaincode confidential to endorsing organizations
- The output of the chaincodes (the transaction proposal responses) are shared with every peer in the channel, whether or not they endorsed the transaction
- 区块提交到 peer 的账本后，peer 会生成相应的事件
    - Block events include the full block content, while block transaction events include summary information only, such as whether each transaction in the block has been validated or invalidated
    - Chaincode events that the chaincode execution has produced can also be published at this time.