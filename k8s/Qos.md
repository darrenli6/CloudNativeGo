# 18. Kubernetes服务质量

## 18.1 QoS原理与使用

### Lab58. Qos

> 节点资源

```bash
$ kubectl get node k8s-worker1 -o yaml | grep -A 13 allo
 `allocatable`:
    cpu: "2"
    ephemeral-storage: "36496716535"
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 3892256Ki
    pods: "110"
 `capacity`:
    cpu: "2"
    ephemeral-storage: 39601472Ki
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 3994656Ki
    pods: "110"
```

> Guaranteed

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: qos-demo
spec:
  containers:
  - name: qos-demo-ctr
    image: nginx
    resources:
      limits:
        memory: "200Mi"
        cpu: "700m"
      requests:
        memory: "200Mi"
        cpu: "700m"
EOF

```

```bash
$ kubectl get pod qos-demo -o yaml | grep qosC
  qosClass: `Guaranteed`
```

> Burstable

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: qos-demo-2
spec:
  containers:
  - name: qos-demo-2-ctr
    image: nginx
    resources:
      limits:
        memory: "200Mi"
      requests:
        memory: "100Mi"
EOF

```

```bash
$ kubectl get pod qos-demo-2 -o yaml | grep qosC
  qosClass: `Burstable`
```

> BestEffort

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: qos-demo-3
spec:
  containers:
  - name: qos-demo-3-ctr
    image: nginx
EOF

```

```bash
$ kubectl get pod qos-demo-3 -o yaml | grep qosC
  qosClass: `BestEffort`
```

> 超出容器的内存限制

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: memory-demo
spec:
  containers:
  - name: memory-demo-ctr
    image: polinux/stress
    resources:
      limits:
        memory: "100Mi"
      requests:
        memory: "50Mi"
    command: ["stress"]
    args: ["--vm", "1", "--vm-bytes", "150M", "--vm-hang", "1"]
EOF

```

```bash
$ kubectl get pod memory-demo
NAME          READY   STATUS      RESTARTS   AGE
memory-demo   0/1    `OOMKilled`  0          20s
```

```bash
$ kubectl describe pod memory-demo | grep -A 2 Last
    Last State:     Terminated
      Reason:       OOMKilled
      Exit Code:    1
```

> 超出节点可用资源

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: memory-demo2
spec:
  containers:
  - name: memory-demo2-ctr
    image: polinux/stress
    resources:
      limits:
        memory: "1000Gi"
      requests:
        memory: "1000Gi"
    command: ["stress"]
    args: ["--vm", "1", "--vm-bytes", "150M", "--vm-hang", "1"]
EOF

```

```bash
$ kubectl get pod memory-demo2
NAME           READY   STATUS    RESTARTS   AGE
memory-demo2   0/1    `Pending`  0          22s
```

```bash
$ kubectl describe pod memory-demo2 | tail -n 1
  Warning  FailedScheduling  10s (x2 over 72s)  default-scheduler  0/3 nodes are available: 1 node(s) had taint {node-role.kubernetes.io/master: }, that the pod didn't tolerate, 2 Insufficient memory.
```


