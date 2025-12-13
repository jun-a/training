package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// NewDBAnswer は新しいデータベース接続を作成する（解答例）
func NewDBAnswer(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続を確認
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// CreateTablesAnswer は必要なテーブルを作成する（解答例）
func (db *DB) CreateTablesAnswer() error {
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

	// jobs テーブルを作成
	if _, err := db.conn.Exec(createJobsTable); err != nil {
		return fmt.Errorf("failed to create jobs table: %w", err)
	}

	// job_skills テーブルを作成
	if _, err := db.conn.Exec(createJobSkillsTable); err != nil {
		return fmt.Errorf("failed to create job_skills table: %w", err)
	}

	return nil
}

// SaveJobAnswer は求人情報をデータベースに保存する（解答例）
func (db *DB) SaveJobAnswer(job Job) error {
	// トランザクションを開始
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // エラー時は自動的にロールバック

	// jobs テーブルに INSERT
	insertJobSQL := `
		INSERT INTO jobs (
			id, title, company_name, company_location, description,
			salary_min, salary_max, salary_currency, employment_type,
			posted_date, application_deadline, is_remote
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.Exec(
		insertJobSQL,
		job.ID,
		job.Title,
		job.Company.Name,
		job.Company.Location,
		job.Description,
		job.Salary.Min,
		job.Salary.Max,
		job.Salary.Currency,
		job.EmploymentType,
		job.PostedDate,
		job.ApplicationDeadline,
		job.IsRemote,
	)
	if err != nil {
		return fmt.Errorf("failed to insert job: %w", err)
	}

	// job_skills テーブルに INSERT
	insertSkillSQL := `INSERT INTO job_skills (job_id, skill) VALUES (?, ?)`
	for _, skill := range job.RequiredSkills.Skills {
		_, err = tx.Exec(insertSkillSQL, job.ID, skill)
		if err != nil {
			return fmt.Errorf("failed to insert skill: %w", err)
		}
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAllJobsAnswer はすべての求人情報を取得する（解答例）
func (db *DB) GetAllJobsAnswer() ([]Job, error) {
	rows, err := db.conn.Query(`
		SELECT id, title, company_name, company_location, description,
		       salary_min, salary_max, salary_currency, employment_type,
		       posted_date, application_deadline, is_remote
		FROM jobs
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Company.Name,
			&job.Company.Location,
			&job.Description,
			&job.Salary.Min,
			&job.Salary.Max,
			&job.Salary.Currency,
			&job.EmploymentType,
			&job.PostedDate,
			&job.ApplicationDeadline,
			&job.IsRemote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		// スキルを取得
		skillRows, err := db.conn.Query(`
			SELECT skill FROM job_skills WHERE job_id = ?
		`, job.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query skills: %w", err)
		}

		var skills []string
		for skillRows.Next() {
			var skill string
			if err := skillRows.Scan(&skill); err != nil {
				skillRows.Close()
				return nil, fmt.Errorf("failed to scan skill: %w", err)
			}
			skills = append(skills, skill)
		}
		skillRows.Close()

		job.RequiredSkills.Skills = skills
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return jobs, nil
}

// GetJobsBySkillAnswer は指定されたスキルを持つ求人を取得する（解答例）
func (db *DB) GetJobsBySkillAnswer(skill string) ([]Job, error) {
	rows, err := db.conn.Query(`
		SELECT DISTINCT j.id, j.title, j.company_name, j.company_location, j.description,
		       j.salary_min, j.salary_max, j.salary_currency, j.employment_type,
		       j.posted_date, j.application_deadline, j.is_remote
		FROM jobs j
		INNER JOIN job_skills js ON j.id = js.job_id
		WHERE js.skill = ?
	`, skill)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs by skill: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Company.Name,
			&job.Company.Location,
			&job.Description,
			&job.Salary.Min,
			&job.Salary.Max,
			&job.Salary.Currency,
			&job.EmploymentType,
			&job.PostedDate,
			&job.ApplicationDeadline,
			&job.IsRemote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		// スキルを取得
		skillRows, err := db.conn.Query(`
			SELECT skill FROM job_skills WHERE job_id = ?
		`, job.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query skills: %w", err)
		}

		var skills []string
		for skillRows.Next() {
			var s string
			if err := skillRows.Scan(&s); err != nil {
				skillRows.Close()
				return nil, fmt.Errorf("failed to scan skill: %w", err)
			}
			skills = append(skills, s)
		}
		skillRows.Close()

		job.RequiredSkills.Skills = skills
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetJobsByRemoteAnswer はリモート可能な求人を取得する（解答例）
func (db *DB) GetJobsByRemoteAnswer() ([]Job, error) {
	rows, err := db.conn.Query(`
		SELECT id, title, company_name, company_location, description,
		       salary_min, salary_max, salary_currency, employment_type,
		       posted_date, application_deadline, is_remote
		FROM jobs
		WHERE is_remote = 1
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query remote jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Company.Name,
			&job.Company.Location,
			&job.Description,
			&job.Salary.Min,
			&job.Salary.Max,
			&job.Salary.Currency,
			&job.EmploymentType,
			&job.PostedDate,
			&job.ApplicationDeadline,
			&job.IsRemote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		// スキルを取得
		skillRows, err := db.conn.Query(`
			SELECT skill FROM job_skills WHERE job_id = ?
		`, job.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query skills: %w", err)
		}

		var skills []string
		for skillRows.Next() {
			var skill string
			if err := skillRows.Scan(&skill); err != nil {
				skillRows.Close()
				return nil, fmt.Errorf("failed to scan skill: %w", err)
			}
			skills = append(skills, skill)
		}
		skillRows.Close()

		job.RequiredSkills.Skills = skills
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// DeleteAllJobsAnswer はすべての求人を削除する（解答例）
func (db *DB) DeleteAllJobsAnswer() error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// job_skills を先に削除（外部キー制約のため）
	if _, err := tx.Exec("DELETE FROM job_skills"); err != nil {
		return fmt.Errorf("failed to delete job_skills: %w", err)
	}

	// jobs を削除
	if _, err := tx.Exec("DELETE FROM jobs"); err != nil {
		return fmt.Errorf("failed to delete jobs: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
