#!/usr/bin/env bash

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm upgrade kube-prometheus prometheus-community/kube-prometheus-stack \
--create-namespace --namespace monitoring \
--install

