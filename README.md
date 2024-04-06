# scgb (self-compiling Go binary)

> [!CAUTION] This is a PoC and not intended for production use.

Install once:

```sh
./scripts/install.sh
```

Enjoy a self-compiling Go binary as you change the source code:

```sh
$ scgb

running main...

$ scgb

2024/04/05 15:57:44 recompiling...
2024/04/05 15:57:44 re-running ~/go/bin/scgb
running changed main...
```

## briefly about how it's different from other solutions

### [fsnotify](https://github.com/fsnotify/fsnotify)

Suitable for long-running processes that need to watch for file changes but not
for one-off command-line tools.

### [minio/selfupdate](https://github.com/minio/selfupdate)

Replaces the binary with a new one taken from a remote source. It's not not
self-compiling.
