#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)

# Install Service mesh
linkerd install --ha --controller-replicas 2 | kubectl apply -f -
linkerd check

# Install Ingress controller
kustomize build ${SCRIPT_DIR}/../bases/traefik | kubectl apply -f -

# Install Redis operator
kustomize build ${SCRIPT_DIR}/../bases/redis-operator | kubectl apply -f -

# Install jager operator
kustomize build ${SCRIPT_DIR}/../bases/jaeger-operator | kubectl apply -f -
