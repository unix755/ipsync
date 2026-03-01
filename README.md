# IPSync

## Features

- Send network information to remote server or file
- Print network information from remote server or file
- Encrypt network information and transmit it to network services
- Periodically update WireGuard endpoint IP

## Example

```sh
ipsync help [command]
```

## Install

```sh
# system is linux(debian,redhat linux,ubuntu,fedora...) and arch is amd64
curl -Lo /usr/local/bin/ipsync https://github.com/unix755/ipsync/releases/latest/download/ipsync-linux-amd64
chmod +x /usr/local/bin/ipsync

# system is freebsd and arch is amd64
curl -Lo /usr/local/bin/ipsync https://github.com/unix755/ipsync/releases/latest/download/ipsync-freebsd-amd64
chmod +x /usr/local/bin/ipsync
```

## Install Service(WebDAV usage examples)

### Linux(systemd)

```sh
curl -Lo "/etc/systemd/system/ipsync.service" "https://github.com/unix755/ipsync/raw/main/service/systemd/send_webdav.service"
systemctl enable ipsync.service && systemctl restart ipsync.service && systemctl status ipsync.service
curl -Lo "/etc/systemd/system/ipsync.timer" "https://github.com/unix755/ipsync/raw/main/service/systemd/ipsync.timer"
systemctl enable ipsync.timer && systemctl restart ipsync.timer && systemctl status ipsync.timer
```

### OpenWRT(init.d)

```sh
curl -Lo "/etc/init.d/ipsync" "https://github.com/unix755/ipsync/raw/main/service/init.d/send_webdav"
chmod +x /etc/init.d/ipsync
service ipsync enable && service ipsync restart
```

### FreeBSD(rc.d)

```sh
curl -Lo "/etc/rc.d/ipsync" "https://github.com/unix755/ipsync/raw/main/service/rc.d/send_webdav"
chmod +x /etc/rc.d/ipsync
service ipsync enable && service ipsync restart
```

## Compile

### How to compile if prebuilt binaries are not found

```sh
git clone https://github.com/unix755/ipsync.git
cd ipsync
export CGO_ENABLED=0
go build -v -trimpath -ldflags "-s -w"
```

## License

- **GPL-3.0 License**
- See `LICENSE` for details
