#!/bin/sh
set -eu
selector="app=devkitd"
kubectl get pods --selector=$selector --field-selector=status.phase=Running -o=go-template --template="{{range .items}}{{.metadata.name}}{{\"\n\"}}{{end}}" --sort-by="{.metadata.name}"
