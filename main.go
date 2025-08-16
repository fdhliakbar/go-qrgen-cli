package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

const (
	defaultSize   = 256
	defaultOutput = "qr.png"
	version       = "1.0.0"
)

type Config struct {
	Text    string
	URL     string
	File    string
	Image   string
	WiFi    string
	VCard   string
	Output  string
	Size    int
	Quiet   bool
	Help    bool
	Version bool
	Quality string
	Batch   bool
	Preview bool
}

func main() {
	config := parseFlags()

	if config.Help {
		showHelp()
		return
	}

	if config.Version {
		showVersion()
		return
	}

	// Validate and get input content
	content, err := getInputContent(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if content == "" {
		fmt.Fprintf(os.Stderr, "Error: No input provided. Use -t, -u, or -f flag.\n")
		showUsage()
		os.Exit(1)
	}

	// Validate size
	if config.Size <= 0 || config.Size > 2048 {
		fmt.Fprintf(os.Stderr, "Error: Size must be between 1 and 2048 pixels\n")
		os.Exit(1)
	}

	// Set quality level
	var recoveryLevel qrcode.RecoveryLevel
	switch strings.ToLower(config.Quality) {
	case "low", "l":
		recoveryLevel = qrcode.Low
	case "medium", "m":
		recoveryLevel = qrcode.Medium
	case "high", "h":
		recoveryLevel = qrcode.High
	case "highest", "hh":
		recoveryLevel = qrcode.Highest
	default:
		recoveryLevel = qrcode.Medium
	}

	// Generate QR code
	if config.Preview {
		showASCIIPreview(content)
	}

	err = generateQRCode(content, config.Output, config.Size, recoveryLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating QR code: %v\n", err)
		os.Exit(1)
	}

	if !config.Quiet {
		fmt.Printf("✅ QR code successfully generated!\n")
		fmt.Printf("📁 Output: %s\n", config.Output)
		fmt.Printf("📏 Size: %dx%d pixels\n", config.Size, config.Size)
		fmt.Printf("📊 Quality: %s\n", config.Quality)
		fmt.Printf("💾 File size: %s\n", getFileSize(config.Output))
	}
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.Text, "text", "", "Text to encode in QR code")
	flag.StringVar(&config.Text, "t", "", "Text to encode in QR code (shorthand)")

	flag.StringVar(&config.URL, "url", "", "URL to encode in QR code")
	flag.StringVar(&config.URL, "u", "", "URL to encode in QR code (shorthand)")

	flag.StringVar(&config.File, "file", "", "File containing text to encode")
	flag.StringVar(&config.File, "f", "", "File containing text to encode (shorthand)")

	flag.StringVar(&config.Image, "image", "", "Image file to encode as base64")
	flag.StringVar(&config.Image, "i", "", "Image file to encode as base64 (shorthand)")

	flag.StringVar(&config.WiFi, "wifi", "", "WiFi credentials: 'SSID:PASSWORD:SECURITY'")
	flag.StringVar(&config.WiFi, "w", "", "WiFi credentials (shorthand)")

	flag.StringVar(&config.VCard, "vcard", "", "vCard file (.vcf) to encode")

	flag.StringVar(&config.Output, "output", defaultOutput, "Output file name")
	flag.StringVar(&config.Output, "o", defaultOutput, "Output file name (shorthand)")

	flag.IntVar(&config.Size, "size", defaultSize, "QR code size in pixels")
	flag.IntVar(&config.Size, "s", defaultSize, "QR code size in pixels (shorthand)")

	flag.StringVar(&config.Quality, "quality", "medium", "Error correction level (low/medium/high/highest)")
	flag.StringVar(&config.Quality, "q", "medium", "Error correction level (shorthand)")

	flag.BoolVar(&config.Batch, "batch", false, "Batch mode - process multiple inputs from file")
	flag.BoolVar(&config.Preview, "preview", false, "Show ASCII QR preview in terminal")
	flag.BoolVar(&config.Quiet, "quiet", false, "Quiet mode - no output messages")
	flag.BoolVar(&config.Help, "help", false, "Show help message")
	flag.BoolVar(&config.Help, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&config.Version, "version", false, "Show version")
	flag.BoolVar(&config.Version, "v", false, "Show version (shorthand)")

	flag.Parse()
	return config
}

func getInputContent(config Config) (string, error) {
	// Priority: batch -> vcard -> wifi -> image -> file -> url -> text
	if config.Batch && config.File != "" {
		return processBatchFile(config.File)
	}

	if config.VCard != "" {
		return readFromFile(config.VCard)
	}

	if config.WiFi != "" {
		return generateWiFiQR(config.WiFi)
	}

	if config.Image != "" {
		return encodeImageToBase64(config.Image)
	}

	if config.File != "" {
		return readFromFile(config.File)
	}

	if config.URL != "" {
		// Check if it's a web URL to fetch content
		if strings.HasPrefix(config.URL, "http://") || strings.HasPrefix(config.URL, "https://") {
			if !isValidURL(config.URL) {
				return "", fmt.Errorf("invalid URL format: %s", config.URL)
			}
			return config.URL, nil
		}
		return config.URL, nil
	}

	if config.Text != "" {
		return config.Text, nil
	}

	return "", nil
}

func readFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("cannot open file %s: %v", filename, err)
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file %s: %v", filename, err)
	}

	return strings.TrimSpace(content.String()), nil
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func generateQRCode(content, outputPath string, size int, recoveryLevel qrcode.RecoveryLevel) error {
	// Check if output file already exists
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("⚠️  File %s already exists. Overwrite? (y/N): ", outputPath)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			return fmt.Errorf("operation cancelled")
		}
	}

	// Create output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("cannot create directory %s: %v", dir, err)
		}
	}

	// Generate QR code
	qr, err := qrcode.New(content, recoveryLevel)
	if err != nil {
		return fmt.Errorf("cannot create QR code: %v", err)
	}

	// Write to file
	err = qr.WriteFile(size, outputPath)
	if err != nil {
		return fmt.Errorf("cannot write QR code to file: %v", err)
	}

	return nil
}

func getFileSize(filename string) string {
	info, err := os.Stat(filename)
	if err != nil {
		return "unknown"
	}

	size := info.Size()
	if size < 1024 {
		return fmt.Sprintf("%d bytes", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
}

func showHelp() {
	fmt.Printf(`🚀 QR Code Generator CLI v%s

USAGE:
    qrgen [OPTIONS]

INPUT OPTIONS:
    -t, --text      Text to encode in QR code
    -u, --url       URL to encode in QR code  
    -f, --file      File containing text to encode
    -i, --image     Image file to encode as base64 data URI
    -w, --wifi      WiFi credentials (SSID:PASSWORD:SECURITY)
    --vcard         vCard file (.vcf) for contact info
    --batch         Batch process multiple inputs from file

OUTPUT OPTIONS:
    -o, --output    Output file name (default: qr.png)
    -s, --size      QR code size in pixels (default: 256)
    -q, --quality   Error correction: low/medium/high/highest (default: medium)
    --preview       Show ASCII QR preview in terminal
    --quiet         Quiet mode - no output messages

GENERAL:
    -h, --help      Show this help message
    -v, --version   Show version

EXAMPLES:
    # Basic text QR
    qrgen -t "Hello World!"
    
    # URL QR with custom size
    qrgen -u "https://github.com/yourusername" -s 512 -o github.png
    
    # WiFi QR code
    qrgen -w "MyWiFi:password123:WPA" -o wifi.png
    
    # Image to base64 QR
    qrgen -i logo.png -o image_qr.png
    
    # Contact info from vCard
    qrgen --vcard contact.vcf -o contact.png
    
    # Batch processing
    qrgen -f urls.txt --batch
    
    # Preview in terminal
    qrgen -t "Preview Test" --preview
    
    # High quality with preview
    qrgen -u "https://important-site.com" -q highest --preview

SUPPORTED FORMATS:
    Input Images: PNG, JPG, JPEG, GIF, WebP
    Output: PNG only
    WiFi Security: WPA, WEP, nopass
    
BATCH FILE FORMAT:
    # Lines starting with # are comments
    https://github.com/user1
    Contact: +1234567890
    https://example.com
    
AUTHOR:
    Generated with ❤️ using Go
`, version)
}

func showVersion() {
	fmt.Printf("qrgen version %s\n", version)
}

func showUsage() {
	fmt.Println("Usage: qrgen [OPTIONS]")
	fmt.Println("Run 'qrgen --help' for more information.")
}

// New functions for enhanced features
func encodeImageToBase64(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot open image file %s: %v", imagePath, err)
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("cannot read image file: %v", err)
	}

	// Get file extension for MIME type
	ext := strings.ToLower(filepath.Ext(imagePath))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	default:
		mimeType = "image/png"
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func generateWiFiQR(wifiConfig string) (string, error) {
	parts := strings.Split(wifiConfig, ":")
	if len(parts) < 2 {
		return "", fmt.Errorf("WiFi format should be 'SSID:PASSWORD' or 'SSID:PASSWORD:SECURITY'")
	}

	ssid := parts[0]
	password := parts[1]
	security := "WPA"

	if len(parts) >= 3 {
		security = strings.ToUpper(parts[2])
	}

	// WiFi QR format: WIFI:T:WPA;S:SSID;P:PASSWORD;H:false;
	wifiQR := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;H:false;", security, ssid, password)
	return wifiQR, nil
}

func processBatchFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("cannot open batch file %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Process each line as separate QR
		outputFile := fmt.Sprintf("batch_%d.png", lineNum)
		err := generateQRCode(line, outputFile, 256, qrcode.Medium)
		if err != nil {
			fmt.Printf("❌ Error processing line %d: %v\n", lineNum, err)
		} else {
			fmt.Printf("✅ Generated: %s\n", outputFile)
		}
		lineNum++
	}

	return fmt.Sprintf("Batch processing completed. Generated %d QR codes.", lineNum-1), nil
}

func showASCIIPreview(content string) {
	fmt.Printf("\n📱 ASCII QR Preview:\n")
	fmt.Println("╭─────────────────────╮")

	// Simple ASCII QR representation
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		fmt.Println("│ Cannot generate preview │")
		fmt.Printf("Content: %s\n", truncateString(content, 50))
		return
	}

	// Get QR bitmap (simplified)
	bitmap := qr.Bitmap()
	size := len(bitmap)

	// Show reduced version for terminal
	step := size / 15 // Reduce to ~15x15 for terminal display
	if step == 0 {
		step = 1
	}

	for i := 0; i < size; i += step {
		fmt.Print("│ ")
		for j := 0; j < size; j += step {
			if bitmap[i][j] {
				fmt.Print("██")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println(" │")
	}

	fmt.Println("╰─────────────────────╯")
	fmt.Printf("Content: %s\n\n", truncateString(content, 50))
}

func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen] + "..."
}

func fetchURLContent(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read response: %v", err)
	}

	return string(body), nil
}
