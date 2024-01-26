# Formatting timestamps via kafka

A basic illustration on how to process data via kafka.

# Prerequisites
[Local dev cluster](https://github.com/DenktMit-eG/local-dev-cluster)

# Development



```shell
kubectl create secret tls root-ca-secret \                              
  --cert="$(mkcert -CAROOT)/rootCA.pem" \
  --key="$(mkcert -CAROOT)/rootCA-key.pem"
secret/root-ca-secret created

```

Additional /etc/hosts entry

```text
127.0.0.1    websocket.local.lgc
```

```shell
# ./mvnw clean spring-boot:build-image 
kind load docker-image ghcr.io/jschlarb/timestamp/producer:0.0.1-SNAPSHOT -n lgc
kind load docker-image ghcr.io/jschlarb/timestamp/formatter:0.0.1-SNAPSHOT -n lgc 
kind load docker-image ghcr.io/jschlarb/timestamp/websocket:0.0.1-SNAPSHOT -n lgc

helm upgrade --install producer charts/base -f charts/base/values.producer.yaml --atomic
helm upgrade --install formatter charts/base -f charts/base/values.formatter.yaml --atomic
helm upgrade --install websocket charts/base -f charts/base/values.websocket.yaml --atomic
```
