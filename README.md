# JS Online Editor

ブラウザ上でJavaScriptコードを編集・実行できるオンラインエディタです。

## アーキテクチャ

```
┌─────────────────┐        ┌──────────────────┐        ┌─────────────────────────┐
│   Frontend      │  HTTP  │   Backend        │  HTTP  │   Sandbox               │
│   (Vue.js)      │ ────▶  │   (Go/Gin)       │ ────▶  │   (Node.js/isolated-vm) │
│                 │ ◀────  │                  │ ◀────  │                         │
│   - Monaco Editor        │   - APIサーバー   │        │   - 隔離されたJS実行     │
│   - UI                   │   - リクエスト転送│        │   - メモリ制限128MB      │
└─────────────────┘        └──────────────────┘        │   - タイムアウト10秒     │
                                                       └─────────────────────────┘
```

## セキュリティ

- **isolated-vm**: V8エンジンの隔離されたコンテキストでユーザーコードを実行
- **NetworkPolicy**: サンドボックスPodは外部ネットワークへのアクセスを禁止
- **リソース制限**: メモリ128MB、タイムアウト10秒
- **非rootユーザー**: すべてのコンテナは非rootユーザーで実行

## ディレクトリ構成

```
.
├── frontend/           # フロントエンド (Vue.js)
│   ├── src/
│   ├── package.json
│   └── ...
├── backend/            # バックエンド (Go)
│   ├── main.go
│   ├── handlers/
│   └── go.mod
├── sandbox/            # サンドボックス (Node.js)
│   ├── src/
│   └── package.json
├── docker/             # Dockerfiles
│   ├── Dockerfile.frontend
│   ├── Dockerfile.backend
│   ├── Dockerfile.sandbox
│   └── build-all.sh
└── k8s/                # Kubernetes設定
    ├── deployment.yaml
    ├── service.yaml
    ├── ingress.yaml
    ├── hpa.yaml
    ├── network-policy.yaml
    └── kustomization.yaml
```

## 開発

### ローカル開発

```bash
# フロントエンド
cd frontend
npm install
npm run dev

# バックエンド
cd backend
go run main.go

# サンドボックス
cd sandbox
npm install
npm start
```

### Dockerビルド

```bash
# すべてのイメージをビルド
./docker/build-all.sh v0.0.1
```

### Kubernetesデプロイ

```bash
# kustomizeでデプロイ
kubectl apply -k k8s/
```

## API

### POST /api/execute

JavaScriptコードを実行します。

**Request:**
```json
{
  "code": "console.log('Hello, World!');",
  "language": "javascript"
}
```

**Response:**
```json
{
  "success": true,
  "output": [
    { "type": "log", "text": "Hello, World!" }
  ]
}
```

## 環境変数

### Backend
- `PORT`: APIサーバーのポート (default: 8081)
- `SANDBOX_SERVICE_URL`: サンドボックスサービスのURL (default: http://sandbox-service:3000)

### Sandbox
- `PORT`: サンドボックスサービスのポート (default: 3000)
