- Orderer 只做排序，不执行交易、维护账本
- Peer 只执行交易、维护账本，不做排序
- 每个 channel 有各自的账本；chaincode 也从属于各自账本
- [x] peer, orderer 可以从属于多个 channel
## 创建网络
### 网络配置 NC
- NC 配置了网络的初始状态下的一套管理权限
- 初始状态下，NC 里配置了只有组织 R4 拥有对网络资源的权限
- `mod_policy` 规定了可以更改 `mod_policy` 本身的一系列组织
- NC 的 `modify_policy` 同时规定了只有 R4 有权限更改网络配置
- 网络演进
    - R4 在 NC 的创建 channel 和集团的策略配置里新增了组织 R1，使得 R1 在 channel 和 集团的上和 R4 具有相同的管理权限
    - 若 R4 把 R1 加到 `mod_policy` 里，则 R1 也可以更改网络配置 NC 里的策略，理论上 R1 甚至可以将 R4 的对 NC 的管理权限移除
        - 实际操作中，R4 通常会将 `mod_policy` 配置成需要 R4 或者所有其它组织都同意才通过对 `mod_policy` 的修改
### Orderer
- 初始状态下，根据 NC，配置了单节点 O4
- 根据 [ ] 事先约定，组织 R4 的管理员配置并启动 O4，并将 O4 部署在组织 R4
    - NC 里配置？
- 以上配置在配置好并启动 O4 后即刻生效
- 共识节点启动后网络即形成
- 尽管共识节点 O4 运行在组织 R4，只有存在网络连接，R1 对 O4 也有和 R4 相同的管理权限
- [ ] 要更改网络配置前要先连接到 orderer
- 排序服务里的每个 orderer 都会记录每个 NC 里的每个 channel
- 排序服务的多个 orderer 可以部署在多个组织
- 应用和 peer 可以通过连接 orderer 来管理网络资源
- 排序服务是网络的**管理节点**，同时也是**交易分发点**
    - 在网络层面上，排序服务为管理网络资源提供入口
    - 在 channel 层面上，排序服务从应用收集得到背书的交易，将交易排序打包成块，分发给对应 channel 中的每个 peer 节点
### Certificate authorities (认证中心)
- 各个组织的 CA 各自给组织内的管理员、节点颁发证书
- 组织颁发的证书 X.509
    1. 用于识别网络成员所属的组织
    2. 用于对交易签名，表明交易结果得到了组织的背书
- 证书和组织的对应关系通过 MSP 实现映射
- NC 通过特定的 MSP 来识别由某 CA 颁发的证书的某些属性，这些属性可以将证书的持有者和其对应的组织关联起来（即通过特定的 MSP 找出组织的成员有哪些）
    - [ ] NC 用特定的 MSP 名称将策略里的规定的对网络资源的权限授予该组织的成员
    - > NC4 can then use this MSP name in policies to grant actors from R4 particular rights over network resources
        - > An example of such a policy is to identify the administrators in R4 who can add new member organizations to the network
- Fabric 自带一个 Fabric-CA
    - 实际应用中组织通常会用自己的 CA
### Channel
- Channel 是集团成员间主要的通信机制
- 一个网络中可以有多个 channel
- C1 由 channel 的配置 CC1 管理，完全独立于网络配置 NC
- C1 为集团 X1 的私密通信提供了通道
- 一个通道可以有任意多个组织与之连接
- 有一个特殊的 system channel，供排序服务使用（ordering service）
    - It behaves in exactly the same way as a regular channel
    - 常规 channel 也成为 application channel
- 一个 channel 的配置独立于 NC 和 其它 channel 的配置
- Peer 根据 channel 配置来控制和应用交互时对 channel 资源的访问
    - 比如，A1 要访问 P1 或 P2 上的合约 S5，P1 或 P2 会根据它保存的 CC1 副本来决定允许 A1 进行哪些操作
#### System channel
- Orderer 之间通过 system channel 形成了一个 mini-blockchain
- Orderer 通过 system channel 分发 network configuration transactions
- 网络配置交易用于各个 orderer 维护一份一致的网络配置的备份
- System channel 的配置是网络配置的一部分
#### Application channel
- 常规 channel 里的各个 peer 维护一份一致的 channel 配置的备份
### Peer
- 组织将 peer 加入 channel
- Peer 维护一本它所在 channel 的账本以供访问
- P1 and O4 通过 C1 通信
- P1 的一个关键配置是由 CA1 颁发的 X.509 **身份**，该身份将 P1 和组织 R1 关联起来
- P1 节点启动后，可以通过 O4 加入 C1；O4 收到加入请求后，查询 CC1 来确定 P1 在 C1 上的权限
    - 如，CC1 决定 P1 是否有对账本 L1 的读、写权限
- 一个 peer 若加入多个 channel，其在不同 channel 下的行为分别受到对应的 channel 配置约束
- 一个组织内的 peer 的发现和通信机制：[gossip protocol](https://hyperledger-fabric.readthedocs.io/en/release-1.4/gossip.html#gossip-protocol)
#### Peer 的类别
- Hyperledger Fabric 中， 所有 peer 都相同，但根据网络配置可以参与多个角色
    1. Committing peer
        - 每个 peer 都是 committing peer
        - 接收区块，校验区块并选择是否将区块追加到自己维护的账本
    2. Endorsing peer
        - 每个安装了合约的 peer 都可以成为 endorsing peer
        - To actually be an endorsing peer, the smart contract on the peer must be [ ] **used by a client application to generate a digitally signed transaction response**
    3. Leader peer
        - 若某个组织在 channel 中有多个 peer，leader peer 负责将交易从 orderer 分发到该组织的 committing peer
        - Peer 可以选择参与静态或动态的领导权选举
            - For the static set, zero or more peers can be **configured** as leaders
            - For the dynamic set, one peer will be elected leader by the set. If a leader peer fails, then the remaining peers will **re-elect** a leader
        - 一个组织里可以有一个或多个 leader peer 连接到排序服务
            - This can help to improve resilience and scalability in large networks which process high volumes of transactions
    4. Anchor peer (optional)
        - 通过锚节点与其他组织的 peer 通信
        - 锚节点的信息在 channel 配置里
        - An organization can have **zero or more** anchor peers defined for it
- 一个 peer 可以同时是以上 4 中角色
- 一个组织中必须有 1,2,3 的 peer 角色
### 应用，合约
- 合约 S5 被安装到 P1
- 组织 R1 的客户端应用 A1 **通过 S5** 来访问 P1 节点上的账本 L1
    - Think of S5 as defining all the common access patterns to the ledger; S5 provides a well-defined set of ways by which the ledger L1 can be queried or updated
- A1 可以通过 C1 与网络上的资源相连（如 P1，O4）
- 客户端应用也有**身份**，与组织绑定
- 一个组织开发的合约可以被该组织所从属的集团的其它组织共享
- 合约必须安装并且初始化
#### 安装合约
- 合约开发完成后，R1 的管理员将其安装到 P1
- 安装完成后，P1 可以看到 S5 的实现逻辑（程序代码）
- 一个组织在 channel 上有多个 peer 时，可以选择**部分节点安装智能合约**，不需要在所有节点都安装
#### 初始化合约
- C1 的其它成员在 S5 被初始化之前感知不到它的存在
- R1 的管理员通过 P1 初始化 S5 后，S5 可以被 A1 调用
- 而后 C1 的其它成员可以访问到 S5，但无法看到内部代码，只有安装该合约的节点（P1）可以看到
- 一个合约的初始化只需要执行一次，所有后来加入 channel 的 peer (安装了该合约的) 就知道合约已可调用
    - Copies of smart contract S5 will usually be identically implemented using the same programming language, but if not, they must be semantically equivalent
#### 背书策略
- 初始化合约后，背书策略也得到了配置
- **合约的背书策略**规定了一个组织的交易必须得到其它哪些组织的认可后才会被写入组织的账本
- 合约初始化的动作将背书策略写入 channel C1 的配置 CC1 中，C1 的所有成员均可访问到
#### 合约调用
- 应用向**背书策略里要求的组织**的 peer 节点发送 transaction proposal 来调用合约
- Transaction proposal 作为智能合约的入参，用于生成经过背书的交易（endorsed transaction response），被 peer 节点返回给应用
- 应用 A1 和 应用 A2 均可以调用 C1 里的合约 S5，通过 P1 或 P2 均可
#### 产生，接收交易
- 安装了合约的 peer 才可以产生交易
- 只有安装了合约的 peer 才可以参与为背书交易的过程
- 所有 peer 都可以验证交易，并选择接接受或拒绝将该交易写进自己维护的账本
### 变更配置
- 配置的改变会被配置交易（configuration transactions）捕获
- Every configuration change results in a new **configuration block transaction** being generated
- 网络和 channel 的配置封装了被网络成员认可的策略
- 网络和 channel 的配置包含了网络和通道的组成，如集团名称、集团包含的组织
- 网络配置和各个 channel 的配置逻辑上只有唯一一份，每个参与网络或 channel 的节点都有一份副本
    - 排序服务里的每个 orderer 节点都会有一份 NC 配置的副本
    - 每个 Peer 会保存所处 channel 的 channel 配置的副本
- 管理员通过提交 configuration transaction 的方式变更网络或 channel 的配置
    - Configuration transaction 必须被 `mod_policy` 里要求的组织签名，因为这些组织要为配置的变更负责