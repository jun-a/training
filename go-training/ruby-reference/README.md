# Ruby言語による参考実装

このディレクトリには、Go言語トレーニング課題と同等の処理をRubyで実装した参考コードが含まれています。

## 環境構築

### 方法1: Gemfileを使用（推奨）

```bash
cd go-training/ruby-reference

# Bundlerをインストール（未インストールの場合）
gem install bundler

# 依存関係をインストール
bundle install
```

### 方法2: 直接gemをインストール

```bash
# Rubyのインストール (バージョン 3.0以上推奨)
# macOS (Homebrew)
brew install ruby

# Ubuntu/Debian
sudo apt-get install ruby-full

# 必要なgemのインストール
gem install sqlite3 nokogiri
```

## ファイル構成

```
ruby-reference/
├── README.md                 # このファイル
├── main.rb                   # レベル1: 基本的なXML読み込み
├── processing.rb             # レベル2: データ正規化・トリム処理
├── diff.rb                   # レベル3: 差分比較
└── database.rb               # レベル4: データベース永続化
```

## 実行方法

### レベル1: 基本的なXML読み込み

```bash
cd go-training/ruby-reference
ruby main.rb
```

### レベル2: データ正規化・トリム処理

```bash
ruby processing.rb
```

### レベル3: 差分比較

```bash
ruby diff.rb
```

### レベル4: データベース永続化

```bash
ruby database.rb
```

## Go言語との対応関係

### データ構造

**Go (構造体):**
```go
type Job struct {
    ID          string
    Title       string
    Company     Company
    // ...
}
```

**Ruby (クラス):**
```ruby
class Job
  attr_accessor :id, :title, :company, # ...
  
  def initialize(attributes = {})
    @id = attributes[:id]
    @title = attributes[:title]
    # ...
  end
end
```

### XMLパース

**Go:**
```go
xml.Unmarshal(data, &jobs)
```

**Ruby:**
```ruby
doc = Nokogiri::XML(file)
doc.xpath('//job').each do |job_node|
  # パース処理
end
```

### 配列操作

**Go (フィルタリング):**
```go
filtered := []Job{}
for _, job := range jobs {
    if condition {
        filtered = append(filtered, job)
    }
}
```

**Ruby (select/reject/map):**
```ruby
filtered = jobs.select { |job| condition }
```

### データベース操作

**Go:**
```go
db, err := sql.Open("sqlite3", dbPath)
rows, err := db.Query("SELECT * FROM jobs")
```

**Ruby:**
```ruby
db = SQLite3::Database.new(db_path)
db.execute("SELECT * FROM jobs") do |row|
  # 処理
end
```

## Rubyの特徴的な機能

### 1. ブロックとイテレータ

```ruby
# 配列の各要素に対して処理
jobs.each do |job|
  puts job.title
end

# フィルタリング
active_jobs = jobs.select { |job| job.status == 'active' }

# 変換
titles = jobs.map { |job| job.title }
```

### 2. 正規表現の組み込みサポート

```ruby
# マッチング
if salary_str =~ /(\d+)万円/
  amount = $1.to_i * 10000
end

# 置換
normalized = salary_str.gsub(/[,円]/, '')
```

### 3. シンボルとハッシュ

```ruby
# ハッシュの作成
job_data = {
  id: 'JOB001',
  title: 'エンジニア',
  salary: { min: 600_0000, max: 1000_0000 }
}

# アクセス
job_data[:id]  # => 'JOB001'
```

### 4. 文字列操作

```ruby
# トリム
"  text  ".strip  # => "text"

# 大文字小文字変換
"Hello".downcase  # => "hello"
"hello".upcase    # => "HELLO"

# 含まれているか
"hello world".include?("world")  # => true
```

### 5. nilとfalseの扱い

```ruby
# nilやfalseは偽、それ以外は真
if value
  # valueがnilやfalse以外の場合
end

# nil安全な呼び出し
job&.company&.name  # jobやcompanyがnilの場合はnil
```

## 学習のポイント

1. **Rubyの文法はより柔軟**: Goよりも簡潔に書けることが多い
2. **動的型付け**: 型宣言が不要だが、その分注意が必要
3. **豊富な組み込みメソッド**: Array, String, Hashなどに便利なメソッドが多数
4. **ブロックの活用**: イテレーションやコールバック処理が直感的
5. **メタプログラミング**: attr_accessorなどで動的にメソッドを生成

## よくある質問

### Q: Goの構造体タグ（`xml:"id"`）に相当するものは？

A: Nokogiriを使う場合はXPathやCSSセレクタで直接要素を取得します。

### Q: Goのエラーハンドリング（`if err != nil`）は？

A: Rubyでは例外（rescue）を使います：

```ruby
begin
  file = File.open(filename)
rescue => e
  puts "エラー: #{e.message}"
end
```

### Q: Goのgoroutineに相当するものは？

A: Rubyでは`Thread`や`Fiber`を使いますが、基本的な課題では不要です。

## 参考資料

- [Ruby公式ドキュメント](https://docs.ruby-lang.org/ja/)
- [Nokogiri (XML解析)](https://nokogiri.org/)
- [sqlite3 gem](https://github.com/sparklemotion/sqlite3-ruby)
