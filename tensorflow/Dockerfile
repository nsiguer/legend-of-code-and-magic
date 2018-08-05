FROM golang:1.9

RUN apt-get update && apt-get install -y \
    curl

RUN TF_TYPE="cpu" TARGET_DIRECTORY='/usr/local' ; curl -L \
     "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-${TF_TYPE}-$(go env GOOS)-x86_64-1.9.0.tar.gz" | \
     tar -C $TARGET_DIRECTORY -xz

RUN go get github.com/tensorflow/tensorflow/tensorflow/go
