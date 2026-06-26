// ============================================================
// antivm.go - Virtual Machine Detection
// ============================================================
package antivm

import (
	"log"
	"os"
	"runtime"
)

var isVM bool

// IsVM returns true if running in a VM
func IsVM() bool {
	return isVM
}

// Check performs VM detection
func Check() {
	if runtime.GOOS == "windows" {
		checkWindowsVM()
	} else {
		checkUnixVM()
	}
}

// checkWindowsVM checks for VM indicators on Windows
func checkWindowsVM() {
	// Check for VM artifacts
	vmIndicators := []string{
		"C:\\Program Files\\VMware",
		"C:\\Program Files\\VirtualBox",
		"C:\\Windows\\System32\\vbox*.dll",
	}
	
	for _, ind := range vmIndicators {
		if _, err := os.Stat(ind); err == nil {
			log.Printf("⚠️ VM indicator found: %s", ind)
			isVM = true
			return
		}
	}
	
	// Check registry (simplified)
	// In production, use windows/registry package
	
	// Check for VM processes
	vmProcesses := []string{
		"vmtoolsd.exe",
		"vboxservice.exe",
	}
	
	for _, proc := range vmProcesses {
		if processExists(proc) {
			log.Printf("⚠️ VM process found: %s", proc)
			isVM = true
			return
		}
	}
}

// checkUnixVM checks for VM indicators on Unix
func checkUnixVM() {
	// Check /sys/class/dmi/id/product_name
	if data, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		if string(data) == "VirtualBox" || string(data) == "VMware" {
			log.Printf("⚠️ VM detected via DMI: %s", string(data))
			isVM = true
			return
		}
	}
}

// processExists checks if a process exists
func processExists(name string) bool {
	// In production, use gopsutil
	return false
}