- The `fabric-ccenv` image which is used to build chaincode, currently includes the github.com/hyperledger/fabric/core/chaincode/shim (“shim”) package.
    - > https://hyperledger-fabric.readthedocs.io/en/release-1.1/releases.html
- `chaintool` is a hyperledger Fabric **chaincode compiler**
    - `chaintool` is a utility to assist in various phases of Hyperledger Fabric chaincode development, such as compilation, test, packaging, and deployment. A chaincode app developer may express the interface to their application in a high-level interface definition language, and `chaintool` will generate (1) chaincode stubs and (2) package the chaincode for convenient deployment.
    - > https://fabric-chaintool.readthedocs.io/en/latest/

- 账本包含不成功的交易。
- peer 节点会在本地维护一个位掩码来隔离有效交易和无效交易。
- 客户端（提交客户端）：提交实际交易调用到背书 peer，广播交易请求到排序服务节点。客户端代表最终用户实体，它必须连接到一个 peer 节点以与区块链交互。客户端可以选择连接任何 peer 节点，创建并调用交易。