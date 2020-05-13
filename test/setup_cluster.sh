#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)
PREREQUISITES=${SCRIPT_DIR}/../../scripts/prerequisites.sh

source $PREREQUISITES kubectl kind kustomize tilt

set -e

${SCRIPT_DIR}/kind-with-registry.sh

echo "waiting for cluster to be ready"
kubectl wait --for=condition=Ready nodes --all --timeout=5m
echo

# Install Redis operator
kustomize build ${SCRIPT_DIR}/../bases/redis-operator | kubectl apply -f -
