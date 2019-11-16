
etcd数据安装
===

## 1 下载

[etcd下载](https://github.com/coreos/etcd/releases/)这里选用的是最新版本v3.4.3

---

## 2 单机版安装

将下载的etcd文件包进行解压，解压后将`etcd`、`etcdctl`二进制文件复制到`/usr/bin`目录,并给权限执行权限

```shell
chmod a+x etcd etcdctl
cp etcd etcdctl /usr/bin/
```

#### 2.1 设置服务文件etcd.service

在`/usr/lib/systemd/system/`目录下创建文件`etcd.service`，内容为：

```systemd
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
TimeoutStartSec=10
Restart=on-failure
WorkingDirectory=/var/lib/etcd/
EnvironmentFile=-/etc/etcd/etcd.conf
ExecStart=/usr/bin/etcd
#--name=${ETCD_NAME} \
#--data-dir=${ETCD_DATA_DIR} \
#--listen-client-urls=${ETCD_LISTEN_CLIENT_URLS} \
#--advertise-client-urls=${ETCD_ADVERTISE_CLIENT_URLS}
#--initial-cluster-state=${ETCD_INITIAL_CLUSTER_STATE} \
#--cert-file=${ETCD_CERT_FILE} \
#--key-file=${ETCD_KEY_FILE} \
#--peer-cert-file=${ETCD_PEER_CERT_FILE} \
#--peer-key-file=${ETCD_PEER_KEY_FILE} \
#--trusted-ca-file=${ETCD_TRUSTED_CA_FILE} \
#--client-cert-auth=${ETCD_CLIENT_CERT_AUTH} \
#--peer-client-cert-auth=${ETCD_PEER_CLIENT_CERT_AUTH} \
#--peer-trusted-ca-file=${ETCD_PEER_TRUSTED_CA_FILE}
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

其中`WorkingDirectory`为`etcd`数据库目录，==需要在etcd启动前创建==  
集群或使用证书取消相关`"#"`,注意`"\"`不要多出来或少掉,并修改`etcd.conf`配置

```shell
  mkdir -p /var/lib/etcd
```

#### 2.2 创建配置文件 `etcd.conf`

在`/etc/etcd/`目录下创建文件`etcd.conf`内容为:

```conf
# [member]
ETCD_NAME=ETCD Server
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
#ETCD_WAL_DIR=""
#ETCD_SNAPSHOT_COUNT="10000"
#ETCD_HEARTBEAT_INTERVAL="100"
#ETCD_ELECTION_TIMEOUT="1000"
#ETCD_LISTEN_PEER_URLS="http://localhost:2380"
ETCD_LISTEN_CLIENT_URLS="http://127.0.0.1:2379"
#ETCD_MAX_SNAPSHOTS="5"
#ETCD_MAX_WALS="5"
#ETCD_CORS=""
#
#[cluster]
#ETCD_INITIAL_ADVERTISE_PEER_URLS="http://localhost:2380"
# if you use different ETCD_NAME (e.g. test), set ETCD_INITIAL_CLUSTER value for this name, i.e. "test=http://..."
#ETCD_INITIAL_CLUSTER="default=http://localhost:2380"
#ETCD_INITIAL_CLUSTER_STATE="new"
#ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_ADVERTISE_CLIENT_URLS="http://127.0.0.1:2379"
#ETCD_DISCOVERY=""
#ETCD_DISCOVERY_SRV=""
#ETCD_DISCOVERY_FALLBACK="proxy"
#ETCD_DISCOVERY_PROXY=""
#
#[proxy]
#ETCD_PROXY="off"
#ETCD_PROXY_FAILURE_WAIT="5000"
#ETCD_PROXY_REFRESH_INTERVAL="30000"
#ETCD_PROXY_DIAL_TIMEOUT="1000"
#ETCD_PROXY_WRITE_TIMEOUT="5000"
#ETCD_PROXY_READ_TIMEOUT="0"
#
#[security]
#ETCD_CERT_FILE="/etc/ssl/etcd/etcd.pem"
#ETCD_KEY_FILE="/etc/ssl/etcd/etcd-key.pem"
#ETCD_CLIENT_CERT_AUTH="true"
#ETCD_TRUSTED_CA_FILE="/etc/ssl/etcd/ca/ca.pem"
#ETCD_PEER_CERT_FILE="/etc/ssl/etcd/etcd.pem"
#ETCD_PEER_KEY_FILE="/etc/ssl/etcd/etcd-key.pem"
#ETCD_PEER_CLIENT_CERT_AUTH="true"
#ETCD_PEER_TRUSTED_CA_FILE="/etc/ssl/etcd/ca/ca.pem"
#
#[logging]
#ETCD_DEBUG="false"
# examples for -log-package-levels etcdserver=WARNING,security=DEBUG
#ETCD_LOG_PACKAGE_LEVELS=""
```

#### 2.3 配置开机启动并运行

```shell
systemctl daemon-reload
systemctl enable etcd.service
systemctl start etcd.service
```

#### 2.4 测试是否安装成功

```shell
etcdctl cluster-health  # v3.3.x
etcdctl endpoint health # v3.4.x
```
