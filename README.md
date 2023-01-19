# Go API Tutorial

本リポジトリは Go 言語を使った API の構築手順を学ぶチュートリアルです。

## Lessons

1. [Lesson1: とにかく API を立ち上げることを目指す](./lesson1/)
   1. Hello world
   1. API サーバの構築
1. Lesson2: パッケージ化 coming soon...
   - 構造体
   - 自作パッケージを読み込む
1. Lesson3: 単体テスト coming soon...
1. [Lesson4: 大量のリクエストを捌けるようにする](./lesson4/)
   1. Nginx でリバースプロキシをする
   1. Nginx でリクエストの振り分けをする
   1. ロードバランシングで負荷分散させる
1. Lesson5: API のリクエスト制御をする coming soon...
   - IP アドレスや Cookie によるアクセス制御
   - Basic 認証
   - A/B Test
1. [Lesson6: データベースと連携する](./lesson6/)
   1. データベースの下準備
   1. API の実装
1. [Lesson7: API 仕様書を書く](./lesson7/)
   1. ソースコードコメントから OpenAPI ドキュメントを自動生成する
   1. 複雑な API のドキュメントを書く
   1. OpenAPI ドキュメントからソースコードを自動生成する
1. Lesson8: スケールアウト coming soon...
   - Kubernetes
1. Lesson9: 監視設定 coming soon...
   - Prometheus, Grafana, Alertmanager

### [番外編(Appendix)](./appendix/)

- GAE へのデプロイ
- ビルド coming soon...
- デーモン化 coming soon...
- 性能試験 coming soon...

## 注意事項

本チュートリアルは version 1.16 以降を想定しています。
1.15 位前を使う場合は、適宜読み替えてください。

また、環境構築は Docker を使って行う想定です。
あらかじめ `docker` および `docker-compose` コマンドが使える状態を用意してください。

Docker を使わない場合は、各章の `Dockerfile` や `docker-compose.yml` を参考に環境構築を行ってください。

## その他のドキュメント

### 公式チュートリアル

https://go-tour-jp.appspot.com/

### コーディング規約

- 公式
  - https://go.dev/doc/effective_go
  - https://github.com/golang/go/wiki/CodeReviewComments
- 日本語訳(非公式)
  - http://go.shibu.jp/effective_go.html
  - https://knsh14.github.io/translations/go-codereview-comments/
  - https://zenn.dev/kenghaya/articles/1b88417b1fa44d
    - あまりオススメしないけど、よく使うものをパッと見たいなら
