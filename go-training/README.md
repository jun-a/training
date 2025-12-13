# Go言語トレーニング課題：求人データ処理システム（実践編）

## 概要
このトレーニング課題では、XMLファイルから求人データを読み込み、データの正規化・トリム処理を行い、新旧データの差分比較を実施し、最終的にデータベースに永続化する一連の処理を実装します。実際の求人検索エンジンで必要となる実践的な処理を学べる課題です。

## 課題で学べること

実際の業務で必要となる以下のスキルを習得できます:

✅ **ファイルの読み込み** - XMLファイルの読み込みとパース
✅ **データの取得や構造化** - XMLデータを構造体に変換
✅ **データハンドリング** - スライス操作、フィルタリング、ソート
✅ **文字列のトリム処理** - 不要な空白や改行の削除
✅ **データの正規化** - 様々な形式の給与データを統一フォーマットに変換
✅ **特定の条件における判定** - 複雑な条件でのデータフィルタリング
✅ **新旧の差分比較** - 2つのデータセット間の変更検出
✅ **データベースへの書き込み** - SQLite3を使った永続化

## ファイル構成

```
go-training/
├── README.md                  # この課題説明ファイル
├── jobs.xml                   # 旧形式の求人データ（8件）
├── jobs_new.xml               # 新形式の求人データ（10件、正規化が必要）
├── jobs_old.xml               # 差分比較用の旧データ（6件）
├── processing.go              # レベル2の課題（データ正規化・トリム処理）
├── diff.go                    # レベル3の課題（差分比較）
├── database.go                # レベル4の課題（データベース永続化）
├── main.go                    # レベル1の課題（基本的なXML読み込み）
└── ruby-reference/            # Ruby言語による参考実装
    ├── README.md              # Ruby実装の説明
    ├── main.rb                # レベル1の参考実装
    ├── processing.rb          # レベル2の参考実装
    ├── diff.rb                # レベル3の参考実装
    └── database.rb            # レベル4の参考実装
```

## 課題内容

### レベル1: 基本的なXML読み込み（基礎）

`main.go` を使用して、基本的なXMLファイルの読み込みを実装します。

**課題:**
- `jobs.xml`（整形済みのデータ）を読み込む
- 基本的なフィルタリング機能を実装

**確認方法:**
```bash
go run main.go
```

---

### レベル2: データ正規化・トリム処理（実践）⭐

`processing.go` にある関数を実装してください。**これが最も重要な課題です！**

#### 1. `LoadJobsFromXML` 関数
XMLファイルから求人データを読み込みます。

**実装のヒント:**
```go
// 1. os.Open() でファイルを開く
// 2. defer file.Close() でファイルを閉じる
// 3. io.ReadAll() でファイルの内容を読み込む
// 4. xml.Unmarshal() でXMLをパース
```

#### 2. `TrimJobData` 関数
求人データから不要な空白や改行を削除します。

**対応すべきデータの問題:**
- タイトルの前後に空白: `"  バックエンドエンジニア（Go言語）  "`
- 企業名に空白: `"  デザインラボ株式会社  "`
- 説明文に余分な改行や空白
- スキル名に空白: `"  Go  "`, `"  PostgreSQL  "`
- 雇用形態に空白: `"正社員  "`
- ステータスに空白: `"  closed  "`

**実装のヒント:**
```go
// - strings.TrimSpace() で前後の空白を削除
// - regexp を使って複数の空白・改行を単一のスペースに置換
// - 各フィールドに対して適切な処理を適用
```

#### 3. `NormalizeSalary` 関数 ⭐⭐⭐
**最重要課題！** 様々な形式の給与データを統一フォーマットに正規化します。

**対応すべき形式（`jobs_new.xml`に含まれています）:**

| 入力形式 | 説明 | 期待される出力 |
|---------|------|--------------|
| `"年収 6,000,000円 〜 10,000,000円"` | カンマ区切り、範囲あり | Min: 6000000, Max: 10000000 |
| `"500万円〜800万円"` | 万円表記、範囲あり | Min: 5000000, Max: 8000000 |
| `"7000000"` | 数値のみ | Min: 7000000, Max: 7000000 |
| `"月給 550,000円 〜 900,000円"` | 月給（年収に変換） | Min: 6600000, Max: 10800000 |
| `"800〜1500万円"` | 万円表記、短縮形 | Min: 8000000, Max: 15000000 |
| `"応相談（想定年収700〜1300万円）"` | 括弧内に範囲 | Min: 7000000, Max: 13000000 |
| `"¥4,500,000 - ¥7,000,000"` | 円マーク、ハイフン区切り | Min: 4500000, Max: 7000000 |
| `"年俸制 8,000,000円〜12,000,000円"` | 年俸制 | Min: 8000000, Max: 12000000 |
| `""` (空文字列) | 給与非公開 | IsPublic: false |
| `"非公開"` | 給与非公開 | IsPublic: false |

**実装のヒント:**
```go
// 1. strings.TrimSpace() で前後の空白を削除
// 2. strings.Contains() で「月給」「万円」「非公開」などをチェック
// 3. カンマを削除: strings.ReplaceAll(s, ",", "")
// 4. regexp で数値を抽出: regexp.MustCompile(`(\d+)`)
// 5. 万円の場合は × 10000
// 6. 月給の場合は × 12
// 7. 数値が1つの場合は Min = Max
// 8. 数値が2つの場合は Min, Max に分割
```

#### 4. `NormalizeRemote` 関数
様々な形式のリモートフラグを正規化します。

**対応すべき形式:**
- `"true"`, `"TRUE"`, `"True"` → `true`
- `"false"`, `"FALSE"`, `"False"` → `false`
- `"yes"`, `"YES"`, `"Yes"` → `true`
- `"no"`, `"NO"`, `"No"` → `false`
- `"1"` → `true`
- `"0"` → `false`
- `"remote"`, `"Remote"`, `"REMOTE"` → `true`

**実装のヒント:**
```go
// 1. strings.TrimSpace() で前後の空白を削除
// 2. strings.ToLower() で小文字に変換
// 3. switch文で分岐処理
```

#### 5. `NormalizeJob` 関数
上記の正規化関数を組み合わせて、生データを正規化します。

#### 6. フィルタリング関数群
- `FilterBySkill`: スキルでフィルタリング（大文字小文字を無視）
- `FilterByRemote`: リモート可能な求人のみ抽出
- `FilterBySalaryRange`: 給与範囲でフィルタリング（非公開は除外）
- `FilterByStatus`: ステータスでフィルタリング

**テスト方法:**
```bash
go run processing.go
```

---

### レベル3: 新旧データの差分比較（応用）⭐⭐

`diff.go` にある関数を実装してください。

#### 差分比較の仕様

`jobs_old.xml` と `jobs_new.xml` を比較して、以下の変更を検出します:

**変更の種類:**
- `added`: 新規追加された求人
- `updated`: 更新された求人
- `deleted`: 削除された求人
- `unchanged`: 変更なしの求人

**検出すべき変更内容:**
- タイトルの変更
- 給与の変更（上がった/下がった）
- スキルの追加・削除
- 締切日の変更
- ステータスの変更（active → closed など）
- 説明文の変更

#### 実装する関数

##### 1. `CompareJobs` 関数
新旧の求人データを比較して差分を抽出します。

**実装のヒント:**
```go
// 1. 新旧の求人をIDでマップ化: map[string]Job
// 2. 新しい求人をループ:
//    - 旧データに存在しない → "added"
//    - 旧データに存在 → 各フィールドを比較
// 3. 旧データにあるが新データにない → "deleted"
```

##### 2. `FindSkillsDiff` 関数
スキルの差分を検出します。

**実装のヒント:**
```go
// 1. oldSkillsとnewSkillsをマップ化: map[string]bool
// 2. newSkillsにあってoldSkillsにない → added
// 3. oldSkillsにあってnewSkillsにない → removed
```

##### 3. フィルタリング関数群
- `FilterByChangeType`: 変更タイプで絞り込み
- `FilterBySalaryIncreased`: 給与が上がった求人のみ
- `FilterBySkillsAdded`: スキルが追加された求人のみ
- `FilterByStatusClosed`: クローズに変更された求人のみ

**実際のデータでの差分例:**

| 求人ID | 変更内容 |
|-------|---------|
| JOB001 | 給与アップ（550〜900万 → 600〜1000万）、PostgreSQLスキル追加、締切延長 |
| JOB002 | 給与アップ（480〜750万 → 500〜800万）、CSSスキル追加 |
| JOB003 | 給与アップ（650万 → 700万）、GCPスキル追加 |
| JOB004 | 新規追加（SREエンジニア） |
| JOB005 | 新規追加（機械学習エンジニア） |
| JOB006 | 給与アップ、プロダクト企画スキル追加、説明文更新 |
| JOB007 | 新規追加（フルスタックエンジニア） |
| JOB008 | 給与アップ、CI/CDスキル追加 |
| JOB009 | 新規追加（DevOpsエンジニア） |
| JOB010 | ステータスがclosedに変更、締切延長 |

**テスト方法:**
```bash
go run processing.go diff.go processing_answer.go
```

---

### レベル4: データベース永続化（発展）

`database.go` にある関数を実装してください。

#### 準備
```bash
go get github.com/mattn/go-sqlite3
go mod init job-training
go mod tidy
```

#### 実装する関数

1. `NewDB`: データベース接続の作成
2. `CreateTables`: テーブルの作成（jobs, job_skills）
3. `SaveJob`: 求人データの保存（トランザクション使用）
4. `GetAllJobs`: すべての求人の取得
5. `GetJobsBySkill`: スキルで絞り込み
6. `GetJobsByRemote`: リモート可能な求人の取得
7. `DeleteAllJobs`: すべての求人の削除

詳細は `database.go` のコメントを参照してください。

---

## 実行方法

### レベル1: 基本的なXML読み込み
```bash
go run main.go
```

### レベル2: データ正規化・トリム処理のテスト
```bash
# テストプログラムを作成して実行
go run processing.go processing_answer.go -test
```

以下のようなテストコードを作成できます:

```go
package main

import "fmt"

func main() {
	// XMLから読み込み
	jobsRaw, err := LoadJobsFromXML("jobs_new.xml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("読み込んだ求人数: %d件\n", len(jobsRaw.JobList))

	// 正規化
	jobs := NormalizeJobs(jobsRaw)

	// 結果を表示
	for _, job := range jobs {
		PrintJob(job)
	}

	// フィルタリングのテスト
	goJobs := FilterBySkill(jobs, "Go")
	fmt.Printf("\nGoスキルを持つ求人: %d件\n", len(goJobs))

	remoteJobs := FilterByRemote(jobs)
	fmt.Printf("リモート可能な求人: %d件\n", len(remoteJobs))

	activeJobs := FilterByStatus(jobs, "active")
	fmt.Printf("アクティブな求人: %d件\n", len(activeJobs))
}
```

### レベル3: 差分比較のテスト
```go
package main

import "fmt"

func main() {
	// 旧データを読み込み
	oldJobsRaw, _ := LoadJobsFromXML("jobs_old.xml")
	oldJobs := NormalizeJobs(oldJobsRaw)

	// 新データを読み込み
	newJobsRaw, _ := LoadJobsFromXML("jobs_new.xml")
	newJobs := NormalizeJobs(newJobsRaw)

	// 差分を比較
	diffs := CompareJobs(oldJobs, newJobs)

	// サマリー表示
	PrintDiffSummary(diffs)

	// 詳細表示
	fmt.Println("\n=== すべての差分 ===")
	for _, diff := range diffs {
		PrintJobDiff(diff)
	}

	// 給与が上がった求人
	fmt.Println("\n=== 給与が上がった求人 ===")
	salaryIncreased := FilterBySalaryIncreased(diffs)
	for _, diff := range salaryIncreased {
		PrintJobDiff(diff)
	}

	// スキルが追加された求人
	fmt.Println("\n=== スキルが追加された求人 ===")
	skillsAdded := FilterBySkillsAdded(diffs)
	for _, diff := range skillsAdded {
		PrintJobDiff(diff)
	}
}
```

---

## 学習のポイント

### レベル2で学べること（最重要）
- **文字列処理の実践**: `strings` パッケージの活用
- **正規表現**: `regexp` パッケージを使ったパターンマッチング
- **データ正規化**: 様々な形式のデータを統一フォーマットに変換
- **エラーハンドリング**: 不正なデータへの対応
- **型変換**: `strconv` パッケージの使用

### レベル3で学べること
- **データ構造の選択**: マップを使った効率的な検索
- **アルゴリズム**: 差分検出のロジック
- **複雑な条件判定**: 多様な変更パターンの検出
- **データ分析**: 変更内容の統計処理

### レベル4で学べること
- **データベース操作**: SQL の基礎
- **トランザクション**: データの一貫性保証
- **JOIN クエリ**: 複数テーブルの結合
- **インデックス**: パフォーマンス最適化

---

## よくある実装のポイント

### 給与の正規化で詰まったら

1. **まずは簡単なケースから実装する**
   ```go
   // ステップ1: 空文字列と「非公開」のチェック
   // ステップ2: カンマの削除
   // ステップ3: 数値の抽出
   // ステップ4: 「万円」のチェックと変換
   // ステップ5: 「月給」のチェックと変換
   ```

2. **正規表現のテスト**
   ```go
   numberRegex := regexp.MustCompile(`(\d+)`)
   matches := numberRegex.FindAllString("年収 600万円〜1000万円", -1)
   fmt.Println(matches) // ["600", "1000"]
   ```

3. **デバッグ出力を活用**
   ```go
   fmt.Printf("入力: %s\n", salaryStr)
   fmt.Printf("抽出された数値: %v\n", matches)
   fmt.Printf("結果: Min=%d, Max=%d\n", min, max)
   ```

### 差分比較で詰まったら

1. **マップの活用**
   ```go
   // IDをキーにしたマップを作成
   oldJobMap := make(map[string]Job)
   for _, job := range oldJobs {
       oldJobMap[job.ID] = job
   }

   // 存在チェックが高速に
   if oldJob, exists := oldJobMap[newJob.ID]; exists {
       // 比較処理
   }
   ```

2. **小さいデータでテスト**
   - まずは2件程度のデータでテストする
   - 各ケース（added, updated, deleted）を個別に確認

---

## 発展課題

基本課題が完了したら、以下にもチャレンジしてみてください:

1. **バリデーション機能**: 不正なデータの検出とエラーレポート
2. **CSV出力**: フィルタリング結果をCSVファイルに出力
3. **JSON API**: net/httpパッケージでREST APIを作成
4. **並行処理**: goroutineを使った高速化
5. **テストコード**: testing パッケージでユニットテスト作成
6. **CLI ツール化**: flag パッケージでコマンドラインツールに
7. **差分レポート**: 変更内容をMarkdown形式でレポート出力

---

## 参考実装について

Go言語での実装に困ったときは、Ruby言語による参考実装を参照してください：

### Ruby言語による参考実装

`ruby-reference/` ディレクトリに、同じ処理を実装したRubyコードがあります：
- `ruby-reference/main.rb` - レベル1の参考実装
- `ruby-reference/processing.rb` - レベル2の参考実装
- `ruby-reference/diff.rb` - レベル3の参考実装
- `ruby-reference/database.rb` - レベル4の参考実装

Ruby実装は完全に動作するコードとして提供されています。Go言語とRubyの違いを学びながら、同じ処理をGo言語で実装してみてください。

詳細は `ruby-reference/README.md` を参照してください。

**ヒント**: まずは自分で考えて実装することをお勧めします。Ruby実装を見る前に、エラーと向き合うことが最大の学びになります！

---

## トラブルシューティング

### XMLのパースに失敗する
- ファイルパスが正しいか確認
- ファイルのエンコーディングがUTF-8か確認
- XMLの構造が正しいか確認

### 給与の正規化がうまくいかない
- デバッグ出力で中間結果を確認
- 正規表現のパターンをテストツールで確認
- 一つずつケースを追加して実装

### 差分比較の結果が合わない
- マップのキー（ID）が正しいか確認
- 比較ロジックを一つずつ確認
- 小さいデータセットでテスト

頑張ってください！実践的なスキルが身につきます！
