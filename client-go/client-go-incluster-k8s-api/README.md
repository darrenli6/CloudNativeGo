
# 集群内访问pod数量

## 编译程序

`GOOS=linux GOARCH=amd64 go build  -o ./app .`

如果出现 `bash: ./app: cannot execute binary file: Exec format error`

## 编译dockerfile 

`docker build -t in-cluster .`


## 加入k8s 

删除pod

`kubectl delete pod test-incluster-client-go`

创建pod
`kubectl apply -f pod.yaml    `

控制台输出 

```
there has 18 pod 
Pod test-node-local-dns not found in default namespace
```