# 🚀 QR Code Generator CLI

A simple, fast, and powerful command-line QR code generator built with Go.

![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Platform](https://img.shields.io/badge/platform-cross--platform-lightgrey)

## ✨ Features

- 📝 Generate QR codes from text, URLs, or files
- 🎨 Customizable size and error correction levels
- 💾 PNG output format
- 🔍 URL validation
- 📁 Automatic directory creation
- ⚡ Fast and lightweight
- 🤫 Quiet mode support
- 🔄 File overwrite protection

## 🚀 Installation

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

### Option 3: Direct install (if published)
```bash
go install github.com/yourusername/go-qrgen-cli@latest
```

## 📖 Usage

### Basic Examples

```bash
# Generate QR from text
qrgen -t "Hello World!"

# Generate QR from URL
qrgen -u "https://github.com/yourusername"

# Generate QR from file
qrgen -f input.txt

# Custom output filename
qrgen -t "My Text" -o my_qr.png

# Large QR code (512x512)
qrgen -t "Big QR Code" -s 512

# High quality error correction
qrgen -u "https://important-site.com" -q highest
```

### Advanced Examples

```bash
# Generate multiple QR codes with different settings
qrgen -u "https://github.com" -o github.png -s 400 -q high
qrgen -f company_info.txt -o company_qr.png -s 600 --quiet

# Batch processing (using shell)
for url in $(cat urls.txt); do
    filename=$(echo $url | sed 's|https://||' | sed 's|/|_|g')
    qrgen -u "$url" -o "qr_${filename}.png" --quiet
done
```

## 🎛️ Command Line Options

| Flag | Long Form | Description | Default |
|------|-----------|-------------|---------|
| `-t` | `--text` | Text to encode in QR code | - |
| `-u` | `--url` | URL to encode in QR code | - |
| `-f` | `--file` | File containing text to encode | - |
| `-o` | `--output` | Output file name | `qr.png` |
| `-s` | `--size` | QR code size in pixels | `256` |
| `-q` | `--quality` | Error correction level | `medium` |
| | `--quiet` | Quiet mode (no output messages) | `false` |
| `-h` | `--help` | Show help message | - |
| `-v` | `--version` | Show version | - |

### Quality Levels

- **low** / **l**: ~7% error correction
- **medium** / **m**: ~15% error correction (default)
- **high** / **h**: ~25% error correction  
- **highest** / **hh**: ~30% error correction

## 📁 Project Structure

```
go-qrgen-cli/
├── main.go          # Main application
├── go.mod           # Go module file
├── go.sum           # Go dependencies checksum
├── README.md        # This file
├── .gitignore       # Git ignore rules
├── LICENSE          # MIT License
├── examples/        # Example QR codes
│   └── sample_qr.png
└── docs/           # Additional documentation
    └── usage.md
```

## 🔧 Development

### Prerequisites
- Go 1.19 or higher
- Git

### Dependencies
```bash
go get github.com/skip2/go-qrcode
```

### Building
```bash
# Build for current platform
go build -o qrgen main.go

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o qrgen-linux main.go
GOOS=windows GOARCH=amd64 go build -o qrgen.exe main.go
GOOS=darwin GOARCH=amd64 go build -o qrgen-mac main.go
```

### Testing
```bash
# Test basic functionality
go run main.go -t "Test QR Code" -o test.png

# Test URL validation
go run main.go -u "https://github.com" -o github_test.png

# Test file input
echo "Hello from file!" > test_input.txt
go run main.go -f test_input.txt -o file_test.png
```

## 📊 Examples Gallery

| Input Type | Command | Output |
|------------|---------|--------|
| Text | `qrgen -t "Hello World!"` | ![Sample QR](examples/sample_qr.png) |
| URL | `qrgen -u "https://github.com"` | QR code linking to GitHub |
| File | `qrgen -f data.txt -s 400` | Large QR from file content |

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. 🍴 Fork the repository
2. 🌟 Create a feature branch (`git checkout -b feature/amazing-feature`)
3. 💻 Commit your changes (`git commit -m 'Add amazing feature'`)
4. 📤 Push to the branch (`git push origin feature/amazing-feature`)
5. 🎯 Open a Pull Request

### Ideas for Contributions
- [ ] SVG output format support
- [ ] Batch processing built-in
- [ ] QR code reading functionality
- [ ] Custom logo embedding
- [ ] Color customization
- [ ] Terminal QR preview (ASCII)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [skip2/go-qrcode](https://github.com/skip2/go-qrcode) - Excellent QR code library
- Go team for the amazing language
- The open-source community

## 📞 Support

- 🐛 [Report Bugs](https://github.com/yourusername/go-qrgen-cli/issues)
- 💡 [Request Features](https://github.com/yourusername/go-qrgen-cli/issues)
- 📖 [Documentation](https://github.com/yourusername/go-qrgen-cli/wiki)

---

**Made with ❤️ and Go**

⭐ **Star this repo if you find it useful!**