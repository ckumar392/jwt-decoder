# JWT Decoder CLI

A command-line tool to decode and display JWT (JSON Web Tokens) in a formatted, colorful manner.

## Features

- 🎨 **Colorful Output** - Easy-to-read colored output for different token parts
- 📅 **Human-readable Timestamps** - Automatically converts Unix timestamps to readable dates
- ⏰ **Expiry Check** - Check if a token is expired with the `-e` flag
- 🔧 **Raw Mode** - Output raw JSON for piping to other tools
- 🚫 **No External Dependencies at Runtime** - Single binary, no runtime dependencies

## Installation

### Using Homebrew (macOS/Linux)

```bash
brew tap ckumar3/tap
brew install jwt-decoder
```

### Using Go

```bash
go install github.com/ckumar3/jwt-decoder@latest
```

### From Source

```bash
git clone https://github.com/ckumar3/jwt-decoder.git
cd jwt-decoder
go build -o jwt-decoder .
```

## Usage

### Basic Usage

```bash
jwt-decoder <token>
```

### Example

```bash
jwt-decoder eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

### Output

```
═══════════════════════════════════════════════════════════════
  HEADER
═══════════════════════════════════════════════════════════════
{
  "alg": "HS256",
  "typ": "JWT"
}

═══════════════════════════════════════════════════════════════
  PAYLOAD
═══════════════════════════════════════════════════════════════
{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022
}

  ── Timestamps (Human Readable) ──
  Issued At: Thu, 18 Jan 2018 01:30:22 UTC

═══════════════════════════════════════════════════════════════
  SIGNATURE
═══════════════════════════════════════════════════════════════
  SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

  ⚠ Note: This tool does not verify the signature.
  Use appropriate libraries to verify token authenticity.
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--raw` | `-r` | Show raw JSON without formatting |
| `--no-color` | `-n` | Disable colored output |
| `--check-expiry` | `-e` | Check if token is expired |
| `--help` | `-h` | Show help message |

### Check Token Expiry

```bash
jwt-decoder -e <token>
```

### Raw Output (for piping)

```bash
jwt-decoder -r <token> | jq .
```

### Version

```bash
jwt-decoder version
```

## Building for Multiple Platforms

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o jwt-decoder-darwin-amd64 .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o jwt-decoder-darwin-arm64 .

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o jwt-decoder-linux-amd64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o jwt-decoder-windows-amd64.exe .
```

## Homebrew Formula

To publish to Homebrew, create a tap repository and add this formula:

```ruby
class JwtDecoder < Formula
  desc "CLI tool to decode and display JWT tokens"
  homepage "https://github.com/ckumar3/jwt-decoder"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ckumar3/jwt-decoder/releases/download/v1.0.0/jwt-decoder-darwin-arm64.tar.gz"
      sha256 "YOUR_SHA256_HERE"
    else
      url "https://github.com/ckumar3/jwt-decoder/releases/download/v1.0.0/jwt-decoder-darwin-amd64.tar.gz"
      sha256 "YOUR_SHA256_HERE"
    end
  end

  on_linux do
    url "https://github.com/ckumar3/jwt-decoder/releases/download/v1.0.0/jwt-decoder-linux-amd64.tar.gz"
    sha256 "YOUR_SHA256_HERE"
  end

  def install
    bin.install "jwt-decoder"
  end

  test do
    system "#{bin}/jwt-decoder", "version"
  end
end
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
