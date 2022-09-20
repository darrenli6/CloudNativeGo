# 14. Kubernetes网络

## 14.1 Kubernetes 网络模型

### Lab45. Node与Pod之间通信实验

```bash
$ kubectl run nginx --image=nginx --port 80

$ kubectl get pod -owide
NAME   READY   STATUS   RESTARTS   AGE     IP             NODE          NOMINATED NODE   READINESS GATES
nginx  1/1     Running  0          27s     172.16.126.3  `k8s-worker2`  <none>           <none>

$ curl 172.16.126.3
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
...输出省略...

$ ip route
...输出省略...
172.16.235.196 dev cali46502c0485c scope link
172.16.235.197 dev cali29fae187fa2 scope link
172.16.235.198 dev cali2bab3c78e13 scope link
...输出省略...
```

### Lab46. Pod与Pod之间通信实验

```bash
$ kubectl run busybox --image=busybox -- sleep 3h

$ kubectl get pod -o wide
NAME       READY   STATUS     RESTARTS     AGE    IP                     NODE
busybox    1/1         Running    0                 27s       172.16.194.71   k8s-worker1
...输出省略
```

```bash
$ kubectl exec -it busybox -- /bin/sh
/ # ping 172.16.126.3
PING 172.16.126.3 (172.16.126.3): 56 data bytes
64 bytes from 172.16.126.3: seq=0 ttl=62 time=0.611 ms
...输出省略...
/ # telnet 172.16.126.3 80
Connected to 172.16.126.3
get
...输出省略...
<hr><center>nginx/1.21.5</center>
</body>
</html>
Connection closed by foreign host
/ # exit
```

### Lab46. 集群外访问-NodePort

```bash
$ kubectl get pod --show-labels
NAME     READY   STATUS  RESTARTS   AGE. LABELS
busybox  1/1        Running  0                17m  run=busybox
nginx    1/1        Running  0                42m  run=nginx
...
```

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Service
metadata:
  name: nginx-access
spec:
  selector:
    run: nginx
  type: NodePort
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30080
EOF

```

```bash
$ kubectl get svc -o wide
NAME          TYPE       CLUSTER-IP     EXTERNAL-IP  PORT(S)           AGE  SELECTOR
kubernetes    ClusterIP  10.96.0.1      <none>       443/TCP           30d   <none>
nginx-access  NodePort   10.110.109.55  <none>       80:`30080`/TCP    68s   run=nginx

$ curl http://k8s-master:30080
...
<h1>Welcome to nginx!</h1>
...

$ ss -ntlp | grep 30080
LISTEN  0       4096              0.0.0.0:30080          0.0.0.0:*
```


