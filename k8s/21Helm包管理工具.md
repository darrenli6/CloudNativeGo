
# 21. Helm包管理工具

## 21.1 Helm简介

## 21.2 使用Helm

### Lab62. 安装、使用Helm

> [install helm](https://helm.sh/docs/intro/install/)

```bash
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
sudo apt-get install apt-transport-https --yes
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm
```

> 使用Helm

```bash
$ helm version

$ helm search hub wordpress
```

```bash
$ helm repo add alirepo https://apphub.aliyuncs.com
$ helm repo update
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "alirepo" chart repository
Update Complete. ⎈Happy Helming!⎈

$ helm repo list
NAME   	URL
alirepo https://apphub.aliyuncs.com

$ helm search repo mysql
NAME                             	CHART VERSION	APP VERSION	DESCRIPTION
alirepo/mysql                    	6.8.0        	8.0.19     	Chart to create a Highly available MySQL cluster
...输出省略...
```

```bash
$ helm install alirepo/mysql --generate-name
```

```bash
$ helm list
NAME            	NAMESPACE	REVISIONUPDATED                                	STATUS  CHART      	APP VERSION
mysql-1653395570	default  	1       2022-05-24 12:32:51.008311838 +0000 UTC	deployedmysql-6.8.0	8.0.19
```

```bash
$ helm uninstall mysql-1653395570
```

## 21.3 Chart简介

### Lab63. Chart

```bash
$ ls ~/.cache/helm/repository

$ tar -xf ~/.cache/helm/repository/mysql-6.8.0.tgz

$ sudo apt install tree

$ tree mysql
mysql
├── Chart.yaml
├── ci
│   └── values-production.yaml
├── files
│   └── docker-entrypoint-initdb.d
│       └── README.md
├── README.md
├── templates
│   ├── _helpers.tpl
│   ├── initialization-configmap.yaml
│   ├── master-configmap.yaml
│   ├── master-statefulset.yaml
│   ├── master-svc.yaml
│   ├── NOTES.txt
│   ├── secrets.yaml
│   ├── servicemonitor.yaml
│   ├── slave-configmap.yaml
│   ├── slave-statefulset.yaml
│   └── slave-svc.yaml
├── values-production.yaml
└── values.yaml

4 directories, 17 files
```

## 21.4 Chart模板的使用

### Lab64. Helm部署mysql

> 1. 默认需要使用pv

```bash
$ helm inspect values alirepo/mysql | less
...
  persistence:
    enabled: true
    mountPath: /bitnami/mysql
    annotations: {}
    accessModes:
      - `ReadWriteOnce`
    size: `8Gi`
...
```

> 2. 创建pv

**[kiosk@192.168.147.128]**

```bash
sudo mkdir -m 777 /nfs_{a..b}
sudo tee /etc/exports <<EOF
/nfs_a *(rw,no_root_squash)
/nfs_b *(rw,no_root_squash)
EOF
sudo apt -y install nfs-kernel-server && \
sudo systemctl enable nfs-server && \
sudo systemctl restart nfs-server && \
showmount -e

```

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pva
spec:
  capacity:
    storage: 8Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  nfs:
    path: /nfs_a
    server: 192.168.147.128
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pvb
spec:
  capacity:
    storage: 8Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  nfs:
    path: /nfs_b
    server: 192.168.147.128
EOF

```

```bash
$ kubectl get pv
NAME   CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
pva    8Gi        RWO            Recycle          Available                                   6s
pvb    8Gi        RWO            Recycle          Available
```

> 3. 安装mysql，并设置密码为 mima

```bash
$ helm install test alirepo/mysql --set root.password=mima

$ kubectl get pod -w
NAME                  READY   STATUS    RESTARTS   AGE
test-mysql-master-0  `1/1`    Running   0          2m35s
test-mysql-slave-0   `1/1`    Running   0          4m4s
<Ctrl-C>

$ helm list
NAME	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART      	APP VERSION
`test`	default  	1       	YYYY-MM-DD 13:49:58.731875229 +0000 UTC	deployed	mysql-6.8.0	8.0.19

$ helm status test
...输出省略...
Tip:

  Watch the deployment status using the command: kubectl get pods -w --namespace default

Services:

  echo Master: test-mysql.default.svc.cluster.local:3306
  echo Slave:  test-mysql-slave.default.svc.cluster.local:3306

Administrator credentials:

  echo Username: root
  echo Password : $(kubectl get secret --namespace default test-mysql -o jsonpath="{.data.mysql-root-password}" | base64 --decode)

To connect to your database:

  1. Run a pod that you can use as a client:

      kubectl run test-mysql-client --rm --tty -i --restart='Never' --image  docker.io/bitnami/mysql:8.0.19-debian-10-r0 --namespace default --command -- bash

  2. To connect to master service (read/write):

      mysql -h test-mysql.default.svc.cluster.local -uroot -p my_database

  3. To connect to slave service (read-only):

      mysql -h test-mysql-slave.default.svc.cluster.local -uroot -p my_database
...
```

> 4. 测试

```bash
$ kubectl get secret --namespace default test-mysql -o jsonpath="{.data.mysql-root-password}" | base64 --decode
`mima`
```

```bash
$ kubectl run test-mysql-client --rm --tty -i --restart='Never' --image  docker.io/bitnami/mysql:8.0.19-debian-10-r0 --namespace default --command -- bash
```

```bash
I have no name!@test-mysql-client:/$ mysql -h test-mysql.default.svc.cluster.local -uroot -p my_database
Enter password: `mima`
mysql> exit
```

```bash
I have no name!@test-mysql-client:/$ mysql -h test-mysql-slave.default.svc.cluster.local -uroot -p my_database
Enter password: `mima`
mysql> exit
```

```bash
I have no name!@test-mysql-client:/$ exit
```



