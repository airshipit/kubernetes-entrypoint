ARG GO_IMAGE=docker.io/golang:1.16.4-buster
ARG FROM=ubuntu:20.04
FROM ${GO_IMAGE} as builder

SHELL [ "/bin/bash", "-cex" ]
ADD . /usr/src/kubernetes-entrypoint
WORKDIR /usr/src/kubernetes-entrypoint
ENV GO111MODULE=on

RUN make get-modules

ARG MAKE_TARGET=build
RUN make ${MAKE_TARGET}



FROM ${FROM} as release

LABEL org.opencontainers.image.authors='airship-discuss@lists.airshipit.org, irc://#airshipit@freenode' \
      org.opencontainers.image.url='https://airshipit.org' \
      org.opencontainers.image.documentation='https://docs.airshipit.org/kubernetes-entrypoint' \
      org.opencontainers.image.source='https://opendev.org/airship/kubernetes-entrypoint' \
      org.opencontainers.image.vendor='The Airship Authors' \
      org.opencontainers.image.licenses='Apache-2.0'

ENV DEBIAN_FRONTEND noninteractive
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

COPY --from=builder /usr/src/kubernetes-entrypoint/bin/kubernetes-entrypoint /usr/local/bin/kubernetes-entrypoint

RUN apt update \
    && apt install -y --no-install-recommends coreutils

USER 65534
ENTRYPOINT [ "/usr/local/bin/kubernetes-entrypoint" ]
