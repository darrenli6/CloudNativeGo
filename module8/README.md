# 作业

现在你对 Kubernetes 的控制面板的工作机制是否有了深入的了解呢？
是否对如何构建一个优雅的云上应用有了深刻的认识，那么接下来用最近学过的知识把你之前编写的 http 以优雅的方式部署起来吧，你可能需要审视之前代码是否能满足优雅上云的需求。
作业要求：编写 Kubernetes 部署脚本将 httpserver 部署到 Kubernetes 集群，以下是你可以思考的维度。

优雅启动
优雅终止
资源需求和 QoS 保证
探活
日常运维需求，日志等级
配置和代码分离






# 完成过程

`cd code`

##  编译go程序

`make build`


## 制作镜像

`make image`

httpserver:v1

## 运行docker

`make rundocker`

杀死docker容器

`make killdocker`


## 编写deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: main-deployment
  labels:
    app: main
spec:
  replicas: 3
  selector:
    matchLabels:
      app: main
  template:
    metadata:
      labels:
        app: main
    spec:
      containers:
      - name: main
        image: httpserver:v1
        volumeMounts:
        - mountPath: /data/logs/
          name: log-volume
        ports:
        - containerPort: 8083
        readinessProbe:
        # 探活
          httpGet:
            path: /healthz
            port: 8083
          initialDelaySeconds: 10
          periodSeconds: 5
        env:
        #环境变量
          - name: GOPORT
            value: "8083"
          - name: GOAPPPATH
            value: "/data/logs/"  
      volumes:
        - name: log-volume
          hostPath:
            # 本机的地址 这里有日志
            path: /Users/darren/go/src/github.com/darrenli6/CloudNativeGo/module8/code/log/
            type: Directory
```

## 创建service 

```
apiVersion: v1
kind: Service
metadata:
  name: main-ing-svc
spec:
  type: NodePort
  ports:
  - port: 8083
  # 因为service通过port 80端口接受，打到targetPort端口 这个targetPort与deployment一致
    targetPort: 8083
    protocol: TCP
  selector:
      # 这里要筛选deployement 所以与上面deployment的标签一致
    app: main
```

得到信息
```
darren@darrendeMacBook-Pro k8s_code % kubectl get svc  
NAME               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
httpd-deployment   NodePort    10.99.10.182    <none>        80:30090/TCP     27d
httpd-svc          ClusterIP   10.103.88.200   <none>        8080/TCP         56d
kubernetes         ClusterIP   10.96.0.1       <none>        443/TCP          224d
main-ing-svc       NodePort    10.107.193.1    <none>        8083:30993/TCP   11m
```

宿主机访问 ：http://127.0.0.1:30993/healthz
