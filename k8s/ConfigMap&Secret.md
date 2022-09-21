
# 16. ConfigMap与Secret

## 16.1 ConfigMap介绍

### Lab51. 创建ConfigMap

```bash
sudo mkdir -p /runfile/configmap

sudo tee /runfile/configmap/game.properties <<EOF
lives=3
enemies.cheat=true
enemies.cheat.level=noGoodRotten
EOF

sudo tee /runfile/configmap/ui.properties <<EOF
color.good=purple
how.nice.to.look=fairlyNice
EOF
```

> - Folder
>   --from-file

```bash
$ kubectl create configmap game-config --from-file=/runfile/configmap

$ kubectl get configmap
NAME               DATA   AGE
`game-config`      2      5s
kube-root-ca.crt   1      31d

$ kubectl describe configmaps game-config 
...输出省略...
Data
====
game.properties:
----
lives=3
enemies.cheat=true
enemies.cheat.level=noGoodRotten

ui.properties:
----
color.good=purple
how.nice.to.look=fairlyNice
...输出省略...
```

> - File: ini
>   --from-file

```bash
$ kubectl create configmap game-config-3 \
  --from-file=/runfile/configmap/game.properties \
  --from-file=/runfile/configmap/ui.properties

$ kubectl get configmap
NAME               DATA   AGE
game-config        2      112s
`game-config-3`    2      12s
kube-root-ca.crt   1      31d

$ kubectl describe configmaps game-config-3
...输出省略...
Data
====
game.properties:
----
lives=3
enemies.cheat=true
enemies.cheat.level=noGoodRotten

ui.properties:
----
color.good=purple
how.nice.to.look=fairlyNice
...输出省略...
```

> --from-literal

```bash
$ kubectl create configmap special-config \
--from-literal=special.how=very \
--from-literal=special.type=charm

$ kubectl get configmaps
NAME               DATA   AGE
game-config        2      2m54s
game-config-3      2      74s
kube-root-ca.crt   1      31d
`special-config`   2      8s

$ kubectl describe configmaps special-config
...输出省略...
Data
====
special.how:
----
very
special.type:
----
charm
...输出省略...
```

> - File: yaml
>   --from-file

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: specialconfig-2
data:
  key1: value1
  pro.property: |
    key2: value2
    key3: value3
EOF

```

```bash
$ kubectl get configmaps
NAME               DATA   AGE
game-config        2      4m
game-config-3      2      2m20s
kube-root-ca.crt   1      31d
special-config     2      74s
`specialconfig-2`  2      15s

$ kubectl describe configmaps specialconfig-2
...输出省略...
Data
====
pro.property:
----
key2: value2
key3: value3

key1:
----
value1
...输出省略...
```

### Lab52. 使用ConfigMap

> - envFrom

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-pod
spec:
  containers:
    - name: cm-container
      image: busybox
      args: [  "/bin/sh", "-c", "sleep 3000" ]
      env:
        - name: special-env
          valueFrom:
            configMapKeyRef:
              name: specialconfig-2
              key: key1
      envFrom:
        - configMapRef:
            name: specialconfig-2
  restartPolicy: Never
EOF

```

```bash
$ kubectl get pods
NAME          READY   STATUS    RESTARTS   AGE
cm-test-pod   1/1     Running   0          26s

$ kubectl exec -it cm-test-pod  -- /bin/sh
/ # env
...输出省略...
key1=value1
pro.property=key2: value2
key3: value3
...输出省略...
/ # exit
```

> - volume

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: cmpod2
spec:
  containers:
  - name: cmpod2
    image: busybox
    args: [ "/bin/sh", "-c", "sleep 3000" ]
    volumeMounts:
    - name: db
      mountPath: "/etc/db"
      readOnly: true
  volumes:
  - name: db
    configMap:
      name: specialconfig-2
EOF

```

```bash
$ kubectl get pods
NAME          READY   STATUS    RESTARTS   AGE
cm-test-pod   1/1     Running   0          6m2s
cmpod2        1/1     Running   0          28s

$ kubectl exec -it cmpod2 -- /bin/sh
/ # ls /etc/db
key1          pro.property
/ # cat /etc/db/key1
value1/ # cat /etc/db/pro.property
key2: value2
key3: value3
# exit
```

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: specialconfig-2
data:
  key1: value-new
  pro.property: |
    key2: value2
    key3: value3
EOF

```

```bash
$ kubectl exec -it cmpod2 -- /bin/sh
/ # cat /etc/db/key1
value-new
# exit
```



## 16.2 Secret介绍

### Lab52. 创建Secret

> - file

```bash
$ echo -n "admin" > username.txt
$ echo -n "mima" > password.txt
```

```bash
$ kubectl create secret generic db-user-pass \
--from-file=./username.txt \
--from-file=./password.txt

$ kubectl get secret
NAME                  TYPE                                  DATA   AGE
`db-user-pass`        Opaque                                2      6s
default-token-7zfdn   kubernetes.io/service-account-token   3      31d

$ kubectl describe secrets db-user-pass
...输出省略...
Data
====
password.txt:  4 bytes
username.txt:  5 bytes
```

> - file: yaml

```bash
$ echo -n "admin" | base64
YWRtaW4=

$ echo -n "mima" | base64
bWltYQ==
```

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: bWltYQ==
EOF

```

```bash
$ kubectl get secrets mysecret
NAME       TYPE     DATA   AGE
mysecret   Opaque   2      13s

$ kubectl describe secrets mysecret
...输出省略...
Data
====
password:  4 bytes
username:  5 bytes
```

### Lab53. 使用Secret

> 挂载全部值

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: spod
spec:
  containers:
  - image: busybox
    name: spod
    args: ["/bin/sh","-c","sleep 3000"]
    volumeMounts:
    - name: secrets
      mountPath: "/etc/secret"
      readOnly: true
  volumes:
  - name: secrets
    secret:
      secretName: mysecret
EOF

```

```bash
$ kubectl get pods spod
NAME   READY   STATUS    RESTARTS   AGE
spod   1/1     Running   0          10s

$ kubectl exec -it spod -- /bin/sh
/ # ls /etc/secret
password  username
/ # cat /etc/secret/username
admin/ # cat /etc/secret/password
mima/ # exit
```

> 挂载指定值

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: spod2
spec:
  containers:
  - image: busybox
    name: spod2
    args: ["/bin/sh","-c","sleep 3000"]
    volumeMounts:
    - name: secrets
      mountPath: "/etc/secret"
      readOnly: true
  volumes:
  - name: secrets
    secret:
      secretName: mysecret
      items:
      - key: password
        path: my-group/my-passwd
EOF

```

```bash
$ kubectl get pod spod2
NAME    READY   STATUS    RESTARTS   AGE
spod2   1/1     Running   0          58s

$ kubectl exec -it spod2 -- /bin/sh
/ # ls /etc/secret
my-group
/ # ls /etc/secret/my-group/
my-passwd
/ # cat /etc/secret/my-group/my-passwd
mima/ # exit
```



