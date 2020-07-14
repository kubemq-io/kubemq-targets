git clone https://github.com/openfaas/faas-netes
kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml
kubectl -n openfaas create secret generic basic-auth --from-literal=basic-auth-user=admin --from-literal=basic-auth-password=password
kubectl apply -f ./yaml
deploy nslookup function


