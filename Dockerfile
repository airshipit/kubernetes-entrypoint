ARG GO_IMAGE=docker.io/golang:1.16.4-buster
ARG RELEASE_IMAGE=ubuntu:focal
FROM ${GO_IMAGE} as builder

SHELL [ "/bin/bash", "-cex" ]
ADD . /usr/src/kubernetes-entrypoint
WORKDIR /usr/src/kubernetes-entrypoint
ENV GO111MODULE=on

RUN make get-modules

ARG MAKE_TARGET=build
RUN make ${MAKE_TARGET}

FROM ${RELEASE_IMAGE} as release
COPY --from=builder /usr/src/kubernetes-entrypoint/bin/kubernetes-entrypoint /usr/local/bin/kubernetes-entrypoint

RUN apt-get update
RUN apt-get install -y --no-install-recommends coreutils

USER 65534
ENTRYPOINT [ "/usr/local/bin/kubernetes-entrypoint" ]
