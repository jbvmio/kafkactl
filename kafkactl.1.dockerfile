FROM alpine:latest
ADD kafkactl /usr/local/bin/kafkactl
ENTRYPOINT ["kafkactl"]