# Process data, train ml model & evaluate ml model


## Source of inspiration
Converted the following jupyter notebook to Python scripts:  
https://www.kaggle.com/ilyapozdnyakov/rain-in-australia-precision-recall-curves-viz#Rain-prediction-in-Australia

---

## Install pre-requisites

### - Install [Python 3](https://www.python.org/downloads/)
Also, install the following *pip* packages:
``` SH
pip install pandas sklearn joblib
# or,
# pip3 install pandas sklearn joblib
```
### - Install [Kaggle CLI](https://github.com/Kaggle/kaggle-api)
#### -- Make sure kaggle is configured with api key in ~/.kaggle/kaggle.json
[Kaggle API Credentials](https://github.com/Kaggle/kaggle-api#api-credentials)
### - Optionally, install [Go](https://golang.org/doc/install)
### - Optionally, install [Docker Desktop](https://www.docker.com/products/docker-desktop) or [Docker](https://docs.docker.com/get-docker/)
### - Optionally, install [Docker Compose](https://docs.docker.com/compose/install/)
### - Optionally, install [Helm](https://helm.sh/docs/intro/install/)
### - Optionally, install [Kubernetes by Docker Desktop](https://docs.docker.com/desktop/kubernetes/) or [Minikube](https://minikube.sigs.k8s.io/docs/start/)

---

## Clone to a local git directory
``` SH
# clone to a local git directory, if not already done so
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git
```


## Steps to manually run python scripts
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the cmd/machine-leaning directory
cd cmd/machine-learning

# download dataset from kaggle
kaggle datasets download jsphyg/weather-dataset-rattle-package -o -f weatherAUS.csv

# unzip and delete the downloaded zip file
unzip -o weatherAUS.csv.zip
rm -f weatherAUS.csv.zip

# set execute permissions to python script files
chmod +x ./*.py

# process data from first argument file saving output in second argument file
./process-data.py ./weatherAUS.csv ./processed-weatherAUS.csv

# train ml model from first argument file and saving model output in second argument file
./train-model.py ./processed-weatherAUS.csv ./ml-model.joblib

# evaluate ml model by reading processed data and model from first and second argument files
./evaluate-model.py ./processed-weatherAUS.csv ./ml-model.joblib
```


## Steps to manually run go workflow
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# optionally, download go module dependencies to local cache
go mod download

# set execute permissions to go main binary
chmod +x cmd/main cmd/main-osx

# run the machine learning workflow
./cmd/main ./cmd/machine-learning/machine-learning.yaml
# or, for mac osx
./cmd/main-osx ./cmd/machine-learning/machine-learning.yaml
```


## Steps to manually run go workflow via docker compose
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# copy kaggle api credentials from ~/.kaggle
cp ~/.kaggle/kaggle.json cmd/

# set execute permissions to go main binary and python scripts
chmod +x cmd/main cmd/machine-learning/*.py

# build docker image
docker build -t roiergasias:latest cmd

# change to the cmd/machine-leaning directory
cd cmd/machine-learning

# run docker compose
docker-compose up
```


## Steps to manually run go workflow via kubernetes helm charts (after trying out via docker compose as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# re-tag local docker image
docker tag roiergasias:latest docker.io/<REPOSITORY>/roiergasias:latest
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
# docker tag roiergasias:latest docker.io/ankursoni/roiergasias:latest

# login to docker hub
docker login

# push the docker image to docker hub
docker push docker.io/<REPOSITORY>/roiergasias:latest
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
# docker push docker.io/ankursoni/roiergasias:latest

# create docker hub registry credentials (for pulling container image pushed previously)
helm upgrade -i --repo https://gabibbo97.github.io/charts imagepullsecrets imagepullsecrets \
  --version 3.0.0 \
  --create-namespace -n roiergasias \
  --set imagePullSecret.registryURL="docker.io" \
  --set imagePullSecret.secretName="container-registry-secret" \
  --set imagePullSecret.username="<USERNAME>" \
  --set imagePullSecret.password="<PASSWORD>"
# where, <USERNAME> and <PASSWORD> are the credentials for login to docker hub

# install helm chart for workflow
helm upgrade -i \
  -n roiergasias \
  roiergasias ./cmd/machine-learning/helm/roiergasias

# browse the pod created by the job
kubectl get pods -n roiergasias

# check the pod logs for the output of the workflow
kubectl logs roiergasias-job-<STRING_FROM_PREVIOUS_STEP> -n roiergasias
```
