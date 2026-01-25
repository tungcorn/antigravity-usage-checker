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
	fmt.Println("🔍 Scanning for Antigravity server...")

	processInfo, err := discovery.FindAntigravityProcess()
	if err != nil {
		fmt.Println("⚠️  Antigravity not running, checking cache...")
		cachedData, cacheErr := cache.LoadLastKnown()
		if cacheErr != nil {
			return fmt.Errorf("Antigravity not running and no cached data available")
		}
		display.ShowUsage(cachedData, outputJSON, true)
		return nil
	}

	fmt.Printf("✅ Found server on port %d (PID: %d)\n", processInfo.ConnectPort, processInfo.PID)

	// Load credentials for token expiry info
	creds, _ := auth.LoadCredentials()

	fmt.Println("📡 Fetching quota data...")

	var quota *api.UsageData

	// Try Google Cloud API first if credentials are valid
	if creds != nil && !creds.IsExpired() {
		googleClient := api.NewGoogleCloudClient(creds.AccessToken)
		quota, err = googleClient.GetUsageData()
		if err != nil {
			fmt.Printf("⚠️  Google Cloud API failed: %v\n", err)
			fmt.Println("📡 Falling back to local server API...")
		} else {
			fmt.Println("✅ Got exact quota from Google Cloud API")
		}
	}

	// Fallback to local server API
	if quota == nil {
		client := api.NewClient(processInfo.ConnectPort, processInfo.CSRFToken, processInfo.HTTPPort)
		quota, err = client.GetUserStatus()
		if err != nil {
			fmt.Printf("⚠️  Local API call failed: %v\n", err)
			cachedData, cacheErr := cache.LoadLastKnown()
			if cacheErr != nil {
				return fmt.Errorf("all API calls failed and no cached data: %v", err)
			}
			display.ShowUsage(cachedData, outputJSON, true)
			return nil
		}
	}

	// Show account info from API response (most accurate source)
	if quota.Email != "" {
		fmt.Printf("👤 Account: %s\n", quota.Email)
	}

	cache.Save(quota)
	display.ShowUsage(quota, outputJSON, false)

	return nil
}

func printHelp() {
	fmt.Printf(`%s v%s

Check Antigravity AI usage quota from the terminal.

USAGE:
    agusage [OPTIONS]
    agu [OPTIONS]      (alias for quick access)

OPTIONS:
    -h, --help      Show help information
    -v, --version   Show version information
    -j, --json      Output in JSON format

EXAMPLES:
    agusage            Check current quota
    agu                Same as above (shorter)
    agusage --json     Output JSON for scripting
`, AppName, Version)
}
