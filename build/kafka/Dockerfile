FROM alpine:latest AS fetcher

ARG scala_version=2.13
ARG kafka_version=3.3.1

RUN apk update
RUN apk add wget

WORKDIR /var/tmp

RUN wget https://dlcdn.apache.org/kafka/$kafka_version/kafka_$scala_version-$kafka_version.tgz

RUN tar -xzf kafka_$scala_version-$kafka_version.tgz

#######################################################################

FROM eclipse-temurin:17.0.3_7-jre

ARG scala_version=2.13
ARG kafka_version=3.3.1

ENV KAFKA_HOME=/opt/kafka

ENV PATH=${PATH}:${KAFKA_HOME}/bin

WORKDIR ${KAFKA_HOME}

COPY --from=fetcher /var/tmp/kafka_$scala_version-$kafka_version ${KAFKA_HOME}

RUN chmod a+x ${KAFKA_HOME}/bin/*.sh

RUN bin/kafka-storage.sh format -t $(bin/kafka-storage.sh random-uuid) -c config/kraft/server.properties

ENTRYPOINT ["/bin/bash"]

CMD ["bin/kafka-server-start.sh", "config/kraft/server.properties"]
