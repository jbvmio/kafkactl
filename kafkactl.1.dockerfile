FROM alpine:latest
ADD dist/linux_amd64/kafkactl /usr/local/bin/kafkactl
ENTRYPOINT ["kafkactl"]