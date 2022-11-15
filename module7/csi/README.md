
# create a host folder
```
 mkdir -p /Users/darren/Documents/project/k8spath/csi
```
# create a file in a folder 

```
echo sh -c "echo 'hello from k8s storage'" > /Users/darren/Documents/project/k8spath/csi/index.html


kubectl create ns module7 

darren@darrendeMacBook-Pro csi % kubectl apply -f pv.yaml -n module7
persistentvolume/task-pv-volume created
darren@darrendeMacBook-Pro csi % kubectl apply -f pvc.yaml -n module7
persistentvolumeclaim/task-pv-claim created
darren@darrendeMacBook-Pro csi % kubectl apply -f pod.yaml -n module7
pod/task-pv-pod created
darren@darrendeMacBook-Pro csi % kubectl apply -f svc.yaml -n module7



```



