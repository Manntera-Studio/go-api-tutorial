# Lesson6 データベースと連携する

REST API の設計原理にもあるようにスケールアウトなどの観点から、基本的に API はステートレス(状態を持たない)である。

そこで、今回は値の実態の保存先としてデータベース(MySQL)を用意し、API を通してデータのやり取りをできる仕組みについて学ぶ。

cf. [docker\-compose で golang と MySQL を繋ぐ](https://zenn.dev/ajapa/articles/443c396a2c5dd1)

## Lesson6-1 データベースの下準備

まずは、データベースのスキーマを設計する。
今回は架空のレースデータを格納するデータベースを想定する。

スキーマが決まったら SQL に書き起こす。
詳細は `mysql/init/10_race.sql` を参照。

`mysql/init` 内の拡張子が `.sql`, `.sh`, `.sql.gz` のファイルはコンテナ起動時に自動的に実行される。
init の中身は上から順に実行されるので順序が大事な場合は、その順になるようにファイル名を命名しておくこと。
シェルスクリプトからは`${環境変数名}` で、`.env` に記述した環境変数を読み込める。

具体的な中身は `mysql/init/race.csv` に記述しておく。

なお、MySQL は version8 からセキュリティ上の観点からデフォルトだとローカルの CSV からデータがインポートできなくなっている。

[docker\-compose ＋ MySQL8（8\.0\.18）で初期データを CSV ロードしようとするとエラー（The used command is not allowed with this MySQL version）に \- Qiita](https://qiita.com/You_name_is_YU/items/6d87f7664c947df84dc1)

そのため、SQL には

```sql
SET　GLOBAL local_infile = 1;
```

を定義しておき、シェルスクリプトから `--local-infile=1` オプションをつけてインポートする。

準備ができたらコンテナを起動する。

```bash
$ docker-compose up -d
```

WSL2 の場合、以下のようなエラーに遭遇することがある。

```bash
$ docker-compose logs mysql
mysql_1   | mysqld: Cannot change permissions of the file 'private_key.pem' (OS errno 1 - Operation not permitted)
mysql_1   | 2020-09-10T05:04:53.449233Z 0 [ERROR] [MY-010295] [Server] Could not set file permission for private_key.pem
mysql_1   | 2020-09-10T05:04:53.449629Z 0 [ERROR] [MY-010119] [Server] Aborting
```

そのため、WSL の設定ファイルでマウントオプションを定義しておく。

```bash
$ sudo vim /etc/wsl.conf
[automount]
options = "metadata"
```

[WSL2 上の Docker で MySQL を構築する際の permissions の対策 \- Qiita](https://qiita.com/n-jun-k2/items/f8f36cebc7312df8bc31)

設定変更出来たらコンテナを再度起動してみる。
このとき、マウントした `mysql/lib` ディレクトリにデータが残っているなら、フォルダごと削除しておくこと。

```bash
$ sudo rm -rf mysql/lib
$ docker-compose up -d
```

コンテナが正常に立ち上がったら、データベースの中を覗いてみる。

```sql
$ docker-compose exec db mysql -utest_user -ppassword -e "select * from race.race;"
mysql: [Warning] Using a password on the command line interface can be insecure.
+----+------------+------------+---------------------+---------------------+-------------+-------+---------------------+
| id | day_raceid | race_date  | start_time          | end_time            | temperature | venue | updated             |
+----+------------+------------+---------------------+---------------------+-------------+-------+---------------------+
|  1 |          1 | 2022-08-01 | 2022-08-01 10:00:00 | 2022-08-01 11:00:00 |        27.5 | tokyo | 2022-08-07 04:21:21 |
|  2 |          2 | 2022-08-01 | 2022-08-01 11:00:00 | 2022-08-01 12:00:00 |        27.5 | tokyo | 2022-08-07 04:21:21 |
|  3 |          1 | 2022-08-01 | 2022-08-01 10:00:00 | 2022-08-01 11:00:00 |        27.5 | chiba | 2022-08-07 04:21:21 |
+----+------------+------------+---------------------+---------------------+-------------+-------+---------------------+
```

ちゃんと CSV で定義した内容が入っていることが確認できる。

## Lesson6-2 API の実装

MySQL と Go による API が起動出来るようになったら、API から MySQL を制御する部分を実装する。

今回は、

1. 取得(`select`) -> GET
1. 追加(`insert`) -> POST
1. 更新(`update`) -> PUT
1. 削除(`delete`) -> DELETE

の 4 機能を実装する。

動作確認については、ターミナルから `curl` コマンドを叩いてもいいのだが、慣れないうちは GUI から操作できる[postman](https://chrome.google.com/webstore/detail/postman/fhbjgbiflinjbdggehcddcbncdddomop/related?hl=ja)を使うと便利である。

### 取得

まずは、取得の動作を確認する。

大まかな処理の流れは、

- `handler.go` でリクエストパラメータを所定のフォーマットにパース
- `race.go`　でパラメータをもとに SQL を組み立て、取得結果を 1 行ずつリストに格納して返却する

である。
このとき、 データベースからの取得結果の受取先として `type Race struct`　のようにデータベースのスキーマに合わせた構造体を用意しておくのがミソである。

使い方は、GET Method で以下の URL にリクエストを送る。
http://localhost:8080/
すると、以下のようなレスポンスが得られる。
このようにデータベースに登録したデータが API 経由で取得できている。

```json
[
  {
    "Id": 1,
    "DayRaceId": 1,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T10:00:00Z",
    "EndTime": "2022-08-01T11:00:00Z",
    "Temperature": 27.5,
    "Venue": "tokyo",
    "Updated": "2022-08-07T15:25:43Z"
  },
  {
    "Id": 2,
    "DayRaceId": 2,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T11:00:00Z",
    "EndTime": "2022-08-01T12:00:00Z",
    "Temperature": 27.5,
    "Venue": "tokyo",
    "Updated": "2022-08-07T15:25:43Z"
  },
  {
    "Id": 3,
    "DayRaceId": 1,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T10:00:00Z",
    "EndTime": "2022-08-01T11:00:00Z",
    "Temperature": 27.5,
    "Venue": "chiba",
    "Updated": "2022-08-07T15:25:43Z"
  }
]
```

また、条件を指定してデータを取得することもできる。
どのようなリクエストパラメータが使えるかは、 `handler.go` の `Select()` 参照。
http://localhost:8080?id=1

```json
[
  {
    "Id": 1,
    "DayRaceId": 1,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T10:00:00Z",
    "EndTime": "2022-08-01T11:00:00Z",
    "Temperature": 27.5,
    "Venue": "tokyo",
    "Updated": "2022-08-07T15:25:43Z"
  }
]
```

http://localhost:8080?date=2022-08-01&stime=10:00:00

```json
[
  {
    "Id": 1,
    "DayRaceId": 1,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T10:00:00Z",
    "EndTime": "2022-08-01T11:00:00Z",
    "Temperature": 27.5,
    "Venue": "tokyo",
    "Updated": "2022-08-07T15:25:43Z"
  },
  {
    "Id": 3,
    "DayRaceId": 1,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T10:00:00Z",
    "EndTime": "2022-08-01T11:00:00Z",
    "Temperature": 27.5,
    "Venue": "chiba",
    "Updated": "2022-08-07T15:25:43Z"
  }
]
```

### 追加

次に API 経由でデータベースに新しいレコードを追加する。

大まかな処理の流れは GET のときと同じである。

- `handler.go` でリクエストパラメータをパース
- `race.go` でクエリを組み立てて、データベースに書き込み

このとき、自動インクリメントされる `id` と デフォルト値で書き込み日時を指定している `updated` に関しては自動入力に任せて書き込みパラメータから外している。
また、必須のパラメータが欠けている場合は Bad Request として返却するようにしている。

実装ができたら、POST Method を選択して各種パラメータを指定する。
ex. localhost:8080?drid=2&date=2022-08-07&stime=12:00:00&venue=osaka&etime=13:00:00

正常に insert できていれば、select をすると insert したデータが取得できるようになっている。

```json
[
  {
    "Id": 1,
    "DayRaceId": 1,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T10:00:00Z",
    "EndTime": "2022-08-01T11:00:00Z",
    "Temperature": 27.5,
    "Venue": "tokyo",
    "Updated": "2022-08-07T15:25:43Z"
  },
  {
    "Id": 2,
    "DayRaceId": 2,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T11:00:00Z",
    "EndTime": "2022-08-01T12:00:00Z",
    "Temperature": 27.5,
    "Venue": "tokyo",
    "Updated": "2022-08-07T15:25:43Z"
  },
  {
    "Id": 3,
    "DayRaceId": 1,
    "RaceDate": "2022-08-01T00:00:00Z",
    "StartTime": "2022-08-01T10:00:00Z",
    "EndTime": "2022-08-01T11:00:00Z",
    "Temperature": 27.5,
    "Venue": "chiba",
    "Updated": "2022-08-07T15:25:43Z"
  },
  {
    "Id": 4,
    "DayRaceId": 1,
    "RaceDate": "2022-08-07T00:00:00Z",
    "StartTime": "2022-08-07T12:00:00Z",
    "EndTime": "2022-08-07T13:00:00Z",
    "Temperature": 0,
    "Venue": "osaka",
    "Updated": "2022-08-07T08:56:07Z"
  },
  {
    "Id": 5,
    "DayRaceId": 2,
    "RaceDate": "2022-08-07T00:00:00Z",
    "StartTime": "2022-08-07T12:00:00Z",
    "EndTime": "2022-08-07T13:00:00Z",
    "Temperature": 0,
    "Venue": "osaka",
    "Updated": "2022-08-07T18:01:24Z"
  }
]
```

### 更新

続いて、登録済みのデータの更新を行う。

処理の流れは、select, insert とほぼ同じ。

例えば以下のデータの Temperature を更新するケースを考える。
http://localhost:8080/?id=5

```json
[
  {
    "Id": 5,
    "DayRaceId": 2,
    "RaceDate": "2022-08-07T00:00:00Z",
    "StartTime": "2022-08-07T12:00:00Z",
    "EndTime": "2022-08-07T13:00:00Z",
    "Temperature": 0,
    "Venue": "osaka",
    "Updated": "2022-08-07T18:01:24Z"
  }
]
```

PUT Method で以下のよう名リクエストを送る。
localhost:8080/5?temperature=20.8

Temperature と Updated が更新されている。

```json
[
  {
    "Id": 5,
    "DayRaceId": 2,
    "RaceDate": "2022-08-07T00:00:00Z",
    "StartTime": "2022-08-07T12:00:00Z",
    "EndTime": "2022-08-07T13:00:00Z",
    "Temperature": 20.8,
    "Venue": "osaka",
    "Updated": "2022-08-07T22:08:31Z"
  }
]
```

### 削除

最後に delete を行う。

流れはこれまで 3 つとほぼ同じ。

例えば、id が 5 番のデータを削除することを考える。
Delete Method にして localhost:8080/5 というリクエストを送る。

```json
{
  "Message": "delete success!"
}
```

というメッセージが返却されていれば、正しくできている。
http://localhost:8080/?id=5 で検索すると確かになくなっている。

```json
null
```

## まとめ

Go の API を通じて MySQL を操作する方法を学んだ。
基本的には、リクエストパラメータを受け取ってパースし、リクエストメソッドに応じた SQL を組み立ててデータベースに送るという流れであった。
API をステートレスにしておくことで、セッション管理が不要になり、シンプルでスケールアウトしやすい形になった。
実際には Lesson4,5 と組み合わせて手前に Nginx を置いたり、API を複数台並べてスケールアウトさせたりする。

今回はあまり気を配らなかったが、ユーザリクエストに応じてデータベースを操作するときは SQL インジェクションに注意しなくてはならない。
[SQL インジェクションとは \| 具体的な攻撃例と対策方法 \- IT を分かりやすく解説](https://medium-company.com/sql%E3%82%A4%E3%83%B3%E3%82%B8%E3%82%A7%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3/)

また、今回は API の設計も甘いので、自分で実装する際は、より見通しの良い、効率的な実装になるよう洗練してほしい。

## References

- [【Golang】MySQL の基本的な操作を行う \- 中堅プログラマーの備忘録](https://www.chuken-engineer.com/entry/2021/09/24/162120)
- [Go database/sql の操作ガイドあったんかい](https://sourjp.github.io/posts/go-db/)
- [Go database/sql チュートリアル 01 \- 概要 · golang\.shop](https://golang.shop/post/go-databasesql-01-overview-ja/)
- [echo の API サーバ実装とエラーハンドリングの落とし穴 \- Qiita](https://qiita.com/usk81/items/5f2bcfe06eb83830ee55)
