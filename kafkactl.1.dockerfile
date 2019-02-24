FROM amd64/alpine:latest
COPY kafkactl /usr/local/bin/kafkactl
ENTRYPOINT ["kafkactl"]