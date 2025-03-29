FROM europe-docker.pkg.dev/kubemq/images/gobuilder as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq-io/kubemq-targets
WORKDIR $GOPATH/github.com/kubemq-io/kubemq-targets
RUN CGO_ENABLED=0 go build  -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.version=$VERSION" -o kubemq-targets-run .
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
MAINTAINER KubeMQ info@kubemq.io
LABEL name="KubeMQ Target Connectors" \
      maintainer="info@kubemq.io" \
      vendor="kubemq.io" \
      version="1.9.0" \
      release="stable" \
      summary="KubeMQ Targets connects KubeMQ Message Broker with external systems and cloud services" \
      description="KubeMQ Targets allows us to build a message-based microservices architecture on Kubernetes with minimal efforts and without developing connectivity interfaces between KubeMQ Message Broker and external systems such as databases, cache, messaging, and REST-base APIs"
COPY licenses /licenses
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN mkdir /kubemq-connector
COPY --from=builder $GOPATH/github.com/kubemq-io/kubemq-targets/kubemq-targets-run ./kubemq-connector
COPY --from=builder $GOPATH/github.com/kubemq-io/kubemq-targets/default_config.yaml ./kubemq-connector/config.yaml
RUN chown -R 1001:root  /kubemq-connector && chmod g+rwX  /kubemq-connector
WORKDIR kubemq-connector
USER 1001
CMD ["./kubemq-targets-run"]
