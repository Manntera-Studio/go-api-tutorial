# Lesson1: とにかくAPIを立ち上げることを目指す

本レッスンでは、なんでもいいからGo言語で記述したAPIサーバを立ち上げ、リクエストを投げられる状態を目指します。

## 環境構築
まずは、コンテナを立ち上げて中に入りましょう。
コンテナ内に入ったら、Go言語がインストールされていることを確認します。
```
$ docker-compose up -d
$ docker-compose exec golang bash
root@7308f6995af6:/go/src/app# go version
go version go1.16.15 linux/amd64
```

Go言語の場合、最初にプロジェクトの初期化を行い、必要なら関連ライブラリのインストールを行います。
今回はすでに初期化済みなので、本手順はスキップします。
初期化済みの場合は、次回実行時に足りないライブラリが自動でインストールされるはずですが、手動インストールが求められた場合は、メッセージに従ってインストールを行ってください。
```
root@7308f6995af6:/go/src/app# cd lesson1.1
root@7308f6995af6:/go/src/app/lesson1.1# go mod init
go: creating new go.mod: module app/lesson1.1
go: to add module requirements and sums:
        go mod tidy
root@7308f6995af6:/go/src/app/lesson1.1# go mod tidy
```

## Lesson1.1: Hello world
まずは、お約束の Hello world を実行しましょう。

```
root@7308f6995af6:/go/src/app/lesson1.1# go run main.go
Hello world!
```

標準出力に `Hello world!` のメッセージが出力されていれば成功です。

コードの解説を簡単にしておきます。

```golang
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
```

まず、冒頭の `package` についてです。
Go言語ではパッケージと呼ばれる単位でコードをまとめています。
すべてのコードは必ず何かのパッケージに属していなければなりません。

また、Go言語では必ず `main` パッケージを作る必要があります。
そして、 `main` パッケージには `main` 関数が定義されている必要があります。
この `main` パッケージの `main` 関数が呼び出しの出発点(エントリーポイント)になります。

続いて `import` についてです。
その名の通り、外部のパッケージを読み込む処理です。
この例では、標準入出力を扱う `fmt` というパッケージを読み込んでいます。
`fmt` は標準パッケージなので、追加のインストールは不要なため `go.mod` には記載されません。

最後に `main` 関数の中身です。
今回は `fmt.Println("Hello world!")` によって、標準出力に文字列を書き込んでいます。

## Lesson1.2: APIサーバの構築
Go言語の大まかな書き方と実行方法がわかったら、さっそくAPIサーバを立ち上げてみましょう。

```
root@7308f6995af6:/go/src/app/lesson1.2# go run main.go

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v3.3.10-dev
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:5000
```
立ち上がったら、適当なブラウザで http://localhost:5000/user にアクセスしてみてください。
次のようなレスポンスが返ってきます。
```json
{
    Name: "",
    Email: ""
}
```
見にくいようでしたら、json整形用のChrome拡張を入れるときれいになります。
https://chrome.google.com/webstore/detail/json-viewer/gbmdgpbipfallnflgajpaliibnhdgobh

今後はパラメータを渡してみましょう。
http://localhost:5000/user?name=Tanaka%20Taro&email=example@gmail.com

すると、指定したパラメータを反映させたレスポンスが返ってきます。
```json
{
    Name: "Tanaka Taro",
    Email: "example@gmail.com"
}
```

中身の説明です。
```golang
package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Name string
	Email string
}

func main() {
	e := echo.New()
	e.GET("/user", show)

	e.Logger.Fatal(e.Start(":5000"))
}

func show(c echo.Context) error {
    name := c.QueryParam("name")
    email := c.QueryParam("email")

	u := new(User)
	u.Name = name
	u.Email = email

	return c.JSON(http.StatusOK, u)
}
```
今回は、`net/http` という標準パッケージと `echo` という外部パッケージを使用します。

`net/http` はHTTPクライアントとサーバーの実装を提供していて、GET、POSTリクエスト、Formデータの送信など色々できます。

cf. https://free-engineer.life/golang-net-http-simple/

[echo](https://github.com/labstack/echo)  はGo言語の代表的なWebフレームワークです。
軽量でシンプルかつハイパフォーマンスなフレームワークです。
`net/http` だけでも書けなくはないのですが、ルーティングなど　`echo` を使うことでよりシンプルに実装できるようになります。

cf. https://osamu-tech-blog.com/go-echo/

続いて構造体です。
```golang
type User struct {
	Name string
	Email string
}
```
C++などにも登場する構造体と同じです。
ここでは、ひとまず自作の型定義のようなものだと思ってもらえば大丈夫です。
詳しくはLesson2で説明します。

次は `main` 関数です。
```golang
func main() {
	e := echo.New()
	e.GET("/user", show)

	e.Logger.Fatal(e.Start(":5000"))
}
```
`e := echo.New()`でechoパッケージのインスタンスの生成を行います。

`e.GET("/user", show)`でリクエストの種類、エンドポイントとその時に呼び出す処理(関数)を定義します。
今回はGETリクエストを想定しています。

`e.Logger.Fatal(e.Start(":5000"))`でリクエストを待機するポートを指定します。

最後に`show`関数です。
```golang
func show(c echo.Context) error {
    name := c.QueryParam("name")
    email := c.QueryParam("email")

	u := new(User)
	u.Name = name
	u.Email = email

	return c.JSON(http.StatusOK, u)
}
```

まず
```golang
    name := c.QueryParam("name")
    email := c.QueryParam("email")
```
でGETリクエストによって渡されたパラメータを取得します。
```golang
	u := new(User)
	u.Name = name
	u.Email = email
```
でUserのインスタンスを生成し、構造体変数に取得したリクエストパラメータを代入します。
```golang
	return c.JSON(http.StatusOK, u)
```
でステータスコードと構造体を付与してjson形式でレスポンスを返すように定義します。

これでひとまずAPIサーバとして動くようにはなりました。
個人で遊びとして使う分にはこれでも十分でしょう。
ただ、実際にサービスに組み込んで使うとなると圧倒的に不十分です。

Lesson2以降では、より高機能で柔軟性の高いAPIに仕上げるために必要な実装方法やサポートするミドルウェアの使い方について学びます。

## 後片付け
使い終わったら、コンテナをシャットダウンさせましょう。
```
$ docker-compose down
```
