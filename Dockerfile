FROM    golang:latest

WORKDIR /
RUN git clone --single-branch --branch dev https://github.com/jbvmio/kafkactl.git
WORKDIR /kafkactl
RUN make docker
ENTRYPOINT [ "/kafkactl" ]
