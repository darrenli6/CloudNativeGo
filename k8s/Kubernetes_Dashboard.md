# 20. Kubernetes Dashboard

## 20.1 Dashboard介绍

### Lab60. [部署 Dashboard UI](https://kubernetes.io/zh/docs/tasks/access-application-cluster/web-ui-dashboard/#%E9%83%A8%E7%BD%B2-dashboard-ui)

<div style="background: #dbfaf4; padding: 12px; line-height: 24px; margin-bottom: 24px;">
<dt style="background: #1abc9c; padding: 6px 12px; font-weight: bold; display: block; color: #fff; margin: -12px; margin-bottom: -12px; margin-bottom: 12px;" >Hint - 提示</dt>
  <li> 国内默认无法访问 https://raw.githubusercontent.com/kubernetes/dashboard/v2.5.0/aio/deploy/recommended.yaml
</div>


> 1. 安装


```bash
*$ kubectl apply -f https://vmcc.xyz:8443/k8s/dashboard/v2.6.1/aio/deploy/recommended.yaml
:<<EOF
namespace/kubernetes-dashboard created
serviceaccount/kubernetes-dashboard created
service/kubernetes-dashboard created
secret/kubernetes-dashboard-certs created
secret/kubernetes-dashboard-csrf created
secret/kubernetes-dashboard-key-holder created
configmap/kubernetes-dashboard-settings created
role.rbac.authorization.k8s.io/kubernetes-dashboard created
clusterrole.rbac.authorization.k8s.io/kubernetes-dashboard created
rolebinding.rbac.authorization.k8s.io/kubernetes-dashboard created
clusterrolebinding.rbac.authorization.k8s.io/kubernetes-dashboard created
deployment.apps/kubernetes-dashboard created
service/dashboard-metrics-scraper created
deployment.apps/dashboard-metrics-scraper created
EOF
```

```bash
$ kubectl -n kubernetes-dashboard get pods -owide
NAME                                         READY   STATUS    RESTARTS   AGE    IP              NODE          NOMINATED NODE   READINESS GATES
dashboard-metrics-scraper-799d786dbf-r7dp9   1/1     Running   0          115s   172.16.194.65   k8s-worker2   <none>           <none>
kubernetes-`dashboard`-546cbc58cd-2lvxt      1/1     `Running`   0          115s   172.16.126.1   k8s-worker2  <none>           <none>
```

> 2. 通过物理机使用==chrome==访问==thisisunsafe==

```bash
*$ kubectl -n kubernetes-dashboard \
  patch svc kubernetes-dashboard -p '{"spec":{"type":"NodePort"}}' 

$ kubectl -n kubernetes-dashboard get svc
NAME                        TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)         AGE
dashboard-metrics-scraper   ClusterIP   10.103.222.187   <none>        8000/TCP        6m2s
kubernetes-dashboard       `NodePort`   10.97.233.31     <none>        443:`32433`/TCP   6m3s
```

> 3. 查看服务帐号kubernetes-dashboard的登陆 token

```bash
针对最新版本，默认未创建token
$ kubectl -n kubernetes-dashboard create token kubernetes-dashboard
:<<EOF
`eyJhbGciOiJSUzI1NiIsImtpZCI6IkN6Zlo5dElKSFAyOXBkM2x0MWxxWU5yZkl1M2JzN1drZ2R1c1F3ZGVfdncifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNjYyMjY3ODE2LCJpYXQiOjE2NjIyNjQyMTYsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsInVpZCI6ImQ0ZjcwN2Q2LTU5YWYtNGZkMC05ZmNkLThjYTQ0MTBlOTcxNCJ9fSwibmJmIjoxNjYyMjY0MjE2LCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZXJuZXRlcy1kYXNoYm9hcmQ6a3ViZXJuZXRlcy1kYXNoYm9hcmQifQ.TV6KHUeqRtLMYRuZ_MNd8cRV6emuXyXNVQgAq7tm2U77VHLJojKGEi5_Ii6Zzg9aTpXjOifY9UwAimrkMD9qZFKXB2dF_PjZAYqMY8lTsM_Uc2NZnUM-BJXNLJzda4I0z4DIhMAg-q3rJMsy-6COJCGTMA0awgJI5hVJE91wgfzD96_rRfRINi0h00O6IMfDxHg2Bh1EfP2Q92NPFVmrnEXViwFCvih4zv1dr9gZmL-9P0vxgsV3g7zhPir7zleTSo7MDsUYUm9SxGSYEDcaoLQnu2ikMjG5U-n5HHp1KHxNdwOXtQWpMy2Ou-h2_usptvDNzB-dAGE2gQACjtECpQ`
EOF
```

> 4. 提权

```bash
$ kubectl -n kubernetes-dashboard \
  describe clusterrole cluster-admin
...输出省略...
PolicyRule:
  Resources  Non-Resource URLs  Resource Names  Verbs
  ---------  -----------------  --------------  -----
  *.*        []                 []              [*]
             [*]                []              [*]

$ kubectl -n kubernetes-dashboard \
describe clusterrole kubernetes-dashboard
...输出省略...
PolicyRule:
  Resources             Non-Resource URLs  Resource Names  Verbs
  ---------             -----------------  --------------  -----
  nodes.metrics.k8s.io  []                 []              [get list watch]
  pods.metrics.k8s.io   []                 []              [get list watch]

*$ kubectl create clusterrolebinding \
      kubernetes-dashboard-cluster-admin \
  --clusterrole=cluster-admin \
  --serviceaccount=kubernetes-dashboard:kubernetes-dashboard

$ kubectl get clusterrolebindings -owide | grep dash
kubernetes-dashboard-cluster-admin  ClusterRole/cluster-admin  6m34s  kubernetes-dashboard/kubernetes-dashboard
kubernetes-dashboard  ClusterRole/kubernetes-dashboard  88m  kubernetes-dashboard/kubernetes-dashboard
```

> 5. 登陆



### Lab61. kubeconfig方式认证

```bash
$ kubectl -n kube-system create sa dashboard-admin

$ kubectl create clusterrolebinding dashboard-admin \
--clusterrole=cluster-admin \
--serviceaccount=kube-system:dashboard-admin

$ SN=$(kubectl -n kube-system get secret | awk '/dashboard-admin/ {print $1}')

$ DASH_TOKEN=$(kubectl -n kube-system get secret ${SN} -o jsonpath={.data.token} | base64 -d)
```

> 创建 dashboard-admin.kubeconfig

```bash
*$ kubectl config set-cluster ck8s \
  --certificate-authority=/etc/kubernetes/pki/ca.crt \
	--embed-certs=true \
	--server=https://192.168.147.128:6443 \
	--kubeconfig=dashboard-admin.kubeconfig

*$ kubectl config set-context dashboard-admin@ck8s \
	--cluster=ck8s \
	--user=dashboard-admin \
	--kubeconfig=dashboard-admin.kubeconfig

*$ kubectl config set-credentials dashboard-admin \
	--token=${DASH_TOKEN} \
	--kubeconfig=dashboard-admin.kubeconfig

*$ kubectl config use-context dashboard-admin@ck8s \
	--kubeconfig=dashboard-admin.kubeconfig
```

```bash
$ vimdiff ~/.kube/config dashboard-admin.kubeconfig
```

物理机

```bash
$ scp kiosk@k8s-master:~/.kube/config .
$ scp kiosk@k8s-master:~/dashboard-admin.kubeconfig .
```



## 20.2 Dashboard功能

## 20.3 Dashboard部署应用