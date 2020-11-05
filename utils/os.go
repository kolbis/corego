package utils

import (
	"os"
	"runtime"
)

// ProcessName ...
func ProcessName() string {
	name, _ := os.Executable()
	return name
}

// HostName will returns the host name reported by the kernel
func HostName() string {
	name, _ := os.Hostname()
	return name
}

// ProcessID will return the process ID reporeted by the kernel
func ProcessID() int {
	return os.Getpid()
}

// OperatingSystem will return the operating system name reporeted by the kernel
func OperatingSystem() string {
	return runtime.GOOS
}
