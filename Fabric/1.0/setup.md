- [x] `e2e_cli_default` 在哪里指定
    - https://docs.docker.com/compose/networking/
- `./network_setup down`; `make clean`; `source download-dockerimages.sh -c x86_64-1.0.0 -f x86_64-1.0.0`
# 搭建 Fabric 1.0 环境
- > Works for 1.0.3
- `network_setup.sh` -> `generateArtifacts.sh` - `docker-compose-e2e.yaml`

```bash
mkdir -p ~/go/src/github.com/hyperledger
cd ~/go/src/github.com/hyperledger
git clone https://github.com/hyperledger/fabric.git
git checkout v1.0.0

cd ~/go/src/github.com/hyperledger/fabric/examples/e2e_cli/
# 下载 Fabric Dockercd 1.0.0 镜像
source download-dockerimages.sh -c x86_64-1.0.0 -f x86_64-1.0.0
```

- `./network_setup.sh up`
	
    ```yaml
    # Error: Error endorsing chaincode: rpc error: code = Unknown desc = Error starting container: API error (404): {"message":"network e2ecli_default not found"}

    # fabric/examples/e2e_cli/base/peer-base.yaml
    - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=e2e_cli_default
    ```

    - 编译生成 Fabric 公私钥、证书的程序（`cryptogen`），程序目录在 `fabric/release/darwin-amd64`
    - 基于 `configtx.yaml` 生成创世区块和通道相关信息，并保存在 `fabric/examples/e2e_cli/channel-artifacts` 文件夹
    - 基于 `crypto-config.yaml` 生成公私钥和证书信息，并保存在 `fabric/examples/e2e_cli/crypto-config` 文件夹中
    - 基于 `docker-compose-cli.yaml` 启动 1 Orderer + 4 Peer + 1 CLI 的 Fabric 容器 
    - 在 CLI 启动的时候，会运行 `fabric/examples/e2e_cli/scripts/script.sh` 文件，这个脚本文件包含了创建 Channel，加入 Channel，安装 `Example02，运行` `Example02` 等功能
- 查询 a 账户余额

    ```bash
    docker exec -it cli bash

    # channel 名 mychannel；chaincode 名 mycc；查询账户 a
    peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
    ```

- a 账户转账 20 给 b 账户

    ```bash
    peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc -c '{"Args":["invoke","a","b","20"]}'

    peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
    ```

- `./network_setup.sh down` 关闭网络
- > http://www.cnblogs.com/studyzy/p/7437157.html
## 详解
- `fabric/common/tools/cryptogen/`
- `fabric/common/configtx/tool/configtxgen`
### 生成公私钥
- Fabric 中有 2 种类型的公私钥和证书，一种是给节点之间通讯安全而准备的 TLS 证书，另一种是用户登录和权限控制的用户证书
    - 这些证书本应由 CA 颁发，但此处的测试环境并没有启用 CA 节点，Fabric 提供了一个工具 `cryptogen`
#### 编译生成 `cryptogen`

```bash
# fabric 目录下
make cryptogen
# 生成 fabric/build/bin/cryptogen

# go build github.com/hyperledger/fabric/vendor/github.com/miekg/pkcs11: invalid flag in #cgo LDFLAGS: -I/usr/local/share/libtool

# ./fabric_v1.0.0/common/tools/cryptogen/go build -tags nopkcs11 main.go 手动 build 并移到目标位置

# 或 export CGO_LDFLAGS_ALLOW=".*"
# https://forum.golangbridge.org/t/invalid-flag-in-cgo-ldflags/10020
```

- `cryptogen` is an utility for generating Hyperledger Fabric key material
    - It is provided as a means of pre-configuring a network for testing purposes
    - It would normally not be used in the operation of a production network
- > https://hyperledger-fabric.readthedocs.io/en/release-1.4/commands/cryptogen.html
#### 配置 `crypto-config.yaml`

```yaml
# fabric/examples/e2e_cli
- Name: Org2
    Domain: org2.example.com
    Template:
      Count: 2
    Users:
      Count: 1
```

- `Name` 和 `Domain` 就是关于这个组织的名字和域名，生成证书时证书内会包含该信息
- `Template.Count=2` 表名要生成 2 套公私钥和证书，一套是 `peer0.org2` 的，一套是 `peer1.org2`
- `Users.Count=1` 是说每个 Template 下面会有几个普通 User（注意，Admin 是 Admin，不包含在这个计数中），此处配置了 1，也就是说我们只需要一个普通用户 `User1@org2.example.com`
    - 可以根据实际需要调整这个配置文件，增删 Org Users 等
#### 生成公私钥和证书

```bash
# fabric/examples/e2e_cli
../../build/bin/cryptogen generate --help
../../build/bin/cryptogen generate --config=./crypto-config.yaml
# default output directory: --output="crypto-config"

tree crypto-config
# 与 crypto-config.yaml 对照查看
```

### 生成创世块和 Channel
#### 编译生成 `configtxgen`

```bash
# fabric 目录下
make configtxgen
# 生成 fabric/build/bin/configtxgen
```

- The `configtxgen` command allows users to create and inspect channel config related artifacts
- The content of the generated artifacts is dictated by the contents of `configtx.yaml`
- > https://hyperledger-fabric.readthedocs.io/en/release-1.4/commands/configtxgen.html
#### 配置 `configtx.yaml`
- `TwoOrgsOrdererGenesis`, `TwoOrgsChannel`
#### 生成创世块

```bash
# fabric/examples/e2e_cli
../../build/bin/configtxgen --help
../../build/bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
# fabric/examples/e2e_cli/configtxgen: Profiles.TwoOrgsOrdererGenesis
```

#### Create the channel transaction artifact

```bash
# fabric/examples/e2e_cli
# 指定 channel 名称
../../build/bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
```

#### Define the anchor peer for orgs on the channel

```bash
# fabric/examples/e2e_cli
../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP

../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP

# tree channel-artifacts
```

### 配置 Fabric 环境
#### 配置 orderer
- Orderer 配置：`fabric/examples/e2e_cli/base/docker-compose-base.yaml`
	
    ```yaml
    services:
        orderer.example.com:
            container_name: orderer.example.com
            image: hyperledger/fabric-orderer
            environment:
                - ORDERER_GENERAL_LOGLEVEL=debug
                - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
                - ORDERER_GENERAL_GENESISMETHOD=file
                # genesis file
                - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
                - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
                - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
                # enabled TLS
                - ORDERER_GENERAL_TLS_ENABLED=true
                - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
                - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
                - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
            working_dir: /opt/gopath/src/github.com/hyperledger/fabric
            command: orderer
            volumes:
                # 映射 genesis file
                - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
                - ../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/var/hyperledger/orderer/msp
                - ../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/:/var/hyperledger/orderer/tls
            ports:
                - 7050:7050
    ```

- Peer 配置：`fabric/examples/e2e_cli/base/docker-compose-base.yaml`, `fabric/examples/e2e_cli/base/peer-base.yaml`
#### 配置 peer

```yaml
# fabric/examples/e2e_cli/base/docker-compose-base.yaml
services:
    peer0.org1.example.com:
        container_name: peer0.org1.example.com
        extends:
            file: peer-base.yaml
            service: peer-base
        environment:
            - CORE_PEER_ID=peer0.org1.example.com
            - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
            - CORE_PEER_CHAINCODELISTENADDRESS=peer0.org1.example.com:7052
            - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.example.com:7051
            - CORE_PEER_LOCALMSPID=Org1MSP
        volumes:
            - /var/run/:/host/var/run/
            - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/fabric/msp
            - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls:/etc/hyperledger/fabric/tls
        ports:
            - 7051:7051
            - 7052:7052
            - 7053:7053
        
# fabric/examples/e2e_cli/base/peer-base.yaml
services:
    peer-base:
        image: hyperledger/fabric-peer
        environment:
            - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
            # the following setting starts chaincode containers on the same
            # bridge network as the peers
            # https://docs.docker.com/compose/networking/
            - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=e2e_cli_default
            #- CORE_LOGGING_LEVEL=ERROR
            - CORE_LOGGING_LEVEL=DEBUG
            - CORE_PEER_TLS_ENABLED=true
            - CORE_PEER_GOSSIP_USELEADERELECTION=true
            - CORE_PEER_GOSSIP_ORGLEADER=false
            - CORE_PEER_PROFILE_ENABLED=true
            - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
            - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
            - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
        working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
        command: peer node start
```

#### 配置 CLI
- CLI 在 Fabric 网络中扮演客户端角色，开发测试的时候可以用 CLI 来代替 SDK，执行各种 SDK 的操作
- **CLI 会和 Peer 相连**，把指令发送给对应的 Peer 执行
- `docker-compose-cli.yaml` 存放了**整个 Fabric Docker 环境的配置**

```yaml
# fabric/examples/e2e_cli/docker-compose-cli.yaml
version: '2'

services:

  orderer.example.com:
    extends:
      file:   base/docker-compose-base.yaml
      service: orderer.example.com
    container_name: orderer.example.com

  peer0.org1.example.com:
    container_name: peer0.org1.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer0.org1.example.com

  peer1.org1.example.com:
    container_name: peer1.org1.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer1.org1.example.com

  peer0.org2.example.com:
    container_name: peer0.org2.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer0.org2.example.com

  peer1.org2.example.com:
    container_name: peer1.org2.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer1.org2.example.com

  cli:
    container_name: cli
    image: hyperledger/fabric-tools
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=DEBUG
      - CORE_PEER_ID=cli
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
      - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    # CLI 启动后会执行如下脚本，手动操作需注释该行
    # command: /bin/bash -c './scripts/script.sh ${CHANNEL_NAME}; sleep $TIMEOUT'
    volumes:
        - /var/run/:/host/var/run/
        # 映射 chaincode
        - ../chaincode/go/:/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/go
        - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
        - ./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/
        - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
    # https://docs.docker.com/compose/compose-file/compose-file-v2/#depends_on
    depends_on:
      - orderer.example.com
      - peer0.org1.example.com
      - peer1.org1.example.com
      - peer0.org2.example.com
      - peer1.org2.example.com
```

- 启动 CLI 默认连接的是 `peer0.org1.example.com:7051`
- 启用 TLS
- 默认以 `Admin@org1.example.com` 身份连接到 Peer 
- `script.sh` 执行 Fabric 环境的初始化和 ChainCode 的安装、运行
### 初始化 Fabric 环境
#### 启动 Fabric 环境的容器

```bash
# fabric/examples/e2e_cli
docker-compose -f docker-compose-cli.yaml up -d
# docker ps
# 1 Orderer + 4 Peer + 1 CLI
```

#### 创建 channel

```bash
docker exec -it cli bash

# critical
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# -o: orderer container name + port
# -c: channel 名在 Create the channel transaction artifact 时指定
peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile $ORDERER_CA
# [channelCmd] readBlock -> DEBU 020 Received block:0
```

#### peer 加入 channel

```bash
# peer0.org1 加入 channel (CLI 默认连接 peer0.org1.example.com:7051)
peer channel join -b mychannel.block
# [channelCmd] executeJoin -> INFO 006 Peer joined the channel!

# peer1.org1 加入 channel
CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer1.org1.example.com:7051

peer channel join -b mychannel.block

# peer0.org2 加入 channel
CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer channel join -b mychannel.block

# peer1.org2 加入 channel
CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer1.org2.example.com:7051

peer channel join -b mychannel.block
```

#### 更新锚节点

```bash
# peer0.org1 设为 org1 的锚节点
CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer0.org1.example.com:7051

peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx --tls true --cafile $ORDERER_CA

# peer0.org2 设为 org2 的锚节点
CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx --tls true --cafile $ORDERER_CA
```

### 安装运行 chaincode
#### Install chaincode

```bash
# 切换到 peer0.org1
CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer0.org1.example.com:7051

peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02
```

- 安装的过程是对指定的代码进行编译打包，并把打包好的文件发送到 Peer，等待接下来的实例化
#### Instantiate chaincode

```bash
# peer0.org1
# takes a while
peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc -v 1.0 -c '{"Args":["init","a","100","b","200"]}' -P "OR('Org1MSP.member','Org2MSP.member')"
# new container: dev-peer0.org1.example.com-mycc-1.0

docker logs -f peer0.org1.example.com

# 2019-06-19 06:25:27.679 UTC [chaincode-platform] generateDockerfile -> DEBU 3fd
# FROM hyperledger/fabric-baseos:x86_64-0.3.1
# ADD binpackage.tar /usr/local/bin
# LABEL org.hyperledger.fabric.chaincode.id.name="mycc" \
#       org.hyperledger.fabric.chaincode.id.version="1.0" \
#       org.hyperledger.fabric.chaincode.type="GOLANG" \
#       org.hyperledger.fabric.version="1.0.0" \
#       org.hyperledger.fabric.base.version="0.3.1"
# ENV CORE_CHAINCODE_BUILDLEVEL=1.0.0
# ENV CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/peer.crt
# COPY peer.crt /etc/hyperledger/fabric/peer.crt
# 2019-06-19 06:25:27.682 UTC [util] DockerBuild -> DEBU 3fe Attempting build with image hyperledger/fabric-ccenv:x86_64-1.0.0
# 2019-06-19 06:25:40.108 UTC [dockercontroller] deployImage -> DEBU 3ff Created image: dev-peer0.org1.example.com-mycc-1.0
# 2019-06-19 06:25:40.109 UTC [dockercontroller] Start -> DEBU 400 start-recreated image successfully
# 2019-06-19 06:25:40.109 UTC [dockercontroller] createContainer -> DEBU 401 Create container: dev-peer0.org1.example.com-mycc-1.0
```

- 实例化链上代码主要是在 Peer 所在的机器上对前面安装好的 chaincode 进行包装，生成 Docker 镜像和容器
#### 发起交易

```bash
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n mycc -c '{"Args":["invoke","a","b","10"]}'
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
```

#### 在其他节点上发起交易

```bash
# 切换到 peer0.org2 安装同一个合约
CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02
```

- 通过 peer0.org2 发起查询交易
	
    ```bash
    peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
    # dev-peer0.org2.example.com-mycc-1.0
    ```

    - `mycc` 已经在 `org1` 上已经实例化，即对应的区块已经生成，所以在 org2 上不能再次初始化
- > https://www.cnblogs.com/studyzy/p/7451276.html
## 升级合约

```bash
peer chaincode install -n mycc -v 2.0 -p github.com/hyperledger/fabric/examples/chaincode/go/checkBill2

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer chaincode upgrade -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n mycc -v 2.0 -c '{"Args":[]}' -P "OR('Org1MSP.member','Org2MSP.member')"

docker exec -it peer0.org1.example.com /bin/bash
ls /var/hyperledger/production/chaincodes
# https://github.com/hyperledger/composer/issues/4043

# https://blog.csdn.net/ASN_forever/article/details/87969949
# https://blog.csdn.net/tiandiwuya/article/details/79418093
```

## 开发模式

```bash
# inside cc container
# cd and build
cd projectLog/ && go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./projectLog

# inside cli
# -p 跟的目录相对于 $GOPATH，是目录
peer chaincode install -p chaincodedev/chaincode/projectLog -n mycc -v 0
# -C: channel name
peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc

# add
# peer chaincode invoke -n mycc -c '{"Args":["add","{\"ProductID\":\"prod8\",\"Step\":0,\"Stage\":2,\"Properties\":[{\"PropertyID\":\"1\",\"PropertyName\":\"资产A\",\"PropertyType\":\"类型A\",\"BuyInfo\":[\"购买记录A\",\"购买记录B\"],\"ReleaseAmount\":[\"释放额度A\",\"释放额度B\"]},{\"PropertyID\":\"2\",\"PropertyName\":\"资产A\",\"PropertyType\":\"类型A\",\"BuyInfo\":[\"购买记录A\",\"购买记录B\"],\"ReleaseAmount\":[\"释放额度A\",\"释放额度B\"]}]}"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["add","{\"ProductID\":\"prod8\",\"Step\":1,\"Stage\":3}"]}' -C myc

# 查产品编号返项目最新记录
peer chaincode query -n mycc -c '{"Args":["query_product","prod8"]}' -C myc

# 基础资产编号查询相关项目记录，每次返一条
# property id, index, cutNum
peer chaincode query -n mycc -c '{"Args":["query_property","2","0","3"]}' -C myc
```

```bash
# fabric-samples/chaincode-docker-devmode
docker-compose -f docker-compose-simple.yaml up 

docker exec -it chaincode bash
cd chaincode_example02/go
go build chaincode_example02.go
CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc2:0 ./chaincode_example02

docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/chaincode_example02/go -n mycc2 -v 0
peer chaincode instantiate -n mycc2 -v 0 -c '{"Args":["init","a","100","b","200"]}' -C myc
peer chaincode invoke -n mycc2 -c '{"Args":["invoke","a","b","1"]}' -C myc
peer chaincode query -n mycc2 -c '{"Args":["query","a"]}' -C myc



cd chaincode_example04
go build chaincode_example04.go
CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=myccp:0 ./chaincode_example04

mkdir -p /opt/gopath/src/github.com/hyperledger/fabric/common/util
cp chaincode/utils.go /opt/gopath/src/github.com/hyperledger/fabric/common/util
peer chaincode install -p chaincodedev/chaincode/chaincode_example04 -n myccp -v 0
peer chaincode instantiate -n myccp -v 0 -c '{"Args":["init","passThgu","1"]}' -C myc
peer chaincode invoke -n myccp -c '{"Args":["invoke","mycc", "passThgu", "1"]}' -C myc




go build reconcile.go && CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc:0 ./reconcile

mkdir -p /opt/gopath/src/github.com/hyperledger/fabric/common/util
cp chaincode/utils.go /opt/gopath/src/github.com/hyperledger/fabric/common/util

peer chaincode install -p chaincodedev/chaincode/reconcile -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc

peer chaincode invoke -C myc -n mycc -c '{"Args":["putICBCUserFlow","icbc","{\"ACC_NO\":\"\",\"ACC_NAME\":\"\",\"CURR_TYPE\":\"\",\"BEGIN_DATE\":\"\",\"END_DATE\":\"\",\"MIN_AMT\":\"\",\"MAX_AMT\":\"888\",\"BANK_TYPE\":\"\",\"NEXT_TAG\":\"\",\"TOTAL_NUM\":\"\",\"DUEBILL_NO\":\"\",\"ACCT_SEQ\":\"\",\"DRCRF\":\"\",\"VOUH_NO\":\"\",\"DEBIT_AMOUNT\":\"\",\"CREDIT_AMOUNT\":\"\",\"BALANCE\":\"\",\"RECIPBK_NO\":\"\",\"RECIPBK_NAME\":\"\",\"RECIPACC_NO\":\"\",\"RECIP_NAME\":\"\",\"SUMMARY\":\"\",\"USE_CN\":\"\",\"POST_SCRIPT\":\"\",\"BUS_CODE\":\"\",\"DATE\":\"\",\"TIME\":\"\",\"REF\":\"\",\"OREF\":\"\",\"EN_SUMMARY\":\"\",\"BUS_TYPE\":\"\",\"VOUH_TYPE\":\"\",\"ADD_INFO\":\"\",\"TOUT_FO\":\"\",\"ONLY_SEQUENCE\":\"0006\",\"AGENT_ACCT_NAME\":\"\",\"AGENT_ACCT_NO\":\"\",\"UP_DTRANF\":\"\",\"VALUE_DATE\":\"\",\"TRX_CODE\":\"\",\"REF1\":\"\",\"OREF1\":\"\",\"CASHF\":\"\",\"BUSI_DATE\":\"\",\"BUSI_TIME\":\"\",\"SEQ_NO\":\"\",\"MG_NO\":\"\",\"MG_ACC_NO\":\"\",\"MG_CURR_TYPE\":\"\",\"CASH_EXF\":\"\"}"]}'

peer chaincode query -C myc -n mycc -c '{"Args":["getICBCUserFlow","icbc", "0006"]}'


peer chaincode invoke -C myc -n mycc -c '{"Args":["putZJUserFlow","zhejin","{\"ID\":\"xxx\",\"INNERBILL_NO\":\"0006\",\"OUT_CODE\":\"aaa\",\"ORDER_TYPE\":\"aaa\",\"FUND_ACCOUNT\":\"aaa\",\"CUSTOMER_NAME\":\"lisan\",\"AMOUNT\":\"88\",\"ORDER_STATUS\":\"normal\",\"DEAL_STATUS\":\"normal\",\"BANK_NO\":\"11111\",\"TRANS_CODE\":\"111\",\"MONEY_TYPE\":\"01\"}"]}'

peer chaincode query -C myc -n mycc -c '{"Args":["getZJUserFlow","zhejin", "0006"]}'

peer chaincode query -C myc -n mycc -c '{"Args":["getZJCheckResult","zhejin", "0006"]}'

# https://github.com/hyperledger/fabric-samples/tree/release/chaincode-docker-devmode
```

## [ ] 创建多通道

```bash
# ./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/yourchannel.tx -channelID yourchannel

../../build/bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel2.tx -channelID channel2
# tree channel-artifacts
```

- > https://blog.csdn.net/mx1222/article/details/82627785
# 多机部署
- https://medium.com/@thotanarendranathreddy/hyperledger-fabric-multi-orgs-multi-nodes-with-kafka-zookeeper-cluster-with-swarm-cluster-5d38be8b1fbc