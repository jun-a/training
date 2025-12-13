#!/usr/bin/env ruby
# encoding: utf-8

# レベル1: 基本的なXML読み込み
# Go言語のmain.goに対応するRuby実装

require 'nokogiri'

# 求人情報を表すクラス
class Job
  attr_accessor :id, :title, :company, :description, :salary, 
                :employment_type, :required_skills, :posted_date, 
                :application_deadline, :is_remote

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

# 給与情報を表すクラス
class Salary
  attr_accessor :min, :max, :currency

  def initialize(min:, max:, currency: 'JPY')
    @min = min
    @max = max
    @currency = currency
  end
end

# XMLファイルから求人データを読み込む
def load_jobs_from_xml(filename)
  jobs = []

  File.open(filename) do |file|
    doc = Nokogiri::XML(file)

    # 各jobノードを処理
    doc.xpath('//job').each do |job_node|
      # 企業情報を取得
      company_node = job_node.at_xpath('company')
      company = Company.new(
        name: company_node.at_xpath('name').text,
        location: company_node.at_xpath('location').text
      )

      # 給与情報を取得
      salary_node = job_node.at_xpath('salary')
      salary = Salary.new(
        min: salary_node.at_xpath('min').text.to_i,
        max: salary_node.at_xpath('max').text.to_i,
        currency: salary_node.at_xpath('currency').text
      )

      # 必要スキルを取得
      skills = job_node.xpath('required_skills/skill').map(&:text)

      # リモート可否を取得
      is_remote = job_node.at_xpath('is_remote').text == 'true'

      # Job オブジェクトを作成
      job = Job.new(
        id: job_node.at_xpath('id').text,
        title: job_node.at_xpath('title').text,
        company: company,
        description: job_node.at_xpath('description').text,
        salary: salary,
        employment_type: job_node.at_xpath('employment_type').text,
        required_skills: skills,
        posted_date: job_node.at_xpath('posted_date').text,
        application_deadline: job_node.at_xpath('application_deadline').text,
        is_remote: is_remote
      )

      jobs << job
    end
  end

  jobs
rescue => e
  puts "エラー: #{e.message}"
  []
end

# 指定されたスキルを持つ求人をフィルタリング
def filter_by_skill(jobs, skill)
  jobs.select do |job|
    job.required_skills.include?(skill)
  end
end

# リモート可能な求人をフィルタリング
def filter_by_remote(jobs)
  jobs.select { |job| job.is_remote }
end

# 指定された給与範囲内の求人をフィルタリング
def filter_by_salary_range(jobs, min_salary, max_salary)
  jobs.select do |job|
    job.salary.min >= min_salary && job.salary.max <= max_salary
  end
end

# 求人情報を整形して出力
def print_job(job)
  puts "====================================="
  puts "求人ID: #{job.id}"
  puts "タイトル: #{job.title}"
  puts "企業名: #{job.company.name}"
  puts "勤務地: #{job.company.location}"
  puts "給与: #{job.salary.min / 10000}万円 〜 #{job.salary.max / 10000}万円"
  puts "雇用形態: #{job.employment_type}"
  puts "必要スキル: #{job.required_skills.join(', ')}"
  puts "リモート: #{job.is_remote}"
  puts "掲載日: #{job.posted_date}"
  puts "応募締切: #{job.application_deadline}"
  puts "説明: #{job.description}"
  puts "=====================================\n\n"
end

# メイン処理
if __FILE__ == $0
  # XMLファイルから求人データを読み込む
  jobs = load_jobs_from_xml('../jobs.xml')

  puts "読み込んだ求人数: #{jobs.length}件\n\n"

  # すべての求人を表示
  puts "=== すべての求人 ==="
  jobs.each { |job| print_job(job) }

  # Goスキルを持つ求人をフィルタリング
  puts "\n=== Goスキルを持つ求人 ==="
  go_jobs = filter_by_skill(jobs, 'Go')
  go_jobs.each { |job| print_job(job) }

  # リモート可能な求人をフィルタリング
  puts "\n=== リモート可能な求人 ==="
  remote_jobs = filter_by_remote(jobs)
  remote_jobs.each { |job| print_job(job) }

  # 給与範囲でフィルタリング（600万円〜1000万円）
  puts "\n=== 給与600万円〜1000万円の求人 ==="
  salary_jobs = filter_by_salary_range(jobs, 6000000, 10000000)
  salary_jobs.each { |job| print_job(job) }
end
