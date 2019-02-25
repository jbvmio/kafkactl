FROM    golang:latest

WORKDIR /
RUN apt-get update && apt-get install -y bc && \
    git clone --single-branch --branch dev https://github.com/jbvmio/kafkactl.git
WORKDIR /kafkactl
RUN make docker
ENTRYPOINT [ "/kafkactl" ]
