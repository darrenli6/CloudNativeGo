
# 蓝绿部署

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-dp-v1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: demo
      version: v1
  template:
    metadata:
      labels:
        app: demo
        version: v1
    spec:
      containers:
      - name: demo
        image: nginx:1.14.2
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: demo-service
spec:
    selector:
        app: demo
        version: v1
    type: loadBalancer    
    ports:
    - port: 8080
        targetPort: 80
        protocol: TCP

 ```       