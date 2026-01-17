package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version  = "1.0.0"
	showRaw  bool
	noColor  bool
	checkExp bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "jwt-decoder [token]",
		Short: "A CLI tool to decode and display JWT tokens",
		Long: `JWT Decoder is a command-line tool that decodes JWT (JSON Web Tokens) 
and displays the header, payload, and signature in a formatted manner.

Example:
  jwt-decoder eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`,
		Args: cobra.ExactArgs(1),
		Run:  decodeToken,
	}

	rootCmd.Flags().BoolVarP(&showRaw, "raw", "r", false, "Show raw JSON without formatting")
	rootCmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "Disable colored output")
	rootCmd.Flags().BoolVarP(&checkExp, "check-expiry", "e", false, "Check if token is expired")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("jwt-decoder version %s\n", version)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func decodeToken(cmd *cobra.Command, args []string) {
	token := strings.TrimSpace(args[0])

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		printError("Invalid JWT format. A JWT should have 3 parts separated by dots.")
		os.Exit(1)
	}

	header, err := decodeSegment(parts[0])
	if err != nil {
		printError(fmt.Sprintf("Failed to decode header: %v", err))
		os.Exit(1)
	}

	payload, err := decodeSegment(parts[1])
	if err != nil {
		printError(fmt.Sprintf("Failed to decode payload: %v", err))
		os.Exit(1)
	}

	signature := parts[2]

	if noColor {
		color.NoColor = true
	}

	printSection("HEADER", header, showRaw)
	fmt.Println()
	printSection("PAYLOAD", payload, showRaw)
	fmt.Println()
	printSignature(signature)

	if checkExp {
		fmt.Println()
		checkTokenExpiry(payload)
	}
}

func decodeSegment(segment string) (map[string]interface{}, error) {
	// Add padding if needed
	segment = strings.TrimRight(segment, "=")
	switch len(segment) % 4 {
	case 2:
		segment += "=="
	case 3:
		segment += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(segment)
	if err != nil {
		// Try with RawURLEncoding
		decoded, err = base64.RawURLEncoding.DecodeString(strings.TrimRight(segment, "="))
		if err != nil {
			return nil, err
		}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(decoded, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func printSection(title string, data map[string]interface{}, raw bool) {
	titleColor := color.New(color.FgCyan, color.Bold)
	keyColor := color.New(color.FgYellow)
	stringColor := color.New(color.FgGreen)
	numberColor := color.New(color.FgMagenta)
	boolColor := color.New(color.FgBlue)

	titleColor.Printf("═══════════════════════════════════════════════════════════════\n")
	titleColor.Printf("  %s\n", title)
	titleColor.Printf("═══════════════════════════════════════════════════════════════\n")

	if raw {
		jsonBytes, _ := json.Marshal(data)
		fmt.Println(string(jsonBytes))
		return
	}

	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	lines := strings.Split(string(prettyJSON), "\n")

	for _, line := range lines {
		// Color the keys and values differently
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			keyColor.Print(parts[0])
			fmt.Print(":")
			if len(parts) > 1 {
				value := strings.TrimSpace(parts[1])
				if strings.HasPrefix(value, "\"") {
					stringColor.Println(parts[1])
				} else if value == "true" || value == "false" || strings.HasPrefix(value, "true,") || strings.HasPrefix(value, "false,") {
					boolColor.Println(parts[1])
				} else if value == "{" || value == "[" || value == "}," || value == "]," || value == "}" || value == "]" {
					fmt.Println(parts[1])
				} else {
					numberColor.Println(parts[1])
				}
			} else {
				fmt.Println()
			}
		} else {
			fmt.Println(line)
		}
	}

	// Print human-readable timestamps for common claims
	printTimestamps(data)
}

func printTimestamps(data map[string]interface{}) {
	timeColor := color.New(color.FgHiBlack, color.Italic)

	timestampClaims := map[string]string{
		"iat": "Issued At",
		"exp": "Expires At",
		"nbf": "Not Before",
	}

	hasTimestamps := false
	for claim := range timestampClaims {
		if _, ok := data[claim]; ok {
			hasTimestamps = true
			break
		}
	}

	if hasTimestamps {
		fmt.Println()
		timeColor.Println("  ── Timestamps (Human Readable) ──")
		for claim, label := range timestampClaims {
			if val, ok := data[claim]; ok {
				if ts, ok := val.(float64); ok {
					t := time.Unix(int64(ts), 0)
					timeColor.Printf("  %s: %s\n", label, t.Format(time.RFC1123))
				}
			}
		}
	}
}

func printSignature(signature string) {
	titleColor := color.New(color.FgCyan, color.Bold)
	sigColor := color.New(color.FgRed)

	titleColor.Printf("═══════════════════════════════════════════════════════════════\n")
	titleColor.Printf("  SIGNATURE\n")
	titleColor.Printf("═══════════════════════════════════════════════════════════════\n")
	sigColor.Printf("  %s\n", signature)

	infoColor := color.New(color.FgHiBlack, color.Italic)
	infoColor.Println("\n  ⚠ Note: This tool does not verify the signature.")
	infoColor.Println("  Use appropriate libraries to verify token authenticity.")
}

func checkTokenExpiry(payload map[string]interface{}) {
	titleColor := color.New(color.FgCyan, color.Bold)
	expiredColor := color.New(color.FgRed, color.Bold)
	validColor := color.New(color.FgGreen, color.Bold)
	infoColor := color.New(color.FgHiBlack)

	titleColor.Printf("═══════════════════════════════════════════════════════════════\n")
	titleColor.Printf("  EXPIRY CHECK\n")
	titleColor.Printf("═══════════════════════════════════════════════════════════════\n")

	exp, ok := payload["exp"]
	if !ok {
		infoColor.Println("  No expiration claim (exp) found in token.")
		return
	}

	expFloat, ok := exp.(float64)
	if !ok {
		infoColor.Println("  Invalid expiration claim format.")
		return
	}

	expTime := time.Unix(int64(expFloat), 0)
	now := time.Now()

	if now.After(expTime) {
		expiredColor.Printf("  ✗ TOKEN EXPIRED\n")
		infoColor.Printf("  Expired: %s\n", expTime.Format(time.RFC1123))
		infoColor.Printf("  Expired: %s ago\n", now.Sub(expTime).Round(time.Second))
	} else {
		validColor.Printf("  ✓ TOKEN VALID\n")
		infoColor.Printf("  Expires: %s\n", expTime.Format(time.RFC1123))
		infoColor.Printf("  Expires in: %s\n", expTime.Sub(now).Round(time.Second))
	}
}

func printError(msg string) {
	errColor := color.New(color.FgRed, color.Bold)
	errColor.Fprintf(os.Stderr, "Error: %s\n", msg)
}
