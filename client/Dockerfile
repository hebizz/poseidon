FROM ubuntu:18.04

RUN apt-get update && apt-get install -y netplan

WORKDIR /etc/poseidon/
COPY main main
COPY config.yaml config.yaml
COPY templates templates

ENTRYPOINT ["./main"]
#docker run -v /etc/:/etc ...
