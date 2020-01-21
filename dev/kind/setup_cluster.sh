#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)
PREREQUISITES=${SCRIPT_DIR}/../../scripts/prerequisites.sh
UPDATE_HOSTS=${SCRIPT_DIR}/../../scripts/update_hosts.sh

source $PREREQUISITES kubectl kind linkerd kustomize

kind create cluster --config=kind-config.yaml

echo "waiting for cluster to be ready"
while [[ $(kubectl get node kind-control-plane --no-headers | sed -e 's/\S\+\s\+\(\S\+\).*/\1/') != "Ready" ]]; do
  echo -n "."
  sleep 2
done
echo

# Create patch and update /etc/hosts for local cluster ip
LOCAL_IP=$(ip route get 1 | sed -n 's/^.*src \([0-9.]*\) .*$/\1/p')

cat <<EOF >kind_ip_patch.yaml
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
                - veidemann.local
EOF

$UPDATE_HOSTS veidemann.local $LOCAL_IP

# Install Service mesh
linkerd install | kubectl apply -f -
linkerd check

# Install Ingress controller
kustomize build ${SCRIPT_DIR}/traefik | kubectl apply -f -

# Install Redis operator
kustomize build ${SCRIPT_DIR}/../../bases/redis-operator | kubectl apply -f -

R=$(firewall-cmd --direct --get-rules ipv4 filter INPUT)
if [ "$R" != "4 -i docker0 -j ACCEPT" ]; then
  sudo firewall-cmd --direct --add-rule ipv4 filter INPUT 4 -i docker0 -j ACCEPT
fi