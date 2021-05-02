# Process data, train ml model & evaluate ml model


## Source of inspiration
Converted the following jupyter notebook to Python scripts:  
https://www.kaggle.com/ilyapozdnyakov/rain-in-australia-precision-recall-curves-viz#Rain-prediction-in-Australia


## Install pre-requisites

### - Install [Python 3](https://www.python.org/downloads/)
Also, install the following *pip* packages:
``` SH
pip install pandas sklearn joblib
# or,
# pip3 install pandas sklearn joblib
```

### - Install [Kaggle CLI](https://github.com/Kaggle/kaggle-api)

### Make sure kaggle is configured with api key in ~/.kaggle/kaggle.json
[Kaggle API Credentials](https://github.com/Kaggle/kaggle-api#api-credentials)


## Steps to run manually go workflow
``` SH
# clone to a local git directory
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

# change to the local git directory
cd kubernetes-operator-roiergasias

# download go module dependencies to local cache
go mod download

# set execute permissions to go main binary
chmod +x cmd/main

# run the machine learning workflow
./cmd/main ./cmd/machine-learning/machine-learning.yaml
# for mac osx
./cmd/main-osx ./cmd/machine-learning/machine-learning.yaml
```


## Steps to run manually python scripts
``` SH
# change to the cmd/machine-leaning directory
cd cmd/machine-learning

# download dataset from kaggle
kaggle datasets download jsphyg/weather-dataset-rattle-package -o -f weatherAUS.csv

# unzip and delete the downloaded zip file
unzip -o weatherAUS.csv.zip
rm -f weatherAUS.csv.zip

# set execute permissions to python script files
chmod +x *.py

# process data from first argument file saving output in second argument file
./process-data.py ./weatherAUS.csv ./processed-weatherAUS.csv

# train ml model from first argument file and saving model output in second argument file
./train-model.py ./processed-weatherAUS.csv ./ml-model.joblib

# evaluate ml model by reading processed data and model from first and second argument files
./evaluate-model.py ./processed-weatherAUS.csv ./ml-model.joblib
```


## Steps to build and run docker image
``` SH
# change to the local git directory
cd kubernetes-operator-roiergasias

# copy kaggle api credentials from ~/.kaggle
cp ~/.kaggle/kaggle.json cmd/

# set execute permissions to go main binary and python scripts
chmod +x cmd/main cmd/machine-learning/*.py

docker build -t docker.io/<REPOSITORY_NAME>/roiergasias:latest ./cmd
# where, <REPOSITORY_NAME> is the docker hub's repository name or username, for e.g.,
# docker build -t docker.io/ankursoni/roiergasias:latest ./cmd

docker run -it --name wf docker.io/<REPOSITORY_NAME>/roiergasias:latest bash
# where, <REPOSITORY_NAME> is the docker hub's repository name or username, for e.g.,
# docker run -it --name wf docker.io/ankursoni/roiergasias:latest bash

# run the machine learning workflow
cd ~
./cmd/main ./cmd/machine-learning/machine-learning.yaml
```