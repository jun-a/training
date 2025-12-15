package main

import (
	"fmt"
)

// database.go - データベース永続化の実装
//
// このファイルでは、以下のベストプラクティスを学びます：
//
// 1. トランザクションの使用
//    - SaveJob内で単一求人をトランザクションで保存
//    - SaveJobs内で複数求人を1つのトランザクションで保存するとさらに効率的
//    - トランザクションはデータの一貫性を保証し、エラー時のロールバックを可能にする
//
// 2. N+1クエリ問題
//    - GetAllJobsでは各求人ごとにスキルを取得すると、N+1クエリ問題が発生
//    - 悪い例: SELECT * FROM jobs; の後、各求人に対して SELECT * FROM job_skills WHERE job_id = ?
//    - 良い方法1: JOINを使って1つのクエリで取得
//    - 良い方法2: すべてのスキルを一度に取得してメモリ上でマッピング
//
// 3. プリペアドステートメント
//    - ループ内でクエリを実行する場合は、db.Prepare()を使用
//    - SQLインジェクション対策とパフォーマンス向上に役立つ
//    - 例: stmt, _ := db.Prepare("INSERT INTO job_skills (job_id, skill) VALUES (?, ?)")
//          defer stmt.Close()
//          for _, skill := range skills {
//              stmt.Exec(jobID, skill)
//          }
//
// 4. エラーハンドリング
//    - すべてのデータベース操作でエラーをチェック
//    - トランザクション使用時は、エラー発生時に必ずRollback()を呼び出す
//
// NOTE: Ruby実装（ruby-reference/database.rb）と比較して、Go実装の違いを学びましょう

// DB はデータベース接続を管理する構造体
// NOTE: この構造体は実装時に database/sql パッケージの sql.DB を使用します
type DB struct {
	conn interface{} // TODO: 実装時に *sql.DB に変更してください
}

// NewDB は新しいデータベース接続を作成する
func NewDB(dbPath string) (*DB, error) {
	// TODO: ここにデータベース接続を作成するコードを実装してください
	// ヒント:
	// 1. sql.Open("sqlite3", dbPath)でデータベースを開く
	// 2. db.Ping()で接続を確認
	// 
	// 必要なインポート:
	// import (
	//     "database/sql"
	//     _ "github.com/mattn/go-sqlite3"
	// )
	return nil, fmt.Errorf("not implemented")
}

// Close はデータベース接続を閉じる
func (db *DB) Close() error {
	// NOTE: この関数は実装済みです
	// TODO: 実装時には conn を *sql.DB にキャストして Close() を呼び出してください
	return nil
}

// CreateTables は必要なテーブルを作成する
func (db *DB) CreateTables() error {
	// TODO: ここにテーブル作成のコードを実装してください
	// 必要なテーブル:
	// 1. jobs - 求人情報
	// 2. job_skills - 求人と必要スキルの中間テーブル
	//
	// NOTE: スキーマには status カラムを含めること（Ruby実装を参照）
	// 
	// 以下のテーブル定義を使用してください:
	// 
	// 必要なインポート:
	// import "database/sql"

	// createJobsTable := `
	// CREATE TABLE IF NOT EXISTS jobs (
	// 	id TEXT PRIMARY KEY,
	// 	title TEXT NOT NULL,
	// 	company_name TEXT NOT NULL,
	// 	company_location TEXT NOT NULL,
	// 	description TEXT,
	// 	salary_min INTEGER,
	// 	salary_max INTEGER,
	// 	salary_currency TEXT,
	// 	salary_is_public BOOLEAN,
	// 	employment_type TEXT,
	// 	posted_date TEXT,
	// 	application_deadline TEXT,
	// 	is_remote BOOLEAN,
	// 	status TEXT
	// );`

	// createJobSkillsTable := `
	// CREATE TABLE IF NOT EXISTS job_skills (
	// 	job_id TEXT NOT NULL,
	// 	skill TEXT NOT NULL,
	// 	PRIMARY KEY (job_id, skill),
	// 	FOREIGN KEY (job_id) REFERENCES jobs(id)
	// );`

	// TODO: db.conn.Exec()を使ってテーブルを作成してください

	return fmt.Errorf("not implemented")
}

// SaveJob は求人情報をデータベースに保存する
func (db *DB) SaveJob(job Job) error {
	// TODO: ここに求人情報を保存するコードを実装してください
	// ヒント:
	// 1. トランザクションを開始 (db.conn.Begin())
	// 2. jobs テーブルに求人情報を INSERT OR REPLACE
	// 3. 既存のスキルを削除 (DELETE FROM job_skills WHERE job_id = ?)
	// 4. job_skills テーブルにスキル情報を INSERT
	// 5. トランザクションをコミット (tx.Commit())
	// エラー時はロールバック (tx.Rollback())
	//
	// NOTE: トランザクションを使用することで、データの一貫性を保証できます
	// NOTE: status フィールドも保存することを忘れずに

	return fmt.Errorf("not implemented")
}

// SaveJobs は複数の求人情報をデータベースに保存する
func (db *DB) SaveJobs(jobs []Job) error {
	// TODO: ここに複数の求人情報を保存するコードを実装してください
	// 
	// 現在の実装: SaveJobを繰り返し呼び出す（各求人ごとにトランザクション）
	// より効率的な実装: 全求人を1つのトランザクションで保存する
	//
	// HINT: SaveJob内でトランザクションを使用しているため、
	// この関数で全体をトランザクションでラップするとさらに効率的になります
	// （Ruby実装を参照）
	
	for _, job := range jobs {
		if err := db.SaveJob(job); err != nil {
			return fmt.Errorf("failed to save job %s: %w", job.ID, err)
		}
	}
	return nil
}

// GetAllJobs はすべての求人情報を取得する
func (db *DB) GetAllJobs() ([]Job, error) {
	// TODO: ここにすべての求人情報を取得するコードを実装してください
	// ヒント:
	// 1. jobs テーブルから求人情報を SELECT
	// 2. 各求人に対してjob_skillsテーブルからスキルを取得
	// 3. Job構造体を組み立てる
	//
	// NOTE: N+1クエリ問題に注意
	// 現在の実装案: 各求人ごとにスキルをクエリ（N+1問題が発生）
	// より良い方法: JOINを使用して1つのクエリで取得する、または
	// すべてのスキルを一度に取得してメモリ上でマッピングする

	return nil, fmt.Errorf("not implemented")
}

// GetJobsBySkill は指定されたスキルを持つ求人を取得する
func (db *DB) GetJobsBySkill(skill string) ([]Job, error) {
	// TODO: ここに指定されたスキルを持つ求人を取得するコードを実装してください
	// ヒント: job_skills テーブルと jobs テーブルを JOIN
	//
	// NOTE: ループ内でクエリを実行する場合は、プリペアドステートメントを使用すること
	// プリペアドステートメントは、SQLインジェクション対策とパフォーマンス向上に役立ちます

	return nil, fmt.Errorf("not implemented")
}

// GetJobsByRemote はリモート可能な求人を取得する
func (db *DB) GetJobsByRemote() ([]Job, error) {
	// TODO: ここにリモート可能な求人を取得するコードを実装してください
	// ヒント: WHERE is_remote = true

	return nil, fmt.Errorf("not implemented")
}

// DeleteAllJobs はすべての求人を削除する（テスト用）
func (db *DB) DeleteAllJobs() error {
	// TODO: ここにすべての求人を削除するコードを実装してください
	// ヒント: DELETE文を実行

	return fmt.Errorf("not implemented")
}
