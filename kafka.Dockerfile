FROM confluentinc/cp-kafka:7.6.0

USER root

WORKDIR /usr/app

RUN yum install -y dos2unix

COPY kafka-create-topics.sh .

RUN dos2unix kafka-create-topics.sh

RUN chmod +x ./kafka-create-topics.sh

CMD ["/bin/bash", "-c", "/etc/confluent/docker/run && ./kafka-create-topics.sh"]
