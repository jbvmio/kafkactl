FROM scratch
ADD kafkactl /usr/bin/kafkactl
ENTRYPOINT ["/kafkactl"]