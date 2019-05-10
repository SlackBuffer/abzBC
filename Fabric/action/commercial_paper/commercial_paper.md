# Tutorial
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/commercial_paper.diagram.1.png)
## 创建网络

```bash
# fabric-samples/basic-network
docker-compose -f docker-compose.yml up
docker ps -a
# wait for Hyperledger Fabric to start
sleep 10

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx
# Join peer0.org1.example.com to the channel
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel join -b mychannel.block

docker network ls
docker network inspect net_basic

docker inspect peer0.org1.example.com
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel create --help
```

- `docker-compose -d` 指后台运行，和 `docker run -d` 中一样
## MagnetoCorp
- `logspout` 将多个输出源的信息聚合到一处，方便从一个窗口查看


    ```bash
    find . -name monitordocker.sh
    # fabric-samples/commercial-paper/organization/magnetocorp/configuration/cli
    # ./monitordocker.sh net_basic <port_number>
    ./monitordocker.sh net_basic
    ```

![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/commercial_paper.diagram.4.png)
- 起 `cliMagnetoCorp` 容器与 `PaperNet` 交互

    ````bash
    # fabric-samples/commercial-paper/organization/magnetocorp/configuration/cli
    docker-compose -f docker-compose.yml up -d
    ````

### 安装合约
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/commercial_paper.diagram.6.png)
- 一个 chaincode 里可以有多个智能合约，chaincode 安装后里面的智能合约就可供网络上的成员使用
- MagnetoCorp 管理员用 `peer chaincode install` 命令将 `papercontract` 从本机拷到目标 peer 容器的文件系统并安装该合约

    ```bash
    docker exec cliMagnetoCorp peer chaincode install -n papercontract -v 0 -p /opt/gopath/src/github.com/contract -l node
    # -p: path to chaincode
    ```

    - 通过 `docker-compose.yml` 里的 `CORE_PEER_ADDRESS=peer0.org1.example.com:7051` 指定 peer
    - `docker-compose.yml` - `./../../../../organization/magnetocorp:/opt/gopath/src/github.com/`

### 初始化合约
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/commercial_paper.diagram.7.png)
- MagnetoCorp 管理员初始化合约后，会有一个新的容器实例来运行该合约

    ```bash
    docker exec cliMagnetoCorp peer chaincode instantiate -n papercontract -v 0 -l node -c '{"Args":["org.papernet.commercialpaper:instantiate"]}' -C mychannel -P "AND ('Org1MSP.member')"
    ```

    - `-P` 指定背书策略
    - [ ] It additionally submits an instantiate transaction to the orderer, which will include the transaction in the next block and distribute it to all peers that have joined `mychannel`, enabling any peer to execute the chaincode in their own isolated chaincode container
    - Instantiate only needs to be issued once for `papercontract` even though typically it is installed on many peers
    - `dev-peer0.org1.example.com-papercontract-0-d96abb966a1ed760663cf0a061700a902284832716c55b4cb05eca53054fe011`: `peer0.org1.example.com` 启动了该容器，合约版本是 0
## 应用依赖

```bash
# fabric-samples/commercial-paper/organization/magnetocorp/application
npm install
```

## 钱包
- As `issue.js` acts on behalf of Isabella, and therefore MagnetoCorp, it will use identity from her wallet that reflects these facts. Need to perform this one-time activity of adding appropriate X.509 credentials to her wallet

    ```bash
    # fabric-samples/commercial-paper/organization/magnetocorp/application
    node addToWallet.js
    # ls ../identity/user/isabella/wallet/
    # ls ../identity/user/isabella/wallet/User1@org1.example.com
    ```

    - `addToWallet.js` is a simple file-copying program. It moves an identity from the basic network sample to Isabella’s wallet
    - 每个身份都各自的目录
    - 目录里的公钥包含在 Isabella 的 X.509 证书里
    - 目录里的 `User1@org.example.com` 证书
        - Contains Isabella’s public key and other X.509 attributes added by the Certificate Authority at certificate creation
        - This certificate is distributed to the network so that different actors at different times can cryptographically verify information created by Isabella’s private key
        - In practice, the certificate file also contains some Fabric-specific metadata such as Isabella’s organization and role
## 发行票据

```bash
# fabric-samples/commercial-paper/organization/magnetocorp/application
node issue.js
```

# DigiBank
- DigiBank 创建与 PaperNet 交互的应用 `cliDigiBank`
- Balaji 作为 DigiBank 的员工，即终端用户，执行买入操作


```bash
./monitordocker.sh net_basic
# fabric-samples/commercial-paper/organization/digibank/configuration/cli
docker-compose -f docker-compose.yml up -d

# fabric-samples/commercial-paper/organization/digibank/application
npm install
node addToWallet.js
node buy.js
node redeem.js
```

<!-- 
peer0.org1.example.com|2019-05-10 09:23:46.494 UTC [endorser] callChaincode -> INFO 053 [mychannel][63388556] Entry chaincode: name:"cscc"
peer0.org1.example.com|2019-05-10 09:23:46.496 UTC [endorser] callChaincode -> INFO 054 [mychannel][63388556] Exit chaincode: name:"cscc"  (2ms)
peer0.org1.example.com|2019-05-10 09:23:46.497 UTC [comm.grpc.server] 1 -> INFO 055 unary call completed grpc.service=protos.Endorser grpc.method=ProcessProposal grpc.peer_address=172.18.0.1:54644 grpc.code=OK grpc.call_duration=6.427586ms
peer0.org1.example.com|2019-05-10 09:23:46.547 UTC [endorser] callChaincode -> INFO 056 [mychannel][28ce101e] Entry chaincode: name:"papercontract"
               couchdb|[notice] 2019-05-10T09:23:46.563490Z nonode@nohost <0.32580.10> 414782444d couchdb:5984 172.18.0.6 undefined GET /mychannel_papercontract/%00org.papernet.commercialpaperlist%00%22MagnetoCorp%22%00%2200001%22%00?attachments=true 200 ok 3
dev-peer0.org1.example.com-papercontract-0|2019-05-10T09:23:46.566Z info [./lib/jsontransactionserializer.js]                info: [mychannel-28ce101e] toBuffer has no schema/lacks sufficient schema to validate against {"timestamp":"2019-05-10T09:23:46.566Z"}
                    peer0.org1.example.com|2019-05-10 09:23:46.568 UTC [endorser] callChaincode -> INFO 057 [mychannel][28ce101e] Exit chaincode: name:"papercontract"  (21ms)
                    peer0.org1.example.com|2019-05-10 09:23:46.568 UTC [comm.grpc.server] 1 -> INFO 058 unary call completed grpc.service=protos.Endorser grpc.method=ProcessProposal grpc.peer_address=172.18.0.1:54644 grpc.code=OK grpc.call_duration=22.667502ms
dev-peer0.org1.example.com-papercontract-0|2019-05-10T09:23:46.567Z info [lib/handler.js]                                    info: [mychannel-28ce101e] Calling chaincode Invoke() succeeded. Sending COMPLETED message back to peer {"timestamp":"2019-05-10T09:23:46.567Z"}
                       orderer.example.com|2019-05-10 09:23:46.603 UTC [comm.grpc.server] 1 -> INFO 00e streaming call completed grpc.service=orderer.AtomicBroadcast grpc.method=Broadcast grpc.peer_address=172.18.0.1:53424 grpc.code=OK grpc.call_duration=4.139294ms
                    peer0.org1.example.com|2019-05-10 09:23:48.606 UTC [gossip.privdata] StoreBlock -> INFO 059 [mychannel] Received block [3] from buffer
                    peer0.org1.example.com|2019-05-10 09:23:48.608 UTC [committer.txvalidator] Validate -> INFO 05a [mychannel] Validated block [3] in 1ms
                                   couchdb|[notice] 2019-05-10T09:23:48.619852Z nonode@nohost <0.32580.10> 8c5d80317b couchdb:5984 172.18.0.6 undefined POST /mychannel_papercontract/_all_docs?include_docs=true 200 ok 10
                                   couchdb|[notice] 2019-05-10T09:23:48.625425Z nonode@nohost <0.32607.10> 725685b5ae couchdb:5984 172.18.0.6 undefined POST /mychannel_lscc/_all_docs?include_docs=true 200 ok 13
                                   couchdb|[notice] 2019-05-10T09:23:48.648235Z nonode@nohost <0.32607.10> fcbcfa00e7 couchdb:5984 172.18.0.6 undefined POST /mychannel_papercontract/_bulk_docs 201 ok 5
                                   couchdb|[notice] 2019-05-10T09:23:48.649495Z nonode@nohost <0.32607.10> 4428e42c6d couchdb:5984 172.18.0.6 undefined POST /mychannel_papercontract/_ensure_full_commit 201 ok 0
                                   couchdb|[notice] 2019-05-10T09:23:48.652961Z nonode@nohost <0.32580.10> 183674ef7c couchdb:5984 172.18.0.6 undefined GET /mychannel_papercontract/_index 200 ok 2
                                   couchdb|[notice] 2019-05-10T09:23:48.653138Z nonode@nohost <0.32607.10> 01cfd97eb8 couchdb:5984 172.18.0.6 undefined GET /mychannel_/statedb_savepoint?attachments=true 200 ok 3
                    peer0.org1.example.com|2019-05-10 09:23:48.660 UTC [kvledger] CommitWithPvtData -> INFO 05b [mychannel] Committed block [3] with 1 transaction(s) in 52ms (state_validation=17ms block_commit=8ms state_commit=24ms)
                                   couchdb|[notice] 2019-05-10T09:23:48.659637Z nonode@nohost <0.32607.10> f9717fdd58 couchdb:5984 172.18.0.6 undefined PUT /mychannel_/statedb_savepoint 201 ok 6
                    peer0.org1.example.com|2019-05-10 09:23:48.666 UTC [comm.grpc.server] 1 -> INFO 05c streaming call completed grpc.service=protos.Deliver grpc.method=DeliverFiltered grpc.peer_address=172.18.0.1:54644 error="context finished before block retrieved: context canceled" grpc.code=Unknown grpc.call_duration=2.069133223s
 -->





























# 开发
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
- A `StateList` doesn’t store anything about an individual state or the total list of states – it delegates all of that to the Fabric state database. This is an important design pattern – it reduces the opportunity for ledger MVCC collisions in Hyperledger Fabric
    - Multi-Version Concurrency Control
## 应用
![](https://hyperledger-fabric.readthedocs.io/en/release-1.4/_images/develop.diagram.3.png)
## Application design elements
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/designelements.html