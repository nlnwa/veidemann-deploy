#!/usr/bin/env bash

# fetch deployment from jaegertracing, see https://github.com/jaegertracing/jaeger-operator
wget -NP vendor/crds https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/crds/jaegertracing.io_jaegers_crd.yaml
wget -NP vendor https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/service_account.yaml
wget -NP vendor https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role.yaml
wget -NP vendor https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role_binding.yaml
wget -NP vendor https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/operator.yaml

wget -NP vendor https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role.yaml
wget -NP vendor https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role_binding.yaml
