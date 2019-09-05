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
<!-- # Blockchain network -->
<!-- - A blockchain network is a technical infrastructure that provides ledger and smart contract services to applications
    - Primarily, smart contracts are used to generate transactions which are subsequently distributed to every peer node in the network where they are immutably recorded on their copy of the ledger
- In most cases, multiple organizations come together as a consortium to form the network and their permissions are determined by a set of policies that are agreed by the consortium when the network is originally configured
    - Moreover, network policies can change over time subject to the agreement of the organizations in the consortium -->
<!-- 
## 1. Creating the network
- The network is formed when an orderer is started
    - The ordering service comprising a single node, O4, is configured according to a network configuration NC4
    - NC4 gives administrative rights to **organization** R4
    - At the network level, certificate authority CA4 is used to **dispense** identities to the administrators and network nodes of the R4 organization
- 根据网络配置 NC4，配置单节点的 orderer O4
- 根据网络配置 NC4，管理员权限授予组织 R4（的用户）
    - NC4 包含 网络初始状态下网络资源的操作权限的 配置
- 根据 [ ] 事先约定，组织 R4 的管理员配置并启动 O4，并将 O4 部署在组织 R4
- 以上配置在配置好并启动共识节点 O4 后即刻生效
- It's helpful to think of the ordering service as the initial administration point for the network
- As agreed beforehand, O4 is initially **configured and started** by an administrator in organization R4, and **hosted** in R4
- The configuration NC4 contains the **policies** that describes the **starting set of administrative capabilities for the network**. Initially this is set to only give R4 rights over the network
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.2.png)
### Certificate authorities
- CA 给对应组织内管理员和节点颁发证书
- 证书用于识别网络成员所属组织；用于签名交易，表示交易结果得到了组织的背书
- 证书和成员组织之间的关系通过 MSP 实现映射
- 网络配置 NC4 通过指定的 MSP 来识别证书，找到 NC4 里与该 MSP 名称对应的权限，将其授予组织 R4 里的角色
- CA is used to issue certificates to administrators and network nodes
    - It dispenses **X.509 certificates** that can be used to **identify** components belonging to organization R4
    - Certificates issued by CAs can also be used to **sign transactions** to indicate that an organization endorses the transaction result - a precondition of it being accepted onto the ledger
- Different components of the blockchain network use certificates to identify themselves to each other as being from a particular organization
    - That's why there's usually more than one CA supporting a blockchain network - different organizations often use different CAs
- CAs are so important that Fabric provides you with a built-in one (Fabric-CA) to help you get going, though in practice, organizations will choose to use their own CA
- The mapping of certificates to member organizations is achieved via a structure called membership services provider MSP

- Network configuration NC4 <mark>uses</mark> a **named MSP** to identify the properties of certificates dispensed by CA4 which associate certificate holders with organization R4
- NC4 can then <mark>use</mark> this MSP name **in policies** to grant actors from R4 particular rights over network resources
    - An example of such a policy is to identify the administrators in R4 who can **add new member organizations to the network**

- X.509 certificates are used in client application transaction proposals and smart transaction response to digitally sign transactions
- Subsequently the network nodes who host copies of the ledger verify that transaction signatures are valid before accepting transactions onto ledger
## 2. Adding network administrators
- 组织 R4 更新网络配置 NC4，把组织 R1 也加为管理员
    - 此后 R1 和 R4 均有更改网络配置的权限
- CA1 加入进来，用于标识来自组织 R1 的用户
- 尽管共识节点 O4 运行在组织 R4，只有存在网络连接，R1 对 O4 也有和 R4 相同的管理权限
    - 此处的 orderer 是单节点，但通常情况下 orderer 会是多节点，可以被配置在不同组织的不同节点上
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.2.1.png)
## 3. Define a consortium
- R1（或 R4）通过更新 NC4 添加组织 R2，R2 可以执行对网络的部分操作（非管理员身份）
- 同时添加了 CA2
- R1 或 R4 创建一个集团 X1，包含组织 R1 和组织 R2
- 集团 X1 的定义存在 NC4 里
- A consortium defines the set of organizations in the network who share a need to transact with one another
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.3.png)
## 4. Creating a channel for a consortium
- Channel 是集团成员间主要的通信机制
- 一个网络中可以有多个 channel
- A channel C1 has been created for R1 and R2 using the consortium definition X1
    - Channel C1 has been connected to the ordering service O4 but that nothing else is attached to it
- C1 由 channel 的配置 CC1 管理，完全独立于网络配置 NC4
- CC1 由 R1 和 R2 管理；R1 和 R2 对 CC1 有相同权限；R4 对 CC1 无权限
- C1 为集团 X1 的私密通信提供了通道
- 一个通道可以有任意多个组织与之连接
- Channels provide privacy from other channels, and from the network
    - Different consortia within the network will have a need for different information and processes to be appropriately shared
- Channels provide an efficient sharing of infrastructure while maintaining data and communications privacy
- There is also a special system channel defined for **use by the ordering service**. It behaves in exactly the same way as a regular channel, which are sometimes called application channels
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.4.png)
## 5. Peers and ledger
- Peers are joined to channels by the organizations that own them
- A peer node P1 has joined the channel C1
- P1’s purpose in the network is purely to host a copy of the ledger L1 for others to access
- Think of L1 as being physically hosted on P1, but logically hosted on the channel C1
- P1 and O4 通过 C1 通信
- P1 的一个关键配置是由 CA1 颁发的 X.509 身份，该身份将 P1 和组织 R1 关联起来
- P1 节点启动后，可以通过 O4 加入 C1；O4 收到加入请求后，查询 CC1 来确定 P1 在 C1 上的权限
    - For example, CC1 determines whether P1 can read and/or write information to the ledger L1
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.5.png)
## 6. Applications and smart contract chaincode
- Now that the channel C1 has a ledger on it, we can start connecting client applications (A1) to consume some of the services provided by workhorse of the ledger, the peer
- 合约 S5 被安装到 P1
- 组织 R1 的客户端应用 A1 **通过 S5** 来访问 P1 节点上的账本 L1
    - Think of S5 as defining all the common access patterns to the ledger; S5 provides a well-defined set of ways by which the ledger L1 can be queried or updated
- A1 可以通过 C1 与网络上的资源相连（如 P1，O4）
- 客户端应用也有与组织绑定的身份
- 一个组织开发的合约可以被该组织从属的集团的其它成员组织共享
- 合约必须安装并初始化
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.6.png)
### Installing a smart contract
- 合约开发完成后，R1 的管理员将其安装到 P1
- 安装完成后，P1 可以看到 S5 的实现逻辑
- 一个组织在 channel 上有多个 peer 时，可以选择**部分节点安装智能合约**，不需要在所有节点都安装
### Instantiating a smart contract
- C1 的其它成员在 S5 被初始化之前感知不到它的存在
- R1 的管理员通过 P1 初始化 S5 后，S5 可以被 A1 调用
- 而后 C1 的其它成员可以访问到 S5，但无法看到内部代码，只有安装该合约的节点（P1）可以看到
### Endorsement policy
- 初始化合约后背书策略也得到配置
- 背书策略描述了一个组织的交易必须得到其它哪些组织的认可后才得以写入其它组织的账本
- 合约初始化的动作将背书策略写入 channel C1 的配置 CC1 中，C1 的所有成员均可访问到
### Invoking a smart contract
- Client applications invoke smart contracts by sending transaction **proposals** to **peers** owned by the organizations specified by the smart contract endorsement policy
- The transaction proposal serves as input to the smart contract, which uses it to generate an **endorsed transaction** response, which is returned by the peer node to the client application
## 7. Network completed
- The network has grown through the addition of infrastructure from organization R2
- Specifically, (an administrator in) R2 has added peer node P2, which hosts a copy of ledger L1, and chaincode S5
- P2 has also joined channel C1, as has application A2
- A2 and P2 are identified using certificates from CA2
- All of this means that both applications A1 and A2 can invoke S5 on C1 either using peer node P1 or P2
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.7.png)
### Generating and accepting transactions
- 安装了合约的 peer 才可以产生交易
- All peer nodes can **validate** and subsequently **accept or reject** transactions onto their copy of the ledger L1
- 只有安装了合约的 peer 才可以参与背书交易的过程
### Types of peers
- In Hyperledger Fabric, while **all peers are the same**, they can **assume multiple roles** depending on how the network is configured
    1. Committing peer
        - 每个 peer 都是 committing peer
        - It **receives** blocks of generated transactions, which are subsequently **validated** before they are committed to the peer node’s copy of the ledger as an append operation
    2. Endorsing peer
        - 每个安装了合约的 peer 都可以成为 endorsing peer
        - To actually be an endorsing peer, the smart contract on the peer must be [ ] **used by a client application to generate a digitally signed transaction response**
        - **合约的**背书策略规定了交易必须在某些组织的 peer 签名过后才可以写入账本
        - An endorsement policy **for a smart contract** identifies the **organizations whose peer** should digitally **sign a generated transaction** before it can be accepted onto a committing peer’s copy of the ledger
    3. Leader peer
        - 若某个组织在 channel 中有多个 peer，leader peer 负责将交易从 orderer 分发到该组织的 committing peer
        - Peer 可以选择参与静态或动态的领导权选举
            - For the static set, zero or more peers can be configured as leaders
            - For the dynamic set, one peer will be elected leader by the set. If a leader peer fails, then the remaining peers will re-elect a leader
        - 一个组织里可以有一个或多个 leader peer 连接到排序服务
            - This can help to improve resilience and scalability in large networks which process high volumes of transactions
    4. Anchor peer
        - 通过锚节点与其他组织的 peer 通信
        - 锚节点的信息在该组织的 channel 配置里
        - An organization can have zero or more anchor peers defined for it
    - 一个 peer 可以同时是以上 4 中角色
    - 一个组织中必须有 1,2,3 的 peer 角色
### Install not instantiate
- R2 上的 S5 不需要被初始化
- 一个合约的初始化只需要执行一次，所有后来加入 channel 的 peer (安装了该合约的) 就知道合约已可调用
    - This fact reflects the fact that ledger L1 and smart contract really exist in a physical manner on the peer nodes, and a logical manner on the channel
    - R2 is merely adding another physical instance of L1 and S5 to the network
- Copies of smart contract S5 will usually be identically implemented using the same programming language, but if not, they must be semantically equivalent
- The careful addition of peers to the network can help support increased throughput, stability, and resilience
- The technical mechanism by which peers within an individual organization efficiently discover and communicate with each other – the [ ] [gossip protocol](https://hyperledger-fabric.readthedocs.io/en/release-1.4/gossip.html#gossip-protocol) – will accommodate a large number of peer nodes in support of such topologies
- The careful use of network and channel policies allow even large networks to be well-governed
    - Organizations are free to add peer nodes to the network so long as they conform to the policies agreed by the network
    - Network and channel policies create the balance between autonomy and control which characterizes a de-centralized network
## 8. Simplifying the visual vocabulary
- The network diagram has been simplified by replacing channel lines with connection points, shown as blue circles which include the channel number
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.8.png)
## 9. Adding another consortium definition
- Define a new consortium, X2, for R2 and R3
- A network administrator from organization R1 or R4 has added a new consortium definition, X2, which includes organizations R2 and R3. This will be used to define a new channel for X2
    - Consortium X2 has been introduced in order to be able to create a new channel for R2 and R3
- A new channel can only be created by those organizations specifically identified in the network configuration policy, NC4, as having the appropriate rights to do so, i.e. R1 or R4
    - This is an example of a policy which separates organizations that can manage resources at the network level versus those who can manage resources at the channel level
- In practice, consortium definition X2 has been added to the network configuration NC4
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.9.pngX)
## 10. Adding a new channel
- A new channel C2 has been created for R2 and R3 using consortium definition X2
- The channel has a channel configuration CC2, completely separate to the network configuration NC4, and the channel configuration CC1
- Channel C2 is managed by R2 and R3 who have equal rights over C2 as defined by a policy in CC2
- R1 and R4 have no rights defined in CC2 whatsoever
- The channel configuration CC2 now contains the policies that govern channel resources, assigning management rights to organizations R2 and R3 over channel C2
- Note how the channel configurations CC1 and CC2 remain completely separate from each other, and completely separate from the network configuration, NC4
- As the network and channels evolve, so will the network and channel configurations
- There is a process by which this is accomplished in a controlled manner – involving configuration transactions which capture the change to these configurations
- Every configuration change results in a new **configuration block transaction** being generated
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.10.png)
### Network and channel configuration
- Network and channel configurations are important because they encapsulate the **policies** agreed by the network members, which provide a shared reference for controlling access to network resources
- Network and channel configurations also contain **facts** about the network and channel composition, such as the name of consortia and its organizations
- When the network is first formed using the ordering service node O4, its behaviour is governed by the network configuration NC4
- The initial configuration of NC4 only contains policies that permit organization R4 to manage network resources
- NC4 is subsequently updated to also allow R1 to manage network resources
- Once this change is made, any administrator from organization R1 or R4 that connects to O4 will have network management rights because that is what the policy in the network configuration NC4 permits
- Internally, each node in the ordering service records each channel in the network configuration, so that there is a record of each channel created, at the network level
    - It means that although ordering service node O4 is the **actor** that created consortia X1 and X2 and channels C1 and C2, the **intelligence** of the network is contained in the network configuration NC4 that O4 is obeying
- As long as O4 behaves as a good actor, and correctly implements the policies defined in NC4 whenever it is dealing with network resources, our network will behave as all organizations have agreed
- In many ways NC4 can be considered more important than O4 because, ultimately, it controls network access
- The same principles apply for channel configurations with respect to peers
    - In our network, P1 and P2 are likewise good actors
    - When peer nodes P1 and P2 are interacting with client applications A1 or A2 they are each using the policies defined within channel configuration CC1 to control access to the channel C1 resources
- If A1 wants to access the smart contract chaincode S5 on peer nodes P1 or P2, each peer node uses its **copy** of CC1 to determine the operations that A1 can perform
- Even though there is logically a single configuration, it is actually replicated and kept consistent by every node that forms the network or channel
- Similarly ordering service node O4 has a copy of the network configuration, but in a multi-node configuration, every ordering service node will have its own copy of the network configuration
- Both network and channel configurations are kept consistent using the same blockchain technology that is used for user transactions – but for **configuration transactions**
- To change a network or channel configuration, an administrator must submit a configuration transaction to change the network or channel configuration
    - It must be signed by the organizations identified in the appropriate policy as being responsible for configuration change
    - This policy is called the `mod_policy`
- The ordering service nodes operate a mini-blockchain, connected via the **system channel**
    - The system channel ordering service nodes distribute **network configuration transactions**
    - These transactions are used to co-operatively maintain a consistent copy of the network configuration at each ordering service node
- In a similar way, peer nodes in an application channel can distribute channel configuration transactions
    - Likewise, these transactions are used to maintain a consistent copy of the channel configuration at each peer node
- Objects like network configurations, that are logically single, turn out to be physically replicated among a set of ordering services nodes
- We also see it with channel configurations, ledgers, and to some extent smart contracts which are installed in multiple places but whose interfaces exist logically at the channel level
- This balance between objects that are logically singular, by being physically distributed is a common pattern in Hyperledger Fabric
- It’s a **pattern** you see repeated time and again in Hyperledger Fabric, and enables Hyperledger Fabric to be both de-centralized and yet manageable at the same time
## 11. Adding another peer
- Client applications A3 can use channel C2 for communication with peer P3 and ordering service O4. Ordering service O4 can make use of the communication services of channels C1 and C2
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.11.png)
## 12. Joining a peer to multiple channels
- R2 is a member of both consortia X1 and X2 by joining it to multiple channels
- Ordering service O4 can make use of the communication services of channels C1 and C2
- Client application A2 can use channel C1 for communication with peers P1 and P2 and channel C2 for communication with peers P2 and P3 and ordering service O4
- Peer node P2 has smart contract S5 installed for channel C1 and smart contract S6 installed for channel C2
- Peer node P2 is a full member of both channels at the same time via different smart contracts for different ledgers
- Peer node P2’s behaviour is controlled very differently depending upon the channel in which it is transacting
    - Specifically, the policies contained in channel configuration CC1 dictate the operations available to P2 when it is transacting in channel C1, whereas it is the policies in channel configuration CC2 that control P2’s behaviour in channel C2
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.12.png)
### The ordering services
- The ordering service node appears to be a centralized component
    - It was used to create the network initially, and connects to every channel in the network
    - Even though we added R1 and R4 to the network configuration policy NC4 which controls the orderer, the node was running on R4’s infrastructure
    - In a world of de-centralization, this looks wrong
- The ordering service can itself too be completely de-centralized
    - An ordering service could be comprised of many individual nodes owned by different organizations
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.15.png)
- The ordering service comprises ordering service nodes O1 and O4. O1 is provided by organization R1 and node O4 is provided by organization R4
- The network configuration NC4 defines network resource permissions for actors from both organizations R1 and R4
- The network configuration policy, NC4, permits R1 and R4 equal rights over network resources
- Client applications and peer nodes from organizations R1 and R4 can manage network resources by connecting to either node O1 or node O4, because both nodes behave the same way, as defined by the policies in network configuration NC4
    - In practice, actors from a particular organization tend to use infrastructure provided by their home organization, but that’s certainly not always the case
### De-centralized transaction distribution
- As well as being the **management point** for the network, the ordering service also provides another key facility – it is the **distribution point** for transactions
- The ordering service is the component which **gathers endorsed transactions** from applications and **orders them into transaction blocks**, which are subsequently **distributed** to every peer node in the channel
- At each of these committing peers, transactions are recorded, whether valid or invalid, and their local copy of the ledger updated appropriately
- When acting at the channel level, O4’s role is to gather transactions and distribute blocks inside channel C1
    - It does this according to the policies defined in channel configuration CC1
- When acting at the network level, O4’s role is to provide a management point for network resources according to the policies defined in network configuration NC4
### Changing policy
- Policy change is managed by a policy within the policy itself
- The `modification policy`, or `mod_policy` for short, is a first class policy within a network or channel configuration that manages change
- When the network was initially set up, only organization R4 was allowed to manage the network
- In practice, this was achieved by making R4 the only organization defined in the network configuration NC4 with permissions to network resources
- Moreover, the `mod_policy` for NC4 only mentioned organization R4 – only R4 was allowed to change this configuration
- We then evolved the network N to also allow organization R1 to administer the network
    - R4 did this by adding R1 to the policies for channel creation and consortium creation
    - Because of this change, R1 was able to define the consortia X1 and X2, and create the channels C1 and C2
    - R1 had equal administrative rights over the channel and consortium policies in the network configuration
- R4 could add R1 to the `mod_policy` such that R1 would be able to manage change of the network policy too
    - This power is much more powerful than the first, because now R1 now has full control over the network configuration NC4
    - This means that R1 can, in principle remove R4’s management rights from the network
        - In practice, R4 would configure the `mod_policy` such that R4 would need to also approve the change, or that all organizations in the `mod_policy` would have to approve the change
- The `mod_policy` behaves like every other policy inside a network or channel configuration; it defines a set of organizations that are allowed to change the `mod_policy` itself
## 13. Network fully formed 
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.14.png)
- The channel configuration of the system channel is part of the network configuration, NC4 -->