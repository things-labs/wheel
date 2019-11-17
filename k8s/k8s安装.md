# 准备

- 硬件
  - 本地ubuntu18.04
  - 虚拟机vmware
- 软件
  - docker-machine
  - etcd
  - kubernetes

## 给节点安装docker

```shell
  docker-machine create --driver generic --generic-ip-address=192.168.1.10 k8s-node1
```

## 安装

### 1. 规划

| 主机名     | 地址         | 角色       | 组件                                                          |
| :--------- | :----------- | :--------- | :------------------------------------------------------------ |
| k8s-master | 192.168.1.12 | k8s-master | etcd、kube-apiserver、kube-controller-manager、kube-scheduler |
| k8s-node-1 | 192.168.1.10 | k8s-node   | kubelet、docker、kube_proxy                                   |
| k8s-node-2 | 192.168.1.11 | k8s-node   | kubelet、docker、kube_proxy                                   |

### 2. 软件下载

#### 2.1 Kubernetes二进制文件下载

[kubernetes下载](https://github.com/kubernetes/kubernetes/releases),从`CHANGELOG`页面下载二进制文件，如图所示为其Linux Server版本

#### 2.2 etcd数据库下载

[etcd下载](https://github.com/coreos/etcd/releases/)这里选用的是最新版本v3.4.3.

### 3. Master安装

#### 3.1 etcd安装

[etcd安装](/k8s/etcd安装.md)

#### 3.2 Master组件

- `kube-apiserver`
- `kube-controller-manager`
- `kube-scheduler`

##### 3.2.1 复制二进制文件到`/usr/bin`目录

将`kube-apiserver`,`kube-controller-manager`,`kube-scheduler`三个可执行文件复制到`/usr/bin`目录,并给执行权限

#### 3.2.2 `kube-apiserver`

在`/usr/lib/systemd/system/`目录下创建文件`kube-apiserver.service`内容为:

```systemd
[Unit]
Description=Kubernetes API Server
After=etcd.service
Wants=etcd.service

[Service]
EnvironmentFile=/etc/kubernetes/apiserver
ExecStart=/usr/bin/kube-apiserver  \
        $KUBE_ETCD_SERVERS \
        $KUBE_API_ADDRESS \
        $KUBE_API_PORT \
        $KUBE_SERVICE_ADDRESSES \
        $KUBE_ADMISSION_CONTROL \
        $KUBE_API_LOG \
        $KUBE_API_ARGS 
Restart=on-failure
Type=notify
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

其中EnvironmentFile为kube-apiserver的配置文件

#### 3.2.3 kube-aouserver配置文件

在`/etc/kubernetes/`目录下创建文件`apiserver`,内容为:

```conf
KUBE_API_ADDRESS="--insecure-bind-address=0.0.0.0"
KUBE_API_PORT="--insecure-port=8080"
KUBE_ETCD_SERVERS="--etcd-servers=http://192.168.1.12:2379"
KUBE_SERVICE_ADDRESSES="--service-cluster-ip-range=172.18.0.0/16"
KUBE_ADMISSION_CONTROL="--admission-control=NamespaceLifecycle,LimitRanger,SecurityContextDeny,ServiceAccount,ResourceQuota"
KUBE_API_LOG="--logtostderr=false --log-dir=/var/log/kubernets/apiserver --v=2"
KUBE_API_ARGS=" "
```

#### 3.2.4 kube-controller-manager

在`/usr/lib/systemd/system/`目录下创建`kube-controller-manager.service`，内容为：

```systemd
[Unit]
Description=Kubernetes Scheduler
After=kube-apiserver.service 
Requires=kube-apiserver.service

[Service]
EnvironmentFile=-/etc/kubernetes/controller-manager.conf
ExecStart=/usr/bin/kube-controller-manager \
        $KUBE_MASTER \
        $KUBE_CONTROLLER_MANAGER_ARGS
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

#### 3.2.5 kube-controller-manager配置文件

在`/etc/kubernetes/`目录下创建文件`controller-manager`,内容为

```conf
KUBE_MASTER="--master=http://192.168.1.12:8080"
KUBE_CONTROLLER_MANAGER_ARGS=" "
```

#### 3.2.6 kube-scheduler

在`/usr/lib/systemd/system/`目录下创建`kube-scheduler.service`，内容为：

```systemd
[Unit]
Description=Kubernetes Scheduler
After=kube-apiserver.service 
Requires=kube-apiserver.service

[Service]
User=root
EnvironmentFile=-/etc/kubernetes/scheduler.conf
ExecStart=/usr/bin/kube-scheduler \
        $KUBE_MASTER \
        $KUBE_SCHEDULER_ARGS
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

#### 3.2.7 kube-scheduler配置文件

在`/etc/kubernetes/`目录下创建文件`scheduler`,内容为

```conf
KUBE_MASTER="--master=http://192.168.1.12:8080"
KUBE_SCHEDULER_ARGS="--logtostderr=true --log-dir=/var/log/kubernetes/scheduler --v=2"
```

#### 3.2.7 各组件开机启动

```shell
systemctl daemon-reload
systemctl enable kube-apiserver.service
systemctl start kube-apiserver.service
systemctl enable kube-controller-manager.service
systemctl start kube-controller-manager.service
systemctl enable kube-scheduler.service
systemctl start kube-scheduler.service
```

#### 3.2.8 检验正确

```shell
kubectl get cs
systemctl status kube-apiserver kube-controller-manager kube-scheduler
```

### 4 Node安装

#### 4.1 组件列表

- `docker`
- `kube-proxy`
- `kubelet`
