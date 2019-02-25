FROM ubuntu:18.04
COPY kafkactl /usr/local/bin/kafkactl
ENTRYPOINT [ "kafkactl" ]