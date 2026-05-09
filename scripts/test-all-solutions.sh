#!/bin/bash
# 批量测试所有 stage 的 solution（Java + Python）
# 用法: ./scripts/test-all-solutions.sh
#
# 分支模型：solution 仓库每种语言一个分支（java / python），
# 脚本通过 git worktree 将各分支 checkout 到临时目录中测试。
# starter 仓库各语言独立 repo（tinynum-python-starter / tinynum-java-starter），
# 测试时从对应 starter repo 的 main 分支复制 test driver 到 solution worktree。

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTER_DIR="$(dirname "$SCRIPT_DIR")"
SOLUTION_DIR="${TESTER_DIR}/../solution"
STARTER_DIR="${TESTER_DIR}/../starter"

# 构建 tester
cd "$TESTER_DIR"
go build -o tinynum-tester .

# Stage 列表（按课程顺序）
STAGES=(
    "storage-and-shape"
    "strides-and-indexing"
    "reshape"
    "transpose"
    "unary-math"
    "binary-ops"
    "broadcasting"
    "sum-and-mean"
    "max-var-and-friends"
    "matmul"
    "slicing-and-views"
    "creation-and-random"
    "join-and-transform"
    "fancy-indexing"
    "capstone-toolkit"
)

# 语言列表
LANGUAGES=("java" "python")

PASSED=0
FAILED=0
SKIPPED=0
TOTAL_TIME=0

echo "=========================================="
echo "  TinyNum Solution Tester"
echo "=========================================="
echo ""

for lang in "${LANGUAGES[@]}"; do
    echo "--- Language: ${lang} ---"
    echo ""

    # 使用 git worktree 将语言分支 checkout 到临时目录
    # 若当前分支已是目标分支，直接使用主目录
    worktree_dir="${SOLUTION_DIR}/.worktree-${lang}"
    use_worktree=true

    current_branch=$(git -C "$SOLUTION_DIR" branch --show-current 2>/dev/null || true)
    if [ "$current_branch" = "$lang" ]; then
        sol_dir="$SOLUTION_DIR"
        use_worktree=false
    else
        if [ -d "$worktree_dir" ]; then
            git -C "$SOLUTION_DIR" worktree remove --force "$worktree_dir" 2>/dev/null || rm -rf "$worktree_dir"
        fi
        if ! git -C "$SOLUTION_DIR" worktree add "$worktree_dir" "$lang" 2>/dev/null; then
            echo "⏭️  [${lang}] SKIPPED - branch not found"
            ((SKIPPED += ${#STAGES[@]}))
            echo ""
            continue
        fi
        sol_dir="$worktree_dir"
    fi

    # 从对应 starter repo 的 main 分支复制 test driver 到 solution worktree
    # （starter 各语言是独立 repo，无需 worktree）
    if [ "$lang" = "java" ]; then
        starter_tests="${STARTER_DIR}/tinynum-java-starter/tests"
        mkdir -p "${sol_dir}/tests"
        cp -f "${starter_tests}/"*.java "${sol_dir}/tests/" 2>/dev/null || true
    elif [ "$lang" = "python" ]; then
        starter_tests="${STARTER_DIR}/tinynum-python-starter/tests"
        mkdir -p "${sol_dir}/tests"
        cp -f "${starter_tests}/"*.py "${sol_dir}/tests/" 2>/dev/null || true
    fi

    for stage in "${STAGES[@]}"; do
        printf "🧪 [%-24s %6s] Testing... " "$stage" "$lang"

        start_time=$(python3 -c 'import time; print(time.time())')

        if ./tinynum-tester -d="$sol_dir" -s="$stage" > /dev/null 2>&1; then
            end_time=$(python3 -c 'import time; print(time.time())')
            elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
            echo "✅ PASSED (${elapsed}s)"
            ((PASSED++))
        else
            end_time=$(python3 -c 'import time; print(time.time())')
            elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
            echo "❌ FAILED (${elapsed}s)"
            ((FAILED++))
        fi

        TOTAL_TIME=$(python3 -c "print(f'{$TOTAL_TIME + $elapsed:.2f}')")
    done

    # 每种语言测试完后立即清理 worktree
    if [ "$use_worktree" = true ]; then
        git -C "$SOLUTION_DIR" worktree remove --force "$worktree_dir" 2>/dev/null || true
    fi

    echo ""
done

echo "=========================================="
echo "  Results: $PASSED passed, $FAILED failed, $SKIPPED skipped"
echo "  Total time: ${TOTAL_TIME}s"
echo "=========================================="

if [ $FAILED -gt 0 ]; then
    exit 1
fi
