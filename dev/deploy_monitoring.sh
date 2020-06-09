#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)

kustomize build $SCRIPT_DIR/bases/kube-prometheus | kubectl apply -f -
