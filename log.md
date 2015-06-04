# ログデータ分析

## Google Analyticsでは見れないデータ
#### Google Analyticsの仕組み
ページに埋め込まれたJavascriptを実行することで、
Googleのサーバーへリクエストが送られ、そのリクエストを解析することで記録している。

そのため下記のような項目は見ることができない

- ステータスコード、メソッドなど、GA上で見れない項目
→ 設定次第で取得することも可能（だと思う）
- 画像やcssなどのページではないファイル
- Javascriptを実行できないブラウザ
→ 仕組み上不可能
JSを実行できないブラウザ = 検索エンジンなどのクローラー



# ログの形式

サーバーで設定された形式

#### 代表的な形式

`Logformat "%h %l %u %t ¥"%r¥" %>s %b ¥"%{Referer}i¥" ¥"%{User-Agent}i¥"" combined`

#### 実際に記録されるログ

```
66.249.73.209 - - [24/Feb/2015:00:00:14 +0900] "GET /biyoshi/feature/opening/?page=23&type=area_list&sort_mode=&dir_name=biyoshi&t_max=10 HTTP/1.1" 301 20 "-" "Mozilla/5.0 (iPhONe; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)" "150900"
```

|記録される内容|説明|BQのカラム名|
|-----|-----|-----|
|66.249.73.209  |   リモートホスト名|remote_hots|
|-  |   クライアントの識別子|identd|
|-  |   認証ユーザー名|user|
|[24/Feb/2015:00:00:14 +0900]   |   時刻|request_time|
 |"GET /biyoshi/feature/（略） HTTP/1.1"   |   リクエストの最初の行|method, url, protocol|
|301    |   レスポンスステータス|status|
|20 |   送信されたバイト数(ヘッダーは含まず)。0バイトの時は「-｣|res_size|
 |"-"   |   リファラー| referer|
|Mozilla/5.0 (iPhONe; （略）   |   ユーザーエージェント|user_agent|
|150900 |   リクエストを処理するのにかかった時間、マイクロ秒単位|required_time|

# ログの量

アクセスログだけでも一日あたり約150万行

・Excel
不可能！

・データベース
不向き！

→__BigQueryを使おう！！__

# ログの分析例

- Googlebot
- 404を吐いてるURLを調べる

などなど

## 課題

- 5/13～5/15の最も遅いページ上位100件

- 各日のログの件数

- Fauraで人気の記事トップ50

- Googlebotを除いた、5/15のPCリジョブのユーザー数

#BQを使わない方法

Hadoopで分散処理とかいろいろあるけれど・・・

## grep

```shell
$grep -i 'googlebot' logdata_0501.txt | wc -l
#GoogleBotからのアクセス数をカウント
```

## awk

```shell
$awk '{
    status[$9] += 1
}
END{
    for(code in status){
        print code,":",status[code]
    }
}' logdata_0501.txt
# 各ステータスコードの数をカウントする
```

awkはテキスト処理に便利！  
ちなみに今回BigQueryへ流し込むためにawkでフォーマットを変換し、１万行ずつ分割  
シェルスクリプトでインサートしてます。

<img src="./img/sed-awk.jpg" style="width: 100px;">
私の席にあるので読みたい人がいたらぜひ
（分析に使うツールは多ければ多いほどいいよね）