# TinyNum Tester

TinyNum 课程自动评测工具。

## 方式一：从源码构建

```bash
git clone https://github.com/tinycs-cn/tinynum-tester
cd tinynum-tester
go build .
./tinynum-tester -s storage-and-shape -d ~/my-solution/java
```

**依赖：** Go 1.24+、Java 21+、python3

## 方式二：Docker 镜像

**快速上手**

```bash
cd ~/my-solution  # 你的解答根目录（包含 java/ 或 python/ 子目录）
docker pull ghcr.io/tinycs/tinynum-tester:latest
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/tinycs/tinynum-tester:latest -s storage-and-shape -d /workspace/java
```

**便捷脚本（推荐）**

在解答根目录创建 `test.sh`：

```bash
#!/bin/bash
LANG=${2:-java}
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/tinycs/tinynum-tester:latest \
  -s "${1:-storage-and-shape}" -d "/workspace/${LANG}"
```

用法：`chmod +x test.sh && ./test.sh broadcasting python`

**本地构建（可选）**

```bash
git clone https://github.com/tinycs-cn/tinynum-tester
cd tinynum-tester
docker build -t my-tester .
# 用法：docker run --rm --user $(id -u):$(id -g) -v ~/my-solution:/workspace my-tester -s storage-and-shape -d /workspace/java
```

## 许可证

MIT
