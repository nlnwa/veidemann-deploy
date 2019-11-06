#!/usr/bin/env bash

cat <<EOF >minikube_ip_patch.yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: veidemann-controller
    spec:
      template:
        spec:
          hostAliases:
            - ip: $(minikube ip)
              hostnames:
                - veidemann.local
EOF

kustomize build . | kubectl apply -f -
