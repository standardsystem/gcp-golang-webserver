# gcp-cloudrun-gcs-golang-webserver
GCS上の指定バケットをWebとして公開する簡易サーバコンテナ

## 環境変数

### BUCKET
コンテンツをホストするGCSのバケットパス

### PORT
Listenするポート

## 認証について
ローカルで実行する場合には、サービスアカウントから生成する認証情報jsonファイルを、ローカルマシンの適当な場所に配置してください。
GCP上で稼働させるときは、VMやCloud Run の稼働サービスアカウントで認証し、このファイルは使わないでください。

ローカルでの実行には、/secrets フォルダに　credential.json というファイル名で運用するのがおすすめです。
環境変数 `GOOGLE_APPLICATION_CREDENTIALS` に認証情報ファイルを指定してください。

/home/username/gcp-cloudrun-gcs-golang-webserver/secrets/ に credential.json を置いた場合の実行例：

```
>docker build . -t easy-web
>docker run -d -p 80:80 -e PORT=80 -v /home/username/gcp-cloudrun-gcs-golang-webserver/secrets/credential.json:/secrets -e BUCKET=golang-easy-webserver-test -e GOOGLE_APPLICATION_CREDENTIALS=/secrets --name web easy-web
```

## credential file について
ターゲットBUCKETに対してのRead権限があるサービスアカウントの認証ファイルを用いてください。