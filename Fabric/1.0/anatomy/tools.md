- `make tools-peer`
    - 基于 `fabric-baseimage` 编译出 `cryptogen` 二进制
    	
        ```bash
        docker run -i --rm  \
        -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric \
        -w /opt/gopath/src/github.com/hyperledger/fabric \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/bin:/opt/gopath/bin \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/cryptogen/pkg:/opt/gopath/pkg \
		hyperledger/fabric-baseimage:x86_64-0.3.1 \
		go install github.com/hyperledger/fabric/common/tools/cryptogen
        # 省略 go install 的 flags
        # go install 生成的二进制 /opt/gopath/bin/cryptogen，映射到 build/docker/bin/cryptogen
        ```
    
    - 基于 `fabric-baseimage` 编译出 `configtxgen` 二进制
    	
        ```bash
        docker run -i --rm \
        -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric \
        -w /opt/gopath/src/github.com/hyperledger/fabric \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/bin:/opt/gopath/bin \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/configtxgen/pkg:/opt/gopath/pkg \
		hyperledger/fabric-baseimage:x86_64-0.3.1 \
		go install github.com/hyperledger/fabric/common/configtx/tool/configtxgen
        # 省略 go install 的 flags
        # go install 生成的二进制 /opt/gopath/bin/configtxgen，映射到 build/docker/bin/configtxgen
        ```
    
    - 基于 `fabric-baseimage` 编译出 `peer` 二进制（同 peer 处）
        - `build/docker/bin/peer`
    - 将 `sampleconfig` 打包
        - `(cd sampleconfig && tar -jc *) > build/sampleconfig.tar.bz2`
    - 拷贝 `peer`, `cryptogen`, `configtxgen`, `sampleconfig.tar.bz2` 到 `build/image/tools/payload`
    	
        ```bash
        cp build/docker/bin/cryptogen build/docker/bin/configtxgen build/docker/bin/peer build/sampleconfig.tar.bz2 build/image/tools/payload
        ```

    - 基于 `fabric-baseimage` 编 tools 镜像
    	
        ```dockerfile
        FROM hyperledger/fabric-baseimage:x86_64-0.3.1
        ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
        VOLUME /etc/hyperledger/fabric
        ADD  payload/sampleconfig.tar.bz2 $FABRIC_CFG_PATH
        COPY payload/cryptogen /usr/local/bin
        COPY payload/configtxgen /usr/local/bin
        COPY payload/peer /usr/local/bin
        ```
    
- `peer` 镜像目录
    - `/usr/local/bin` 存放二进制 `peer`, `cryptogen`, `configtxgen`
    - `/etc/hyperledger/fabric` (<mark>`FABRIC_CFG_PATH`</mark>) 存放 sampleconfig
    	
        ```
        ├── configtx.yaml
        ├── core.yaml
        ├── msp
        │   ├── admincerts
        │   │   └── admincert.pem
        │   ├── cacerts
        │   │   └── cacert.pem
        │   ├── config.yaml
        │   ├── keystore
        │   │   └── key.pem
        │   ├── signcerts
        │   │   └── peer.pem
        │   └── tlscacerts
        │       └── cert.pem
        └── orderer.yaml
        ```