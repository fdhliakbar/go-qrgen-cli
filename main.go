package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	Output  string
	Size    int
	Quiet   bool
	Help    bool
	Version bool
	Quality string
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

	flag.StringVar(&config.Output, "output", defaultOutput, "Output file name")
	flag.StringVar(&config.Output, "o", defaultOutput, "Output file name (shorthand)")

	flag.IntVar(&config.Size, "size", defaultSize, "QR code size in pixels")
	flag.IntVar(&config.Size, "s", defaultSize, "QR code size in pixels (shorthand)")

	flag.StringVar(&config.Quality, "quality", "medium", "Error correction level (low/medium/high/highest)")
	flag.StringVar(&config.Quality, "q", "medium", "Error correction level (shorthand)")

	flag.BoolVar(&config.Quiet, "quiet", false, "Quiet mode - no output messages")
	flag.BoolVar(&config.Help, "help", false, "Show help message")
	flag.BoolVar(&config.Help, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&config.Version, "version", false, "Show version")
	flag.BoolVar(&config.Version, "v", false, "Show version (shorthand)")

	flag.Parse()
	return config
}

func getInputContent(config Config) (string, error) {
	// Priority: file -> url -> text
	if config.File != "" {
		return readFromFile(config.File)
	}

	if config.URL != "" {
		if !isValidURL(config.URL) {
			return "", fmt.Errorf("invalid URL format: %s", config.URL)
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

OPTIONS:
    -t, --text      Text to encode in QR code
    -u, --url       URL to encode in QR code  
    -f, --file      File containing text to encode
    -o, --output    Output file name (default: qr.png)
    -s, --size      QR code size in pixels (default: 256)
    -q, --quality   Error correction level: low/medium/high/highest (default: medium)
    --quiet         Quiet mode - no output messages
    -h, --help      Show this help message
    -v, --version   Show version

EXAMPLES:
    # Generate QR from text
    qrgen -t "Hello World!"
    
    # Generate QR from URL
    qrgen -u "https://github.com/yourusername" -o github.png
    
    # Generate QR from file with custom size
    qrgen -f input.txt -s 512 -o large_qr.png
    
    # High quality QR code
    qrgen -t "Important Data" -q highest -o important.png
    
    # Quiet mode
    qrgen -u "https://example.com" --quiet

SUPPORTED FORMATS:
    Output: PNG (only)
    
AUTHOR:
    Generated with ❤️ using Go
`, version)
}

func showVersion() {
	fmt.Printf("qrgen version %s\n", version)
}

func showUsage() {
	fmt.Println("Usage: qrgen -t \"text\" OR qrgen -u \"url\" OR qrgen -f \"file.txt\"")
	fmt.Println("Run 'qrgen --help' for more information.")
}
