<!-- 
```bash
cat /etc/os-release
# Ubuntu 16.04.6 LTS"
```
 -->
# 环境
## Golang

```bash
wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.12.6.linux-amd64.tar.gz

# .bash_profile
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin
# https://apple.stackexchange.com/questions/42537/why-must-i-source-bashrc-every-time-i-open-terminal-for-aliases-to-work
# https://golang.org/doc/install#tarball
# https://github.com/golang/go/wiki/SettingGOPATH
```

## Docker, docker-compose

```bash
# https://docs.docker.com/install/linux/docker-ce/ubuntu/
# Docker version 18.09.7, build 2d0083d

# https://docs.docker.com/compose/install/
# docker-compose version 1.24.0, build 0aa59064
```

## cURL

```bash
# curl 7.47.0 (x86_64-pc-linux-gnu) libcurl/7.47.0 GnuTLS/3.4.10 zlib/1.2.8 libidn/1.32 librtmp/2.3
# Protocols: dict file ftp ftps gopher http https imap imaps ldap ldaps pop3 pop3s rtmp rtsp smb smbs smtp smtps telnet tftp
# Features: AsynchDNS IDN IPv6 Largefile GSS-API Kerberos SPNEGO NTLM NTLM_WB SSL libz TLS-SRP UnixSockets
```

## `make`

```bash
apt install make
```

# Issues
- `make cryptogen`: `vendor/github.com/miekg/pkcs11/pkcs11.go:29:18: fatal error: ltdl.h: No such file or directory` 
    - `apt install libltdl-dev`
    - > `apt-get install gcc`