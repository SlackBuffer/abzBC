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