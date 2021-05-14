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
### - Optionally, install [Kubernetes by Docker Desktop](https://docs.docker.com/desktop/kubernetes/) or [Minikube](https://minikube.sigs.k8s.io/docs/start/)
### - Optionally, install [Helm](https://helm.sh/docs/intro/install/)

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
go mod tidy && go mod download

# set execute permissions to go main binary
chmod +x cmd/main cmd/main-osx

# run the machine learning workflow
./cmd/main ./cmd/machine-learning/machine-learning.yaml
# or, for mac osx
./cmd/main-osx ./cmd/machine-learning/machine-learning.yaml
```


## Steps to build docker image
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# copy kaggle api credentials from ~/.kaggle
cp ~/.kaggle/kaggle.json cmd/

# set execute permissions to go main binary and python scripts
chmod +x cmd/main cmd/machine-learning/*.py

# build docker image
docker build -t roiergasias:latest cmd
```


## Steps to manually run go workflow via docker compose (after building docker image as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the cmd/machine-leaning directory
cd cmd/machine-learning

# run docker compose
docker-compose up

# clean up docker compose
docker-compose down
```


## Steps to push docker image to docker hub
``` SH
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
```
> NOTE: Make sure you have changed the above mentioned docker hub repository as **private** because it contains your kaggle api key credentials


## Steps to create kubernetes secret for docker hub credentials (after pushing docker image to docker hub as mentioned above)
``` SH
# create docker hub registry credentials (for pulling container image pushed previously)
helm upgrade -i --repo https://gabibbo97.github.io/charts imagepullsecrets imagepullsecrets \
  --version 3.0.0 \
  --create-namespace -n roiergasias \
  --set imagePullSecret.registryURL="docker.io" \
  --set imagePullSecret.secretName="container-registry-secret" \
  --set imagePullSecret.username="<USERNAME>" \
  --set imagePullSecret.password="<PASSWORD>"
# where, <USERNAME> and <PASSWORD> are the credentials for login to docker hub
```


## Steps to manually run go workflow via kubernetes helm charts (after creating kubernetes secret for docker hub credentials as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the cmd/machine-leaning directory
cd cmd/machine-learning

# create a new helm chart values override file: ./helm/roiergasias-job/values-secret.yaml
cp ./helm/roiergasias-job/values.yaml ./helm/roiergasias-job/values-secret.yaml

# update the values in ./helm/roiergasias-job/values-secret.yaml using nano or vi editor
nano ./helm/roiergasias-job/values-secret.yaml
# update "image" to be "docker.io/<REPOSITORY>/roiergasias:latest
#   where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
#          "docker.io/ankursoni/roiergasias:latest"
# update "hostPath" to be the full path of the local git clone directory + "/cmd/machine-learning", for e.g.,
#   "/Users/ankursoni/go/src/github.com/ankursoni/kubernetes-operator-roiergasias/cmd/machine-learning"

# output helm chart template for roiergasias job
helm template \
  -n roiergasias \
  -f ./helm/roiergasias-job/values-secret.yaml \
  roiergasias-job ./helm/roiergasias-job >machine-learning-job.yaml

# explore the contents of the machine-learning-job.yaml
cat machine-learning-job.yaml

# apply the manifest
kubectl apply -f machine-learning-job.yaml

# browse pod created by the job
kubectl get pods -n roiergasias

# check pod logs for the output
kubectl logs roiergasias-job-<STRING_FROM_PREVIOUS_STEP> -n roiergasias

# delete the manifest
kubectl delete -f machine-learning-job.yaml
```


## Steps to manually run go workflow via kubernetes operator (after pushing the docker image to docker hub and creating kubernetes secret for docker hub credentials as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the operator directory
cd operator

# optionally, download go module dependencies to local cache
go mod tidy && go mod download

# build the operator docker image and push to docker hub
make docker-build docker-push IMG=docker.io/<REPOSITORY>/roiergasias:operator
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
# make docker-build docker-push IMG=docker.io/ankursoni/roiergasias:operator

# deploy the operator to kubernetes
make deploy IMG=docker.io/<REPOSITORY>/roiergasias:operator
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
make deploy IMG=docker.io/ankursoni/roiergasias:operator

# change to the cmd/machine-leaning directory
cd ../cmd/machine-learning

# create a new helm chart values override file: ./helm/roiergasias-workflow/values-secret.yaml
cp ./helm/roiergasias-workflow/values.yaml ./helm/roiergasias-workflow/values-secret.yaml

# update the values in ./helm/roiergasias-workflow/values-secret.yaml using nano or vi editor
nano ./helm/roiergasias-workflow/values-secret.yaml
# update "image" to be "docker.io/<REPOSITORY>/roiergasias:latest
#   where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
#          "docker.io/ankursoni/roiergasias:latest"
# update "hostPath" to be the full path of the local git clone directory + "/cmd/machine-learning", for e.g.,
#   "/Users/ankursoni/go/src/github.com/ankursoni/kubernetes-operator-roiergasias/cmd/machine-learning"

# output helm chart template for roiergasias workflow
helm template \
  -n roiergasias \
  -f ./helm/roiergasias-workflow/values-secret.yaml \
  roiergasias-workflow ./helm/roiergasias-workflow >machine-learning-workflow.yaml

# explore the contents of the machine-learning-workflow.yaml
cat machine-learning-workflow.yaml

# apply the manifest
kubectl apply -f machine-learning-workflow.yaml

# browse pod created by the job
kubectl get pods -n roiergasias

# check pod logs for the output
kubectl logs roiergasias-workflow-<STRING_FROM_PREVIOUS_STEP> -n roiergasias

# delete the manifest
kubectl delete -f machine-learning-workflow.yaml

# change to the operator directory
cd ../../operator

# undeploy the operator
make undeploy
```


## Steps to provision AWS infrastructure
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the infra/aws directory
cd infra/aws

# create a new terraform values override file: ./values-secret.tfvars
cp ./values.tfvars ./values-secret.tfvars

# substitute the value for <PREFIX> by replacing PLACEHOLDER in the following command:
# PLACEHOLDER e.g. "roiergasias" or "workflow" etc.
sed -i 's|<PREFIX>|PLACEHOLDER|g' ./values-secret.tfvars

# substitute the value for <ENVIRONMENT> by replacing PLACEHOLDER in the command
# PLACEHOLDER e.g. "demo" or "play" or "poc" or "dev" or "test" etc.
sed -i 's|<ENVIRONMENT>|PLACEHOLDER|g' ./values-secret.tfvars

# substitute the value for <REGION> by replacing PLACEHOLDER in the command
# PLACEHOLDER e.g. "ap-southeast-2" for Sydney or "ap-southeast-1" for Singapore or "us-east-1" for North Virginia etc.
# Browse https://aws.amazon.com/about-aws/global-infrastructure/regions_az/ for more regions
# run this to know more: "aws ec2 describe-regions -o table"
sed -i 's|<REGION>|PLACEHOLDER|g' ./values-secret.tfvars

# substitute the value for <NODE_COUNT> by replacing PLACEHOLDER in the command
# PLACEHOLDER e.g. 1
sed -i 's|<NODE_COUNT>|PLACEHOLDER|g' ./values-secret.tfvars

# verify the ./values-secret.tfvars file by displaying its content
cat ./values-secret.tfvars

# output should be something like this
prefix="roiergasias"
environment="demo"
region="ap-southeast-2"
node_count=1

# if there is a correction needed then use text editor 'nano' to update the file and then press ctrl+x after you are done editing
nano ./values-secret.tfvars

# initialise terraform providers
terraform init

# execute infrastructure provisioning command
terraform apply -var-file=values-secret.tfvars

# get kubectl credentials
aws eks update-kubeconfig --region <REGION> --name <PREFIX>-<ENVIRONMENT>-eks01
# for e.g., aws eks update-kubeconfig --region ap-southeast-2 --name roiergasias-demo-eks01

# patch coredns to use fargate
kubectl patch deployment coredns -n kube-system --type json \
-p='[{"op": "remove", "path": "/spec/template/metadata/annotations/eks.amazonaws.com~1compute-type"}]'
```


## Steps to automatically run go workflow via kubernetes operator (after provisioning AWS infrastructure as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the cmd/machine-leaning directory
cd ../cmd/machine-learning
```