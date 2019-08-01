- `make peer-orderer`
    - 基于 `fabric-baseimage` 编译出 `orderer` 二进制（给 orderer-image）
    	
        ```bash
        # host 机的 fabric 源代码在编镜像的时候会映射进去
        docker run -i --rm \
        -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric \
        -w /opt/gopath/src/github.com/hyperledger/fabric \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/bin:/opt/gopath/bin \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/orderer/pkg:/opt/gopath/pkg \
		hyperledger/fabric-baseimage:x86_64-0.3.1 \
		go install github.com/hyperledger/fabric/orderer
        # 省略 go install 的 flags
        # go install 生成的二进制 /opt/gopath/bin/orderer，映射到 build/docker/bin/orderer
        ```
    
    - 将 `sampleconfig` 打包（给 orderer-image）
        - `(cd sampleconfig && tar -jc *) > build/sampleconfig.tar.bz2`
    - 拷贝 `orderer` 二进制和 `sampleconfig.tar.bz2` 到 `build/image/orderer/payload`
    	
        ```bash
        cp build/docker/bin/orderer build/sampleconfig.tar.bz2 build/image/orderer/payload
        ```

    - 基于 `fabric-baseos` 编 orderer 镜像
    	
        ```dockerfile
        # Dockerfile 所在目录下有 payload 目录 (Sending build context to Docker daemon  xx.xxMB)
        
        FROM hyperledger/fabric-baseos:x86_64-0.3.1
        ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
        RUN mkdir -p /var/hyperledger/production $FABRIC_CFG_PATH
        COPY payload/orderer /usr/local/bin
        ADD payload/sampleconfig.tar.bz2 $FABRIC_CFG_PATH/

        EXPOSE 7050
        
        ```

- `orderer` 镜像目录
    - `/usr/local/bin` 存放二进制 `orderer`
    - `/etc/hyperledger/fabric` (`FABRIC_CFG_PATH`) 存放 sampleconfig
    	
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

    - `/var/hyperledger/production` ([ ] 容器映射用)