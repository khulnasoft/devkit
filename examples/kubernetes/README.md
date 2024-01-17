# Kubernetes manifests for DevKit

This directory contains Kubernetes manifests for `Pod`, `Deployment` (with `Service`), `StatefulSet`, and `Job`.
* `Pod`: good for quick-start
* `Deployment` + `Service`: good for random load balancing with registry-side cache
* `StateFulset`: good for client-side load balancing, without registry-side cache
* `Job`: good if you don't want to have daemon pods

Using Rootless mode (`*.rootless.yaml`) is recommended because Rootless mode image is executed as non-root user (UID 1000) and doesn't need `securityContext.privileged`.
See [`../../docs/rootless.md`](../../docs/rootless.md).

See also ["Building Images Efficiently And Securely On Kubernetes With DevKit" (KubeCon EU 2019)](https://kccnceu19.sched.com/event/MPX5).

## `Pod`

```console
$ kubectl apply -f pod.rootless.yaml
$ buildctl \
  --addr kube-pod://devkitd \
  build --frontend dockerfile.v0 --local context=/path/to/dir --local dockerfile=/path/to/dir
```

If rootless mode doesn't work, try `pod.privileged.yaml`.

:warning: `kube-pod://` connection helper requires Kubernetes role that can access `pods/exec` resources. If `pods/exec` is not accessible, use `Service` instead (See below).

## `Deployment` + `Service`

Setting up mTLS is highly recommended.

`./create-certs.sh SAN [SAN...]` can be used for creating certificates.
```console
$ ./create-certs.sh 127.0.0.1
```

The daemon certificates is created as `Secret` manifest named `devkit-daemon-certs`.
```console
$ kubectl apply -f .certs/devkit-daemon-certs.yaml
```

Apply the `Deployment` and `Service` manifest:
```console
$ kubectl apply -f deployment+service.rootless.yaml
$ kubectl scale --replicas=10 deployment/devkitd
```

Run `buildctl` with TLS client certificates:
```console
$ kubectl port-forward service/devkitd 1234
$ buildctl \
  --addr tcp://127.0.0.1:1234 \
  --tlscacert .certs/client/ca.pem \
  --tlscert .certs/client/cert.pem \
  --tlskey .certs/client/key.pem \
  build --frontend dockerfile.v0 --local context=/path/to/dir --local dockerfile=/path/to/dir
```

## `StatefulSet`
`StatefulSet` is useful for consistent hash mode.

```console
$ kubectl apply -f statefulset.rootless.yaml
$ kubectl scale --replicas=10 statefulset/devkitd
$ buildctl \
  --addr kube-pod://devkitd-4 \
  build --frontend dockerfile.v0 --local context=/path/to/dir --local dockerfile=/path/to/dir
```

See [`./consistenthash`](./consistenthash) for how to use consistent hashing.

## `Job`

```console
$ kubectl apply -f job.rootless.yaml
```

To push the image to the registry, you also need to mount `~/.docker/config.json`
and set `$DOCKER_CONFIG` to `/path/to/.docker` directory.
