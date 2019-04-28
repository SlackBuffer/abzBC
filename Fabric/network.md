# Blockchain network
- A blockchain network is a technical infrastructure that provides ledger and smart contract services to applications
    - Primarily, smart contracts are used to generate transactions which are subsequently distributed to every peer node in the network where they are immutably recorded on their copy of the ledger
- In most cases, multiple organizations come together as a consortium to form the network and their permissions are determined by a set of policies that are agreed by the consortium when the network is originally configured
    - Moreover, network policies can change over time subject to the agreement of the organizations in the consortium
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
- 尽管共识节点 O4 运行在组织 R4，只有存在网络简介，R1 对 O4 也有和 R4 相同的管理权限
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
- 只能合约开发完成后，R1 的管理员将其安装到 P1
- 安装完成后，P1 可以看到 S5 的实现逻辑
- 一个组织在 channel 上有多个 peer 时，可以选择部分节点安装智能合约，不需要在所有节点都安装
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
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/network.diagram.7.png)