#!/bin/bash
#
# Deploy SSE Go - Radar Futebol
# Uso: ./deploy.sh
#

set -e

export PATH=$PATH:/usr/local/go/bin

echo "=========================================="
echo "  DEPLOY SSE GO - $(date '+%Y-%m-%d %H:%M:%S')"
echo "=========================================="

cd /var/www/radarfutebol-sse

echo ""
echo "[1/4] Git pull..."
git pull origin master

echo ""
echo "[2/4] Build (otimizado para producao)..."
# -ldflags="-s -w" remove debug symbols, reduz tamanho do binario
go build -ldflags="-s -w" -o bin/radarfutebol-sse ./cmd/main.go

echo ""
echo "[3/4] Restart systemd..."
systemctl restart sse-go

echo ""
echo "[4/4] Verificando..."
sleep 2
systemctl status sse-go --no-pager | head -10
echo ""
curl -s http://localhost:3005/sse/health

echo ""
echo ""
echo "=========================================="
echo "  DEPLOY CONCLUIDO!"
echo "=========================================="
