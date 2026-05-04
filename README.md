# 🐳 mini-docker

A minimal container runtime built from scratch in Go using Linux namespaces, chroot, and cgroups v2.

## What This Project Does

`mini-docker` runs commands inside isolated containers using the same Linux kernel features that power real Docker:

- **PID Namespace** — container gets its own process tree, starts at PID 1
- **UTS Namespace** — container gets its own hostname
- **Mount Namespace** — container gets its own filesystem view
- **chroot** — filesystem isolation using Alpine Linux rootfs
- **cgroups v2** — memory (256MB) and PID (20) resource limits

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
(empty — host filesystem not visible)
#ARCHITECHTURE
mydocker run /bin/sh
       │
       ▼
  [run() — host side]
  Creates new namespaces:
  • CLONE_NEWPID  → isolated process tree
  • CLONE_NEWUTS  → isolated hostname
  • CLONE_NEWNS   → isolated mounts
  Applies cgroup v2 limits:
  • memory.max = 256MB
  • pids.max   = 20
       │
       ▼
  [child() — runs INSIDE namespaces]
  • syscall.Chroot("rootfs")  → Alpine filesystem
  • Mount /proc               → ps/top work correctly
  • exec user command         → /bin/sh
  #PROJECT STRUCTURE
  mini-docker/
├── cmd/mydocker/
│   └── main.go              # CLI entry point + namespace setup
├── internal/
│   └── cgroups/
│       └── cgroups.go       # cgroup v2 resource limiting
├── rootfs/                  # Alpine Linux root filesystem
├── go.mod
└── README.md
