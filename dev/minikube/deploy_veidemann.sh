#!/usr/bin/env bash

SCRIPT_DIR=$(dirname $0)

kustomize build $SCRIPT_DIR | kubectl apply -f -
