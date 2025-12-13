# Go言語トレーニング課題：求人データ処理システム

## 概要
このトレーニング課題では、XMLファイルから求人データを読み込み、フィルタリング処理を行い、データベースに永続化する一連の処理を実装します。実際の求人検索エンジンの基本的な機能を模擬した実践的な課題です。

## 課題内容

### レベル1: XML読み込みとフィルタリング（基礎）
`main.go` にある以下の関数を実装してください。

#### 1. `LoadJobsFromXML` 関数
XMLファイルから求人データを読み込む関数を実装してください。

**実装のヒント:**
```go
// 1. os.Open() でファイルを開く
// 2. defer file.Close() でファイルを閉じる
// 3. io.ReadAll() でファイルの内容を読み込む
// 4. xml.Unmarshal() でXMLをパースして Jobs 構造体に変換
```

**確認方法:**
```bash
go run main.go
```
「読み込んだ求人数: 8件」と表示されればOKです。

#### 2. `FilterBySkill` 関数
指定されたスキルを持つ求人をフィルタリングする関数を実装してください。

**実装のヒント:**
```go
// 1. jobs をループで回す
// 2. 各求人の RequiredSkills.Skills をチェック
// 3. 指定されたスキルが含まれていれば結果に追加
```

#### 3. `FilterByRemote` 関数
リモート可能な求人をフィルタリングする関数を実装してください。

**実装のヒント:**
```go
// 1. jobs をループで回す
// 2. IsRemote が true の求人を抽出
```

#### 4. `FilterBySalaryRange` 関数
指定された給与範囲内の求人をフィルタリングする関数を実装してください。

**実装のヒント:**
```go
// 1. jobs をループで回す
// 2. 求人の Salary.Min が minSalary 以上かつ
//    Salary.Max が maxSalary 以下のものを抽出
```

### レベル2: データベース永続化（応用）
`database.go` にある以下の関数を実装してください。

#### 準備
SQLite3 ドライバをインストール:
```bash
go get github.com/mattn/go-sqlite3
```

go.mod の初期化（未実施の場合）:
```bash
go mod init job-training
go mod tidy
```

#### 1. `NewDB` 関数
データベース接続を作成する関数を実装してください。

**実装のヒント:**
```go
// 1. sql.Open("sqlite3", dbPath) でデータベースを開く
// 2. db.Ping() で接続を確認
// 3. DB 構造体を返す
```

#### 2. `CreateTables` 関数
必要なテーブルを作成する関数を実装してください。

**実装のヒント:**
```go
// 1. db.conn.Exec() を使って CREATE TABLE 文を実行
// 2. jobs テーブルと job_skills テーブルを作成
```

#### 3. `SaveJob` 関数
求人情報をデータベースに保存する関数を実装してください。

**実装のヒント:**
```go
// 1. db.conn.Begin() でトランザクションを開始
// 2. jobs テーブルに INSERT
// 3. job_skills テーブルにスキル情報を INSERT
// 4. tx.Commit() でコミット
// エラー時は tx.Rollback() でロールバック
```

#### 4. `GetAllJobs` 関数
データベースからすべての求人を取得する関数を実装してください。

**実装のヒント:**
```go
// 1. SELECT 文で jobs テーブルから求人を取得
// 2. 各求人について job_skills テーブルからスキルを取得
// 3. Job 構造体に組み立てる
```

#### 5. `GetJobsBySkill` 関数
指定されたスキルを持つ求人を取得する関数を実装してください。

**実装のヒント:**
```go
// 1. jobs と job_skills テーブルを JOIN
// 2. WHERE 句でスキルを絞り込む
```

#### 6. `GetJobsByRemote` 関数
リモート可能な求人を取得する関数を実装してください。

**実装のヒント:**
```go
// 1. WHERE is_remote = 1 で絞り込む
```

#### 7. `DeleteAllJobs` 関数
すべての求人を削除する関数を実装してください。

**実装のヒント:**
```go
// 1. トランザクションを開始
// 2. job_skills を先に DELETE（外部キー制約のため）
// 3. jobs を DELETE
// 4. コミット
```

## データベース処理のテスト例

以下のようなテストプログラムを作成して動作確認できます:

```go
package main

import (
	"fmt"
	"log"
)

func testDatabase() {
	// XMLから読み込み
	jobs, err := LoadJobsFromXML("jobs.xml")
	if err != nil {
		log.Fatal(err)
	}

	// データベース接続
	db, err := NewDB("jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブル作成
	if err := db.CreateTables(); err != nil {
		log.Fatal(err)
	}

	// データ保存
	if err := db.SaveJobs(jobs.JobList); err != nil {
		log.Fatal(err)
	}

	// データ取得
	savedJobs, err := db.GetAllJobs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("保存された求人数: %d件\n", len(savedJobs))

	// スキルで絞り込み
	goJobs, err := db.GetJobsBySkill("Go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Goスキルを持つ求人数: %d件\n", len(goJobs))
}
```

## ファイル構成

```
go-training/
├── README.md              # この課題説明ファイル
├── jobs.xml               # 求人データ（8件の求人情報）
├── main.go                # レベル1の課題（XMLの読み込みとフィルタリング）
├── database.go            # レベル2の課題（データベース永続化）
├── main_answer.go         # レベル1の解答例
└── database_answer.go     # レベル2の解答例
```

## 学習のポイント

### レベル1で学べること
- Go言語での構造体の定義と使用
- XMLパッケージを使ったデータの読み込み
- スライスの操作とループ処理
- エラーハンドリングの基礎

### レベル2で学べること
- database/sqlパッケージの使用方法
- SQLite3との連携
- トランザクション処理
- SQLクエリの実行とデータの取得
- データベース設計の基礎

## 求人データについて

`jobs.xml` には以下の8件の求人データが含まれています:

1. **JOB001**: バックエンドエンジニア（Go言語） - 東京都渋谷区
2. **JOB002**: フロントエンドエンジニア - 大阪府大阪市北区
3. **JOB003**: データエンジニア - 東京都港区
4. **JOB004**: SREエンジニア - 福岡県福岡市中央区
5. **JOB005**: 機械学習エンジニア - 東京都千代田区
6. **JOB006**: プロダクトマネージャー - 東京都新宿区
7. **JOB007**: フルスタックエンジニア（契約社員） - 神奈川県横浜市
8. **JOB008**: QAエンジニア - 愛知県名古屋市

各求人には以下の情報が含まれています:
- 求人ID
- タイトル
- 企業名・所在地
- 説明文
- 給与範囲
- 雇用形態
- 必要スキル
- 掲載日・応募締切
- リモート可否

## 解答例について

実装に困ったときは `main_answer.go` と `database_answer.go` を参照してください。ただし、まずは自分で考えて実装することをお勧めします。

## 実行方法

### レベル1の実行
```bash
cd go-training
go run main.go
```

### レベル2の実行（データベース処理を含む場合）
```bash
# 必要なパッケージをインストール
go mod init job-training
go get github.com/mattn/go-sqlite3
go mod tidy

# 実行
go run main.go database.go
```

## 発展課題

基本課題が完了したら、以下のような機能追加にもチャレンジしてみてください:

1. **複数スキルでのフィルタリング**: AND条件、OR条件の両方に対応
2. **ソート機能**: 給与順、掲載日順などでソート
3. **全文検索**: 説明文からキーワード検索
4. **集計機能**: 都道府県別の求人数、スキル別の平均給与など
5. **JSON出力**: フィルタリング結果をJSONで出力
6. **Webサーバー化**: net/httpパッケージを使ってREST APIを作成

頑張ってください！
