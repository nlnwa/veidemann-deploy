#!/usr/bin/env bash

# fetch deployment from jaegertracing, see https://github.com/jaegertracing/jaeger-operator
wget -P deploy/crds https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/crds/jaegertracing.io_jaegers_crd.yaml
wget -P deploy https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/service_account.yaml
wget -P deploy https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role.yaml
wget -P deploy https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role_binding.yaml
wget -P deploy https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/operator.yaml
