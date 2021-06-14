# roi ergasias
> [roí ergasías](https://translate.google.com/?sl=en&tl=el&text=workflow&op=translate) as pronounced in greek means workflow.

This **kubernetes operator** is meant to address a fundamental requirement of any data science / machine learning project running their pipelines on Kubernetes - which is to quickly provision a declarative data pipeline (on demand) for their various project needs using simple kubectl commands. Basically, implementing the concept of **No Ops**.

[![Go Reference](https://pkg.go.dev/badge/github.com/ankursoni/kubernetes-operator-roiergasias.svg)](https://pkg.go.dev/github.com/ankursoni/kubernetes-operator-roiergasias)

&#x1F534; NOT OFFICIALLY RELEASED YET - first version that supports split workflow jobs to be launched in June 2021.  
> MAIN BRANCH WORKS CORRECTLY AT THE MOMENT

## Run "Hello world" workflow locally
``` SH
# clone to a local git directory, if not already done so
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

# change to the local git directory
cd kubernetes-operator-roiergasias

# set execute permissions to go main binary
chmod +x cmd/main cmd/main-osx

# run the hello world workflow
./cmd/main ./cmd/hello-world/hello-world.yaml
# or, for mac osx
./cmd/main-osx ./cmd/hello-world/hello-world.yaml
```
![hello-world](docs/images/hello-world.png)


## Run "Hello world" workflow via operator in kubernetes
### - Install [Helm](https://helm.sh/docs/intro/install/)
### - Optionally, install [Kubernetes by Docker Desktop](https://docs.docker.com/desktop/kubernetes/) or [Minikube](https://minikube.sigs.k8s.io/docs/start/)

``` SH
# install roiergasias operator
helm install --repo https://github.com/ankursoni/kubernetes-operator-roiergasias/raw/main/operator/helm/ \
  --version 0.1.0 \
  roiergasias-operator roiergasias-operator

# write the following yaml file
cat <<EOF>hello-world-manifest.yaml
apiVersion: batch.ankursoni.github.io/v1
kind: Workflow
metadata:
  name: roiergasias-kubernetes
spec:
  workflowYAML:
    name: hello-world
    yaml: |
      version: 1.0

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
              set-environment:
                - greeting: "Warm greetings!"

        - sequential:
            - print:
                - "{{env:welcome}}"
            - print:
                - "{{env:greeting}}"
      
  jobTemplate:
    spec:
      backoffLimit: 3
      template:
        spec:
          imagePullSecrets:
            - name: container-registry-secret
          containers:
            - name: roiergasias
              image: docker.io/ankursoni/roiergasias-operator:workflow
              command: ["./cmd/main", "./cmd/hello-world/hello-world.yaml"]
              volumeMounts:
                # volume - 'yaml' is automatically created by the operator using a generated configMap
                - name: yaml
                  mountPath: /root/cmd/hello-world
              resources:
                requests:
                  memory: "100Mi"
                  cpu: "100m"
                limits:
                  memory: "200Mi"
                  cpu: "200m"
          restartPolicy: Never
EOF

# apply the manifest
kubectl apply -f hello-world-manifest.yaml

# delete the manifest
kubectl delete -f hello-world-manifest.yaml

# uninstall the operator
helm uninstall roiergasias-operator
```


## Run "Machine learning" workflow locally
Follow this [README](cmd/machine-learning/README.md)


## Run "Machine learning" workflow in AWS
![topology](docs/images/aws-topology.png)
Follow this [README](cmd/machine-learning/README.md#process-data-train-ml-model--evaluate-ml-model-in-aws)


## Install Roiergasias operator
``` SH
# install the operator
helm install --repo https://github.com/ankursoni/kubernetes-operator-roiergasias/raw/main/operator/helm/ \
  --version v0.1.0 \
  roiergasias-operator roiergasias-operator

# uninstall the operator
helm uninstall roiergasias-operator
```
