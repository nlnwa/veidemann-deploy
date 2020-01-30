#!/usr/bin/env bash

set -e

OPSYS=$(uname)
KUSTOMIZE_VERSION=v3.5.4
KIND_VERSION=0.7.0
LINKERD_VERSION=stable-2.6.1
KUBECTL_VERSION=v1.16.4
MINIKUBE_VERSION=v1.6.2
VEIDEMANNCTL_VERSION=0.3.3

function check_cmd() {
  local CMD=$1
  local VER=$2
  W=$(which ${CMD} 2>/dev/null)
  if [ $? -ne 0 ]; then
    ask "${CMD} not found. Would you like to install ${CMD} ${VER}? (y/n)"
    return $?
  fi

  CURRENT_VER=$(eval $(printf "%s %s" "$W" "$3"))

  if ! [ "$CURRENT_VER" = "${VER}" ]; then
    ask "${CMD} version is ${CURRENT_VER}, but we would like to have ${VER}. Would you like to install ${CMD} ${VER}? (y/n)"
    return $?
  fi
  return 0
}

function ask() {
  read -p "${1}" -n 1 -r
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo
    return 0
  fi
  echo
  return 1
}

for CMD in "$@"; do
  case $CMD in
  kubectl)
    check_cmd kubectl $KUBECTL_VERSION 'version --client --short | sed -e "s/Client Version: \(.*\)/\1/"'
    INSTALL=$?
    if [ $INSTALL -ne 0 ]; then
      echo "Installing Kubectl"
      curl -LO https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl
      chmod +x ./kubectl
      sudo mv ./kubectl /usr/local/bin/kubectl
    fi
    ;;
  kind)
    check_cmd kind $KIND_VERSION '-q version'
    INSTALL=$?
    if [ $INSTALL -ne 0 ]; then
      echo "Installing Kind"
      curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v${KIND_VERSION}/kind-${OPSYS}-amd64
      chmod +x ./kind
      sudo mv ./kind /usr/local/bin/kind
    fi
    ;;
  kustomize)
    check_cmd kustomize $KUSTOMIZE_VERSION 'version --short | sed -e "s/{kustomize\/\(\S\+\).*/\1/"'
    INSTALL=$?
    if [ $INSTALL -ne 0 ]; then
      echo "Installing Kustomize"
      curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_${OPSYS}_amd64.tar.gz |
        tar xz
      chmod +x ./kustomize
      sudo mv ./kustomize /usr/local/bin/kustomize
    fi
    ;;
  linkerd)
    check_cmd linkerd $LINKERD_VERSION 'version --client --short'
    INSTALL=$?
    if [ $INSTALL -ne 0 ]; then
      echo "Installing Linkerd"
      curl -sL https://run.linkerd.io/install | sh
    fi
    ;;
  minikube)
    check_cmd minikube $MINIKUBE_VERSION 'version | grep version | sed -e "s/.*version: \(.*\)/\1/"'
    INSTALL=$?
    if [ $INSTALL -ne 0 ]; then
      echo "Installing Minikube"
      curl -LO https://storage.googleapis.com/minikube/releases/${MINIKUBE_VERSION}/minikube-linux-amd64 &&
        sudo install minikube-linux-amd64 /usr/local/bin/minikube
    fi
    ;;
  veidemannctl)
    check_cmd veidemannctl $VEIDEMANNCTL_VERSION '--version | sed -e "s/.*version: \(.*\), Go version.*/\1/"'
    INSTALL=$?
    if [ $INSTALL -ne 0 ]; then
      echo "Installing Veidemannctl"
      curl -L#o veidemannctl https://github.com/nlnwa/veidemannctl/releases/download/${VEIDEMANNCTL_VERSION}/veidemannctl_${VEIDEMANNCTL_VERSION}_linux_amd64 &&
        sudo install veidemannctl /usr/local/bin/veidemannctl &&
        sudo sh -c "/usr/local/bin/veidemannctl completion > /etc/bash_completion.d/veidemannctl"
    fi
    ;;
  esac
done
echo
