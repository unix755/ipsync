# IPSync

## Features

- Send network information to remote server or file
- Receive network information from remote server or file
- Encrypted or decrypted transmission of network information
- Show local network information
- Periodically update WireGuard endpoint IP

## Example

```sh
# Show local network information
ipsync show
## Show send preload
ipsync show -p

# Send local network information to a file
ipsync send file -filepath="./home.json"
## Send local network information to a file and encrypt the file
ipsync send file -filepath="./home.json" -encryption_key="admin123"
## Send local network information to a file and encrypt the file every 5 min
ipsync send file -filepath="./home.json" -encryption_key="admin123" -interval="5m"

# Send local network information to s3 server
ipsync send s3 -endpoint="https://s3.amazonaws.com" -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Send local network information to minio s3 server
ipsync send s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Send local network information to minio s3 server and encrypt the file
ipsync send s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123"
## Send local network information to minio s3 server and encrypt the file every 5 min
ipsync send s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123" -interval="5m"

# Send local network information to webdav server
ipsync send webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json"
## Send local network information to webdav server and encrypt the file
ipsync send webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123"
## Send local network information to webdav server and encrypt the file every 5 min
ipsync send webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123" -interval="5m"


# Receive local network information from a file
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" file -filepath="./home.json"
## Receive local network information from a file and decrypt the file
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" file -filepath="./home.json" -encryption_key="admin123"
## Receive local network information from a file and decrypt the file every 5 min
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" -interval="5m" file -filepath="./home.json" -encryption_key="admin123"

# Receive local network information from s3 server
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" s3 -endpoint="https://s3.amazonaws.com" -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Receive local network information from minio s3 server
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json"
## Receive local network information from minio s3 server and decrypt the file
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123"
## Receive local network information from minio s3 server and decrypt the file every 5 min
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" -interval="5m" s3 -endpoint="http://192.168.1.185:9000" -path_style -access_key_id="admin" -secret_access_key="adminadmin" -bucket="storage" -object_path="home.json" -encryption_key="admin123"

# Receive local network information from webdav server
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json"
## Receive local network information from webdav server and decrypt the file
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123"
## Receive local network information from webdav server and decrypt the file every 5 min
ipsync receive -remote_interface="pppoe-wan" -wg_interface="wg0" -interval="5m" webdav -endpoint="http://192.168.1.2/" -filepath="/dav/home.json" -encryption_key="admin123"

# Decrypt a encrypted file
ipsync decrypt -filepath "./home.json" -encryption_key="admin123"
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
curl -Lo "/etc/systemd/system/ipsync.timer" "https://github.com/unix755/ipsync/raw/main/service/systemd/send.timer"
systemctl enable ipsync.timer && systemctl restart ipsync.timer && systemctl status ipsync.timer
```

### Alpine Linux(openrc)

```sh
curl -Lo "/etc/init.d/ipsync" "https://github.com/unix755/ipsync/raw/main/service/openrc/send_webdav"
chmod +x /etc/init.d/ipsync
rc-update add ipsync && rc-service ipsync restart && rc-service ipsync status
```

### FreeBSD(rc.d)

```sh
mkdir /usr/local/etc/rc.d/
curl -Lo "/usr/local/etc/rc.d/ipsync" "https://github.com/unix755/ipsync/raw/main/service/rc.d/send_webdav"
chmod +x /usr/local/etc/rc.d/ipsync
service ipsync enable && service ipsync restart && service ipsync status
```

### OpenWRT(init.d)

```sh
curl -Lo "/etc/init.d/ipsync" "https://github.com/unix755/ipsync/raw/main/service/init.d/send_webdav"
chmod +x /etc/init.d/ipsync
service ipsync enable && service ipsync restart && service ipsync status
```

## Compile

### How to compile if prebuilt binaries are not found

```sh
git clone https://github.com/unix755/ipsync.git
cd ipsync
export CGO_ENABLED=0
go build -v -trimpath -ldflags "-s -w"
```

### For mipsle router

```sh
git clone https://github.com/unix755/ipsync.git
cd ipsync
export GOOS=linux
export GOARCH=mipsle
export GOMIPS=softfloat
export CGO_ENABLED=0
go build -v -trimpath -ldflags "-s -w"
```

## License

- **GPL-3.0 License**
- See `LICENSE` for details
