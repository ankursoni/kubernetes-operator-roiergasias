FROM python:3.9.2

RUN curl -fsSL https://golang.org/dl/go1.16.3.linux-amd64.tar.gz -o go1.16.3.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.3.linux-amd64.tar.gz && rm go1.16.3.linux-amd64.tar.gz && \
    export PATH=$PATH:/usr/local/go/bin && echo PATH=$PATH:/usr/local/go/bin >> /root/.bashrc && \
    pip install kaggle pandas sklearn joblib

COPY ./kaggle.json /root/.kaggle/
COPY ./main /root/cmd/

COPY ./machine-learning/*.yaml /root/cmd/machine-learning/
COPY ./machine-learning/*.py /root/cmd/machine-learning/
