// ============================================================
// system.go - System Information Collection
// ============================================================
package system

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type SystemInfo struct {
	Hostname      string   `json:"hostname"`
	Username      string   `json:"username"`
	OS            string   `json:"os"`
	Arch          string   `json:"arch"`
	CPU           string   `json:"cpu"`
	Memory        string   `json:"memory"`
	Disk          string   `json:"disk"`
	IP            string   `json:"ip"`
	MAC           string   `json:"mac"`
	Uptime        string   `json:"uptime"`
	InstalledApps []string `json:"installed_apps"`
	Processes     []string `json:"processes"`
	Timestamp     string   `json:"timestamp"`
}

var systemInfo SystemInfo

// Run collects and saves system information
func Run() {
	log.Println("🖥️ Collecting system information...")
	
	systemInfo = SystemInfo{
		Hostname:      getHostname(),
		Username:      getUsername(),
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		CPU:           getCPUInfo(),
		Memory:        getMemoryInfo(),
		Disk:          getDiskInfo(),
		IP:            getIPAddress(),
		MAC:           getMACAddress(),
		Uptime:        getUptime(),
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
		InstalledApps: getInstalledApps(),
		Processes:     getRunningProcesses(),
	}
	
	saveSystemInfo()
}

func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}

func getUsername() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERNAME")
	}
	return os.Getenv("USER")
}

func getCPUInfo() string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "cpu", "get", "name")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/cpuinfo")
		if err == nil {
			return string(data)
		}
	}
	
	return "unknown"
}

func getMemoryInfo() string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "memorychip", "get", "capacity")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/meminfo")
		if err == nil {
			return string(data)
		}
	}
	
	return "unknown"
}

func getDiskInfo() string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "logicaldisk", "get", "size,freespace,caption")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	if runtime.GOOS == "linux" {
		cmd := exec.Command("df", "-h")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	return "unknown"
}

func getIPAddress() string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("ipconfig")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	if runtime.GOOS == "linux" {
		cmd := exec.Command("ifconfig")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	return "unknown"
}

func getMACAddress() string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("getmac")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	if runtime.GOOS == "linux" {
		cmd := exec.Command("ip", "link")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	return "unknown"
}

func getUptime() string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "os", "get", "lastbootuptime")
		output, err := cmd.Output()
		if err == nil {
			return string(output)
		}
	}
	
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/uptime")
		if err == nil {
			return string(data)
		}
	}
	
	return "unknown"
}

func getInstalledApps() []string {
	var apps []string
	
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "product", "get", "name")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.Contains(line, "Name") && !strings.Contains(line, "Caption") {
					apps = append(apps, line)
				}
			}
		}
	}
	
	return apps
}

func getRunningProcesses() []string {
	var processes []string
	
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.Contains(line, "Image Name") && !strings.Contains(line, "====") {
					processes = append(processes, line)
				}
			}
		}
	}
	
	if runtime.GOOS == "linux" {
		cmd := exec.Command("ps", "aux")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					processes = append(processes, line)
				}
			}
		}
	}
	
	return processes
}

func saveSystemInfo() {
	data, err := json.MarshalIndent(systemInfo, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal system info: %v", err)
		return
	}
	
	if err := os.WriteFile("system_info.json", data, 0644); err != nil {
		log.Printf("❌ Failed to save system info: %v", err)
		return
	}
	
	log.Printf("✅ System info saved to system_info.json")
}