#!/usr/bin/env bash

set -e

SCRIPT_DIR=$(dirname $0)
PREREQUISITES=${SCRIPT_DIR}/../../scripts/prerequisites.sh
source $PREREQUISITES kubectl minikube kustomize

minikube delete
