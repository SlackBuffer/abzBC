# [cURL](https://curl.haxx.se/download.html)
- Mac uses 7.63.0, Ubuntu 16.04 uses 7.64.1
# Docker

```bash
docker --version
docker-compose --version
```

# Golang
- (Install and upgrade)[https://golang.org/doc/install#install]

    ```bash
    # lsb_release -a
    uname -r

    # ubuntu
    ## uninstall
    # sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.12.5.linux-amd64.tar.gz
    # put in ~/.zshrc
    export PATH=$PATH:/usr/local/go/bin # for just Golang, only this one is enough
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    source ~/.zshrc
    # source ~/.profile

    go version
    ```

- https://github.com/golang/go/wiki/SettingGOPATH
# Shadowsocks
- https://jingsam.github.io/2016/05/08/setup-shadowsocks-http-proxy-on-ubuntu-server.html

```json
// /etc/shadowsocks.json
{
    "server":"ip_address", // 必须是 IP
    "server_port":port_number,
    "local_address": "127.0.0.1",
    "local_port":1080,
    "password":"your_password",
    "timeout":300,
    "method":"aes-256-cfb",
    "fast_open": false
}

// /etc/polipo/config
logSyslog = true
logFile = /var/log/polipo/polipo.log

proxyAddress = "0.0.0.0"
socksParentProxy = "127.0.0.1:1080"
socksProxyType = socks5

chunkHighMark = 50331648
objectHighMark = 16384

serverMaxSlots = 64
serverSlots = 16
serverSlots1 = 32
```

1.（若没启动）`sudo sslocal -c /etc/shadowsocks.json -d start`
2.（若没启动）`sudo /etc/init.d/polipo restart`
1. `export http_proxy="http://127.0.0.1:8123/"`
2. `curl www.google.com`, `curl ip.gs`, `curl ip.sb`
# NodeJS
- https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-16-04
- https://github.com/nvm-sh/nvm#manual-install
# Python (Ubuntu16.04)
- Python 2.7 is required
# Fabric scripts
- Latest production release: `curl -sSL http://bit.ly/2ysbOFE | bash -s` or `curl https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s`
    - `binaryDownload()` 将二进制（`bin`）和配置（`/config`）下载下来放入 `fabric-sample` 目录，自行拷入即可，不需通过脚本
- https://github.com/hyperledger/fabric-samples
# bat
- https://github.com/sharkdp/bat#installation

    ```bash
    wget https://github.com/sharkdp/bat/releases/download/v0.10.0/bat_0.10.0_amd64.deb
    sudo dpkg -i bat_0.10.0_amd64.deb
    ```

- https://github.com/sharkdp/bat#highlighting-theme