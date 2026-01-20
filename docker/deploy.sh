#!/bin/bash

# 本番環境デプロイスクリプト
# 使用方法: ./deploy.sh [version]
# 例: ./deploy.sh v0.0.1

set -e

VERSION=${1:-v0.0.1}
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "=========================================="
echo "Deploying JS Online Editor - Version: $VERSION"
echo "=========================================="

# 1. Dockerイメージをビルド
echo ""
echo "[1/4] Building Docker images..."
echo "  - Building frontend..."
docker build -t online-editer-frontend:$VERSION -f docker/Dockerfile.frontend .

echo "  - Building backend..."
docker build -t online-editer-backend:$VERSION -f docker/Dockerfile.backend .

echo "  - Building sandbox..."
docker build -t online-editer-sandbox:$VERSION -f docker/Dockerfile.sandbox .

# 2. イメージをtarファイルに保存
echo ""
echo "[2/4] Saving images to tar files..."
mkdir -p docker/images
docker save online-editer-frontend:$VERSION > docker/images/frontend.tar
docker save online-editer-backend:$VERSION > docker/images/backend.tar
docker save online-editer-sandbox:$VERSION > docker/images/sandbox.tar

# 3. MicroK8sにイメージをインポート
echo ""
echo "[3/4] Importing images to MicroK8s..."
microk8s ctr images import docker/images/frontend.tar
microk8s ctr images import docker/images/backend.tar
microk8s ctr images import docker/images/sandbox.tar

# 4. Kubernetesにデプロイ
echo ""
echo "[4/4] Deploying to Kubernetes..."
microk8s kubectl apply -k k8s/

echo ""
echo "=========================================="
echo "Deployment completed successfully!"
echo "=========================================="
echo ""
echo "Deployed images:"
echo "  - online-editer-frontend:$VERSION"
echo "  - online-editer-backend:$VERSION"
echo "  - online-editer-sandbox:$VERSION"
echo ""
echo "Check deployment status:"
echo "  microk8s kubectl get pods"
echo "  microk8s kubectl get services"
