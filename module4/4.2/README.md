
课后练习 4.2
启动一个 Envoy Deployment。
要求 Envoy 的启动配置从外部的配置文件 Mount 进 Pod。
进入 Pod 查看 Envoy 进程和配置。
更改配置的监听端口并测试访问入口的变化。
通过非级联删除的方法逐个删除对象。


```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: envoy
  name: envoy
spec:
  replicas: 1
  selector:
    matchLabels:
      run: envoy
  template:
    metadata:
      labels:
        run: envoy
    spec:
      containers:
      - image: envoyproxy/envoy-dev
        name: envoy
        volumeMounts:
        - name: envoy-config
          mountPath: "/etc/envoy"
          readOnly: true
      volumes:
      - name: envoy-config
        configMap:
          name: envoy-config
```



# 启动 Envoy Deployment

```
root@k8s-node1:~# kubectl create configmap envoy-config --from-file=envoy.properties
configmap/envoy-config created
root@k8s-node1:~# kubect get configmap
kubect: command not found
root@k8s-node1:~# kubectl get configmap
NAME               DATA   AGE
envoy-config       1      12s
kube-root-ca.crt   1      4d11h
root@k8s-node1:~# ls
envoy.properties  overlay

# 外部mount进去
root@k8s-node1:~# cat envoy.properties 


hostPath.path=/k8s/etc/
hostPath.type=DirectoryOrCreate


root@k8s-master:~/yaml# kubectl apply -f envoy.yaml 
deployment.apps/envoy created
root@k8s-master:~/yaml# kubectl get ns 
NAME               STATUS   AGE
calico-apiserver   Active   4d11h
calico-system      Active   4d11h
default            Active   4d11h
kube-node-lease    Active   4d11h
kube-public        Active   4d11h
kube-system        Active   4d11h
tigera-operator    Active   4d11h
root@k8s-master:~/yaml# kubectl get deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
envoy              0/1     1            0           35s
nginx-deployment   1/1     1            1           3d12h
root@k8s-master:~/yaml# kubectl get deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
envoy              0/1     1            0           44s
nginx-deployment   1/1     1            1           3d12h

```

等待启动





# 进入 Pod 查看 Envoy 进程和配置。

```
root@k8s-master:~#   kubectl get pod
NAME                                READY   STATUS    RESTARTS      AGE
envoy-67b686dfb7-wrtdt              1/1     Running   2 (18m ago)   139m
nginx-deployment-6799fc88d8-h8tkv   1/1     Running   0             3d15h
root@k8s-master:~# kubectl exec -it envoy-67b686dfb7-wrtdt   sh


# pwd
/etc/envoy
# ls
envoy.properties
# cat envoy.properties


hostPath.path=/k8s/etc/
hostPath.type=DirectoryOrCreate


# ps -ef
UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 03:34 ?        00:00:00 /bin/sh -ce sleep 3600
root           6       1  0 03:34 ?        00:00:00 sleep 3600
root           7       0  0 03:53 pts/0    00:00:00 sh
root          21       7  0 03:56 pts/0    00:00:00 ps -ef


```


# 通过非级联删除的方法逐个删除对象。




```
root@k8s-master:~# kubectl get pod 
NAME                                READY   STATUS    RESTARTS      AGE
envoy-67b686dfb7-wrtdt              1/1     Running   2 (24m ago)   145m
nginx-deployment-6799fc88d8-h8tkv   1/1     Running   0             3d15h
root@k8s-master:~# kubectl get deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
envoy              1/1     1            1           145m
nginx-deployment   1/1     1            1           3d15h
 
root@k8s-master:~# kubectl delete deployment envoy 
deployment.apps "envoy" deleted
root@k8s-master:~# kubectl get pod
NAME                                READY   STATUS        RESTARTS      AGE
envoy-67b686dfb7-wrtdt              1/1     Terminating   2 (24m ago)   145m
nginx-deployment-6799fc88d8-h8tkv   1/1     Running       0             3d15h
root@k8s-master:~# kubectl get deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   1/1     1            1           3d15h
```

