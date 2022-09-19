# 13. Pod健康检查

## 13.2 使用探针

### Lab41. livenessProbe-exec
```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
  name: liveness-exec
spec:
  containers:
  - name: liveness
    args:
    - /bin/sh
    - -c
    - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
    image: busybox
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5
EOF

```

```bash
$ kubectl get pod liveness-exec
NAME            READY   STATUS    RESTARTS     AGE
liveness-exec   1/1     Running  `1`(5s ago)   80s

$ kubectl describe pod liveness-exec

$ kubectl describe pod liveness-exec | grep -i liveness:.*exec
    Liveness:       exec [cat /tmp/healthy] delay=5s timeout=1s period=5s #success=1 #failure=3
```

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
# name: liveness-exec
  name: liveness-exec3
spec:
  containers:
  - name: liveness
    args:
    - /bin/sh
    - -c
    - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
    image: busybox
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5
      # 增加 1 行
      timeoutSeconds: 3
EOF

```

```bash
$ kubectl describe pod liveness-exec3 | grep -i liveness:.*exec
    Liveness:       exec [cat /tmp/healthy] delay=5s timeout=`3s` period=5s #success=1 #failure=3
```

### Lab42. livenessProbe-http

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
  name: liveness-http
spec:
  containers:
  - name: liveness
    image: mirrorgooglecontainers/liveness
    args:
    - /server
    livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
        httpHeaders:
        - name: X-Custom-Header
          value: Awesome
      initialDelaySeconds: 3
      periodSeconds: 3
EOF

```

```bash
$ kubectl get pod liveness-http
NAME            READY   STATUS    RESTARTS   AGE
liveness-http   1/1     Running   0          41s

$ kubectl describe pod liveness-http | grep -i liveness:.*http
    Liveness:      `http-get` http://:8080/healthz delay=3s timeout=1s period=3s #success=1 #failure=3
```

### Lab43. livenessProbe-tcp

```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: ubuntu
  labels:
    app: ubuntu
spec:
  containers:
  - name: ubuntu
    image: ubuntu
    args:
    - /bin/sh
    - -c
    - apt-get update && apt-get -y install openbsd-inetd telnetd && /etc/init.d/openbsd-inetd start; sleep 30000
    livenessProbe:
      tcpSocket:
        port: 23
      initialDelaySeconds: 60
      periodSeconds: 20
EOF

```

```bash
$ kubectl get pod ubuntu
NAME     READY   STATUS    RESTARTS   AGE
ubuntu   1/1     Running   0          2m32s

$ kubectl describe pod ubuntu | grep -i live
    Liveness:      `tcp-socket` :23 delay=60s timeout=1s period=20s #success=1 #failure=3
```

## 13.3 使用就绪探针

### Lab44. readinessProbe-exec

```yaml
kubectl apply -f- <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpd-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: httpd
  template:
    metadata:
      labels:
        app: httpd
    spec:
      containers:
      - name: httpd
        image: httpd
        ports:
        - containerPort: 80
        readinessProbe:
          exec:
            command:
            - cat
            - /usr/local/apache2/htdocs/index.html
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: httpd-svc
spec:
  selector:
    app: httpd
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 80
EOF

```

```bash
$ kubectl get endpoints httpd-svc
NAME        ENDPOINTS                                            AGE
httpd-svc   172.16.126.57:80,172.16.126.58:80,172.16.126.59:80   97s

$ kubectl get pod -l app=httpd
NAME                              READY   STATUS    RESTARTS   AGE
httpd-deployment-c459d5dd-6pxl5   1/1     Running   0          84s
httpd-deployment-c459d5dd-79nvd   1/1     Running   0          84s
httpd-deployment-c459d5dd-vmlh2   1/1     Running   0          84s
```

```bash
$ kubectl exec -it httpd-deployment-c459d5dd-6pxl5 -- /bin/sh
# rm /usr/local/apache2/htdocs/index.html
# exit
```

```bash
$ kubectl get endpoints httpd-svc
NAME        ENDPOINTS                           AGE
httpd-svc   172.16.126.57:80,172.16.126.58:80   5m12s

$ kubectl get pods -l app=httpd
NAME                              READY   STATUS    RESTARTS   AGE
httpd-deployment-c459d5dd-6pxl5   0/1     Running   0          5m42s
httpd-deployment-c459d5dd-79nvd   1/1     Running   0          5m42s
httpd-deployment-c459d5dd-vmlh2   1/1     Running   0          5m42s

$ kubectl describe pod httpd-deployment-c459d5dd-6pxl5
...输出省略...
  Warning  Unhealthy  58s (x22 over 2m33s)  kubelet            Readiness probe failed: cat: /usr/local/apache2/htdocs/index.html: No such file or directory
```


