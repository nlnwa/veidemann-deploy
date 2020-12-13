#!/usr/bin/env bash

sed -i "s|value: .*$|value: \"${SKIP_AUTHENTICATION:-false}\"|" kustomization.yaml

SCRIPT_DIR=$(dirname $0)

# Basic safeguard to ensure that we're working in a development-context
if ! [[ $(kubectl config current-context) = *minikube* ]]; then echo "WARNING: Not a dev context, use: kubectl config use-context minikube";  exit 1; fi

kustomize build $SCRIPT_DIR | ${SCRIPT_DIR}/../../scripts/versions.sh > versions.json

kustomize build $SCRIPT_DIR | kubectl apply -f -
