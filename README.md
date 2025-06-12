---
title: GitLab から GitHub への移行
source: 
author: 
published: 
created: 2025-06-12
description: GitLab から GitHub への移行に関するメモと GitHub Actions の概要
tags: [GitHub, GitLab, CI/CD, GitHub Actions]
---
#### GitHubの移行理由
プラン: GitHub Enterprise,Copilot Business

- セキュリティ強化
	- GitHub Advanced Security
		- CodeQL
			- コードの静的解析
		- シークレットスキャン
			- コードに含まれるシークレットの検知
		- Dependency Review
			- 依存関係のあるパッケージの脆弱性のスキャン
- コスト

#### GitLab -> GitHub へのリポジトリ移行

```sh
git clone --mirror {repo url}
cd {プロジェクト}
git remote set-url origin {new repo url}
git push --mirror
```


#### GitLab -> GitHub へのリポジトリ移行に伴う変更点
- CIのファイルの変更
- E2Eテストのコード変更



####  CIのファイルの変更
GitHubでCIを実行したい場合は、`.github/workflows`配下にCIのファイルを配置しています。
GitHubのCIはGitHub Actionsと呼びます。

### GitHub Actions ワークフロー管理のツリーマップ

```mermaid
flowchart TD
    A[GitHub Actions 管理] --> B[ワークフローファイル]
    B --> BA[.github/workflows/*.yml]
    BA --> BA1[基本構造]
    BA1 --> BA1A[name: ワークフロー名]
    BA1 --> BA1B[on: トリガー定義]
    BA1 --> BA1C[jobs: ジョブ定義]
    
    
    BA1B --> BB1[push]
    BA1B --> BB2[pull_request]
    BA1B --> BB3[schedule]
    BA1B --> BB4[workflow_dispatch]

    A --> C[CI環境変数]
    
    C --> E[シークレット管理]
    E --> EA[リポジトリシークレット]
    E --> EB[環境シークレット]
    E --> EC[Organization シークレット]

    C --> F[変数の管理]
    F --> FA[リポジトリ変数]
    F --> FB[環境変数]
    F --> FC[Organization 変数]
    
```

### GitHub Actions ワークフロー

```mermaid
flowchart TD
A[GitHub Actions 管理] --> B[ワークフローファイル]
B --> BA[.github/workflows/*-deploy.yml]
B --> BB[.github/workflows/unit-test.yml]
B --> BC[.github/workflows/e2e-test.yml]
B --> BD[.github/workflows/api-test.yml]
B --> BE[.github/workflows/dependency-review.yml]
B --> BF[.github/workflows/base64.yml]

BA --> BBA[コンテナイメージをビルドするアクション]
BB --> BBB[単体テストアクション]
BC --> BBC[E2Eテストアクション]
BD --> BBD[APIテストアクション]
BE --> BBE[パッケージの脆弱性チェックアクション]
BF --> BBF[コードにbase64が含まれていないかをチェックするアクション]

```

```yaml
name: CI/CD パイプライン
on:
  push:
    branches: [ main, develop ]
jobs:
  build:
    runs-on: vm-runner
    steps:         
      - name: 依存関係インストール
        run: npm ci
        
      - name: リント実行
        run: npm run lint
        
      - name: テスト実行
        run: npm test
        
      - name: ビルド実行
        run: npm run build
```

##### GitHub Actions (コンテナイメージのビルド)
```mermaid
sequenceDiagram
    開発者 ->> GitHub: ソースコードをプッシュ
    GitHub ->> GitHub Actions: プッシュをトリガーにCIを実行
    GitHub Actions ->> ECR: コンテナイメージをビルド、プッシュ
    ECR--> 開発者: CIの実行ログ
    ECR ->> Kubernetes: Kubernetesにデプロイ
    Kubernetes->開発者: 各環境(dev,stg,prod)のシステムを確認
```

##### GitHub Actions (単体テスト)
```mermaid
sequenceDiagram
    開発者 ->> GitHub: プルリクエストの作成
    GitHub ->> GitHub Actions: プルリクエストをトリガーにCIを実行
    GitHub Actions ->> QAシステム: テスト結果を登録
    GitHub Actions--> 開発者: CIの実行ログ
    QAシステム-> 開発者: テスト結果の確認
```

##### GitHub Actions (APIテスト,E2E)
```mermaid
sequenceDiagram
    開発者 ->> Kubernetes: 最新のコンテナイメージのデプロイ
    Kubernetes ->> ArgoCD: ArgoCDのステータスがHealth
    ArgoCD ->> GitHub: GitHubのAPIを実行
    GitHub ->> GitHub Actions: GitHubのAPI経由でCIを実行
    GitHub Actions ->> QAシステム: テスト結果を登録
    GitHub Actions--> 開発者: CIの実行ログ
    QAシステム-> 開発者: テスト結果の確認
```


####  E2Eテストのコード変更
GitLab RunnerのCIでのみ使用できる環境変数をGitHub用に変更した。
