package cgroups

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// CgroupConfig holds resource limits for a container.
type CgroupConfig struct {
	// MaxMemory is the memory limit in bytes.
	// e.g., 268435456 = 256MB
	MaxMemory int64

	// MaxPIDs is the maximum number of processes allowed.
	MaxPIDs int
}

const (
	// cgroupRoot is where cgroup v2 is mounted on modern Linux systems
	cgroupRoot = "/sys/fs/cgroup"

	// Our container's cgroup name
	cgroupName = "mini-docker"
)

// Apply creates a cgroup for the container and applies resource limits.
// pid is the container process ID on the HOST (not inside the container).
func Apply(pid int, cfg CgroupConfig) error {
	// Path to our cgroup directory
	cgroupPath := filepath.Join(cgroupRoot, cgroupName)

	// Create the cgroup directory
	// The kernel automatically creates control files inside it
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("creating cgroup dir: %w", err)
	}

	// Apply memory limit
	if cfg.MaxMemory > 0 {
		memLimit := strconv.FormatInt(cfg.MaxMemory, 10)
		if err := writeFile(cgroupPath, "memory.max", memLimit); err != nil {
			return fmt.Errorf("setting memory limit: %w", err)
		}
		fmt.Printf("[mydocker] Memory limit: %d bytes\n", cfg.MaxMemory)
	}

	// Apply PID limit
	if cfg.MaxPIDs > 0 {
		if err := writeFile(cgroupPath, "pids.max", strconv.Itoa(cfg.MaxPIDs)); err != nil {
			return fmt.Errorf("setting pid limit: %w", err)
		}
		fmt.Printf("[mydocker] PID limit: %d\n", cfg.MaxPIDs)
	}

	// Add the container process to this cgroup
	// Writing the PID to cgroup.procs moves it into the cgroup
	if err := writeFile(cgroupPath, "cgroup.procs", strconv.Itoa(pid)); err != nil {
		return fmt.Errorf("adding process to cgroup: %w", err)
	}

	fmt.Printf("[mydocker] Applied cgroup limits to PID %d\n", pid)
	return nil
}

// Cleanup removes the cgroup after the container exits.
// Always call this to avoid leaving stale cgroups.
func Cleanup() error {
	cgroupPath := filepath.Join(cgroupRoot, cgroupName)
	if err := os.Remove(cgroupPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing cgroup: %w", err)
	}
	return nil
}

// writeFile writes value to a cgroup control file.
func writeFile(cgroupPath, filename, value string) error {
	path := filepath.Join(cgroupPath, filename)
	return os.WriteFile(path, []byte(value), 0700)
}
