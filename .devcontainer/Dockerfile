# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.148.1/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.15, 1.14
ARG VARIANT="1"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends protobuf-compiler

# [Optional] Uncomment the next line to use go get to install anything else you need
RUN su vscode -c "go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 \
  && go install entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc@master"

ENV PROTOBUF_VER=3.17.3
ENV PROTOBUF_ZIP=protoc-${PROTOBUF_VER}-linux-x86_64.zip
RUN curl -OL https://github.com/google/protobuf/releases/download/v${PROTOBUF_VER}/${PROTOBUF_ZIP} \
  && unzip $PROTOBUF_ZIP -d protoc3 \
  && mv protoc3/bin/* /usr/local/bin/ \
  && mv protoc3/include/* /usr/local/include/ \
  && rm -rf protoc3 $PROTOBUF_ZIP
