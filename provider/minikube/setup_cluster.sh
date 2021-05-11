#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)
PREREQUISITES=${SCRIPT_DIR}/../../scripts/prerequisites.sh
UPDATE_HOSTS=${SCRIPT_DIR}/../../scripts/update_hosts.sh

source $PREREQUISITES kubectl minikube helm linkerd kustomize veidemannctl

set -e

minikube start # --cpus 4 --memory 12000 --driver docker

# Create patch and update /etc/hosts for local cluster ip
HOSTNAME=veidemann.test
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
        - ip: ${LOCAL_IP}
          hostnames:
            - ${HOSTNAME}
EOF

$UPDATE_HOSTS $LOCAL_IP veidemann.test linkerd.veidemann.test

echo "Waiting for nodes to be ready"
kubectl wait --for=condition=Ready nodes --all --timeout=5m
echo

kubectl config set contexts.minikube.namespace veidemann

# Install linkerd
LINKERD_SERVER_VERSION=$(linkerd version | tail -1 | awk '{print $3}')
if [ "$LINKERD_SERVER_VERSION" = "unavailable" ]; then
  linkerd check --pre
  linkerd install | kubectl apply -f -
fi
linkerd check

# Install traefik
kustomize build ${SCRIPT_DIR}/../../dev/ingress-traefik | kubectl apply -f -

# Install linkerd ingressroute
kustomize build ${SCRIPT_DIR}/../../dev/linkerd | kubectl apply -f -

# Install redis-operator
kustomize build ${SCRIPT_DIR}/../../dev/redis-operator | kubectl apply -f -


# Install jaeger-operator
kustomize build ${SCRIPT_DIR}/../../dev/observability | kubectl apply -f -

# Give kubernetes time to install jaeger-operator CRD
sleep 1

# Install jaeger
kustomize build ${SCRIPT_DIR}/../../dev/observability/jaeger | kubectl apply -f -

# Install cert manager
${SCRIPT_DIR}/../../dev/cert-manager/install_cert_manager.sh

# Install scylla-operator
${SCRIPT_DIR}/../../dev/scylla-operator/install_scylla_operator.sh
