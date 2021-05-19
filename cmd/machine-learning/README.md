# Process data, train ml model & evaluate ml model locally


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

# set execute permissions to go main binary
chmod +x cmd/main

# build docker image
docker build -t roiergasias:local cmd
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
docker tag roiergasias:local docker.io/<REPOSITORY>/roiergasias:local
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
# docker tag roiergasias:local docker.io/ankursoni/roiergasias:local

# login to docker hub
docker login

# push the docker image to docker hub
docker push docker.io/<REPOSITORY>/roiergasias:local
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
# docker push docker.io/ankursoni/roiergasias:local
```
> NOTE: Make sure you have changed the above mentioned docker hub repository as **private** because it contains your kaggle api key credentials


## Steps to create kubernetes secret for docker hub credentials (after pushing docker image to docker hub as mentioned above)
``` SH
# create docker hub registry credentials (for pulling docker image pushed previously)
helm upgrade -i --repo https://gabibbo97.github.io/charts imagepullsecrets imagepullsecrets \
  --version 3.0.0 \
  --create-namespace -n roiergasias \
  --set imagePullSecret.registryURL="docker.io" \
  --set imagePullSecret.secretName="container-registry-secret" \
  --set imagePullSecret.username="<USERNAME>" \
  --set imagePullSecret.password="<PASSWORD>"
# where, <USERNAME> and <PASSWORD> are the credentials for login to docker hub
```


## Steps to manually run go workflow via kubernetes job (after creating kubernetes secret for docker hub credentials as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the cmd/machine-leaning directory
cd cmd/machine-learning

# create a new helm chart values override file: ./helm/roiergasias-job/values-secret.yaml
cp ./helm/roiergasias-job/values.yaml ./helm/roiergasias-job/values-secret.yaml

# update the values in ./helm/roiergasias-job/values-secret.yaml using nano or vi editor
nano ./helm/roiergasias-job/values-secret.yaml
# update "image" to be "docker.io/<REPOSITORY>/roiergasias:local"
#   where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
#          "docker.io/ankursoni/roiergasias:local"
# update "hostPath" to be the full path of the local git clone directory + "/cmd/machine-learning", for e.g.,
#   "/Users/ankursoni/go/src/github.com/ankursoni/kubernetes-operator-roiergasias/cmd/machine-learning"

# output helm chart template for roiergasias job
helm template \
  -n roiergasias \
  -f ./helm/roiergasias-job/values-secret.yaml \
  roiergasias-job ./helm/roiergasias-job >machine-learning-job-manifest.yaml

# explore the contents of the machine-learning-job-manifest.yaml
cat machine-learning-job-manifest.yaml

# apply the manifest
kubectl apply -f machine-learning-job-manifest.yaml

# browse pod created by the job
kubectl get pods -n roiergasias

# check pod logs for the output
kubectl logs roiergasias-job-<STRING_FROM_PREVIOUS_STEP> -n roiergasias

# delete the manifest
kubectl delete -f machine-learning-job-manifest.yaml
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
# make deploy IMG=docker.io/ankursoni/roiergasias:operator

# change to the cmd/machine-leaning directory
cd ../cmd/machine-learning

# create a new helm chart values override file: ./helm/roiergasias-workflow/values-secret.yaml
cp ./helm/roiergasias-workflow/values.yaml ./helm/roiergasias-workflow/values-secret.yaml

# update the values in ./helm/roiergasias-workflow/values-secret.yaml using nano or vi editor
nano ./helm/roiergasias-workflow/values-secret.yaml
# update "image" to be "docker.io/<REPOSITORY>/roiergasias:local"
#   where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
#          "docker.io/ankursoni/roiergasias:local"
# update "hostPath" to be the full path of the local git clone directory + "/cmd/machine-learning", for e.g.,
#   "/Users/ankursoni/go/src/github.com/ankursoni/kubernetes-operator-roiergasias/cmd/machine-learning"

# output helm chart template for roiergasias workflow
helm template \
  -n roiergasias \
  -f ./helm/roiergasias-workflow/values-secret.yaml \
  roiergasias-workflow ./helm/roiergasias-workflow >machine-learning-workflow-manifest.yaml

# explore the contents of the machine-learning-workflow-manifest.yaml
cat machine-learning-workflow-manifest.yaml

# apply the manifest
kubectl apply -f machine-learning-workflow-manifest.yaml

# browse pod created by the job
kubectl get pods -n roiergasias

# check pod logs for the output
kubectl logs roiergasias-workflow-<STRING_FROM_PREVIOUS_STEP> -n roiergasias

# delete the manifest
kubectl delete -f machine-learning-workflow-manifest.yaml

# change to the operator directory
cd ../../operator

# undeploy the operator
make undeploy
```

---

# Process data, train ml model & evaluate ml model in AWS


## Install pre-requisites

### - Install [Kaggle CLI](https://github.com/Kaggle/kaggle-api)
#### -- Make sure kaggle is configured with api key in ~/.kaggle/kaggle.json
[Kaggle API Credentials](https://github.com/Kaggle/kaggle-api#api-credentials)
### - Install [AWS CLI v2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
### - Install [Helm](https://helm.sh/docs/intro/install/)
### - Install [Kubectl](https://kubernetes.io/docs/tasks/tools/)
### - Optionally, install [Go](https://golang.org/doc/install)

---

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


## Steps to build docker image for AWS
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# copy kaggle api credentials from ~/.kaggle
cp ~/.kaggle/kaggle.json cmd/

# copy aws cli credentials from ~/.aws
cp -RP ~/.aws cmd/.aws/

# set execute permissions to go main binary
chmod +x cmd/main

# build docker image
docker build -t roiergasias:aws -f cmd/Dockerfile-AWS cmd
```


## Steps to push docker image for AWS to docker hub
``` SH
# re-tag local docker image
docker tag roiergasias:aws docker.io/<REPOSITORY>/roiergasias:aws
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
# docker tag roiergasias:aws docker.io/ankursoni/roiergasias:aws

# login to docker hub
docker login

# push the docker image to docker hub
docker push docker.io/<REPOSITORY>/roiergasias:aws
# where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
# docker push docker.io/ankursoni/roiergasias:aws
```
> NOTE: Make sure you have changed the above mentioned docker hub repository as **private** because it contains your kaggle api key credentials and aws cli credentials


## Steps to create kubernetes secret for docker hub credentials (after pushing docker image for AWS as mentioned above)
``` SH
# create docker hub registry credentials (for pulling docker image pushed previously)
helm upgrade -i --repo https://gabibbo97.github.io/charts imagepullsecrets imagepullsecrets \
  --version 3.0.0 \
  --create-namespace -n roiergasias \
  --set imagePullSecret.registryURL="docker.io" \
  --set imagePullSecret.secretName="container-registry-secret" \
  --set imagePullSecret.username="<USERNAME>" \
  --set imagePullSecret.password="<PASSWORD>"
# where, <USERNAME> and <PASSWORD> are the credentials for login to docker hub
```


## Steps to manually run go workflow via kubernetes operator (after provisioning AWS infrastructure and pushing docker image for AWS as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the operator directory
cd operator

# deploy the operator to kubernetes
make deploy IMG=docker.io/ankursoni/roiergasias-operator:latest

# change to the cmd/machine-leaning directory
cd ../cmd/machine-learning

# upload the workflow yaml and python script files
# assumes 'roiergasias' as <PREFIX> and 'demo' as <ENVIRONMENT> values
aws s3 cp process-data.py s3://roiergasias-demo-s3b01/
aws s3 cp train-model.py s3://roiergasias-demo-s3b01/
aws s3 cp evaluate-model.py s3://roiergasias-demo-s3b01/

# create a new helm chart values override file: ./helm/roiergasias-aws-manual/values-secret.yaml
cp ./helm/roiergasias-aws-manual/values.yaml ./helm/roiergasias-aws-manual/values-secret.yaml

# update the values in ./helm/roiergasias-aws-manual/values-secret.yaml using nano or vi editor
nano ./helm/roiergasias-aws-manual/values-secret.yaml
# update "image" to be "docker.io/<REPOSITORY>/roiergasias:local"
#   where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
#          "docker.io/ankursoni/roiergasias:local"
# update "s3URI" to be "s3://<PREFIX>-<ENVIRONMENT>-s3b01/",
#   where, <PREFIX> and <ENVIRONMENT> were set in the infra/aws/values-secret.tfvars, for e.g.,
#          "s3://roiergasias-demo-s3b01/"
# update "enablePersistentVolume" to be either 0 (default) or 1 to turn OFF or ON the persistent volume from elastic file system (EFS)
#   it gives persistence to the data written by steps in the workflow. Regardless of value, each sequential task syncs up to the S3 at the end of each stage.
# update "efsId" by running the command and copying the second value from output:
#   aws --region <REGION> efs describe-file-systems --query 'FileSystems[*].[Name, FileSystemId]' --output text | grep <PREFIX>-<ENVIRONMENT>-efs01, for e.g.,
# aws --region ap-southeast-2 efs describe-file-systems --query 'FileSystems[*].[Name, FileSystemId]' --output text | grep roiergasias-demo-efs01

# output helm chart template for roiergasias aws manual
helm template \
  -n roiergasias \
  -f ./helm/roiergasias-aws-manual/values-secret.yaml \
  roiergasias-aws-manual ./helm/roiergasias-aws-manual >machine-learning-aws-manual-manifest.yaml

# explore the contents of the machine-learning-aws-manual-manifest.yaml
cat machine-learning-aws-manual-manifest.yaml

# apply the manifest
kubectl apply -f machine-learning-aws-manual-manifest.yaml

# browse pod created by the job
kubectl get pods -n roiergasias

# check pod logs for the output
kubectl logs roiergasias-aws-<STRING_FROM_PREVIOUS_STEP> -n roiergasias

# delete the manifest
kubectl delete -f machine-learning-aws-manual-manifest.yaml

# change to the operator directory
cd ../../operator

# undeploy the operator
make undeploy
```


## Steps to automatically run go workflow via kubernetes operator (after provisioning AWS infrastructure and pushing docker image for AWS as mentioned above)
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the operator directory
cd operator

# install the operator in kubernetes
helm install roiergasias-operator ./helm/roiergasias-operator

# change to the cmd/machine-leaning directory
cd ../cmd/machine-learning

# upload the workflow yaml and python script files
# assumes 'roiergasias' as <PREFIX> and 'demo' as <ENVIRONMENT> values
aws s3 cp process-data.py s3://roiergasias-demo-s3b01/
aws s3 cp train-model.py s3://roiergasias-demo-s3b01/
aws s3 cp evaluate-model.py s3://roiergasias-demo-s3b01/

# create a new helm chart values override file: ./helm/roiergasias-aws-auto/values-secret.yaml
cp ./helm/roiergasias-aws-auto/values.yaml ./helm/roiergasias-aws-auto/values-secret.yaml

# update the values in ./helm/roiergasias-aws-auto/values-secret.yaml using nano or vi editor
nano ./helm/roiergasias-aws-auto/values-secret.yaml
# update "image" to be "docker.io/<REPOSITORY>/roiergasias:local"
#   where, <REPOSITORY> is the docker hub repository name or docker hub username, for e.g.,
#          "docker.io/ankursoni/roiergasias:local"
# update "s3URI" to be "s3://<PREFIX>-<ENVIRONMENT>-s3b01/",
#   where, <PREFIX> and <ENVIRONMENT> were set in the infra/aws/values-secret.tfvars, for e.g.,
#          "s3://roiergasias-demo-s3b01/"
# update "enablePersistentVolume" to be either 0 (default) or 1 to turn OFF or ON the persistent volume from elastic file system (EFS)
#   it gives persistence to the data written by steps in the workflow. Regardless of value, each sequential task syncs up to the S3 at the end of each stage.
# update "efsId" by running the command and copying the second value from output:
#   aws --region <REGION> efs describe-file-systems --query 'FileSystems[*].[Name, FileSystemId]' --output text | grep <PREFIX>-<ENVIRONMENT>-efs01, for e.g.,
# aws --region ap-southeast-2 efs describe-file-systems --query 'FileSystems[*].[Name, FileSystemId]' --output text | grep roiergasias-demo-efs01

# output helm chart template for roiergasias aws auto
helm template \
  -n roiergasias \
  -f ./helm/roiergasias-aws-auto/values-secret.yaml \
  roiergasias-aws-auto ./helm/roiergasias-aws-auto >machine-learning-aws-auto-manifest.yaml

# explore the contents of the machine-learning-aws-auto-manifest.yaml
cat machine-learning-aws-auto-manifest.yaml

# apply the manifest
kubectl apply -f machine-learning-aws-auto-manifest.yaml

# browse pod created by the job
kubectl get pods -n roiergasias

# check pod logs for the output
kubectl logs roiergasias-aws-<STRING_FROM_PREVIOUS_STEP> -n roiergasias

# delete the manifest
kubectl delete -f machine-learning-aws-auto-manifest.yaml

# uninstall the operator
helm uninstall roiergasias-operator
```


## Steps to de-provision AWS infrastructure
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# change to the infra/aws directory
cd infra/aws

# make sure your s3 bucket is empty
# assumes 'roiergasias' as <PREFIX> and 'demo' as <ENVIRONMENT> values
aws s3 rm s3://roiergasias-demo-s3b01 --recursive

# execute infrastructure de-provisioning command
terraform destroy -var-file=values-secret.tfvars
# sometimes it fails the first time. So, after a delay of 5-10 mins you can repeat the above command as many times until it succeeds
```