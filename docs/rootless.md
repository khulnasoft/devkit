# Rootless mode

Rootless mode allows running DevKit daemon as a non-root user.

## Distribution-specific hint
Using Ubuntu kernel is recommended.
### Container-Optimized OS from Google
Make sure to have an `emptyDir` volume below:
```yaml
spec:
  containers:
    - name: devkitd
      volumeMounts:
        # Dockerfile has `VOLUME /home/user/.local/share/devkit` by default too,
        # but the default VOLUME does not work with rootless on Google's Container-Optimized OS
        # as it is mounted with `nosuid,nodev`.
        # https://github.com/khulnasoft/devkit/issues/879#issuecomment-1240347038
        - mountPath: /home/user/.local/share/devkit
          name: devkitd
  volumes:
    - name: devkitd
      emptyDir: {}
```

See also the [example manifests](#Kubernetes).

### Bottlerocket OS

Needs to set the max user namespaces to a positive integer, through the [API settings](https://github.com/bottlerocket-os/bottlerocket#kernel-settings):

```toml
[settings.kernel.sysctl]
"user.max_user_namespaces" = "16384"
```

See [`../examples/eksctl/bottlerocket.yaml`](../examples/eksctl/bottlerocket.yaml) for an example to configure a Node Group in EKS.

<details>
<summary>Old distributions</summary>

<p>

### Debian GNU/Linux 10
Add `kernel.unprivileged_userns_clone=1` to `/etc/sysctl.conf` (or `/etc/sysctl.d`) and run `sudo sysctl -p`.
This step is not needed for Debian GNU/Linux 11 and later.

### RHEL/CentOS 7
Add `user.max_user_namespaces=28633` to `/etc/sysctl.conf` (or `/etc/sysctl.d`) and run `sudo sysctl -p`.
This step is not needed for RHEL/CentOS 8 and later.

### Fedora, before kernel 5.13
You may have to disable SELinux, or run DevKit with `--oci-worker-snapshotter=fuse-overlayfs`.

</p>
</details>

## Known limitations
* Using the `overlayfs` snapshotter requires kernel >= 5.11 or Ubuntu kernel.
  On kernel >= 4.18, the `fuse-overlayfs` snapshotter is used instead of `overlayfs`.
  On kernel < 4.18, the `native` snapshotter is used.
* Network mode is always set to `network.host`.

## Running DevKit in Rootless mode (OCI worker)

[RootlessKit](https://github.com/rootless-containers/rootlesskit/) needs to be installed.

```console
$ rootlesskit devkitd
```

```console
$ buildctl --addr unix:///run/user/$UID/devkit/devkitd.sock build ...
```

To isolate DevKit daemon's network namespace from the host (recommended):
```console
$ rootlesskit --net=slirp4netns --copy-up=/etc --disable-host-loopback devkitd
```

## Running DevKit in Rootless mode (containerd worker)

[RootlessKit](https://github.com/rootless-containers/rootlesskit/) needs to be installed.

Run containerd in rootless mode using rootlesskit following [containerd's document](https://github.com/containerd/containerd/blob/main/docs/rootless.md).

```
$ containerd-rootless.sh
```

Then let devkitd join the same namespace as containerd.

```
$ containerd-rootless-setuptool.sh nsenter -- devkitd --oci-worker=false --containerd-worker=true --containerd-worker-snapshotter=native
```

## Troubleshooting

### Error related to `overlayfs`
Try running `devkitd` with `--oci-worker-snapshotter=fuse-overlayfs`:

```console
$ rootlesskit devkitd --oci-worker-snapshotter=fuse-overlayfs
```

### Error related to `fuse-overlayfs`
Try running `devkitd` with `--oci-worker-snapshotter=native`:

```console
$ rootlesskit devkitd --oci-worker-snapshotter=native
```

### Error related to `newuidmap` or `/etc/subuid`
See https://rootlesscontaine.rs/getting-started/common/subuid/

### Error `Options:[rbind ro]}]: operation not permitted`
Make sure to mount an `emptyDir` volume on `/home/user/.local/share/devkit` .

### Error `fork/exec /proc/self/exe: no space left on device` with `level=warning msg="/proc/sys/user/max_user_namespaces needs to be set to non-zero."`
Run `sysctl -w user.max_user_namespaces=N` (N=positive integer, like 63359) on the host nodes.

See [`../examples/kubernetes/sysctl-userns.privileged.yaml`](../examples/kubernetes/sysctl-userns.privileged.yaml).

### Error `mount proc:/proc (via /proc/self/fd/6), flags: 0xe: operation not permitted`
This error is known to happen when DevKit is executed in a container without the `--oci-worker-no-sandbox` flag.
Make sure that `--oci-worker-no-process-sandbox` is specified (See [below](#docker)).

## Containerized deployment

### Kubernetes
See [`../examples/kubernetes`](../examples/kubernetes).

### Docker

```console
$ docker run \
  --name devkitd \
  -d \
  --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined \
  --device /dev/fuse \
  khulnasoft/devkit:rootless --oci-worker-no-process-sandbox
$ buildctl --addr docker-container://devkitd build ...
```

If you don't mind using `--privileged` (almost safe for rootless), the `docker run` flags can be shorten as follows:

```console
$ docker run --name devkitd -d --privileged khulnasoft/devkit:rootless
```

#### About `--device /dev/fuse`
Adding `--device /dev/fuse` to the `docker run` arguments is required only if you want to use `fuse-overlayfs` snapshotter.

#### About `--oci-worker-no-process-sandbox`

By adding `--oci-worker-no-process-sandbox` to the `devkitd` arguments, DevKit can be executed in a container without adding `--privileged` to `docker run` arguments.
However, you still need to pass `--security-opt seccomp=unconfined --security-opt apparmor=unconfined` to `docker run`.

Note that `--oci-worker-no-process-sandbox` allows build executor containers to `kill` (and potentially `ptrace` depending on the seccomp configuration) an arbitrary process in the DevKit daemon container.

To allow running rootless `devkitd` without `--oci-worker-no-process-sandbox`, run `docker run` with `--security-opt systempaths=unconfined`. (For Kubernetes, set `securityContext.procMount` to `Unmasked`.)

The `--security-opt systempaths=unconfined` flag disables the masks for the `/proc` mount in the container and potentially allows reading and writing dangerous kernel files, but it is safe when you are running `devkitd` as non-root.

### Change UID/GID

The `khulnasoft/devkit:rootless` image has the following UID/GID configuration:

Actual ID (shown in the host and the DevKit daemon container)| Mapped ID (shown in build executor containers)
----------|----------
1000      | 0
100000    | 1
...       | ...
165535    | 65536

```
$ docker exec devkitd id
uid=1000(user) gid=1000(user)
$ docker exec devkitd ps aux
PID   USER     TIME   COMMAND
    1 user       0:00 rootlesskit devkitd --addr tcp://0.0.0.0:1234
   13 user       0:00 /proc/self/exe devkitd --addr tcp://0.0.0.0:1234
   21 user       0:00 devkitd --addr tcp://0.0.0.0:1234
   29 user       0:00 ps aux
$ docker exec cat /etc/subuid
user:100000:65536
```

To change the UID/GID configuration, you need to modify and build the DevKit image manually.
```
$ vi Dockerfile
$ make images
$ docker run ... khulnasoft/devkit:local-rootless ...
```

