# roi ergasias
> [roí ergasías](https://translate.google.com/?sl=en&tl=el&text=workflow&op=translate) as pronounced in greek means workflow.

This **kubernetes operator** is meant to address a fundamental requirement of any data science / machine learning project running their pipelines on Kubernetes - which is to quickly provision a declarative data pipeline (on demand) for their various project needs using simple kubectl commands. Basically, implementing the concept of **No Ops**.

[![Go Reference](https://pkg.go.dev/badge/github.com/ankursoni/kubernetes-operator-roiergasias.svg)](https://pkg.go.dev/github.com/ankursoni/kubernetes-operator-roiergasias)

&#x1F534; NOT OFFICIALLY RELEASED YET - first version that supports split workflow jobs to be launched in July 2021.  
> MAIN BRANCH WORKS CORRECTLY AT THE MOMENT


## Repository map
```
 📌 -----------------------> you are here
┬
├── cmd    ----------------> contains go main starting point for roiergasias workflow cli
│   ├── linux   -----------> contains linux amd64 executable for roiergasias workflow cli
│   └── osx   -------------> contains mac-osx amd64 executable for roiergasias workflow cli
├── docs   ----------------> contains documentation / images
├── examples  
│   ├── hello-world   -----> contains both single node and multi node split workflow example
│   ├── machine-learning
│   │   ├── aws   ---------> contains multi node split workflow in 2 node groups example
│   │   └── local   -------> contains single node workflow example
├── infra   ---------------> contains terraform scripts for infrastructure as code
│   └── aws
├── operator   ------------> contains kubernetes operator code for roiergasias workflow
│   ├── api
│   ├── config
│   ├── controllers
│   ├── hack
│   └── helm   ------------> contains kubernetes operator helm chart repository
└── pkg   -----------------> contains go packages for roiergasias workflow engine
    ├── lib
    ├── mocks
    ├── steps
    ├── tasks
    └── workflow
```


## Run "Hello world" workflow locally
``` SH
# clone to a local git directory, if not already done so
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

# change to the local git directory
cd kubernetes-operator-roiergasias

# set execute permissions to go binary
chmod +x cmd/linux/roiergasias cmd/osx/roiergasias

# run the hello world workflow
./cmd/linux/roiergasias run -f ./examples/hello-world/hello-world.yaml
# or, for mac osx
./cmd/osx/roiergasias run -f ./examples/hello-world/hello-world.yaml
```
![hello-world](docs/images/hello-world.png)


## Run "Hello world" workflow via operator in kubernetes
### - Install [Helm](https://helm.sh/docs/intro/install/)
### - Optionally, install [Kubernetes by Docker Desktop](https://docs.docker.com/desktop/kubernetes/) or [Minikube](https://minikube.sigs.k8s.io/docs/start/)

``` SH
# install roiergasias operator
helm install --repo https://github.com/ankursoni/kubernetes-operator-roiergasias/raw/main/operator/helm/ \
  --version 0.1.1 \
  roiergasias-operator roiergasias-operator

# read the following example hello-world-kubernetes.yaml file
cat examples/hello-world/hello-world-kubernetes.yaml
---
apiVersion: batch.ankursoni.github.io/v1
kind: Workflow
metadata:
  name: roiergasias-demo
spec:
  workflowYAML:
    name: hello-world
    yaml: |
      version: 0.1

      environment:
        - welcome: "Welcome to the demo workflow!"

      task:
        - sequential:
            - print:
                - "Hello"
                - "World!"
            - print:
                - "Hi"
                - "Universe!"
            - environment:
                - greeting: "Warm greetings!"

        - sequential:
            - print:
                - "{{env:welcome}}"
            - execute:
                - "echo {{env:greeting}}"
            - environment:
                - greeting: "Warm greetings again!"

        - sequential:
            - execute:
                - "echo {{env:greeting}}"

  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: roiergasias
              image: docker.io/ankursoni/roiergasias-operator:workflow
              command: ["/root/roiergasias", "run", "--file=/root/hello-world/hello-world.yaml"]
              volumeMounts:
                # volume - 'yaml' is automatically created by the operator using a generated configMap
                - name: yaml
                  mountPath: /root/hello-world
---

# apply the manifest
kubectl apply -f examples/hello-world/hello-world-kubernetes.yaml

# browse workflow created by the manifest
kubectl get workflow

# browse job created by the workflow
kubectl get job

# browse pod created by the job
kubectl get pod

# check pod logs for the output and wait till it is completed
kubectl logs roiergasias-demo-<STRING_FROM_PREVIOUS_STEP>

# delete the manifest
kubectl delete -f examples/hello-world/hello-world-kubernetes.yaml

# delete the roiergasias namespace (optional)
kubectl delete ns roiergasias

# uninstall the operator (optional)
helm uninstall roiergasias-operator
```


## Run "Machine learning" workflow locally
Follow this [README](examples/machine-learning/local/README.md)


## Run "Machine learning" workflow in AWS
![topology](docs/images/aws-topology.png)
Follow this [README](examples/machine-learning/aws/README.md)


## Install Roiergasias operator
``` SH
# install the operator
helm install --repo https://github.com/ankursoni/kubernetes-operator-roiergasias/raw/main/operator/helm/ \
  --version v0.1.1 \
  roiergasias-operator roiergasias-operator

# uninstall the operator
helm uninstall roiergasias-operator
```
