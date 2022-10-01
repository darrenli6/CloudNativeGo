# 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化

## 构建文件
```
GOOS=linux GOARCH=amd64 go build .
```

## httpserver 容器化

```
docker build -t httpserver . 
```

## 启动容器 

```
docker run -p 8081:8081 --name httpserver -itd httpserver
```



# 将镜像推送至 docker 官方镜像仓库

- 修改tag

```
 docker tag httpserver:latest darren94me/httpserver:v1
```

- push image

```
darren@darrendeMBP module3 %  docker push darren94me/httpserver:v1
The push refers to repository [docker.io/darren94me/httpserver]
5f70bf18a086: Pushed 
f468fbcda9cb: Layer already exists 
03c3d223c3ca: Layer already exists 
98fcde1c6e0f: Layer already exists 
v1: digest: sha256:3471ee85537da6732cd786b40dee9deb6eaf74a821494a466782293868cdcea1 size: 1359  
```



# 通过 docker 命令本地启动 httpserver

```
docker run -p 8081:8081 --name httpserver -itd httpserver
```

# 通过 nsenter 进入容器查看 IP 配置


- 进入容器，查看ip 
```
docker exec -it httpserver  bash

root@b4654315ca63:/# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: tunl0@NONE: <NOARP> mtu 1480 qdisc noop state DOWN group default qlen 1000
    link/ipip 0.0.0.0 brd 0.0.0.0
3: ip6tnl0@NONE: <NOARP> mtu 1452 qdisc noop state DOWN group default qlen 1000
    link/tunnel6 :: brd :: permaddr 8a50:7f7a:6215::
122: eth0@if123: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:07 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.7/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever


 172.17.0.7
        
```

- 到宿主机查找PID 

```
darren@darrendeMBP module3 %  docker inspect httpserver | grep Pid
            "Pid": 8752,
            "PidMode": "",
            "PidsLimit": null,


nsenter -t 8752 -n  ip a             
```

 
