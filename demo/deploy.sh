#!/bin/bash

# echo "Creating certificates"
# openssl req -nodes -new -x509 -keyout certs/ca.key -out certs/ca.crt -subj "/CN=Etimo Admission Controller"
# openssl genrsa -out certs/admission-tls.key 2048
# openssl req -new -key certs/admission-tls.key -subj "/CN=admission-server.default.svc" | openssl x509 -req -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/admission-tls.crt

# echo "Creating k8s Secret"

echo "Creating k8s admission deployment"
kubectl apply -f deployment.yaml
kubectl apply -f webhooks.yaml

kubectl apply -f deployments/ingress.yaml
sleep 5
kubectl delete ingress etimo-ingress-fail