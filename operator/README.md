# Roiergasias Kubernetes operator


## Install Roiergasias operator in Kubernetes
```shell
# install the operator
helm install --repo https://github.com/ankursoni/kubernetes-operator-roiergasias/raw/main/operator/helm/ \
  --version v0.1.2 \
  roiergasias-operator roiergasias-operator

# uninstall the operator
helm uninstall roiergasias-operator
```


## Getting started with Roiergasias workflow
For workflow yaml file syntax and command syntax, follow this [README](../cmd/README.md)


## Repository map
```text
┬
├── ...
│   📌 --------------------> you are here
├── operator   ------------> contains kubernetes operator code for roiergasias workflow
│   ├── api
│   ├── config
│   ├── controllers
│   ├── hack
│   └── helm   ------------> contains kubernetes operator helm chart repository
└── ...
```