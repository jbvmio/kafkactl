FROM alpine:latest
ADD kafkactl.linux /usr/local/bin/kafkactl
ENTRYPOINT ["kafkactl"]