# QR Code Generator CLI

A simple, fast, and powerful command-line QR code generator built with Go.
<img src="https://i.pinimg.com/1200x/55/a9/ab/55a9aba97d1e214f849ab2e55f3dabff.jpg" alt="Go Banner"/>

![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Platform](https://img.shields.io/badge/platform-cross--platform-lightgrey)

## Features

- Generate QR codes from text, URLs, or files
- **NEW!** Encode images as base64 data URIs
- **NEW!** Generate WiFi QR codes for easy connection
- **NEW!** Create QR codes from vCard contact files
- **NEW!** ASCII QR preview in terminal
- **NEW!** Batch processing multiple inputs
- Customizable size and error correction levels
- PNG output format
- URL validation
- Automatic directory creation
- Fast and lightweight
- Quiet mode support
- File overwrite protection

## Installation

### Option 1: Install from source
```bash
git clone https://github.com/yourusername/go-qrgen-cli.git
cd go-qrgen-cli
go mod tidy
go install
```

### Option 2: Build locally
```bash
git clone https://github.com/yourusername/go-qrgen-cli.git
cd go-qrgen-cli
go mod tidy
go build -o qrgen main.go
```

## Usage

### Basic Examples

```bash
# Generate QR from text
qrgen -t "Hello World!"

# Generate QR from URL
qrgen -u "https://github.com/yourusername"

# Generate QR from file
qrgen -f input.txt

# NEW! WiFi QR code
qrgen -w "MyWiFi:password123:WPA" -o wifi.png

# NEW! Image to base64 QR
qrgen -i company_logo.png -o logo_qr.png

# NEW! Preview in terminal
qrgen -t "Preview Test" --preview

# Custom output filename and size
qrgen -t "My Text" -o my_qr.png -s 512

# High quality error correction
qrgen -u "https://important-site.com" -q highest
```

### Advanced Examples

```bash
# Contact QR from vCard file
qrgen --vcard contact.vcf -o contact_qr.png

# Batch processing URLs
echo "https://github.com" > urls.txt
echo "https://google.com" >> urls.txt
qrgen -f urls.txt --batch

# Large WiFi QR with preview
qrgen -w "CoffeeShop:free123:WPA" -s 600 --preview -o cafe_wifi.png

# Process image with high quality
qrgen -i screenshot.png -q highest -s 800 -o screenshot_qr.png

# Silent batch processing
qrgen -f company_urls.txt --batch --quiet
```

## Command Line Options

| Flag | Long Form | Description | Default |
|------|-----------|-------------|---------|
| `-t` | `--text` | Text to encode in QR code | - |
| `-u` | `--url` | URL to encode in QR code | - |
| `-f` | `--file` | File containing text to encode | - |
| `-i` | `--image` | **NEW!** Image file to encode as base64 | - |
| `-w` | `--wifi` | **NEW!** WiFi credentials (SSID:PASSWORD:SECURITY) | - |
| | `--vcard` | **NEW!** vCard file for contact info | - |
| | `--batch` | **NEW!** Batch process multiple inputs | `false` |
| `-o` | `--output` | Output file name | `qr.png` |
| `-s` | `--size` | QR code size in pixels | `256` |
| `-q` | `--quality` | Error correction level | `medium` |
| | `--preview` | **NEW!** Show ASCII QR preview | `false` |
| | `--quiet` | Quiet mode (no output messages) | `false` |
| `-h` | `--help` | Show help message | - |
| `-v` | `--version` | Show version | - |


## Examples

<p>You can test this QR code, which will lead to github.com/fdhliakbar</p>
<div style="display: flex; justify-content:center">
<img src="./github.png" alt="Github QR Code"/>
</div>

---
