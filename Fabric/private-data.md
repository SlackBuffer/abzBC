- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/private-data-arch.html

- 一个 channel 中的部分组织可以通过使用私密数据将一部分分数据隔离到单独的数据库，这部分数据与 channel 的账本隔离开，只有有权的组织才能访问
- 将合约的数据进一步用通用的加密算法（如 AES）加密，再发送给排序服务，这样写入到账本的数据只有有对应私钥的一方才能解密
<!-- - When a subset of organizations on that channel need to keep their transaction data confidential, a private data collection (**collection**) is used to segregate this data in a private database, logically separate from the channel ledger, accessible only to the authorized subset of organizations
    - Channels keep transactions private from the broader network whereas collections keep data private between subsets of organizations on the channel -->
<!-- - To further obfuscate the data, values within chaincode can be encrypted (in part or in total) using common cryptographic algorithms such as AES before sending transactions to the ordering service and appending blocks to the ledger
    - Once encrypted data has been written to the ledger, it can be decrypted only by a user in possession of the corresponding key that was used to generate the cipher text -->