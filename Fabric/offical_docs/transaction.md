- 交易包含每个背书节点的签名，被提交给排序服务
- 交易包括读到的键值对的版本（read set）和写入的键值对（write set）
- 交易由 orderer 打成区块并分发给 peer
- 根据背书策略，指定的 peer 校验交易（为交易背书）
- 往区块链追加区块之前，会做版本检查，确保读到的资产状态在执行合约期间未被改变
- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/txflow.html
- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/arch-deep-dive.html#swimlane