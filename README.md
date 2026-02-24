# Netinfo

## Features

- Send network information to a remote server or file
- Show local network information

## Example

```sh
# Show local network information
netinfo show
## Show send preload
netinfo show -p

# Send local network information to a file
netinfo send file -filepath="./home.json"
## Send local network information to a file and encrypt the file
netinfo send file -filepath="./home.json" -encryption_key="admin123"
## Send local network information to a file and encrypt the file every 5 min
netinfo send file -filepath="./home.json" -encryption_key="admin123" -interval="5m"

# Send local network information to s3 server
netinfo send s3 -endpoint="https://s3.amazonaws.com" -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Send local network information to minio s3 server
netinfo send s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Send local network information to minio s3 server and encrypt the file
netinfo send s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123"
## Send local network information to minio s3 server and encrypt the file every 5 min
netinfo send s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123" -interval="5m"

# Send local network information to webdav server
netinfo send webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json"
## Send local network information to webdav server and encrypt the file
netinfo send webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123"
## Send local network information to webdav server and encrypt the file every 5 min
netinfo send webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123" -interval="5m"


# Receive local network information from a file
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" file -filepath="./home.json"
## Receive local network information from a file and decrypt the file
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" file -filepath="./home.json" -encryption_key="admin123"
## Receive local network information from a file and decrypt the file every 5 min
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" -interval="5m" file -filepath="./home.json" -encryption_key="admin123"

# Receive local network information from s3 server
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" s3 -endpoint="https://s3.amazonaws.com" -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Receive local network information from minio s3 server
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Receive local network information from minio s3 server and decrypt the file
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123"
## Receive local network information from minio s3 server and decrypt the file every 5 min
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" -interval="5m" s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123"

# Receive local network information from webdav server
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json"
## Receive local network information from webdav server and decrypt the file
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123"
## Receive local network information from webdav server and decrypt the file every 5 min
netinfo receive -remote_interface="pppoe-wan" -wg_interface="wg0" -interval="5m" webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123"

# Decrypt a encrypted file
netinfo decrypt -filepath "./home.json" -encryption_key="admin123"
```

## Install

```sh
# system is linux(debian,redhat linux,ubuntu,fedora...) and arch is amd64
curl -Lo /usr/local/bin/netinfo https://github.com/unix755/netinfo/releases/latest/download/netinfo-linux-amd64
chmod +x /usr/local/bin/netinfo

# system is freebsd and arch is amd64
curl -Lo /usr/local/bin/netinfo https://github.com/unix755/netinfo/releases/latest/download/netinfo-freebsd-amd64
chmod +x /usr/local/bin/netinfo
```

## Install Service(WebDAV usage examples)

### Linux(systemd)

```sh
curl -Lo "/etc/systemd/system/netinfo.service" "https://github.com/unix755/netinfo/raw/main/service/systemd/netinfo_sender_webdav.service"
systemctl enable netinfo.service && systemctl restart netinfo.service && systemctl status netinfo.service
curl -Lo "/etc/systemd/system/netinfo.timer" "https://github.com/unix755/netinfo/raw/main/service/systemd/netinfo_sender.timer"
systemctl enable netinfo.timer && systemctl restart netinfo.timer && systemctl status netinfo.timer
```

### Alpine Linux(openrc)

```sh
curl -Lo "/etc/init.d/netinfo" "https://github.com/unix755/netinfo/raw/main/service/openrc/netinfo_sender_webdav"
chmod +x /etc/init.d/netinfo
rc-update add netinfo && rc-service netinfo restart && rc-service netinfo status
```

### FreeBSD(rc.d)

```sh
mkdir /usr/local/etc/rc.d/
curl -Lo "/usr/local/etc/rc.d/netinfo" "https://github.com/unix755/netinfo/raw/main/service/rc.d/netinfo_sender_webdav"
chmod +x /usr/local/etc/rc.d/netinfo
service netinfo enable && service netinfo restart && service netinfo status
```

### OpenWRT(init.d)

```sh
curl -Lo "/etc/init.d/netinfo" "https://github.com/unix755/netinfo/raw/main/service/init.d/netinfo_sender_webdav"
chmod +x /etc/init.d/netinfo
service netinfo enable && service netinfo restart && service netinfo status
```

## Compile

### How to compile if prebuilt binaries are not found

```sh
git clone https://github.com/unix755/netinfo.git
cd netinfo
export CGO_ENABLED=0
go build -v -trimpath -ldflags "-s -w"
```

### For mipsle router

```sh
git clone https://github.com/unix755/netinfo.git
cd netinfo
export GOOS=linux
export GOARCH=mipsle
export GOMIPS=softfloat
export CGO_ENABLED=0
go build -v -trimpath -ldflags "-s -w"
```

## License

- **GPL-3.0 License**
- See `LICENSE` for details
