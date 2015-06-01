新卒研修 - マーケ/データ分析偏 (編集中

#SQLおさらい
BigQueryの前に・・・

## データベースの構造
#### データベース
テーブルの集合

#### テーブル
レコードの集合

#### レコード
一つ一つのデータ
テーブルの各行

#### カラム（フィールド）
データを構成する情報
テーブル、レコードの各列

#### ・・・Excelで言うと
データベース = ファイル自体
テーブル = シート
レコード = 行(1～)
カラム = 列(A～)


## 主な命令文
- UPDATE - データの更新
- INSERT - データの作成（挿入)
- DELETE - データの削除
→覚えなくてよし！(マーケでは)

- SELECT
データの取得
→__覚えよう！！！__

## 実際のSQL文を簡単に復習
- 基本

```mysql
SELECT # 表示するカラムを指定する
    id, in_day, shop_id
FROM # 使用するテーブル指定する
    scout_log
WHERE # 抽出する条件
    in_day BETWEEEN '2015-01-01' AND '2015-05-01'
;
```

- 結合

```mysql
SELECT
    sl.id, sl.in_day, m.id, m.name, m.city
FROM
    scout_log AS sl
        INNER JOIN members AS m
            ON sl.usr_id = m.id
   # scout_logのusr_idとmembersのidをキーに内部結合する
WHER
    sl.in_day BETWEEEN '2015-04-01' AND '2015-05-01'
;
```

- 集計

```mysql
SELECT
    m.city, count(sl.id) # 各都道府県ごとに応募数を表示する
FROM
    scout_log AS sl
        INNER JOIN members AS m
            ON sl.usr_id = m.id
WHERE
    sl.in_day BETWEEEN '2015-04-01' AND '2015-05-01'
GROUP BY
    m.city # m.cityの値でグループ化
;
```

## 実際に叩いてみよう(演習)

#### 現行リジョブDBについて
各テーブルでstatusというカラムが使われている
　status = 1 ・・・使用されているデータ
　status = 0 ・・・削除されたデータ

→不要なレコードをテーブル上から物理的に削除するのではなく、
statusというフラグで論理的に削除して（存在しないように見せかけて）いる。
ので、原則WHERE句にstatus = 1を指定する

※membersなど一部のテーブルではstatusの意味合いが変わる
(members.status = 0 ・・・会員情報の削除ではなく、退会済という意味)

※__リニューアルで大幅に刷新される予定__

#### よく使う関数・式・演算子など

BETWEEEN
LIKE
IN
COUNT
SUM
DISTINCT
NOW
DATE_FORMAT
FIND_IN_SET
INTERVAL
CASE
EXISTS
UNION (ALL)
LEFT(RIGHT) JOIN
INNER JOIN
GROUP BY
HAVING
GROUP_CONCAT

#### 例題

- 都道府県ごとの求職者数
- 都道府県ごとの求職者数を、登録中/退会済みで分けて算出
- 都道府県ごとの平均求職者数 
- 会社が契約中、東京都、の株式会社で、現在表示されている案件名
- 2014年以降の各月応募数、採用数、採用率
- 2015年以降に登録したユーザーのうち、登録から1日以内に応募したユーザー（重複して表示させない）
- 0
- 0



#### 暗黙の型変換

```mysql
# status は整数型
status = 1   #整数で指定
status = '1' #文字列として指定
status = 1.0 #浮動小数点として指定
````

上記は全て正しく評価される。
自動的に必要な型(≠適切な型)へ変換してくれるため

しかし・・・
パフォーマンスが落ちたり、意図しない結果になることがあるので、
カラムのデータ構造を理解し、極力正しい型で指定する

```mysql
'rejob' < 'jigen'
# エラーが出そう・・・でも通ります。つまり
# 元のカラムにどういうデータが入ってるかわかっていないと結果が予測できない
```

しかし効果的に使えば便利！
例えば
`WHERE timestamp BETWEEEN '2015-01-01' AND '2015-02-01'`

```mysql
timestampはDATETIME型なので
'2015-01-01' => '2015-01-01 00:00:00'
'2015-02-01' => '2015-02-01 00:00:00'
へ自動的に変換される
```

つまり下記を表す簡略的な書き方になる
```
'2015-01-01 00:00:00' ≦ timestamp ≦ '2015-01-31 23:59:59'
# 厳密には 「<= '2015-02-01 00:00:00'」なので
# '2015-02-01 00:00:00' ジャストのデータがあれば含まれてしまう
# 
```

下記は全て等価
```mysql
timestamp BETWEEEN '2015-01-01' AND '2015-02-01'

timestamp >= '2015-01-01'
    AND timestamp <= '2015-02-01'

timestamp >= '2015-01-01 00:00:00'
    AND timestamp <= '2015-02-01 00:00:00'
```


#### 実装依存な書き方




## NULLのおはなし

NULLの場合は is NULL


#### よくある落とし穴

各サロンの応募数を抽出したい！

```mysql
SELECT
    sd.cus_id, sd.name, count(sl.id)
FROM
    shop_data AS sd LEFT JOIN scout_log AS sl
        ON sd.cus_id = sl.shop_id
WHERE
    sl.status = 1
    AND sd.status = 1
GROUP BY
    sd.cus_id
;
```

→　__× bad __

#### なぜ？

LEFT JOIN・・・ 右側のテーブルに存在しないデータは__NULLが返る__

|sd.cus_id|sd.name|sl.id|
|------|------|------|
|1|エステサロンA|NULL|
|2|美容室B|121034|
|2|美容室B|124054|
|3|ネイルサロンC|125612|

サロンB,Cは応募データが存在するので応募ID(sl.id)が返される
サロンAは応募が0件=データが存在しないのでNULLが返る

上記をサロンごとにGROUP BYで集約し、COUNTで集計すると

|sd.cus_id|sd.name|__count(sl.id)__|
|------|------|------|
|1|エステサロンA|0|
|2|美容室B|2|
|3|ネイルサロンC|1|

こうなるはず
(COUNTはNULL以外の値の数を数える = Nullだけの場合は0)

だけど実際は

|sd.cus_id|sd.name|__count(sl.id)__|
|------|------|------|
|2|美容室B|2|
|3|ネイルサロンC|1|

どこに原因があるんでしょう？

```mysql
WHERE
    sl.status = 1
 ```

__ココ！！！！__

#### NULLとは

状態が不明、存在しない、定義されていない

箱があった時、下記は明確に区別される
・箱の中身が空であることがわかっている
・箱の中身がわからない = 何か入っているかもしれないし、入っていないかもしれない

下記はいずれも不明(unknown)が返される
- 1 > NULL
- 10 < NULL
- 14 = NULL
- 200 <> NULL
- NULL = NULL
中身のわからない箱Aと、中身のわからない箱Bは、中身が同じであるか？
→わからない

NULLであることを判定するには、SQLではISを使用

```mysql
# boxがNullである時
box = NULL #unknown
box is NULL #True
```


#### 正解例

```msyql
# 不正解
WHERE
    sl.status = 1
# status が1 → 含む
# status が0 → 含まない
# status がNULL → 含まない
```

```mysql
# 正解
WHERE
    (sl.status = 1 OR sl.status IS NULL)
    AND sd.status = 1
GROUP BY
    sd.cus_id
```

## 最後に

興味があれば
インデックス
実装依存な書き方
なんかも

#### SQLを覚えるには？
__とにかく叩け！！！！__








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

代表的な形式

`Logformat "%h %l %u %t ¥"%r¥" %>s %b ¥"%{Referer}i¥" ¥"%{User-Agent}i¥"" combined`

実際に記録されるログ

```
66.249.73.209 - - [24/Feb/2015:00:00:14 +0900] "GET /biyoshi/feature/opening/?page=23&type=area_list&sORt_mode=&dir_name=biyoshi&t_max=10 HTTP/1.1" 301 20 "-" "Mozilla/5.0 (iPhONe; CPU iPhONe OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) VersiON/6.0 Mobile/10A5376e Safari/8536.25 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)" "150900"
```

|記録される内容|説明|
|-----|-----|
|66.249.73.209  |   リモートホスト名|
|-  |   クライアントの識別子|
|-  |   認証ユーザー名|
|[24/Feb/2015:00:00:14 +0900]   |   時刻|
 |"GET /biyoshi/feature/（略） HTTP/1.1"   |   リクエストの最初の行|
|301    |   レスポンスステータス|
|20 |   送信されたバイト数(ヘッダーは含まず)。0バイトの時は「-｣|
 |"-"   |   リファラー|
|Mozilla/5.0 (iPhONe; （略）   |   ユーザーエージェント|
|150900 |   リクエストを処理するのにかかった時間、マイクロ秒単位|

# ログの量

アクセスログだけでも一日あたり150万行

・Excel
不可能！

・データベース
不向き！

→__BigQueryを使おう！！__


#Google BigQuery

## BigQueryとは

Googleが提供するSQLライクな構文で問い合わせができる、ビッグデータ解析サービス。
インフラ構築、保守、などのもろもろの面倒事が不要で、低価格で運用ができる。
__かつ、高速！__

TB級の大規模なDBのフルテーブルスキャンも数秒でこなす
・・・5000台のマシンを使って並列処理している（ものすごい力技）

## SQLとの違い

- データのネスト
- データのupdateができない
- javascriptでUDFが使える

## おいくら？

### 課金形態

- 使用ストレージ容量
$0.02 / 1GB (1月)

- クエリでスキャンするデータ量
$5 / 1TB
(毎月1TB分は無料)

- StreamingInsert
100,000行につき$0.01
(バルクインサートは無料)

### 課金対象

BQはフルスキャン
= WHEREで対象の行を絞っても課金対象は変わらない
→テーブルの行数　×　使用したカラム　= スキャンするデータ量

なので基本的にテーブルを月単位・日付単位などで分割して保存しておく


## 使ってみようBQ

### 基本形

```
# データセット名を付け加える以外は基本的にSQLと同じ
SELECT
    カラム名
FROM
    [データセット名.テーブル名]
WHERE
    条件
;
```

### UNION

テーブル名をカンマ区切りで指定するだけでUNION ALLとして機能する

```
SELECT
    remote_host
FROM
    [Rejob.access_log_20150513],[Rejob.access_log_20150516]
LIMIT 1000;
```

### TABLE_DATE_RANGE

前述の通り、日付単位などで分割されているため、期間が長くなると指定するのが大変

```
# access_log_で始まるテーブル名で日付の期間指定
# テーブル名は'YYYYMMDD'で付けておく
SELECT
    remote_host, request_time
FROM (
    TABLE_DATE_RANGE(
        Rejob.access_log_,
        TIMESTAMP('2015-05-01'),
        TIMESTAMP('2015-05-10')
    )
)
```

### TABLE_QUERY

正規表現などでの指定もできる

```
SELECT
    request_time
FROM
    TABLE_QUERY(
        Rejob,
        'REGEXP_MATCH(table_id, r"2015051[67$]")'
    )
;
```


EACH
JOIN

### 正規表現系

REGEXP_MATCH
REGEXP_EXTRACT
REGEXP_REPLACE

### WITHIN、FLATTEN

ネストされたフィールドに対して使う
*今回用意したデータセットでは使わないので割愛

### その他

その他にもいろいろな機能や関数があるので、詳しくはGoogleのリファレンスを参照
https://cloud.google.com/bigquery/query-reference

#BQを使わない方法

## grep

```
$grep -i 'googlebot' logdata_0501.txt | wc -l
#GoogleBotからのアクセス数をカウント
```

## awk

```
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

その他Hadoopで分散処理とか・・・

# ログの分析例

- Googlebot
- 404を吐いてるURLを調べる
などなど