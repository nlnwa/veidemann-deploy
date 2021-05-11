#!/usr/bin/env bash

# helm repo add scylla-operator https://storage.googleapis.com/scylla-operator-charts/stable
# helm repo update

SCRIPT_DIR=$(dirname $0)

helm upgrade scylla scylla-operator/scylla \
--install \
-f ${SCRIPT_DIR}/values.yaml \
--namespace veidemann
