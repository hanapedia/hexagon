# Kafka Strimzi Operator installation
[Using Helm](https://strimzi.io/docs/operators/latest/deploying.html#deploying-cluster-operator-helm-chart-str)

[Deploy Kafka cluster](https://strimzi.io/docs/operators/latest/deploying.html#deploying-kafka-cluster-str)

Install:
```sh
helm install strimzi-cluster-operator oci://quay.io/strimzi-helm/strimzi-kafka-operator \
    -n kafka --create-namespace
```

Deploy Kafka Cluster:
```sh
kubectl apply -f kafka.yaml -n kafka
```
`kafka.yaml` contains manifest to deploy single node persistent kafka cluster.
