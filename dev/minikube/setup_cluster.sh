#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)
PREREQUISITES=${SCRIPT_DIR}/../../scripts/prerequisites.sh
UPDATE_HOSTS=${SCRIPT_DIR}/../../scripts/update_hosts.sh

source $PREREQUISITES kubectl minikube linkerd kustomize veidemannctl

EXISTS=$(minikube status --format='{{.Host}}')
if [ -n "$EXISTS" ]; then
  echo "Cluster is already set up and has status: ${EXISTS}"
  exit
fi

minikube addons disable ingress
minikube start # --cpus 2 --memory 8096 --vm-driver kvm2

# Create patch and update /etc/hosts for local cluster ip
LOCAL_IP=$(minikube ip)

cat <<EOF >minikube_ip_patch.yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: veidemann-controller
    spec:
      template:
        spec:
          hostAliases:
            - ip: $(minikube ip)
              hostnames:
                - veidemann.local
EOF

$UPDATE_HOSTS veidemann.local $LOCAL_IP

# Install Service mesh
linkerd install | kubectl apply -f -
linkerd check

# Install Ingress controller
kustomize build ${SCRIPT_DIR}/../bases/traefik | kubectl apply -f -

# Install Redis operator
kustomize build ${SCRIPT_DIR}/../../bases/redis-operator | kubectl apply -f -

# Install jager operator
kustomize build ${SCRIPT_DIR}/../../bases/jaeger-operator | kubectl apply -f -
