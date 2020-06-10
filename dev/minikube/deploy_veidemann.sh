#!/usr/bin/env bash

if [ "$1" = true ]; then
  SKIP_AUTHENTICATION="\"true\""
else
  SKIP_AUTHENTICATION="\"false\""
fi

sed -i "s|value: \".*\"|value: $SKIP_AUTHENTICATION|" kustomization.yaml


SCRIPT_DIR=$(dirname $0)

kustomize build $SCRIPT_DIR | kubectl apply -f -
