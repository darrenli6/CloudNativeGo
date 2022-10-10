[toc]

# Docker核心技术 一

## linux对Namespace的操作方法

- clone 

在创建新进程的系统调用时，可以通过flags参数指定需要新的namespace类型

- setns 
  
系统调用可以让调用进程加入已存在的namespace

- unshare
  
系统调用可以让调用进程加入新的namespace

## 隔离性 namespace 

![image](./imgs/namespace.png)


