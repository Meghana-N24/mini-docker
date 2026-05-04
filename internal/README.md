# 🐳 mini-docker

A minimal container runtime built from scratch in Go, 
implementing core Docker concepts using Linux primitives.

## What This Project Does

`mini-docker` runs a command inside an isolated environment 
using the same Linux kernel features that power real Docker:

- **Namespaces** — process, mount, and hostname isolation
- **chroot** — filesystem isolation using Alpine Linux rootfs  
- **cgroups v2** — memory and PID resource limits
- **CLI interface** — `mydocker run <command>`

## Demo

```bash
$ sudo ./mydocker run /bin/sh
[mydocker] Starting container...
[mydocker] Host PID: 1234
[mydocker] Memory limit: 256 MB
[mydocker] PID limit: 20
[mydocker] Container PID: 1
/ # hostname
mini-docker-container
/ # ps aux
PID   USER     COMMAND
    1 root     /proc/self/exe child /bin/sh
    6 root     /bin/sh
    7 root     ps aux
/ # ls /home
(empty — host filesystem is not visible)

mydocker run /bin/sh
       │
       ▼
  [run() function]
  Creates new namespaces:
  • CLONE_NEWPID  → isolated process tree
  • CLONE_NEWUTS  → isolated hostname
  • CLONE_NEWNS   → isolated mounts
       │
       ▼
  Applies cgroup limits:
  • memory.max = 256MB
  • pids.max   = 20
       │
       ▼
  [child() function] ← runs INSIDE namespaces
  • syscall.Chroot("rootfs")  → Alpine filesystem
  • Mount /proc               → ps/top work correctly
  • exec user command         → /bin/sh

  mini-docker/
├── cmd/mydocker/
│   └── main.go              # CLI entry point
├── internal/
│   ├── container/           # container lifecycle
│   └── cgroups/
│       └── cgroups.go       # resource limiting
├── rootfs/                  # Alpine Linux filesystem
└── README.md