- `make peer-docker`
    - 基于 `fabric-baseimage` 编译出 `peer` 二进制（给 peer-image）
    	
        ```bash
        # host 机的 fabric 源代码在编镜像的时候会映射进去
        docker run -i --rm \
        -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric \
        -w /opt/gopath/src/github.com/hyperledger/fabric \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/bin:/opt/gopath/bin \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/peer/pkg:/opt/gopath/pkg \
		hyperledger/fabric-baseimage:x86_64-0.3.1 \
		go install -ldflags "-X github.com/hyperledger/fabric/common/metadata.Version=1.0.0 -X github.com/hyperledger/fabric/common/metadata.BaseVersion=0.3.1 -X github.com/hyperledger/fabric/common/metadata.BaseDockerLabel=org.hyperledger.fabric -X github.com/hyperledger/fabric/common/metadata.DockerNamespace=hyperledger -X github.com/hyperledger/fabric/common/metadata.BaseDockerNamespace=hyperledger -linkmode external -extldflags '-static -lpthread'" github.com/hyperledger/fabric/peer

        # go install 生成的二进制 /opt/gopath/bin/peer，映射到 build/docker/bin/peer
        ```
    
        - `go install` 生成的二进制保存在 `$GOPATH/bin`，compiled packages 保存在 `$GOPATH/pkg`
        - `ldflags` 参数传给 `fabric/common/metadata` 包
            - > `LDFLAGS`: Extra flags to give to compilers when they are supposed to invoke the linker `ld`, such as `-L`. Libraries (`-lfoo`) should be added to the `LDLIBS` variable instead.
    - 将 `sampleconfig` 打包（给 peer-image）
        - `(cd sampleconfig && tar -jc *) > build/sampleconfig.tar.bz2`
    - 拷贝 `peer` 二进制和 `sampleconfig.tar.bz2` 到 `build/image/peer/payload`
    	
        ```bash
        cp build/docker/bin/peer build/sampleconfig.tar.bz2 build/image/peer/payload
        ```
    
    - 基于 `fabric-baseimage` 编 gotools（给 ccenv-image）
    	
        ```bash
        mkdir -p build/docker/gotools/bin build/docker/gotools/obj

        # host 机的 fabric 源代码在编镜像的时候会映射进去
        docker run -i --rm \
        -v /Users/slackbuffer/go/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric \
        -w /opt/gopath/src/github.com/hyperledger/fabric \
		-v /Users/slackbuffer/go/src/github.com/hyperledger/fabric/build/docker/gotools:/opt/gotools \
		-w /opt/gopath/src/github.com/hyperledger/fabric/gotools \
		hyperledger/fabric-baseimage:x86_64-0.3.1 \
		make install BINDIR=/opt/gotools/bin OBJDIR=/opt/gotools/obj # 指定 obj 目录，最终会映射到 fabric/build/docker/gotools/obj
        # ...fabric-baseimage 内 make 得到指定的 gotools...
        cp /opt/gotools/obj/gopath/bin/protoc-gen-go /opt/gotools/bin # 映射到 build/docker/gotools/bin/protoc-gen-go
        ```
    
    - 获取 `chaintool`，保存为 `build/bin/chaintool` 并赋予可执行权限（给 ccenv-image）
    - 打包得到 `goshim.tar.bz2`（给 ccenv-image）
    	
        ```makefile
        PROJECT_NAME   = hyperledger/fabric
        PKGNAME = github.com/$(PROJECT_NAME)
        GOSHIM_DEPS = $(shell ./scripts/goListFiles.sh $(PKGNAME)/core/chaincode/shim)
        build/goshim.tar.bz2: $(GOSHIM_DEPS)
	        echo "Creating $@"
	        tar -jhc -C $(GOPATH)/src $(patsubst $(GOPATH)/src/%,%,$(GOSHIM_DEPS)) > $@
        # build/goshim.tar.bz2
        ```
    
    - 将 `protoc-gen-go`, `chaintool`, `goshim.tar.bz2` 拷贝到 `build/image/ccenv/payload`

        ```bash
        mkdir -p build/image/ccenv/payload
        cp build/docker/gotools/bin/protoc-gen-go build/bin/chaintool build/goshim.tar.bz2 build/image/ccenv/payload
        ```

    - 基于 `fabric-baseimage` 编 `ccenv-image`
    	
        ```dockerfile
        COPY payload/chaintool payload/protoc-gen-go /usr/local/bin/
        # ADD 会解压并 un-tar
        ADD payload/goshim.tar.bz2 $GOPATH/src/
        RUN mkdir -p /chaincode/input /chaincode/output
        ```
    
        - [ ] 为何要编 ccenv-image 镜像？
            - 生成 chaincode 容器时会用到，生成成功后移除
    - 基于 `fabric-baseos` 编 peer 镜像
    	
        ```dockerfile
        # Dockerfile 所在目录下有 payload 目录 (Sending build context to Docker daemon  xx.xxMB)

        FROM hyperledger/fabric-baseos:x86_64-0.3.1
        ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
        RUN mkdir -p /var/hyperledger/production $FABRIC_CFG_PATH
        COPY payload/peer /usr/local/bin
        # ADD 会解压并 un-tar
        ADD  payload/sampleconfig.tar.bz2 $FABRIC_CFG_PATH
        ```
    
- `peer` 镜像目录
    - `/usr/local/bin` 存放二进制 `peer`
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

    - `/var/hyperledger/production` ([ ] 容器映射用)
- `make peer` 应该是为进程起 peer 节点用的
    - 基于 host 机编二进制 peer
- > Dockerfile source
	
    ```makefile
    build/image/%/Dockerfile: images/%/Dockerfile.in
        @cat $< \
            | sed -e 's/_BASE_NS_/$(BASE_DOCKER_NS)/g' \
            | sed -e 's/_NS_/$(DOCKER_NS)/g' \
            | sed -e 's/_BASE_TAG_/$(BASE_DOCKER_TAG)/g' \
            | sed -e 's/_TAG_/$(DOCKER_TAG)/g' \
            > $@
        @echo LABEL $(BASE_DOCKER_LABEL).version=$(PROJECT_VERSION) \\>>$@
        @echo "     " $(BASE_DOCKER_LABEL).base.version=$(BASEIMAGE_RELEASE)>>$@
    ```

- Linux 目录
    - `/etc`
        - The `/etc` directory contains all the **system-wide configuration files**. It also contains a collection of **shell scripts that start each of the system services at boot time**
        - Everything in this directory should be ***readable text***
        - > While everything in `/etc` is interesting, here are some all-time favorites: `/etc/crontab`, a file that defines when automated jobs will run; `/etc/fstab`, a table of storage devices and their associated mount points; and `/etc/passwd`, a list of the user accounts
    - `/usr/local`
        - The `/usr/local` tree is where programs that are not included with your distribution but are intended for **system-wide use** are installed
        - Programs compiled from **source code** are normally installed in `/usr/local/bin`
            - > `/usr/local/lib`
        - On a newly installed Linux system, this tree exists, but it will be empty until the system administrator puts something in it
    - `/var`
        - The `/var` directory tree is where data that is likely to **change** is stored. Various databases, spool files, user mail, and so forth, are located here
    - `/opt`
        - The `/opt` directory is used to install **“optional” software**. This is mainly used to hold commercial software products that might be installed on the system

- 指定 peer 节点发送交易后，若该 peer 节点安装有对应 chaincode，则会拉起 chaincode 容器。


Instantiating chaincode on org2/peer2...
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
CORE_PEER_LOCALMSPID=Org2MSP
CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
CORE_PEER_TLS_ENABLED=true
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
CORE_PEER_ID=cli
CORE_LOGGING_LEVEL=DEBUG
CORE_PEER_ADDRESS=peer0.org2.example.com:7051
2019-08-02 02:47:39.460 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
2019-08-02 02:47:39.483 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
2019-08-02 02:47:39.503 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 003 Using default escc
2019-08-02 02:47:39.503 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 004 Using default vscc
Error: error endorsing chaincode: rpc error: code = Unknown desc = Error starting container: pull access denied for hyperledger/fabric-yxbaseos, repository does not exist or may require 'docker login'