#!/usr/bin/env bash

set -e

SCRIPT_DIR=$(dirname $0)
PREREQUISITES=${SCRIPT_DIR}/../../scripts/prerequisites.sh
source $PREREQUISITES kubectl kind kustomize

kind delete cluster
