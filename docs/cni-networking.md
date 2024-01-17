# CNI networking

It can be useful to use a bridge network for your builder if for example you
encounter a network port contention during multiple builds. If you're using
the DevKit image, CNI is not [(yet)](https://github.com/khulnasoft/devkit/issues/28)
available in it.

But you can create your own DevKit image with CNI support:

```dockerfile
ARG DEVKIT_VERSION=v0.9.3
ARG CNI_VERSION=v1.0.1

FROM --platform=$BUILDPLATFORM alpine AS cni-plugins
RUN apk add --no-cache curl
ARG CNI_VERSION
ARG TARGETOS
ARG TARGETARCH
WORKDIR /opt/cni/bin
RUN curl -Ls https://github.com/containernetworking/plugins/releases/download/$CNI_VERSION/cni-plugins-$TARGETOS-$TARGETARCH-$CNI_VERSION.tgz | tar xzv

FROM khulnasoft/devkit:${DEVKIT_VERSION}
ARG DEVKIT_VERSION
RUN apk add --no-cache iptables
COPY --from=cni-plugins /opt/cni/bin /opt/cni/bin
ADD https://raw.githubusercontent.com/khulnasoft/devkit/${DEVKIT_VERSION}/hack/fixtures/cni.json /etc/devkit/cni.json
```

Here we use the [CNI config for integration tests in DevKit](../hack/fixtures/cni.json),
but feel free to use your own config.
