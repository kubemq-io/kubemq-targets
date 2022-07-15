helm repo add bitnami https://charts.bitnami.com/bitnami
helm install rabbitmq bitnami/rabbitmq --set auth.username=rabbitmq --set auth.password=rabbitmq --set service.type=LoadBalancer
helm delete rabbitmq
kubectl port-forward --namespace default svc/rabbitmq 5672:5672
kubectl port-forward --namespace default svc/rabbitmq 15672:15672

kubectl create secret generic rabbitmq-certificates --from-file=./ca.crt --from-file=./tls.crt --from-file=./tls.key
helm install rabbitmq bitnami/rabbitmq --set auth.username=rabbitmq --set auth.password=rabbitmq --set auth.tls.enabled=true --set auth.tls.existingSecret=rabbitmq-certificates
