package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tungcorn/antigravity-usage-checker/internal/api"
	"github.com/tungcorn/antigravity-usage-checker/internal/auth"
	"github.com/tungcorn/antigravity-usage-checker/internal/cache"
	"github.com/tungcorn/antigravity-usage-checker/internal/discovery"
	"github.com/tungcorn/antigravity-usage-checker/internal/display"
)

// Version is set by ldflags at build time
var Version = "dev"

const AppName = "Antigravity Usage Checker"

func main() {
	// Define flags
	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(showVersion, "v", false, "Show version information (shorthand)")
	
	outputJSON := flag.Bool("json", false, "Output in JSON format")
	flag.BoolVar(outputJSON, "j", false, "Output in JSON format (shorthand)")
	
	// Custom usage message
	flag.Usage = printHelp
	
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s v%s\n", AppName, Version)
		return
	}

	// Run the main check
	if err := run(*outputJSON); err != nil {
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
		if creds.IsExpired() {
			fmt.Println("⚠️  Credentials loaded but expired")
		} else {
			fmt.Printf("✅ Credentials loaded (expires in %d min)\n", creds.ExpiresInMinutes())
		}
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

Check Antigravity AI usage quota from the terminal.

USAGE:
    agusage [OPTIONS]

OPTIONS:
    -h, --help      Show help information
    -v, --version   Show version information
    -j, --json      Output in JSON format

EXAMPLES:
    agusage            Check current quota
    agusage --json     Output JSON for scripting
`, AppName, Version)
}
