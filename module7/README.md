
## kube-schedule 

kube-schedule 负责调度pod到集群内部节点上，它监听kube-apiserver ,查询还有分配pod的node
根据调度策略为这些pod分配节点
考虑诸多的因素
- 公平调度
- 资源高效利用
- Qos


## 调度器

分为两个阶段：predicate 和 priority
predicate: 过滤不符合条件的节点 
priority：优先级排序


