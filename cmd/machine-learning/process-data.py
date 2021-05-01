#!/usr/bin/env python
# coding: utf-8

# # Rain prediction in Australia
# **Task type:** Processing data

import sys
import pandas as pd

print(f'Running script: {sys.argv[0]} with input data file: {sys.argv[1]}')

df = pd.read_csv(sys.argv[1])

zeros_cnt = df.isnull().sum().sort_values(ascending=False)
percent_zeros = (df.isnull().sum() / df.isnull().count()).sort_values(ascending=False)
missing_data = pd.concat([zeros_cnt, percent_zeros], axis=1, keys=['Total', 'Percent'])
missing_data

dropList = list(missing_data[missing_data['Percent'] > 0.15].index)
dropList
df.drop(dropList, axis=1, inplace=True)

print(f'Saving processed data to file: {sys.argv[2]}')
df.to_csv(sys.argv[2])
