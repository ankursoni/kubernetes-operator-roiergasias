FROM python:3.9.2

RUN pip install kaggle pandas sklearn joblib

WORKDIR /root
COPY ./kaggle.json ./.kaggle/
COPY ./main ./cmd/
