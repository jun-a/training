#!/usr/bin/env ruby
# encoding: utf-8

# レベル2: データ正規化・トリム処理
# Go言語のprocessing.goに対応するRuby実装

require 'nokogiri'

# 正規化後の求人情報を表すクラス
class Job
  attr_accessor :id, :title, :company, :description, :salary, 
                :employment_type, :required_skills, :posted_date, 
                :application_deadline, :is_remote, :status

  def initialize(attributes = {})
    @id = attributes[:id]
    @title = attributes[:title]
    @company = attributes[:company]
    @description = attributes[:description]
    @salary = attributes[:salary]
    @employment_type = attributes[:employment_type]
    @required_skills = attributes[:required_skills] || []
    @posted_date = attributes[:posted_date]
    @application_deadline = attributes[:application_deadline]
    @is_remote = attributes[:is_remote]
    @status = attributes[:status]
  end
end

# 企業情報を表すクラス
class Company
  attr_accessor :name, :location

  def initialize(name:, location:)
    @name = name
    @location = location
  end
end

# 給与情報を表すクラス（正規化後）
class Salary
  attr_accessor :min, :max, :currency, :is_public

  def initialize(min: 0, max: 0, currency: 'JPY', is_public: true)
    @min = min
    @max = max
    @currency = currency
    @is_public = is_public
  end
end

# XMLファイルから生データを読み込む
class JobRaw
  attr_accessor :id, :title, :company_name, :company_location, :description,
                :salary_raw, :employment_type, :required_skills, :posted_date,
                :application_deadline, :is_remote_raw, :status

  def initialize(attributes = {})
    @id = attributes[:id]
    @title = attributes[:title]
    @company_name = attributes[:company_name]
    @company_location = attributes[:company_location]
    @description = attributes[:description]
    @salary_raw = attributes[:salary_raw]
    @employment_type = attributes[:employment_type]
    @required_skills = attributes[:required_skills] || []
    @posted_date = attributes[:posted_date]
    @application_deadline = attributes[:application_deadline]
    @is_remote_raw = attributes[:is_remote_raw]
    @status = attributes[:status]
  end
end

# XMLファイルから求人データを読み込む
def load_jobs_from_xml(filename)
  file = File.open(filename)
  doc = Nokogiri::XML(file)
  file.close

  jobs_raw = []

  doc.xpath('//job').each do |job_node|
    # 生データとして読み込む
    job_raw = JobRaw.new(
      id: job_node.at_xpath('id')&.text,
      title: job_node.at_xpath('title')&.text,
      company_name: job_node.at_xpath('company/name')&.text,
      company_location: job_node.at_xpath('company/location')&.text,
      description: job_node.at_xpath('description')&.text,
      salary_raw: job_node.at_xpath('salary')&.text,
      employment_type: job_node.at_xpath('employment_type')&.text,
      required_skills: job_node.xpath('required_skills/skill').map(&:text),
      posted_date: job_node.at_xpath('posted_date')&.text,
      application_deadline: job_node.at_xpath('application_deadline')&.text,
      is_remote_raw: job_node.at_xpath('is_remote')&.text,
      status: job_node.at_xpath('status')&.text
    )

    jobs_raw << job_raw
  end

  jobs_raw
rescue => e
  puts "エラー: #{e.message}"
  []
end

# 求人データのトリム処理
def trim_job_data(job_raw)
  # 前後の空白を削除
  job_raw.title = job_raw.title&.strip
  job_raw.company_name = job_raw.company_name&.strip
  job_raw.company_location = job_raw.company_location&.strip
  job_raw.employment_type = job_raw.employment_type&.strip
  job_raw.status = job_raw.status&.strip

  # 説明文は複数の空白・改行を単一のスペースに置換してからトリム
  if job_raw.description
    job_raw.description = job_raw.description.gsub(/\s+/, ' ').strip
  end

  # スキルの各要素をトリム
  job_raw.required_skills = job_raw.required_skills.map(&:strip)

  job_raw
end

# 給与データの正規化
def normalize_salary(salary_str)
  return Salary.new(is_public: false) if salary_str.nil? || salary_str.strip.empty?

  salary_str = salary_str.strip

  # 非公開の場合
  if salary_str.include?('非公開') || salary_str.include?('応相談')
    return Salary.new(is_public: false)
  end

  # カンマと円マークを削除
  normalized = salary_str.gsub(/[,¥円]/, '')

  # 全角の波ダッシュと全角チルダを統一
  normalized = normalized.gsub(/[〜～]/, '~')

  # 数値を抽出（正規表現）
  numbers = normalized.scan(/(\d+)/).flatten.map(&:to_i)

  return Salary.new(is_public: false) if numbers.empty?

  # 万円かどうかをチェック
  is_man_yen = salary_str.include?('万円') || salary_str.include?('万')
  
  # 月給かどうかをチェック
  is_monthly = salary_str.include?('月給') || salary_str.include?('月額')

  # 数値を変換
  if is_man_yen
    numbers = numbers.map { |n| n * 10000 }
  end

  if is_monthly
    numbers = numbers.map { |n| n * 12 }
  end

  # MinとMaxを設定
  if numbers.length == 1
    min_salary = max_salary = numbers[0]
  elsif numbers.length >= 2
    min_salary = numbers[0]
    max_salary = numbers[1]
  else
    return Salary.new(is_public: false)
  end

  Salary.new(min: min_salary, max: max_salary, currency: 'JPY', is_public: true)
end

# リモートフラグの正規化
def normalize_remote(remote_str)
  return false if remote_str.nil? || remote_str.strip.empty?

  normalized = remote_str.strip.downcase

  case normalized
  when 'true', 'yes', '1', 'remote'
    true
  when 'false', 'no', '0'
    false
  else
    false
  end
end

# 生データを正規化してJobオブジェクトに変換
def normalize_job(job_raw)
  # トリム処理
  job_raw = trim_job_data(job_raw)

  # 給与を正規化
  salary = normalize_salary(job_raw.salary_raw)

  # リモートフラグを正規化
  is_remote = normalize_remote(job_raw.is_remote_raw)

  # 企業情報を作成
  company = Company.new(
    name: job_raw.company_name || '',
    location: job_raw.company_location || ''
  )

  # Job オブジェクトを作成
  Job.new(
    id: job_raw.id || '',
    title: job_raw.title || '',
    company: company,
    description: job_raw.description || '',
    salary: salary,
    employment_type: job_raw.employment_type || '',
    required_skills: job_raw.required_skills,
    posted_date: job_raw.posted_date || '',
    application_deadline: job_raw.application_deadline || '',
    is_remote: is_remote,
    status: job_raw.status || 'active'
  )
end

# 複数の求人データを正規化
def normalize_jobs(jobs_raw)
  jobs_raw.map { |job_raw| normalize_job(job_raw) }
end

# 指定されたスキルを持つ求人をフィルタリング（大文字小文字を無視）
def filter_by_skill(jobs, skill)
  skill_lower = skill.downcase
  jobs.select do |job|
    job.required_skills.any? { |s| s.downcase == skill_lower }
  end
end

# リモート可能な求人をフィルタリング
def filter_by_remote(jobs)
  jobs.select { |job| job.is_remote }
end

# 指定された給与範囲内の求人をフィルタリング
def filter_by_salary_range(jobs, min_salary, max_salary)
  jobs.select do |job|
    # 給与が公開されている場合のみ
    job.salary.is_public && 
      job.salary.min >= min_salary && 
      job.salary.max <= max_salary
  end
end

# ステータスで求人をフィルタリング
def filter_by_status(jobs, status)
  jobs.select { |job| job.status.downcase == status.downcase }
end

# 求人情報を整形して出力
def print_job(job)
  puts "====================================="
  puts "求人ID: #{job.id}"
  puts "タイトル: #{job.title}"
  puts "企業名: #{job.company.name}"
  puts "勤務地: #{job.company.location}"
  
  if job.salary.is_public
    puts "給与: #{job.salary.min / 10000}万円 〜 #{job.salary.max / 10000}万円"
  else
    puts "給与: 非公開"
  end
  
  puts "雇用形態: #{job.employment_type}"
  puts "必要スキル: #{job.required_skills.join(', ')}"
  puts "リモート: #{job.is_remote}"
  puts "ステータス: #{job.status}"
  puts "掲載日: #{job.posted_date}"
  puts "応募締切: #{job.application_deadline}"
  puts "説明: #{job.description}"
  puts "=====================================\n\n"
end

# メイン処理
if __FILE__ == $0
  # XMLから読み込み
  jobs_raw = load_jobs_from_xml('../jobs_new.xml')
  puts "読み込んだ求人数（生データ）: #{jobs_raw.length}件\n\n"

  # 正規化
  jobs = normalize_jobs(jobs_raw)
  puts "正規化後の求人数: #{jobs.length}件\n\n"

  # 結果を表示
  puts "=== 正規化後のすべての求人 ==="
  jobs.each { |job| print_job(job) }

  # フィルタリングのテスト
  puts "\n=== Goスキルを持つ求人 ==="
  go_jobs = filter_by_skill(jobs, 'Go')
  puts "該当件数: #{go_jobs.length}件"
  go_jobs.each { |job| print_job(job) }

  puts "\n=== リモート可能な求人 ==="
  remote_jobs = filter_by_remote(jobs)
  puts "該当件数: #{remote_jobs.length}件"
  remote_jobs.each { |job| print_job(job) }

  puts "\n=== アクティブな求人 ==="
  active_jobs = filter_by_status(jobs, 'active')
  puts "該当件数: #{active_jobs.length}件"
  active_jobs.each { |job| print_job(job) }

  puts "\n=== 給与600万円〜1200万円の求人 ==="
  salary_jobs = filter_by_salary_range(jobs, 6000000, 12000000)
  puts "該当件数: #{salary_jobs.length}件"
  salary_jobs.each { |job| print_job(job) }
end
