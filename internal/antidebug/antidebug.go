// ============================================================
// antidebug.go - Anti-Debugging Detection
// ============================================================
package antidebug

import (
	"log"
	"os"
	"runtime"
)

var isDebugged bool

// IsDebugged returns true if a debugger is detected
func IsDebugged() bool {
	return isDebugged
}

// Run starts anti-debugging checks
func Run() {
	if runtime.GOOS == "windows" {
		checkWindowsDebugger()
	} else {
		checkUnixDebugger()
	}
}

// checkWindowsDebugger checks for debuggers on Windows
func checkWindowsDebugger() {
	// Check for debugger processes
	debuggers := []string{
		"ollydbg.exe",
		"x64dbg.exe",
		"x32dbg.exe",
		"windbg.exe",
		"ida.exe",
		"idag.exe",
		"idaw.exe",
		"idaldr.exe",
		"processhacker.exe",
		"procexp.exe",
	}
	
	for _, proc := range debuggers {
		if processExists(proc) {
			log.Printf("⚠️ Debugger detected: %s", proc)
			isDebugged = true
			return
		}
	}
	
	// Check for debugger environment
	if os.Getenv("PYTHONDEBUG") != "" || os.Getenv("PYTHONHASHSEED") != "" {
		isDebugged = true
		return
	}
}

// checkUnixDebugger checks for debuggers on Unix
func checkUnixDebugger() {
	debuggers := []string{
		"gdb",
		"lldb",
		"strace",
		"dtrace",
	}
	
	for _, proc := range debuggers {
		if processExists(proc) {
			log.Printf("⚠️ Debugger detected: %s", proc)
			isDebugged = true
			return
		}
	}
}

// processExists checks if a process exists (simplified)
func processExists(name string) bool {
	// In production, use gopsutil or syscall
	return false
}