#!/bin/bash

# create topics
kafka-topics --create --topic $TOPIC_NAME --bootstrap-server kafka:9092
echo "topic $TOPIC_NAME was created"
