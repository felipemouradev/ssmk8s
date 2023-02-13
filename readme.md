Role for execution SSM K8S Service

```
kubectl apply -f permissions.yaml
./configure.sh
export KUBECONFIG=./kubeconfig
go run main.go YOUR_SSM_PATH YOUR_APPLICATION 
```