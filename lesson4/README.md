# Lesson4

Lesson1 で API を立ち上げた。
自分一人が趣味で使う分にはこれでも十分だが、ビジネス用途でこれをそのまま使うとサービス要件を満たすのは難しい。

そこで、大抵の場合、 API の手前に Web サーバを置く。

[図解で解説！！　 Apache、Tomcat ってなんなの？ \- Qiita](https://qiita.com/tanayasu1228/items/11e22a18dbfa796745b5)

今回は最近主流になってきている Nginx を使って学習する。

Apache と Nginx の違いに関するより詳しい解説は以下を参考に
[Nginx と Apache って何が違うの？？ \- Qiita](https://qiita.com/hiroaki-u/items/f2455d62f8a4017663cb)

VSCode で Nginx の設定ファイルを書くときは拡張機能を使うと便利である。
https://marketplace.visualstudio.com/items?itemName=raynigon.nginx-formatter

## Lesson 4-1 Nginx でリバースプロキシをする

Nginx 経由で API にアクセスするようにする。
API は Lesson2 で作ったものと同じである。

```bash
$ docker-compose up -d
```

でコンテナを立ち上げたら、まずは、API に直接アクセスしてみる。
8080 ポートで立ち上げているので http://localhost:8080 でアクセスできる。

今度は Nginx 経由でアクセスしてみる。
80 ポートなので http://locahost でアクセスできる。
ユーザ的には違うポートにアクセスしたのに、同じ見た目のページにアクセスできている。

これは Nginx の機能の 1 つであるリバースプロキシによるもの。
[リバースプロキシとは？仕組みをわかりやすく解説 \- カゴヤのサーバー研究室](https://www.kagoya.jp/howto/it-glossary/web/reverse-proxy/)

`proxy/default.conf`で以下のように 80 ポートで受け付けたリクエストを 8080 ポートにリダイレクトさせている。

```conf
server {
    listen 80;
    server_name 127.0.0.1;

    location / {
        proxy_pass http://app:8080;
    }
}
```

また、以下のようにユーザのアクセスログを取ることもできる。

```bash
$ docker-compose logs -f
proxy    | 172.20.0.1 - - [02/Aug/2022:16:16:51 +0000] "GET / HTTP/1.1" 200 23 "-" "Mozilla/X.X (Windows NT XX.X; WinXX; x64) AppleWebKit/XXX.XX (KHTML, like Gecko) Chrome/XXX.X.X.X Safari/XXX.XX" "-" "0.000"
```

アクセスログを残しておくことでアクセスの多いリクエストやレスポンスに時間がかかっているリクエスト、40X,50X になったリクエスト、外部からの攻撃などいろいろな用途の分析の手掛かりになる。

## Lesson 4-2 Nginx でリクエストの振り分けをする

大量のリクエストがある API だと、すべてのリクエストを 1 台のサーバで捌こうとするとパンクしてしまう。
そのため、API の手前にロードバランサーと呼ばれる負荷分散装置を置く。

[ロードバランサーとは？意味・定義 \| IT トレンド用語 \| NTT コミュニケーションズ](https://www.ntt.com/bizon/glossary/j-r/load-balancer.html)

今回は Nginx の機能を使ってロードバランシングをする。

構成は 2 つの Nginx の手前にロードバランサーとしての Nginx を置いたものである。
cf. [Docker でロードバランサ・アプリケーションサーバ・DB サーバの環境構築 \- A Memorandum](https://blog1.mammb.com/entry/2019/11/01/215930)

まずはコンテナを立ち上げる。

```bash
$ docker-compose up -d
```

http://localhost で手前のロードバランサーとしての Nginx 経由で後ろの Nginx にアクセスしてみる。
何度かページをリロードすると「Hello nginx1」と「Hello nginx2」の 2 種類のメッセージが表示されていることがわかる。

何度かページをリロードすると「Hello nginx1」のときは後ろ側の 1 番機の Nginx に「Hello nginx2」のときは 2 番機の Nginx にアクセスしている。

ユーザとしては同じ URL でアクセスしているつもりだが、裏側では 2 台に冗長化させた片方のサーバにアクセスするよう振り分けられている。

この仕組みを応用すると AB テストができるようになる。
詳しくは Lesson5 で扱う。

## Lesson 4-3 ロードバランシングで負荷分散させる

ロードバランシングによって負荷軽減されているかを確かめる。
構成は Lesson4-1 とほぼ同じで API を 2 台に増やしたものである。

まずはコンテナを立ち上げる。

```bash
$ docker-compose up -d
```

http://localhost?name=hoge&email=hoge@example.com にアクセスと API にアクセスできていることを確かめる。

ログを見ると app1 または app2 のいずれかにリバースプロキシされている。

```bash
$ docker-compose logs -f
proxy    | 172.22.0.1 - - [02/Aug/2022:16:48:00 +0000] "GET /?name=hoge&email=hoge@example.com HTTP/1.1" 200 43 "-" "Mozilla/X.X (Windows NT XX.X; Win64; x64) AppleWebKit/XXX.XX (KHTML, like Gecko) Chrome/XXX.0.0.0 Safari/XXX.XX" "-" "0.001"
app1     | 2022/08/02 16:48:04 params: hoge,hoge@example.com
proxy    | 172.22.0.1 - - [02/Aug/2022:16:48:04 +0000] "GET /?name=hoge&email=hoge@example.com HTTP/1.1" 200 43 "-" "Mozilla/X.X (Windows NT XX.X; Win64; x64) AppleWebKit/XXX.XX (KHTML, like Gecko) Chrome/XXX.0.0.0 Safari/XXX.XX" "-" "0.000"
app2     | 2022/08/02 16:48:06 params: hoge,hoge@example.com
proxy    | 172.22.0.1 - - [02/Aug/2022:16:48:06 +0000] "GET /?name=hoge&email=hoge@example.com HTTP/1.1" 200 43 "-" "Mozilla/X.X (Windows NT XX.X; Win64; x64) AppleWebKit/XXX.XX (KHTML, like Gecko) Chrome/XXX.0.0.0 Safari/XXX.XX" "-" "0.001"
app1     | 2022/08/02 16:48:07 params: hoge,hoge@example.com
proxy    | 172.22.0.1 - - [02/Aug/2022:16:48:07 +0000] "GET /?name=hoge&email=hoge@example.com HTTP/1.1" 200 43 "-" "Mozilla/X.X (Windows NT XX.X; Win64; x64) AppleWebKit/XXX.XX (KHTML, like Gecko) Chrome/XXX.0.0.0 Safari/XXX.XX" "-" "0.001"
```

また、 http://localhost:8080?name=hoge&email=hoge@example.com や http://localhost:8081?name=hoge&email=hoge@example.com のように直接個別の API にもアクセスできるようにしてある。

これに対して負荷試験を行う。
今回は負荷試験ツールとして Taurus(トーラス)を使用する。
https://gettaurus.org/

日本語で Taurus の使い方が読みたければ、以下が参考になる。
[手軽に負荷テストができるツール「Taurus」がスゴい](https://zenn.dev/tonchan1216/articles/11afd147ea3dd2734315)

まずは、app1 に直接リクエストを飛ばしてみる。

```bash
$ cd taurus
$ sh run.sh d
```

負荷設定は以下のとおりである。

```yml
execution:
  - concurrency: 10 # 同時接続数
    ramp-up: 10s # ramp-up(concurrency数に到達するまでの)時間
    hold-for: 2m # 試験時間
    throughput: 130 # 最大スループット(RPS: Request per second)
    scenario: quick-test # シナリオ名

scenarios:
  quick-test:
    requests:
      - http://host.docker.internal:8080?name=hoge&email=hoge@example.com
```

後半負荷に耐えられず一部のリクエストが捌き切れていない。

今後はロードバランサー経由で負荷をかけてみる。

```bash
$ sh run.sh
```

先ほどと違ってすべてのリクエストを捌き切れている。
また、CPU や Memory にもちょっとだけ優しくなっている。

あまり差が見られなければ負荷を増やしたり、試験時間を延ばすとわかりやすいかもしれない。

## その他補足資料

- [nginx\.conf の中身を理解したいので一つずつ調べました \- Qiita](https://qiita.com/na0kiB/items/a8b081fe30ff1c6d99a9)
- [Docker Compose でカーネルパラメータを設定する – ymyzk’s blog](https://blog.ymyzk.com/2017/01/docker-compose-sysctls/)
- [Docker コンテナ内の net\.core\.somaxconn をいじる方法 \- Qiita](https://qiita.com/ma2shita/items/f1a68a3f909c5cee7869)
- [負荷が低いのにアクセスを捌けきれない時の対応 \- Carpe Diem](https://christina04.hatenablog.com/entry/2016/12/31/124314)
- [Taurus を用いて JMeter のレポートを作成する \- Qiita](https://qiita.com/hmsnakr/items/a64507c31d1365dd6bb0)
