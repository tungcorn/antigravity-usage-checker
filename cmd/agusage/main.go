package main

import (
	"fmt"
	"os"

	"github.com/TungCorn/antigravity-usage-checker/internal/api"
	"github.com/TungCorn/antigravity-usage-checker/internal/auth"
	"github.com/TungCorn/antigravity-usage-checker/internal/cache"
	"github.com/TungCorn/antigravity-usage-checker/internal/discovery"
	"github.com/TungCorn/antigravity-usage-checker/internal/display"
)

const (
	Version = "0.1.0"
	AppName = "Antigravity Usage Checker"
)

func main() {
	// Parse command line flags
	showVersion := false
	outputJSON := false
	
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-v", "--version":
			showVersion = true
		case "-j", "--json":
			outputJSON = true
		case "-h", "--help":
			printHelp()
			return
		}
	}

	if showVersion {
		fmt.Printf("%s v%s\n", AppName, Version)
		return
	}

	// Run the main check
	if err := run(outputJSON); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(outputJSON bool) error {
	// Step 1: Find Antigravity process and extract info
	fmt.Println("🔍 Scanning for Antigravity server...")
	
	processInfo, err := discovery.FindAntigravityProcess()
	if err != nil {
		// Try cache fallback
		fmt.Println("⚠️  Antigravity not running, checking cache...")
		cachedData, cacheErr := cache.LoadLastKnown()
		if cacheErr != nil {
			return fmt.Errorf("Antigravity not running and no cached data available")
		}
		display.ShowUsage(cachedData, outputJSON, true)
		return nil
	}
	
	fmt.Printf("✅ Found server on port %d (PID: %d)\n", processInfo.ConnectPort, processInfo.PID)
	
	// Step 2: Load OAuth credentials (optional, for future use)
	creds, _ := auth.LoadCredentials()
	if creds != nil {
		fmt.Printf("✅ Credentials loaded (expires in %d min)\n", creds.ExpiresInMinutes())
	}
	
	// Step 3: Call API to get quota
	fmt.Println("📡 Fetching quota data...")
	
	client := api.NewClient(processInfo.ConnectPort, processInfo.CSRFToken, processInfo.HTTPPort)
	quota, err := client.GetUserStatus()
	if err != nil {
		// Try cache fallback on API error
		fmt.Printf("⚠️  API call failed: %v\n", err)
		cachedData, cacheErr := cache.LoadLastKnown()
		if cacheErr != nil {
			return fmt.Errorf("API call failed and no cached data: %v", err)
		}
		display.ShowUsage(cachedData, outputJSON, true)
		return nil
	}
	
	// Save to cache for future fallback
	cache.Save(quota)
	
	// Step 4: Display result
	display.ShowUsage(quota, outputJSON, false)
	
	return nil
}

func printHelp() {
	fmt.Printf(`%s v%s

Kiểm tra usage quota của Antigravity AI từ terminal.

USAGE:
    agusage [OPTIONS]

OPTIONS:
    -h, --help      Hiển thị trợ giúp
    -v, --version   Hiển thị phiên bản
    -j, --json      Xuất định dạng JSON

EXAMPLES:
    agusage            Kiểm tra quota hiện tại
    agusage --json     Xuất JSON để script
`, AppName, Version)
}
