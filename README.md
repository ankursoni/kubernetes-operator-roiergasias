# roi ergasias
> [roí ergasías](https://translate.google.com/?sl=en&tl=el&text=workflow&op=translate) as pronounced in greek means workflow.

This **kubernetes operator** is meant to address a fundamental requirement of any data science / machine learning project running their pipelines on Kubernetes - which is to quickly provision a declarative data pipeline (on demand) for their various project needs using simple kubectl commands. Basically, implementing the concept of **No Ops**.

---

## Install basic pre-requisites
### - Install [Docker](https://docs.docker.com/get-docker/)
### - Install [Docker Compose](https://docs.docker.com/compose/install/)
### - Optionally, install [Go](https://golang.org/doc/install)

---

## Run "Hello world" workflow locally
``` SH
# clone to a local git directory
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

# change to the local git directory
cd kubernetes-operator-roiergasias

# download go module dependencies to local cache
go mod download

# set execute permissions to go main binary
chmod +x cmd/main cmd/main-osx

# run the hello world workflow
./cmd/main ./cmd/hello-world/hello-world.yaml
# for mac osx
./cmd/main-osx ./cmd/hello-world/hello-world.yaml
```

## Run "Machine learning" workflow locally
Follow this [README](cmd/machine-learning/README.md)