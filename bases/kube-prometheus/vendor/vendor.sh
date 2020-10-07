#!/usr/bin/env sh

# This script vendors the generated manifests and kustomization file frome
# github.com/prometheus-operator/kube-prometheus

# clone kube-prometheus repository into kube-prometheus folder
git clone https://github.com/prometheus-operator/kube-prometheus --depth=1 --single-branch

# copy manifests into manifests folder
cp -r kube-prometheus/manifests .

# copy kustomization
cp kube-prometheus/kustomization.yaml .

# delete repository
rm -rf kube-prometheus/
