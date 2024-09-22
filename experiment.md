## Configure Minikube

```
minikube start
minikube config set memory 4000
minikube delete
minikube start
minikube addons enable metrics-server
minikube dashboard
```

visit [dashboard](http://127.0.0.1:51466/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#/pod?namespace=memallocdemo) to inspect memory and cpu usage.

## Deploy Two Memoryhoggers

```
kubectl apply -f experiment1.yaml
kubens memallocdemo
```

## Allocate 3GB to pod 1

```
kubectl port-forward svc/memallocdemo-1 8081:80
curl -X POST localhost:8081/allocate/3000
```

## Allocate 3GB to pod 2

```
kubectl port-forward svc/memallocdemo-2 8081:80
curl -X POST localhost:8081/allocate/2000
{"message":"Allocating 2000MB."}
```

This OOMKills pod 1! Why?

## Allocate 3GB to pod 1, then free it

```
kubectl port-forward svc/memallocdemo-1 8081:80
curl -X POST localhost:8081/allocate/3000
curl -X POST localhost:8081/free
```

## Allocate 3GB to pod 2

```
kubectl port-forward svc/memallocdemo-2 8081:80
curl -X POST localhost:8081/allocate/3000
{"message":"Allocating 3000MB."}
```

This works, confirming that the scheduler frees up the memory to the memory pool once an application stops using it.