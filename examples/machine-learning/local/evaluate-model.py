#!/usr/bin/env python
# coding: utf-8

# # Rain prediction in Australia
# **Task type:** Load and evaluate model
# **Metrics:** Accuracy

import sys
from joblib import load
import pandas as pd
from sklearn import preprocessing
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline
from sklearn.model_selection import cross_val_score
from sklearn.metrics import accuracy_score, f1_score

print(f'Running script: {sys.argv[0]} with input processed data file: {sys.argv[1]} and input model file: {sys.argv[2]}')

df = pd.read_csv(sys.argv[1])
df['Location'].unique()
df.drop(['Date'], axis=1, inplace=True)
df.drop(['Location'], axis=1, inplace=True)

ohe = pd.get_dummies(data=df, columns=['WindGustDir','WindDir9am','WindDir3pm'])
ohe['RainToday'] = df['RainToday'].astype(str)
ohe['RainTomorrow'] = df['RainTomorrow'].astype(str)
lb = preprocessing.LabelBinarizer()
ohe['RainToday'] = lb.fit_transform(ohe['RainToday'])
ohe['RainTomorrow'] = lb.fit_transform(ohe['RainTomorrow'])
ohe = ohe.dropna()
y = ohe['RainTomorrow']
X = ohe.drop(['RainTomorrow'], axis=1)

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.3, random_state=0)

pipe = load(sys.argv[2])

print(f'Score: {pipe.score(X_train, y_train)}')

print(f'Cross validation score: {cross_val_score(pipe, X, y, cv=3)}')

y_pred = pipe.predict(X_test)
print(f'Accuracy score: {accuracy_score(y_test, y_pred)}')
print(f'F1 score: {f1_score(y_test, y_pred)}')
