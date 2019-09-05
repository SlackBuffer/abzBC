- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/private-data-arch.html

- 一个 channel 中的部分组织可以通过使用私密数据将一部分分数据隔离到单独的数据库，这部分数据与 channel 的账本隔离开，只有有权的组织才能访问
    - Use cases in which you want all channel participants to see a transaction while keeping a portion of the data private
- v1.2 起支持隐私数据集
    - Allow a defined subset of organizations on a channel the ability to endorse, commit, or query private data without having to create a separate channel
- 隐私数据集包含 2 个部分
    1. 隐私数据
        - 通过 gossip 协议点对点地发送给有权限查看的组织
        - 数据存放在有权限查看的 peer 的 private state database (也被称作 "side" database)，可以被有权限的组织的 peer 的 chaincode 访问到
        - orderer 不参与到此过程，隐私数据对其也不可见
        - 需要在每个相关 peer 上配置 `CORE_PEER_GOSSIP_EXTERNALENDPOINT` 设置锚节点信息来实现跨组织的通信
    2. 隐私数据的 hash
        - A hash of that data, which is endorsed, ordered, and written to the ledgers of every peer on the channel. The hash serves as evidence of the transaction and is used for state validation and can be used for audit purposes
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/PrivateDataConcept-2.png)
- 私有数据可以被清除
- 将合约的数据进一步用通用的加密算法（如 AES）加密，再发送给排序服务，这样写入到账本的数据只有有对应私钥的一方才能解密
<!-- - When a subset of organizations on that channel need to keep their transaction data confidential, a private data collection (**collection**) is used to segregate this data in a private database, logically separate from the channel ledger, accessible only to the authorized subset of organizations
    - Channels keep transactions private from the broader network whereas collections keep data private between subsets of organizations on the channel -->
<!-- - To further obfuscate the data, values within chaincode can be encrypted (in part or in total) using common cryptographic algorithms such as AES before sending transactions to the ordering service and appending blocks to the ledger
    - Once encrypted data has been written to the ledger, it can be decrypted only by a user in possession of the corresponding key that was used to generate the cipher text -->