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


## 创建ingress

```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: main-ingress-test
  annotations:
    kubernetes.io/ingress.class: "nginx"   
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - http:
      paths:       
      - path: /hello
        pathType: Prefix
        backend:
          service:
          #这里是service的信息
            name: main-ing-svc
            port:
              number: 8083
```

```
darren@darrendeMacBook-Pro k8s_code % kubectl get ingress -owide
NAME                CLASS    HOSTS   ADDRESS   PORTS   AGE
main-ingress-test   <none>   *                 80      14s
```


### 高可用

```
darren@darrendeMacBook-Pro resume % kubectl delete pod main-deployment-84649c9fb6-dj4b9
pod "main-deployment-84649c9fb6-dj4b9" deleted
darren@darrendeMacBook-Pro resume % kubectl get pod                                    
NAME                               READY   STATUS    RESTARTS         AGE
busybox                            1/1     Running   223 (119m ago)   56d
hppod                              1/1     Running   83 (99m ago)     56d
main-deployment-84649c9fb6-bxj57   1/1     Running   0                29m
main-deployment-84649c9fb6-d82sg   1/1     Running   0                29s
main-deployment-84649c9fb6-vkx7c   0/1     Running   0                2s
nginx                              1/1     Running   6 (2d8h ago)     56d
redis-v2-7897dd7df5-svg5z          1/1     Running   1 (2d8h ago)     11d
test-incluster-client-go           1/1     Running   5 (2d8h ago)     29d
darren@darrendeMacBook-Pro resume % curl http://127.0.0.1:30993/healthz
ok

```