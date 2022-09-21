# 17. StatefulSet管理与使用

## 17.1 使用StatefulSet

### Lab54. StatefulSet

```bash
sudo mkdir /etc/exports.d /nfs{1..3} && \
sudo tee /etc/exports.d/s.exports <<EOF
/nfs1 *(rw,no_root_squash)
/nfs2 *(rw,no_root_squash)
/nfs3 *(rw,no_root_squash)
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
  name: mypv1
spec:
  storageClassName: my-sc
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  nfs:
    path: /nfs1
    server: 192.168.147.128
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mypv2
spec:
  storageClassName: my-sc
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  nfs:
    path: /nfs2
    server: 192.168.147.128
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mypv3
spec:
  storageClassName: my-sc
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  nfs:
    path: /nfs3
    server: 192.168.147.128
EOF

```

```bash
$ kubectl get pv
NAME    CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
mypv1   1Gi        RWO            Recycle          Available           my-sc                   41s
mypv2   1Gi        RWO            Recycle          Available           my-sc                   41s
mypv3   1Gi        RWO            Recycle          Available           my-sc                   10s
```

```yaml
kubectl apply -f- <<EOF
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  selector:
    matchLabels:
      app: nginx
  serviceName: nginx
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: stor
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: stor
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: "my-sc"
      resources:
        requests:
          storage: 1Gi
EOF

```

```bash
$ watch -n 1 kubectl get pod
NAME    READY   STATUS    RESTARTS   AGE
`web-0` 1/1     Running   0          79s
`web-1` 1/1     Running   0          55s
`web-2` 1/1     Running   0          30s
<Ctrl-C>

$ kubectl get pvc
NAME         STATUS   VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
stor-web-0   Bound    mypv1    1Gi        RWO            my-sc          112s
stor-web-1   Bound    mypv2    1Gi        RWO            my-sc          88s
stor-web-2   Bound    mypv3    1Gi        RWO            my-sc          63s
```

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
  - port: 80
    name: web
  clusterIP: None
  selector:
    app: nginx
EOF

```

```bash
$ kubectl get svc nginx
NAME    TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
nginx   ClusterIP   None         <none>        80/TCP    15s
```

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: clientpod
spec:
  containers:
    - name: clientpod
      image: busybox:1.28.3
      args:
      - /bin/sh
      - -c
      - sleep 3h
EOF

```

```bash
$ kubectl exec -it clientpod -- /bin/sh
/ # nslookup nginx
...输出省略...
Name:      nginx
Address 1: 172.16.194.66 web-2.nginx.default.svc.cluster.local
Address 2: 172.16.126.1 web-0.nginx.default.svc.cluster.local
Address 3: 172.16.194.65 web-1.nginx.default.svc.cluster.local
/ # nslookup web-0.nginx
...输出省略...
Name:      web-0.nginx
Address 1: 172.16.126.1 web-0.nginx.default.svc.cluster.local

/# exit

```

### Lab55. StatefulSet的故障处理

```bash
$ sudo touch /nfs1/newfile

$ kubectl exec -it web-1 -- ls /usr/share/nginx/html
newfile

$ kubectl exec -it clientpod -- nslookup web-1.nginx
...输出省略...
Address 3: 172.16.194.`65` web-1.nginx.default.svc.cluster.local

$ kubectl delete pod web-1
pod "web-1" deleted

$ kubectl exec -it clientpod -- nslookup web-1.nginx
...输出省略...
Name:      web-1.nginx
Address 1: 172.16.194.`67` web-1.nginx.default.svc.cluster.local

$ k exec -it web-1 -- ls /usr/share/nginx/html
newfile
```

### Lab56. 扩缩容和升级

```bash
$ watch -n 1 kubectl get pod
```

新开个终端

```bash
$ kubectl scale statefulset web --replicas=1

$ kubectl scale statefulset web --replicas=3
```

### Lab57. Pod管理策略

```bash
kubectl delete statefulsets web
kubectl delete pvc stor-web-0 stor-web-1 stor-web-2

$ kubectl get svc

$ kubectl get pv
```

```yaml
kubectl apply -f- <<EOF
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  selector:
    matchLabels:
      app: nginx
  serviceName: nginx
  # 增加 1 行
  podManagementPolicy: "Parallel"
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: stor
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: stor
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: "my-sc"
      resources:
        requests:
          storage: 1Gi
EOF

```

```bash
$ watch -n 1 kubectl get pod
```

新开个终端

```bash
$ kubectl scale statefulset web --replicas=1

$ kubectl scale statefulset web --replicas=3
```


