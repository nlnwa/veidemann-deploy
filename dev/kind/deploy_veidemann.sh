#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)

# Basic safeguard to ensure that we're working in a development-context
if ! [[ $(kubectl config current-context) = *kind* ]]; then echo "WARNING: Not a dev context, use: kubectl config use-context kind-kind ";  exit 1; fi

kustomize build $SCRIPT_DIR | kubectl apply -f -
