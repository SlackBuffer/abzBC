![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/develop.diagram.1.png)
- 5 月底 MagnetoCorp 发行 5M 的商业票据，6 个月后赎回；DigiBank 以 4.94M 买入，6 个月后 MagnetoCorp 以 5M 的价格赎回
- 6 月底 MagnetoCorp 发行 5M 的商业票据；BigFund 以 4.94M 买入
- RateM 是风险评级机构
## 流程分析
### 票据生命周期（世界状态）
- 发行

    ```
    Issuer = MagnetoCorp
    Paper = 00001
    Owner = MagnetoCorp
    Issue date = 31 May 2020
    Maturity = 30 November 2020
    Face value = 5M USD
    Current state = issued
    ```

- 买入

    ```
    Issuer = MagnetoCorp
    Paper = 00001

    Owner = DigiBank

    Issue date = 31 May 2020
    Maturity date = 30 November 2020
    Face value = 5M USD

    Current state = trading
    ```

- 兑换

    ```
    Issuer = MagnetoCorp
    Paper = 00001
    Owner = MagnetoCorp
    Issue date = 31 May 2020
    Maturity date = 30 November 2020
    Face value = 5M USD

    Current state = redeemed
    ```

    - The value of `Owner` of a paper can be used to perform access control on the redeem transaction, by comparing the Owner against the identity of the transaction creator
    - Fabric supports this through the `getCreator()` [chaincode API](https://github.com/hyperledger/fabric-chaincode-node/blob/master/fabric-shim/lib/stub.js#L293)
        - If golang is used as a chaincode language, the [client identity chaincode library](https://github.com/hyperledger/fabric/blob/master/core/chaincode/shim/ext/cid/README.md) can be used to retrieve additional attributes of the transaction creator
### 交易（交易记录组成的区块链）
- 发行

    ```
    Txn = issue
    Issuer = MagnetoCorp
    Paper = 00001
    Issue time = 31 May 2020 09:00:00 EST
    Maturity date = 30 November 2020
    Face value = 5M USD
    ```

    - MagnetoCorp 需要对此交易（发行）签名
- 买入

    ```
    Txn = buy
    Issuer = MagnetoCorp
    Paper = 00001
    Current owner = MagnetoCorp
    New owner = DigiBank
    Purchase time = 31 May 2020 10:00:00 EST
    Price = 4.94M USD
    ```

    - 买卖双方都需要对此交易（买入）签名
- 兑换

    ```
    Txn = redeem
    Issuer = MagnetoCorp
    Paper = 00001
    Current owner = HedgeMatic
    Redeem time = 30 Nov 2020 12:00:00 EST
    ```

    - 买卖双方都需要对此交易（兑换）签名
## 数据结构
- A state will have a combination of properties that uniquely identify it in a given context – it’s **key**
    - > The key for a PaperNet commercial paper is formed by a concatenation of the `Issuer` and `paper` properties; so for MagnetoCorp’s first paper, it’s `MagnetoCorp00001`
    - When a unique key is not available from the available set of properties, an application-determined unique key is specified as an input to the transaction that creates the state. This unique key is usually with some form of UUID
- 通过 state key 可以识别一个票据
- Fabric 要求每个账本里的状态都有唯一的 key
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/develop.diagram.7.png)
- Each paper in the list is represented by a vector state, with a unique composite key formed by the concatenation of `org.papernet.paper`, `Issuer` and `Paper` properties
    - It allows us to examine any state vector in the ledger to determine which list it’s in, without reference to a separate list
    - Hyperlegder Fabric internally uses a [concurrency control mechanism](https://hyperledger-fabric.readthedocs.io/en/release-1.4/arch-deep-dive.html#the-endorsing-peer-simulates-a-transaction-and-produces-an-endorsement-signature) to update a ledger, such that keeping papers in separate state vectors vastly reduces the opportunity for shared-state collisions
        - Such collisions require transaction re-submission, complicate application design, and decrease performance
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/develop.diagram.8.png)
## 智能合约
- The same version of the smart contract must be used by all applications connected to the network so that they jointly implement the same shared business processes and data