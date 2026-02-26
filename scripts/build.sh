#!/usr/bin/env bash
set -e

# ==========================
# 用法
# ==========================
usage() {
  echo "Usage: $0 --version vX.Y.Z"
  exit 1
}

# ==========================
# 参数解析
# ==========================
VERSION=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      VERSION="$2"
      shift 2
      ;;
    *)
      usage
      ;;
  esac
done

if [[ -z "$VERSION" ]]; then
  echo "[ERROR] version is required"
  usage
fi

# 去掉 v 前缀
VERSION="${VERSION#v}"

echo "[INFO] Build version: $VERSION"

# ==========================
# 路径
# ==========================
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/www"
BACKEND_DIR="$PROJECT_ROOT/backend"
BIN_DIR="$PROJECT_ROOT/bin/$VERSION"
CONFIG_TEMPLATE="$BACKEND_DIR/config.yaml.template"

MAIN_GO="$BACKEND_DIR/cmd/main.go"

# ==========================
# 前端打包
# ==========================
echo "[STEP] Build frontend"

cd "$FRONTEND_DIR"
rm -rf dist
npm install
npm run build

# 同步到 backend/web/dist（给 go:embed 用）
echo "[STEP] Sync frontend dist"
rm -rf "$BACKEND_DIR/web/dist"
mkdir -p "$BACKEND_DIR/web"
cp -r dist "$BACKEND_DIR/web/"

# ==========================
# 后端打包
# ==========================
mkdir -p "$BIN_DIR"
cd "$BACKEND_DIR"

# ---------- macOS arm64 ----------
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 Building: macOS arm64"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

GOOS=darwin GOARCH=arm64 \
go build \
  -o "$BIN_DIR/goinsight-osx-arm64" \
  -ldflags "-X main.Version=$VERSION" \
  "$MAIN_GO"

FILE_SIZE=$(du -h "$BIN_DIR/goinsight-osx-arm64" | cut -f1)
echo "✓ Binary built: goinsight-osx-arm64 ($FILE_SIZE)"

cd "$BIN_DIR"
cp "$CONFIG_TEMPLATE" .
tar -czf "goinsight-osx-arm64-v$VERSION.tar.gz" \
  "goinsight-osx-arm64" \
  "config.yaml.template" 2>/dev/null

# 生成 SHA256
if command -v shasum &> /dev/null; then
    shasum -a 256 "goinsight-osx-arm64-v$VERSION.tar.gz" > "goinsight-osx-arm64-v$VERSION.tar.gz.sha256"
elif command -v sha256sum &> /dev/null; then
    sha256sum "goinsight-osx-arm64-v$VERSION.tar.gz" > "goinsight-osx-arm64-v$VERSION.tar.gz.sha256"
fi

rm -f "config.yaml.template"
PKG_SIZE=$(du -h "goinsight-osx-arm64-v$VERSION.tar.gz" | cut -f1)
echo "✓ Package created: goinsight-osx-arm64-v$VERSION.tar.gz ($PKG_SIZE)"
echo "✓ Checksum: goinsight-osx-arm64-v$VERSION.tar.gz.sha256"

# ---------- linux amd64 ----------
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 Building: Linux amd64"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

cd "$BACKEND_DIR"
GOOS=linux GOARCH=amd64 \
go build \
  -o "$BIN_DIR/goinsight-linux-amd64" \
  -ldflags "-X main.Version=$VERSION" \
  "$MAIN_GO"

FILE_SIZE=$(du -h "$BIN_DIR/goinsight-linux-amd64" | cut -f1)
echo "✓ Binary built: goinsight-linux-amd64 ($FILE_SIZE)"

cd "$BIN_DIR"
cp "$CONFIG_TEMPLATE" .
tar -czf "goinsight-linux-amd64-v$VERSION.tar.gz" \
  "goinsight-linux-amd64" \
  "config.yaml.template" 2>/dev/null

# 生成 SHA256
if command -v sha256sum &> /dev/null; then
    sha256sum "goinsight-linux-amd64-v$VERSION.tar.gz" > "goinsight-linux-amd64-v$VERSION.tar.gz.sha256"
elif command -v shasum &> /dev/null; then
    shasum -a 256 "goinsight-linux-amd64-v$VERSION.tar.gz" > "goinsight-linux-amd64-v$VERSION.tar.gz.sha256"
fi

rm -f "config.yaml.template"
PKG_SIZE=$(du -h "goinsight-linux-amd64-v$VERSION.tar.gz" | cut -f1)
echo "✓ Package created: goinsight-linux-amd64-v$VERSION.tar.gz ($PKG_SIZE)"
echo "✓ Checksum: goinsight-linux-amd64-v$VERSION.tar.gz.sha256"

# ==========================
# 完成
# ==========================
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Build Completed Successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📁 Output directory: $BIN_DIR"
echo ""
echo "📦 Build artifacts:"
ls -lh "$BIN_DIR" | grep -E '\.tar\.gz$|\.sha256$' | awk '{printf "   %s  %s\n", $5, $9}'
echo ""
echo "✓ Version: $VERSION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
