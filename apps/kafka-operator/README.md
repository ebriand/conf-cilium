# Kafka deployment

## Deploy the kafka operator

kubectl apply -f strimzi-cluster-operator-0.11.1.yaml

## Deploy a kafka cluster

kubectl apply -f kafka-ephemeral.yaml

## Test

When all pods (zookeeper and kafka) are ready, you can try sending and receiving messages.
Sending will offer you a prompt to send messages, each time you hit `Return` a new message will be delivered
Receiving will display all messages from the beginning.

Sending:

```shell
kubectl -n default run kafka-producer -ti --image=strimzi/kafka:0.11.1-kafka-2.1.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --broker-list my-cluster-kafka-bootstrap:9092 --topic my-topic
```

Receiving:

```shell
kubectl -n default run kafka-consumer -ti --image=strimzi/kafka:0.11.1-kafka-2.1.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning
```
