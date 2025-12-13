#!/usr/bin/env ruby
# encoding: utf-8

# レベル4: データベース永続化
# Go言語のdatabase.goに対応するRuby実装

require 'sqlite3'
require_relative 'processing'

# データベースクラス
class JobDatabase
  attr_reader :db

  def initialize(db_path)
    @db = SQLite3::Database.new(db_path)
    @db.results_as_hash = true
  end

  # データベースを閉じる
  def close
    @db.close if @db
  end

  # テーブルを作成
  def create_tables
    # jobsテーブル
    @db.execute <<-SQL
      CREATE TABLE IF NOT EXISTS jobs (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        company_name TEXT NOT NULL,
        company_location TEXT NOT NULL,
        description TEXT,
        salary_min INTEGER,
        salary_max INTEGER,
        salary_currency TEXT,
        salary_is_public BOOLEAN,
        employment_type TEXT,
        posted_date TEXT,
        application_deadline TEXT,
        is_remote BOOLEAN,
        status TEXT
      );
    SQL

    # job_skillsテーブル
    @db.execute <<-SQL
      CREATE TABLE IF NOT EXISTS job_skills (
        job_id TEXT NOT NULL,
        skill TEXT NOT NULL,
        PRIMARY KEY (job_id, skill),
        FOREIGN KEY (job_id) REFERENCES jobs(id)
      );
    SQL
  end

  # 求人を保存（トランザクション使用）
  def save_job(job)
    @db.transaction do
      # jobsテーブルに挿入
      @db.execute(
        <<-SQL,
          INSERT OR REPLACE INTO jobs (
            id, title, company_name, company_location, description,
            salary_min, salary_max, salary_currency, salary_is_public,
            employment_type, posted_date, application_deadline, is_remote, status
          ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        SQL
        job.id,
        job.title,
        job.company.name,
        job.company.location,
        job.description,
        job.salary.min,
        job.salary.max,
        job.salary.currency,
        job.salary.is_public ? 1 : 0,
        job.employment_type,
        job.posted_date,
        job.application_deadline,
        job.is_remote ? 1 : 0,
        job.status
      )

      # 既存のスキルを削除
      @db.execute('DELETE FROM job_skills WHERE job_id = ?', job.id)

      # 新しいスキルを挿入
      job.required_skills.each do |skill|
        @db.execute(
          'INSERT INTO job_skills (job_id, skill) VALUES (?, ?)',
          job.id,
          skill
        )
      end
    end
  rescue SQLite3::Exception => e
    puts "求人の保存に失敗: #{e.message}"
    raise
  end

  # 複数の求人を保存
  def save_jobs(jobs)
    jobs.each { |job| save_job(job) }
  end

  # すべての求人を取得
  def get_all_jobs
    jobs = []

    @db.execute('SELECT * FROM jobs') do |row|
      # スキルを取得
      skills = @db.execute(
        'SELECT skill FROM job_skills WHERE job_id = ?',
        row['id']
      ).map { |skill_row| skill_row['skill'] }

      # 企業情報を作成
      company = Company.new(
        name: row['company_name'],
        location: row['company_location']
      )

      # 給与情報を作成
      salary = Salary.new(
        min: row['salary_min'],
        max: row['salary_max'],
        currency: row['salary_currency'],
        is_public: row['salary_is_public'] == 1
      )

      # Jobオブジェクトを作成
      job = Job.new(
        id: row['id'],
        title: row['title'],
        company: company,
        description: row['description'],
        salary: salary,
        employment_type: row['employment_type'],
        required_skills: skills,
        posted_date: row['posted_date'],
        application_deadline: row['application_deadline'],
        is_remote: row['is_remote'] == 1,
        status: row['status']
      )

      jobs << job
    end

    jobs
  rescue SQLite3::Exception => e
    puts "求人の取得に失敗: #{e.message}"
    []
  end

  # スキルで絞り込み
  def get_jobs_by_skill(skill)
    jobs = []

    # JOINクエリを使用して効率的に取得
    query = <<-SQL
      SELECT DISTINCT jobs.* FROM jobs
      INNER JOIN job_skills ON jobs.id = job_skills.job_id
      WHERE job_skills.skill = ?
    SQL

    @db.execute(query, skill) do |row|
      # スキルを取得
      skills = @db.execute(
        'SELECT skill FROM job_skills WHERE job_id = ?',
        row['id']
      ).map { |skill_row| skill_row['skill'] }

      # 企業情報を作成
      company = Company.new(
        name: row['company_name'],
        location: row['company_location']
      )

      # 給与情報を作成
      salary = Salary.new(
        min: row['salary_min'],
        max: row['salary_max'],
        currency: row['salary_currency'],
        is_public: row['salary_is_public'] == 1
      )

      # Jobオブジェクトを作成
      job = Job.new(
        id: row['id'],
        title: row['title'],
        company: company,
        description: row['description'],
        salary: salary,
        employment_type: row['employment_type'],
        required_skills: skills,
        posted_date: row['posted_date'],
        application_deadline: row['application_deadline'],
        is_remote: row['is_remote'] == 1,
        status: row['status']
      )

      jobs << job
    end

    jobs
  rescue SQLite3::Exception => e
    puts "求人の取得に失敗: #{e.message}"
    []
  end

  # リモート可能な求人を取得
  def get_jobs_by_remote
    jobs = []

    @db.execute('SELECT * FROM jobs WHERE is_remote = 1') do |row|
      # スキルを取得
      skills = @db.execute(
        'SELECT skill FROM job_skills WHERE job_id = ?',
        row['id']
      ).map { |skill_row| skill_row['skill'] }

      # 企業情報を作成
      company = Company.new(
        name: row['company_name'],
        location: row['company_location']
      )

      # 給与情報を作成
      salary = Salary.new(
        min: row['salary_min'],
        max: row['salary_max'],
        currency: row['salary_currency'],
        is_public: row['salary_is_public'] == 1
      )

      # Jobオブジェクトを作成
      job = Job.new(
        id: row['id'],
        title: row['title'],
        company: company,
        description: row['description'],
        salary: salary,
        employment_type: row['employment_type'],
        required_skills: skills,
        posted_date: row['posted_date'],
        application_deadline: row['application_deadline'],
        is_remote: row['is_remote'] == 1,
        status: row['status']
      )

      jobs << job
    end

    jobs
  rescue SQLite3::Exception => e
    puts "求人の取得に失敗: #{e.message}"
    []
  end

  # すべての求人を削除（テスト用）
  def delete_all_jobs
    @db.transaction do
      @db.execute('DELETE FROM job_skills')
      @db.execute('DELETE FROM jobs')
    end
    puts "すべての求人を削除しました"
  rescue SQLite3::Exception => e
    puts "削除に失敗: #{e.message}"
    raise
  end
end

# メイン処理
if __FILE__ == $0
  # データベースファイル名（環境変数から取得、デフォルトはjobs.db）
  db_path = ENV['DB_PATH'] || 'jobs.db'
  
  # データベースを作成
  db = JobDatabase.new(db_path)

  begin
    # テーブルを作成
    db.create_tables
    puts "テーブルを作成しました\n\n"

    # 既存データを削除
    db.delete_all_jobs
    puts "既存データを削除しました\n\n"

    # XMLから求人データを読み込み
    jobs_raw = load_jobs_from_xml('../jobs_new.xml')
    jobs = normalize_jobs(jobs_raw)
    puts "#{jobs.length}件の求人を読み込みました\n\n"

    # データベースに保存
    db.save_jobs(jobs)
    puts "#{jobs.length}件の求人をデータベースに保存しました\n\n"

    # すべての求人を取得
    puts "=== データベースから取得したすべての求人 ==="
    all_jobs = db.get_all_jobs
    puts "取得件数: #{all_jobs.length}件\n"
    all_jobs.each { |job| print_job(job) }

    # Goスキルを持つ求人を取得
    puts "\n=== Goスキルを持つ求人（データベースから取得） ==="
    go_jobs = db.get_jobs_by_skill('Go')
    puts "取得件数: #{go_jobs.length}件\n"
    go_jobs.each { |job| print_job(job) }

    # リモート可能な求人を取得
    puts "\n=== リモート可能な求人（データベースから取得） ==="
    remote_jobs = db.get_jobs_by_remote
    puts "取得件数: #{remote_jobs.length}件\n"
    remote_jobs.each { |job| print_job(job) }

  ensure
    # データベースを閉じる
    db.close
    puts "\nデータベース接続を閉じました"
  end
end
