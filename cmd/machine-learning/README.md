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


## Steps to manually run go workflow
``` SH
# clone to a local git directory, if not already done so
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

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


## Steps to run go workflow via docker compose
``` SH
# clone to a local git directory, if not already done so
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

# change to the local git directory
cd kubernetes-operator-roiergasias

# copy kaggle api credentials from ~/.kaggle
cp ~/.kaggle/kaggle.json .

# set execute permissions to go main binary and python scripts
chmod +x ../main ./*.py

# build docker image
docker build -t roiergasias:latest .

# change to the cmd/machine-leaning directory
cd cmd/machine-learning

# run docker compose
docker-compose up
```


## Steps to manually run python scripts
``` SH
# clone to a local git directory, if not already done so
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

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
chmod +x *.py

# process data from first argument file saving output in second argument file
./process-data.py ./weatherAUS.csv ./processed-weatherAUS.csv

# train ml model from first argument file and saving model output in second argument file
./train-model.py ./processed-weatherAUS.csv ./ml-model.joblib

# evaluate ml model by reading processed data and model from first and second argument files
./evaluate-model.py ./processed-weatherAUS.csv ./ml-model.joblib
```
