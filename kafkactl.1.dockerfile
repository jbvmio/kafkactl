FROM scratch
COPY kafkactl /
ENTRYPOINT ["/kafkactl"]