# DevKit

DevKit is a concurrent, cache-efficient, and Dockerfile-agnostic builder toolkit.

Report issues at https://github.com/khulnasoft/devkit

Join `#devkit` channel on [Docker Community Slack](https://dockr.ly/comm-slack)

# Tags

### Latest stable release

- [`v0.10.0`, `latest`](https://github.com/khulnasoft/devkit/blob/v0.10.0/Dockerfile)

- [`v0.10.0-rootless`, `rootless`](https://github.com/khulnasoft/devkit/blob/v0.10.0/Dockerfile) (see [`docs/rootless.md`](https://github.com/khulnasoft/devkit/blob/master/docs/rootless.md) for usage)

### Development build from master branch

- [`master`](https://github.com/khulnasoft/devkit/blob/master/Dockerfile)

- [`master-rootless`](https://github.com/khulnasoft/devkit/blob/master/Dockerfile)


Binary releases and changelog can be found from https://github.com/khulnasoft/devkit/releases

# Usage


To run daemon in a container:

```bash
docker run -d --name devkitd --privileged khulnasoft/devkit:latest
export DEVKIT_HOST=docker-container://devkitd
buildctl build --help
```

See https://github.com/khulnasoft/devkit#devkit for general DevKit usage instructions


## Docker Buildx

[Buildx](https://github.com/docker/buildx) uses the latest stable image by default. To set a custom DevKit image version use `--driver-opt`:

```bash
docker buildx create --driver-opt image=khulnasoft/devkit:master --use
```


## Rootless

For Rootless deployments, see [`docs/rootless.md`](https://github.com/khulnasoft/devkit/blob/master/docs/rootless.md)


## Kubernetes

For Kubernetes deployments, see [`examples/kubernetes`](https://github.com/khulnasoft/devkit/tree/master/examples/kubernetes)


## Daemonless

To run the client and an ephemeral daemon in a single container ("daemonless mode"):

```bash
docker run \
    -it \
    --rm \
    --privileged \
    -v /path/to/dir:/tmp/work \
    --entrypoint buildctl-daemonless.sh \
    khulnasoft/devkit:master \
        build \
        --frontend dockerfile.v0 \
        --local context=/tmp/work \
        --local dockerfile=/tmp/work
```

Rootless mode:

```bash
docker run \
    -it \
    --rm \
    --security-opt seccomp=unconfined \
    --security-opt apparmor=unconfined \
    -e DEVKITD_FLAGS=--oci-worker-no-process-sandbox \
    -v /path/to/dir:/tmp/work \
    --entrypoint buildctl-daemonless.sh \
    khulnasoft/devkit:master-rootless \
        build \
        --frontend \
        dockerfile.v0 \
        --local context=/tmp/work \
        --local dockerfile=/tmp/work
```
