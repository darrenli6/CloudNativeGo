

## 部署etcd
```
create ns etcd

kubectl apply -f sts.yaml -n etcd 

kubectl apply -f svc.yaml -n etcd 

```


kubectl logs etcd-0 -n etcd 
