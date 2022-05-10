#!/bin/bash

docker build . -t etimodanielwinther/etimo-kubernetes-admissioncontroller
docker tag etimodanielwinther/etimo-kubernetes-admissioncontroller etimodanielwinther/etimo-kubernetes-admissioncontroller:latest
docker push etimodanielwinther/etimo-kubernetes-admissioncontroller:latest
