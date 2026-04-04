# TinyNum Tester

Automated testing tool for the TinyNum course.

## Option 1: Build from Source

```bash
git clone https://github.com/bootcraft-cn/tinynum-tester
cd tinynum-tester
go build .
./tinynum-tester -s storage-and-shape -d ~/my-solution/java
```

**Dependencies:** Go 1.24+, Java 21+, python3

## Option 2: Docker Image

**Quick Start**

```bash
cd ~/my-solution  # your solution root (contains java/ or python/)
docker pull ghcr.io/bootcraft-cn/tinynum-tester:latest
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/bootcraft-cn/tinynum-tester:latest -s storage-and-shape -d /workspace/java
```

**Simplified script (recommended)**

Create `test.sh` in your solution root:

```bash
#!/bin/bash
LANG=${2:-java}
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/bootcraft-cn/tinynum-tester:latest \
  -s "${1:-storage-and-shape}" -d "/workspace/${LANG}"
```

Usage: `chmod +x test.sh && ./test.sh broadcasting python`

**Local build (optional)**

```bash
git clone https://github.com/bootcraft-cn/tinynum-tester
cd tinynum-tester
docker build -t my-tester .
# Usage: docker run --rm --user $(id -u):$(id -g) -v ~/my-solution:/workspace my-tester -s storage-and-shape -d /workspace/java
```

## License

MIT
