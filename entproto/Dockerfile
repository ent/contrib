FROM golang:1.17.3
RUN apt update && apt install unzip
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip && \
 unzip -o protoc-3.17.3-linux-x86_64.zip -d ./proto && \
 chmod 755 -R ./proto/bin && \
 cp ./proto/bin/protoc /usr/bin/ && \
 cp -R ./proto/include/* /usr/include/
RUN  go get google.golang.org/protobuf/cmd/protoc-gen-go && \
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
