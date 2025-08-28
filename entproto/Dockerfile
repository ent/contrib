FROM golang:1.23.4
RUN apt update && apt install unzip
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_64.zip && \
 unzip -o protoc-3.19.4-linux-x86_64.zip -d ./proto && \
 chmod 755 -R ./proto/bin && \
 cp ./proto/bin/protoc /usr/bin/ && \
 cp -R ./proto/include/* /usr/include/
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0 \
&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
