1. Download kind - to run kubernetes on docker - [https://kind.sigs.k8s.io/]
2. make sure that the ca crt and key extracted and available in the current directory
3. run in command line
```
        kubectl create secret generic rabbitmq-certificates --from-file=./ca.crt --from-file=./tls.crt --from-file=./tls.key
```
4. run in command line
```
        helm repo add bitnami https://charts.bitnami.com/bitnami
        helm install rabbitmq bitnami/rabbitmq --set auth.username=rabbitmq --set auth.password=rabbitmq --set auth.tls.enabled=true --set auth.tls.existingSecret=rabbitmq-certificates
```
5. wait until it will completed and you will see rabbitmq running in the cluster
6. expose rabbitmq service in the cluster with port-forwarding
```
        kubectl port-forward svc/rabbitmq 5671:5671
```
then you can connect to rabbitmq





