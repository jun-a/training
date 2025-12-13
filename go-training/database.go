package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DB はデータベース接続を管理する構造体
type DB struct {
	conn *sql.DB
}

// NewDB は新しいデータベース接続を作成する
func NewDB(dbPath string) (*DB, error) {
	// TODO: ここにデータベース接続を作成するコードを実装してください
	// ヒント:
	// 1. sql.Open("sqlite3", dbPath)でデータベースを開く
	// 2. db.Ping()で接続を確認
	return nil, fmt.Errorf("not implemented")
}

// Close はデータベース接続を閉じる
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// CreateTables は必要なテーブルを作成する
func (db *DB) CreateTables() error {
	// TODO: ここにテーブル作成のコードを実装してください
	// 必要なテーブル:
	// 1. jobs - 求人情報
	// 2. companies - 企業情報
	// 3. job_skills - 求人と必要スキルの中間テーブル

	createJobsTable := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		company_name TEXT NOT NULL,
		company_location TEXT NOT NULL,
		description TEXT,
		salary_min INTEGER,
		salary_max INTEGER,
		salary_currency TEXT,
		employment_type TEXT,
		posted_date TEXT,
		application_deadline TEXT,
		is_remote BOOLEAN
	);`

	createJobSkillsTable := `
	CREATE TABLE IF NOT EXISTS job_skills (
		job_id TEXT NOT NULL,
		skill TEXT NOT NULL,
		PRIMARY KEY (job_id, skill),
		FOREIGN KEY (job_id) REFERENCES jobs(id)
	);`

	// TODO: db.conn.Exec()を使ってテーブルを作成してください

	return fmt.Errorf("not implemented")
}

// SaveJob は求人情報をデータベースに保存する
func (db *DB) SaveJob(job Job) error {
	// TODO: ここに求人情報を保存するコードを実装してください
	// ヒント:
	// 1. トランザクションを開始 (db.conn.Begin())
	// 2. jobs テーブルに求人情報を INSERT
	// 3. job_skills テーブルにスキル情報を INSERT
	// 4. トランザクションをコミット (tx.Commit())
	// エラー時はロールバック (tx.Rollback())

	return fmt.Errorf("not implemented")
}

// SaveJobs は複数の求人情報をデータベースに保存する
func (db *DB) SaveJobs(jobs []Job) error {
	// TODO: ここに複数の求人情報を保存するコードを実装してください
	// ヒント: SaveJobを繰り返し呼び出す
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

	return nil, fmt.Errorf("not implemented")
}

// GetJobsBySkill は指定されたスキルを持つ求人を取得する
func (db *DB) GetJobsBySkill(skill string) ([]Job, error) {
	// TODO: ここに指定されたスキルを持つ求人を取得するコードを実装してください
	// ヒント: job_skills テーブルと jobs テーブルを JOIN

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
