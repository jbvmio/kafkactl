FROM alpine:latest
COPY kafkactl /usr/local/bin/kafkactl
ENTRYPOINT ["kafkactl"]