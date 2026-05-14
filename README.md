<p align="center">
  <img src="https://i.ibb.co.com/TMKpr41z/mirage.png" alt="mirage" border="0">
</p>
<h1 align="center">Mirage</h1>
<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-blue.svg" alt="Language Go">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Release-v1.0.0-orange.svg" alt="Release">
</p>

**Mirage** is a blazingly fast, lightweight, and pipeline-friendly CNAME detector built for bug hunters and network engineers. Designed with concurrency in mind, it rapidly resolves thousands of subdomains to reveal their underlying Canonical Names (CNAME), making it an essential tool for discovering **Subdomain Takeovers**.

## Features

- **Blazingly Fast:** Uses Go routines for high concurrency checking.
- **Pipeline-Friendly:** Seamlessly integrates with other tools like `subfinder`, `httpx`, etc.
- **Smart Input Filtering:** Automatically strips `http://`, `https://`, and URL paths from messy input files.
- **Custom DNS Resolvers:** Bypass ISP rate limits by specifying your own resolver (e.g., `8.8.8.8`).
- **Zero Dependencies:** Compiled as a single, statically-linked binary.

## Installation

### Easy Install (Recommended)
You can easily install Mirage by running the following one-liner in your terminal. This script will download the source, build it, and move it to your system path.

```bash
curl -sSL [https://raw.githubusercontent.com/openverselabs/mirage/main/install.sh](https://raw.githubusercontent.com/openverselabs/mirage/main/install.sh) | bash
```

### Manual Install (Go required)

If you prefer to build it yourself:

```bash
git clone [https://github.com/openverselabs/mirage.git](https://github.com/openverselabs/mirage.git)
cd mirage
go build -ldflags="-s -w" -o mirage main.go
sudo mv mirage /usr/local/bin/

```

## Usage

Mirage can be used by reading from a file or directly from standard input (stdin).

**1. Basic Usage (Reading from a file)**

```bash
mirage -l subdomains.txt -o results.txt

```

**2. Pipelining from other tools (e.g., Subfinder)**

```bash
subfinder -d target.com -silent | mirage

```

**3. Silent Mode (Ninja Mode)**
Hides the banner and only prints the results, perfect for chaining commands.

```bash
cat dirty_urls.txt | mirage -silent -c 200 > cnames.txt

```

**4. Using Custom DNS and Timeout**

```bash
echo "www.github.com" | mirage -r 1.1.1.1 -t 3

```

## Flags

| Flag | Description | Default |
| --- | --- | --- |
| `-c` | Maximum concurrency (number of workers) | `100` |
| `-l` | File containing the list of domains to check | `""` |
| `-o` | File to write the output to | `""` |
| `-r` | Custom DNS resolver (e.g., `8.8.8.8`) | `System Default` |
| `-t` | DNS Timeout in seconds | `5` |
| `-silent` | Show only results in the output (hides banner) | `false` |


## License and Contributions

* **License**: Distributed under the MIT License.
* **Contributing**: Pull requests are welcome. For major changes, please open an issue first to discuss the proposed updates.
