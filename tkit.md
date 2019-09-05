# Linux
- `printenv | less`; `echo $CHANNEL_NAME`
- [tcptrack](https://www.howtoforge.com/tracking_tcp_connectios_with_tcptrack)
    - `sudo apt-get install tcptrack`
    - `sudo tcptrack -i ens3`
# Docker
- `docker stop $(docker ps -aq)`
- `docker rm -f $(docker ps -a | grep -v learn-cli | grep -oE [0-9a-f]{12})`