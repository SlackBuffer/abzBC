# DEV

```bash
#### re-make images
make peer-docker-clean && make tools-docker-clean
make peer-docker && make tools-docker

#### re-build just peer
## fabric/peer
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags nopkcs11
docker cp peer peer0.org1.example.com:/usr/local/bin 
docker start peer0.org1.example.com

## fabric
make peer-docker
docker cp build/docker/bin/peer peer0.org2.example.com:/usr/local/bin 
docker start peer0.org2.example.com


env CHANNEL_NAME=mychannel TIMEOUT=10000 docker-compose -f docker-compose-cli.yaml restart peer0.org1.example.com


Building build/docker/bin/peer
docker run -i --rm -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric -w /opt/gopath/src/github.com/hyperledger/fabric \
                -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/bin:/opt/gopath/bin \
                -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/peer/pkg:/opt/gopath/pkg \
                hyperledger/fabric-baseimage:x86_64-0.3.1 \
                go install -ldflags "-X github.com/hyperledger/fabric/common/metadata.Version=1.0.0 -X github.com/hyperledger/fabric/common/metadata.BaseVersion=0.3.1 -X github.com/hyperledger/fabric/common/metadata.BaseDockerLabel=org.hyperledger.fabric -X github.com/hyperledger/fabric/common/metadata.DockerNamespace=hyperledger -X github.com/hyperledger/fabric/common/metadata.BaseDockerNamespace=hyperledger -linkmode external -extldflags '-static -lpthread'" github.com/hyperledger/fabric/peer
touch build/docker/bin/peer  
mkdir -p build/image/peer/payload
cp build/docker/bin/peer build/sampleconfig.tar.bz2 build/image/peer/payload



# Redirecting standard output and standard error to one file
docker logs -f --details peer0.org1.example.com >& peer0.log &

docker logs -f --details peer0.org1.example.com 2> error.log > output.log

docker inspect --format='{{.LogPath}}' peer0.org1.example.com
## for linux
# /var/lib/docker/containers/728640fbfb24ef5ee3a8c0879cff115c3c261ba410f53e4cb5758f137c9c68b8/728640fbfb24ef5ee3a8c0879cff115c3c261ba410f53e4cb5758f137c9c68b8-json.log
## for mac
# screen ~/Library/Containers/com.docker.docker/Data/vms/0/tty
# cat /var/lib/docker/containers/728640fbfb24ef5ee3a8c0879cff115c3c261ba410f53e4cb5758f137c9c68b8/728640fbfb24ef5ee3a8c0879cff115c3c261ba410f53e4cb5758f137c9c68b8-json.log
# or grep
# ctrl+a, ctrl+\
```

- > https://stackoverflow.com/questions/41144589/how-to-redirect-docker-logs-to-a-single-file/48345336#48345336, https://stackoverflow.com/questions/7901517/how-to-redirect-stderr-and-stdout-to-different-files-in-the-same-line-in-script
- > https://stackoverflow.com/questions/4847691/how-do-i-get-out-of-a-screen-without-typing-exit
- > https://www.gnu.org/software/screen/manual/screen.html
- > [Go 交叉编译](https://blog.csdn.net/panshiqu/article/details/53788067)
# OP

```bash
cp build/docker/bin/peer build/sampleconfig.tar.bz2 build/image/peer/payload
Building docker peer-image
docker build -t hyperledger/fabric-peer build/image/peer
Sending build context to Docker daemon  25.28MB
COPY payload/peer /usr/local/bin

# fabric-tools (cli) 里用 peer 执行操作
cp build/docker/bin/cryptogen build/docker/bin/configtxgen build/docker/bin/peer build/sampleconfig.tar.bz2 build/image/tools/payload
Building docker tools-image
docker build -t hyperledger/fabric-tools build/image/tools
Sending build context to Docker daemon  52.06MB
COPY payload/peer /usr/local/bin

docker exec -it cli bash
which peer
# /usr/local/bin/peer
```

- cli 容器的启动命令是执行 `script.sh` 脚本
- [x] 替换 peer 二进制；停止容器后再启动
- [ ] 指定某个 peer 发送交易
- [ ] 再看程序对配置文件查找路径；容器内
- [ ] query 走 invoke 否确认
- "Building docker peer-image" 和 "Building docker orderer-image" 里有 `/var/hyperledger/production` 目录出现（`RUN mkdir -p /var/hyperledger/production $FABRIC_CFG_PATH`）
- `$FABRIC_CFG_PATH` (Fabric config path, 在 cli, peer, orderer 容器内): `/etc/hyperledger/fabric`
    - peer 的在 `docker-compose-base.yaml` (映射进去) 和 `peer-base.yaml` (环境变量的形式，`printenv | grep -i core`)
    - cli, orderer 的在 `images/%/Dockerfile.in` 里定义，编镜像后生效
    	
        ```yaml
        ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
        VOLUME /etc/hyperledger/fabric
        ADD payload/sampleconfig.tar.bz2 $FABRIC_CFG_PATH

        # tar tjf build/image/orderer/payload/sampleconfig.tar.bz2
        ```
    
        - > `ADD`: If `<src>` is a local tar archive in a recognized compression format (identity, gzip, bzip2 or xz) then it is unpacked as a directory
        - > [Use volumes](https://docs.docker.com/storage/volumes/)
- `/usr/local/lib` for 国密库，.so 等
# 书
- 《Hyperledger Fabric源代码分析与深入解读》_蔡亮.pdf Fabric 1.0
    - csdn 博客出处
- 《区块链开发实战：Hyperledger Fabric关键技术与案例分析》_冯翔.pdf Fabric 1.1
    - cc 容器外运行，IDE 调试 cc p156
    - 进程起节点 p196
    - 非初始化组织的加入 p228
- 《区块链网络构建和应用：基于超级账本Fabric的商业实践》_陆平.pdf
    - 深蓝居多机部署 Fabric 1.0
    - 密码学，数字证书，分布式
- 《白话区块链》_蒋勇.pdf
    - 密码算法，共识算法
    - 扩容，侧链，闪电网络，多链
    - 搭微链
- HyperLedger Fabric开发实战-快速掌握区块链技术_www.java1234.com.pdf Fabric 1.1
    - 内网外网部署
- 区块链原理、设计与应用 Fabric 1.0
    - 分布式
    - 密码学，PKI
    - 侧链，闪电网络
    - 进程起节点
- 深度探索区块链：Hyperledger技术与应用.pdf Fabric 1.0