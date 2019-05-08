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
    export PATH=$PATH:/usr/local/go/bin
    source ~/.profile

    go version
    ```

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
3. `export http_proxy="http://127.0.0.1:8123/"`
4. `curl www.google.com`, `curl ip.gs`, `curl ip.sb`