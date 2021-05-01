#!/usr/bin/env python
# coding: utf-8

# # Rain prediction in Australia
# **Task type:** Training and saving model
# **ML algorithm used:** Random Forest Classifier

import sys
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline
from sklearn import preprocessing
from sklearn.preprocessing import StandardScaler
from joblib import dump

print(f'Running script: {sys.argv[0]} with input processed data file: {sys.argv[1]}')

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

# **Please uncomment this part of code to use grid search for hyperparameter tuning for the model. The model below uses the outcome of the GridSearch operation with best parameters.**
# from sklearn.model_selection import GridSearchCV
# param_grid = {
#    'n_estimators': [100, 200],
#    'max_features': ['auto'],
#    'max_depth' : [4,5,8,10],
#    'criterion' :['gini', 'entropy']
# }
# RFC = RandomForestClassifier()
# cv_RFC = GridSearchCV(estimator=RFC, param_grid=param_grid, cv=2)
# cv_RFC.fit(X_train, y_train)
# cv_RFC.best_params_
# sorted(zip(cv_RFC.best_estimator_.feature_importances_,ohe.columns))

pipe = Pipeline([('scaler', StandardScaler()), ('RFC', RandomForestClassifier(criterion='gini', 
                                                                              max_depth=10, 
                                                                              max_features='auto',
                                                                              n_estimators=200))])
pipe.fit(X_train, y_train)

print(f'Saving trained model to file: {sys.argv[2]}')
dump(pipe, sys.argv[2])

# **Please uncomment this part of code to evaluate the model here only
# from sklearn.model_selection import cross_val_score
# from sklearn.metrics import accuracy_score, f1_score
# print(f'Score: {pipe.score(X_train, y_train)}')

# print(f'Cross validation score: {cross_val_score(pipe, X, y, cv=3)}')

# y_pred = pipe.predict(X_test)
# print(f'Accuracy score: {accuracy_score(y_test, y_pred)}')
# print(f'F1 score: {f1_score(y_test, y_pred)}')
