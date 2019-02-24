FROM alpine:latest
COPY dist/linux_amd64/kafkactl /usr/local/bin/kafkactl
ENTRYPOINT ["kafkactl"]