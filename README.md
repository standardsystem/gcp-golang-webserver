# gcp-golang-webserver

GCS上の指定バケットをWebとして公開する簡易サーバコンテナ

## 環境変数

### BUCKET
コンテンツをホストするGCSのバケットパス

### PORT
Listenするポート

## セットアップ手順

以下の手順で、Google Cloud 環境にウェブサーバーを構築できます。

### 1. プロジェクトの設定
```bash
export PROJECT_ID="あなたのプロジェクトID"
export BUCKET_NAME="あなたのバケット名"
export REGION="asia-northeast1"

gcloud config set project $PROJECT_ID
```

### 2. GCS バケットの作成とコンテンツのアップロード
```bash
# バケットの作成
gcloud storage buckets create gs://$BUCKET_NAME --location=$REGION

# インデックスファイルのアップロード (例)
echo "<h1>Hello from GCS</h1>" > index.html
gcloud storage cp index.html gs://$BUCKET_NAME/index.html
```

### 3. Cloud Run へのデプロイ (IAP 有効化)
初回デプロイ時に `--iap` フラグを指定することで、ロードバランサなしで認証を有効にした状態でデプロイできます。

```bash
# ビルドとデプロイ (IAPを直接有効化)
gcloud beta run deploy easy-webserver \
  --source . \
  --region $REGION \
  --set-env-vars BUCKET=$BUCKET_NAME \
  --iap \
  --no-allow-unauthenticated
```

### 4. 権限の設定 (重要)
Cloud Run のデフォルトのサービスアカウントが GCS バケットを読み取れるように権限を付与します。
```bash
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format='value(projectNumber)')

gcloud storage buckets add-iam-policy-binding gs://$BUCKET_NAME \
  --member="serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role="roles/storage.objectViewer"

# 署名付きURL生成（大容量ファイル用）のための権限付与
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role="roles/iam.serviceAccountTokenCreator"
```

## ローカルでの実行

ローカルで実行する場合には、サービスアカウントから生成する認証情報jsonファイルを、ローカルマシンの適当な場所に配置してください。
環境変数 `GOOGLE_APPLICATION_CREDENTIALS` にそのファイルを指定します。

```bash
docker build . -t easy-web
docker run -d -p 8080:8080 \
  -e PORT=8080 \
  -e BUCKET=あなたのバケット名 \
  -v /path/to/credential.json:/secrets/credential.json \
  -e GOOGLE_APPLICATION_CREDENTIALS=/secrets/credential.json \
  --name web easy-web
```

## credential file について
ターゲットBUCKETに対してのRead権限があるサービスアカウントの認証ファイルを用いてください。

## ユーザー認証 (IAP)

このサーバーは Google Cloud Identity-Aware Proxy (IAP) による認証に対応しています。IAPを有効にすることで、IAMに登録されたユーザーのみにアクセスを制限できます。

### 5. IAP の権限設定 (gcloud コマンド)

デプロイ時に `--iap` を指定した後は、アクセスを許可するユーザーに対して権限を付与するだけです。

```bash
# ユーザーにアクセス権限を付与
export USER_EMAIL="user@example.com"
gcloud iap web add-iam-policy-binding \
  --member="user:$USER_EMAIL" \
  --role="roles/iap.httpsResourceAccessor"
```

> [!TIP]
> 既存のサービスに対して後から IAP を有効にする場合は、以下のコマンドを使用します：
> `gcloud beta run services update easy-webserver --region $REGION --iap`

### トラブルシューティング: エラーコード 11

IAP 有効化後に「Error code 11」が表示される場合、OAuth クライアントの設定が正しく初期化されていない可能性があります。

以下の手順を試してください：

1. **OAuth 同意画面を確認**:
   - [Google Cloud Console > API とサービス > OAuth 同意画面](https://console.cloud.google.com/apis/credentials/consent) を開き、設定が完了していることを確認してください。
   - 「内部」または「外部」のユーザータイプが適切に選択されているか確認します。
2. **IAP の無効化と再有効化**:
   - IAP の設定が反映されるまで時間がかかる場合があります。一度 IAP をオフにして数分待ち、再度オンにすることで解消することが多いです。
   ```bash
   # 一度オフにする
   gcloud beta run services update easy-webserver --region $REGION --no-iap
   # 5分ほど待ってから再度オンにする
   gcloud beta run services update easy-webserver --region $REGION --iap
   ```
4. **IAP サービスエージェントの作成**:
   - `Service account ...@gcp-sa-iap.iam.gserviceaccount.com does not exist` というエラーが出る場合、以下のコマンドで IAP 用のサービスエージェントを手動で作成する必要があります。
   ```bash
   gcloud beta services identity create --service=iap.googleapis.com --project=$PROJECT_ID
   ```

