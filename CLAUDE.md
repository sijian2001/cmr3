# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリのコードを操作する際のガイダンスを提供します。

## プロジェクト構成

フロントエンドとバックエンドが明確に分離された CRM（顧客関係管理）システムです：

### フロントエンド (Vue.js + TypeScript)
- **場所**: `frontend/crm3/`
- **フレームワーク**: Vue.js 3 with TypeScript
- **状態管理**: Pinia
- **ルーティング**: Vue Router
- **開発サーバー**: http://localhost:5173/
- **ビルドツール**: Vite

### バックエンド (Go + Echo)
- **場所**: `backend/crm3_api/`
- **言語**: Go 1.25.1
- **フレームワーク**: Echo v4
- **ORM**: GORM (READMEに記載されているが、現在のmain.goには未実装)
- **APIサーバー**: http://localhost:1323/

### データベース
- **種類**: MariaDB
- **データベース名**: crm3
- **ユーザー名**: crm3user
- **パスワード**: 1234

## 開発コマンド

### フロントエンドコマンド (`frontend/crm3/`から実行)
```bash
npm run dev          # 開発サーバーを起動
npm run build        # 本番用ビルド
npm run type-check   # TypeScript型チェック
npm run lint         # ESLintによる自動修正
npm run format       # Prettierによるコード整形
npm run test:unit    # Vitestユニットテスト実行
npm run test:e2e     # PlaywrightE2Eテスト実行
npm run preview      # 本番ビルドプレビュー
```

### バックエンドコマンド (`backend/crm3_api/`から実行)
```bash
go run main.go       # APIサーバーを起動
go build            # バイナリをビルド
go mod tidy         # 依存関係をクリーンアップ
```

## 主要機能（計画中/開発中）

システムは以下の機能を扱うように設計されています：
- **認証システム**: ユーザー登録、ログイン、ログアウト、セッション管理
- **顧客管理**: 検索とページネーション機能付きのCRUD操作
- **製品管理**: 在庫管理機能付きのCRUD操作
- **店舗管理**: ステータス管理機能付きのCRUD操作
- **店員管理**: 店舗との関連付け機能付きのCRUD操作

## 開発ノート

- バックエンドは現在、最小限のルートを持つ基本的なEchoセットアップ (main.go:34-36)
- フロントエンドはComposition APIを使用したモダンなVue.js 3を使用
- テスト環境には、ユニットテスト（Vitest）とE2Eテスト（Playwright）の両方が含まれる
- コード整形はPrettierとESLintで処理
- TypeScriptは開発とビルドプロセスの両方で設定済み