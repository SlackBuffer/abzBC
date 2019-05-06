- 公钥机制（Public key infrastructure）用于为组织，网络成员，终端用户生成密码学证书
- 身份校验发生在交易流的每一步
    - 访问权限实现在网络的每个层级，从排序服务到 channel
    - Payload 反复签名、验证，通过认证后作为交易 proposal 传递给各个组件
<!-- - There're also ongoing identity verifications happening in all directions of the transaction flow
    - Access control lists are implemented on hierarchical layers of the network (ordering service down to channels)
        - For example, a specific user ID could be permitted to invoke a chaincode application, but be blocked from deploying new chaincode
    - Payloads are repeatedly signed, verified, and authenticated as a transaction proposal passes through the different architectural components -->
<!-- - Public key infrastructure is used to generate cryptographic certificates which are tied to organizations, network components, and end users or client applications
    - As a result, data access control can be manipulated and governed on the broader network and on channel levels -->
- 区块链网络有各种参与者，包括 peer，orderer，客户端应用，管理员等
    <!-- - Actors — active elements inside or outside a network able to consume services -->
    - 参与者指区块链网络内外可以消费区块链提供的服务的元素
- 每个参与者的数字身份都封装在 X.509 数字证书里
- 身份用于确定参与者在网络中的权限
    <!-- - Determine the exact permissions over resources and access to information that actors have in a blockchain network -->
- 身份和其伴随的属性的集合称为 principal
    - principal 和 ID 类似但更灵活
    <!-- - Principals are just like userIDs or groupIDs, but a little more flexible because they can include a wide range of properties of an actor’s identity, such as the actor’s organization, organizational unit, role or even the actor’s specific identity -->
    - Principals 可以理解为决定权限的属性
    <!-- - When we talk about principals, they are the properties which determine their permissions -->
- 身份必须来自可信机构才可核实，Fabric 通过 MSP 实现
- PKI 提供身份，MSP 将该身份匹配到参与到网络的组织中
    - > 超市购物结账时，需要有一张有效的且超级同意接收的信用卡，信用卡有效，但是超市对其不认可（如与该发卡行没有合作），无法完成交易
    <!-- - A PKI provides a list of identities, and an MSP says which of these are members of a given organization that participates in the networ -->
    - > PKI 对应信用卡发卡行，分发不同类型的可核实的身份；MSP 对应超市支持的发卡行，决定哪些身份是网络的可信参与者
    - MSP 将可核查的身份映射为区块链网络的成员（参与者）
- 成员通过 MSP 参与到网络中
- An MSP is a component that defines the rules that govern the valid identities for this organization
- The default MSP implementation in Fabric uses X.509 certificates as identities, adopting a traditional Public Key Infrastructure (PKI) hierarchical model
## Public Key Infrastructure (PKI)
- PKI 生成可校验的数字身份
- PKI 提供网络中的安全通行，是 HTTPS 里的 S
- 组成
    1. Digital certificates
    2. Public and private keys
    3. Certificate authorities (CAs)
    4. Certificate revocation lists
### 数字证书
- 数字证书是一份包含证书持有者属性的文档
- 最常见的类型是兼容 X.509 标准的证书，此类证书将参与方的身份详情加密其中
    <!-- - The most common type of certificate is the one compliant with the X.509 standard, which allows the encoding of a party’s identifying details in its structure -->
- 描述一个叫做 Mary Morris 参与方的数字证书
    - Mary 是证书的主题（`SUBJECT` of the certificate）
    - `Subject` 字段包含 Mary 的关键信息
    - 证书包含 Mary 的公钥
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/identity.diagram.8.png)
- Mary 的属性用密码学的方式被记录下来，篡改属性会使得证书失效
- 密码学技术使得 Mary 可以出示她的证书来证实自身的身份，前提是其它参与方信任 Mary 的证书颁发方（CA）
- 只要 CA 自身的私钥是安全的，其它组织读到由该 CA 颁发的证书时就能确信该证书的信息未被篡改
- 可以把 Mary 的 X.509 想象成无法串改的数字身份证
### 公私钥
- 公钥是公之于众的
- 私钥加密的信息（生成信息的数字签名）只有用对应的公钥才能解（原信息）；若信息被篡改了也会变得不可解
    <!-- - The private key can be used to produce a signature on a message that only the corresponding public key can match, and only on the same message -->
    - 验证了信息的来源且信息未被篡改
### 认证中心
- 认证中心为成员颁发数字证书
- 在于网络的通信中，成员用数字证书标识身份（authenticate themselves）
- 流行的 CA：GoDaddy, Symantec (originally Verisign), GeoTrust, DigiCert, Comodo
- CA 对它颁发的证书、参与者的公钥、一些可选信息签名；信任该 CA 的参与者（知道该 CA 的公钥）可以通过验证该 CA 颁发的证书的签名，解出并信任证书的里内容
- [ ] CA 自身也有证书，他方用该 CA 的公钥校验 CA 的证书通过，且校验该 CA 颁发的证书也通过，则证明该 CA 颁发的证书确实使用该 CA 的私钥签名的
    - [ ] 自己为自己签名
- CA 分 root CA (Symantec, Geotrust, etc) 和 intermediate CA
- 中间 CA  的证书由根 CA 或者其它 CA 颁发
- 通过多级 CA 组成的信任链条可以追溯到根 CA，同时避免暴露根 CA
    - 根 CA 暴露会波及整个链条
- 多级 CA 的模式提供了更好的可扩展性
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/identity.diagram.1.png)
- Fabric CA 是內建的私有根 CA
    - 默认不能提供浏览器端的 SSL 证书
    <!-- - Because Fabric CA is a custom CA targeting the Root CA needs of Fabric, it is inherently not capable of providing SSL certificates for general/automatic use in browsers -->
### 吊销证书名单
- CRL 是已失效的证书
    - 可理解成丢失的信用卡
- 第三方校验某一方的身份时，首先检查发行该证书的 CA 的 CRL
- 证书被吊销和证书过期不同
## Membership 成员资格
- MSP 将 PKI 生成的身份表示为区块链网络中的可信成员
- MSP 将身份映射为角色
    - 将数字身份映射到所属的组织，决定在组织中担任的角色和拥有的权限
- MSP 识别出哪些根 CA、中间 CA 可以定义组织的成员
    1. MSP 列出组织成员的身份
    2. MSP 识别出哪些 CA 有创建成员（为其提供有限的身份）的权限
    3. 结合 1 和 2
- MSP 不仅可以记录出网络参与者或是 channel 的成员
- An MSP can identify specific roles an actor might play either within the scope of the organization the MSP represents (e.g., admins), or as members of a sub-organization group, and sets the basis for defining access privileges in the context of a network and channel (e.g., channel admins, readers, writers)
- MSP 的配置通过 channel MSP 的形式广播到所有相关的 channel
- 除了 channel MSP 外，peer，orderer，客户端也维护了一份本地 MSP，以认证（authenticate）通道以外的信息、定义特定成员的权限
- 通常情况一个组织用一个 MSP 管理成员，所以组织的名称也通常包含在其对应的 MSP 名称中
    - 组织名 `ORG1`，MSP 名 `ORG1-MSP`
- 多 channel 情况下用多个 MSP 表示不同的组织关系
    - 组织名 `ORG2`，`ORG2-MSP-NATIONAL`，`ORG2-MSP-GOVERNMENT`
- 一个组织可划分为多个 organizational units (OUs) 组织单位，每个单位负责不同职责
    - CA 颁发的 X.509 证书的包含 `OU` 字段
- OU 有时还被用作区分同一个集团下不同组织的
    - 此种情况下不同组织用相同的根 CA 和中间 CA，用 OU 字段标识每个组织的成员
- 分 channel MSP 和本地 MSP
- 为客户端（用户）、节点（peer，orderer）定义了本地 MSP
    - 节点的本地 MSP 定义了节点的权限
    - The local MSPs of the users allow the user side to authenticate itself in its transactions as a member of a channel (e.g. in chaincode transactions), or as the owner of a specific role into the system (an org admin, for example, in configuration transactions)
- 每个节点和用户必须有本地的 MSP，该 MSP 觉得管理权和参与权
- 每个参与 channel 的组织都只需有对应的 MSP，该 MSP 用于认证 channel 的参与者
    - This means that if an organization wishes to join the channel, an MSP incorporating the chain of trust for the organization’s members would need to be included in the channel configuration. Otherwise transactions originating from this organization’s identities will be rejected
- An administrator `B` connects to the peer with an identity issued by `RCA1` and stored in their local MSP. When `B` tries to install a smart contract on the peer, the peer checks its local MSP, `ORG1-MSP`, to verify that the identity of `B` is indeed a member of `ORG1`. A successful verification will allow the install command to complete successfully. Subsequently, `B` wishes to instantiate the smart contract on the channel. Because this is a channel operation, all organizations on the channel must agree to it. Therefore, the peer must check the MSPs of the channel before it can successfully commit this command. (Other things must happen too, but concentrate on the above for now.)
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/membership.diagram.4.png)
- Local MSPs are only defined on the file system of the node or user to which they apply. Therefore, physically and logically there is only one local MSP per node or user
- However, as channel MSPs are available to all nodes in the channel, they are logically defined once in the channel configuration. However, a channel MSP is also instantiated on the file system of every node in the channel and kept synchronized via consensus. So while there is a copy of each channel MSP on the local file system of every node, logically a channel MSP resides on and is maintained by the channel or the network
- 可以认为存在多层 MSP，需要分别为网络、channel、peer、orderer、用户顶一个 MSP
- peer 和 orderer 的 MSP 是保存在本地的，channel 和网络的 MSP 是相关成员间共享的
- 网络 MSP
    - 定义网络的成员，包括参与的组织和管理员权限的授权
- Channel MSP
- Peer MSP
    - 在 peer 节点上安装合约需要校验该 peer 的 MSP
- Orderer MSP
## MSP 结构
- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/msp.html
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/membership.diagram.5.png)
1. Root CAs
    - 包含此 MSP 所代表的组织所信任的 root CA 的自签名的 X.509 证书（至少要有一个）
    - 其它证书从此证书衍生出来
2. Intermediate CAs
    - 可能为空
3. Organizational units
    - 位于 `$FABRIC_CFG_PATH/msp/config.yaml`
4. Administrators
    - 是一组 X.509 证书
5. Revoked certificates
    - 存的的身份信息（Subject Key Identifier (SKI) 和 Authority Access Identifier (AKI)）而非身份本身
6. Node identity
    - Node identity, in combination to the content of `KeyStore`, would allow the node to authenticate itself in the messages that is sends to other participants of its channels and network
    - 对于基于 X.509 的身份，node identity 是一个 X.509 的证书
    - > This is the certificate a peer places in a transaction proposal response, for example, to indicate that the peer has endorsed it — which can subsequently be checked against the resulting transaction’s endorsement policy at validation time
    - 对于本地 MSP 是必须的，且一个节点只能有一个 X.509 证书
    - channel MSP 不需要
7. `KeyStore` for private key
    - 包含节点的 signing key
    - 用于对数据签名
        - For example to sign a transaction proposal response, as part of the endorsement phase
    - 对于本地 MSP 是必须的，且必须只包含一个 private key
    - channel MSP 不需要
8. TLS root CA
    - 包含此 MSP 所代表的组织所信任的参与 TLS 通信的 root CA 的自签名的 X.509 证书（至少要有一个）
        - An example of a TLS communication would be when a peer needs to connect to an orderer so that it can receive ledger updates
9.  TLS intermediate CA