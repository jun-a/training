#!/usr/bin/env ruby
# encoding: utf-8

# レベル3: 差分比較
# Go言語のdiff.goに対応するRuby実装

require_relative 'processing'

# 求人の差分情報を表すクラス
class JobDiff
  attr_accessor :job_id, :change_type, :old_job, :new_job, :changes

  def initialize(job_id:, change_type:, old_job: nil, new_job: nil, changes: [])
    @job_id = job_id
    @change_type = change_type  # "added", "updated", "deleted", "unchanged"
    @old_job = old_job
    @new_job = new_job
    @changes = changes
  end
end

# スキルの差分情報を表すクラス
class SkillsDiff
  attr_accessor :added, :removed

  def initialize(added: [], removed: [])
    @added = added
    @removed = removed
  end
end

# 新旧の求人データを比較
def compare_jobs(old_jobs, new_jobs)
  diffs = []

  # 旧求人をマップ化（ID -> Job）
  old_jobs_map = old_jobs.each_with_object({}) { |job, hash| hash[job.id] = job }

  # 新求人をマップ化（ID -> Job）
  new_jobs_map = new_jobs.each_with_object({}) { |job, hash| hash[job.id] = job }

  # 新しい求人をループ
  new_jobs.each do |new_job|
    old_job = old_jobs_map[new_job.id]

    if old_job.nil?
      # 新規追加
      diffs << JobDiff.new(
        job_id: new_job.id,
        change_type: 'added',
        new_job: new_job,
        changes: ['新規追加']
      )
    else
      # 既存求人：変更を比較
      changes = []

      # タイトルの変更
      if old_job.title != new_job.title
        changes << "タイトル変更: '#{old_job.title}' → '#{new_job.title}'"
      end

      # 給与の変更
      if old_job.salary.is_public && new_job.salary.is_public
        if old_job.salary.min != new_job.salary.min || old_job.salary.max != new_job.salary.max
          old_range = "#{old_job.salary.min / 10000}〜#{old_job.salary.max / 10000}万円"
          new_range = "#{new_job.salary.min / 10000}〜#{new_job.salary.max / 10000}万円"
          changes << "給与変更: #{old_range} → #{new_range}"
        end
      end

      # スキルの差分
      skills_diff = find_skills_diff(old_job.required_skills, new_job.required_skills)
      if skills_diff.added.any?
        changes << "スキル追加: #{skills_diff.added.join(', ')}"
      end
      if skills_diff.removed.any?
        changes << "スキル削除: #{skills_diff.removed.join(', ')}"
      end

      # 締切日の変更
      if old_job.application_deadline != new_job.application_deadline
        changes << "締切日変更: #{old_job.application_deadline} → #{new_job.application_deadline}"
      end

      # ステータスの変更
      if old_job.status != new_job.status
        changes << "ステータス変更: #{old_job.status} → #{new_job.status}"
      end

      # 説明文の変更
      if old_job.description != new_job.description
        changes << "説明文が更新されました"
      end

      # 変更があれば "updated"、なければ "unchanged"
      change_type = changes.any? ? 'updated' : 'unchanged'
      diffs << JobDiff.new(
        job_id: new_job.id,
        change_type: change_type,
        old_job: old_job,
        new_job: new_job,
        changes: changes
      )
    end
  end

  # 削除された求人を検出
  old_jobs.each do |old_job|
    unless new_jobs_map.key?(old_job.id)
      diffs << JobDiff.new(
        job_id: old_job.id,
        change_type: 'deleted',
        old_job: old_job,
        changes: ['削除された求人']
      )
    end
  end

  diffs
end

# スキルの差分を検出
def find_skills_diff(old_skills, new_skills)
  old_skills_set = old_skills.map(&:downcase).to_set
  new_skills_set = new_skills.map(&:downcase).to_set

  # 追加されたスキル
  added = new_skills.select { |skill| !old_skills_set.include?(skill.downcase) }

  # 削除されたスキル
  removed = old_skills.select { |skill| !new_skills_set.include?(skill.downcase) }

  SkillsDiff.new(added: added, removed: removed)
end

# 変更タイプでフィルタリング
def filter_by_change_type(diffs, change_type)
  diffs.select { |diff| diff.change_type == change_type }
end

# 給与が上がった求人のみフィルタリング
def filter_by_salary_increased(diffs)
  diffs.select do |diff|
    next false unless diff.change_type == 'updated'
    next false unless diff.old_job && diff.new_job
    next false unless diff.old_job.salary.is_public && diff.new_job.salary.is_public

    # 最低給与または最高給与が上がった場合
    diff.new_job.salary.min > diff.old_job.salary.min ||
      diff.new_job.salary.max > diff.old_job.salary.max
  end
end

# スキルが追加された求人のみフィルタリング
def filter_by_skills_added(diffs)
  diffs.select do |diff|
    next false unless diff.change_type == 'updated'
    
    skills_diff = find_skills_diff(diff.old_job.required_skills, diff.new_job.required_skills)
    skills_diff.added.any?
  end
end

# クローズに変更された求人のみフィルタリング
def filter_by_status_closed(diffs)
  diffs.select do |diff|
    next false unless diff.change_type == 'updated'
    next false unless diff.old_job && diff.new_job

    diff.old_job.status != 'closed' && diff.new_job.status == 'closed'
  end
end

# 差分のサマリーを表示
def print_diff_summary(diffs)
  added_count = diffs.count { |d| d.change_type == 'added' }
  updated_count = diffs.count { |d| d.change_type == 'updated' }
  deleted_count = diffs.count { |d| d.change_type == 'deleted' }
  unchanged_count = diffs.count { |d| d.change_type == 'unchanged' }

  puts "====================================="
  puts "差分サマリー"
  puts "====================================="
  puts "新規追加: #{added_count}件"
  puts "更新: #{updated_count}件"
  puts "削除: #{deleted_count}件"
  puts "変更なし: #{unchanged_count}件"
  puts "====================================="
  puts
end

# 求人の差分を表示
def print_job_diff(diff)
  puts "-------------------------------------"
  puts "求人ID: #{diff.job_id}"
  puts "変更タイプ: #{diff.change_type}"

  case diff.change_type
  when 'added'
    puts "新規追加された求人"
    puts "タイトル: #{diff.new_job.title}"
    puts "企業名: #{diff.new_job.company.name}"
    if diff.new_job.salary.is_public
      puts "給与: #{diff.new_job.salary.min / 10000}〜#{diff.new_job.salary.max / 10000}万円"
    end
  when 'deleted'
    puts "削除された求人"
    puts "タイトル: #{diff.old_job.title}"
    puts "企業名: #{diff.old_job.company.name}"
  when 'updated'
    puts "更新された求人"
    puts "タイトル: #{diff.new_job.title}"
    puts "変更内容:"
    diff.changes.each do |change|
      puts "  - #{change}"
    end
  when 'unchanged'
    puts "変更なし"
    puts "タイトル: #{diff.new_job.title}"
  end

  puts "-------------------------------------"
  puts
end

# メイン処理
if __FILE__ == $0
  # 旧データを読み込み
  old_jobs_raw = load_jobs_from_xml('../jobs_old.xml')
  old_jobs = normalize_jobs(old_jobs_raw)
  puts "旧データ: #{old_jobs.length}件\n\n"

  # 新データを読み込み
  new_jobs_raw = load_jobs_from_xml('../jobs_new.xml')
  new_jobs = normalize_jobs(new_jobs_raw)
  puts "新データ: #{new_jobs.length}件\n\n"

  # 差分を比較
  diffs = compare_jobs(old_jobs, new_jobs)

  # サマリー表示
  print_diff_summary(diffs)

  # すべての差分を表示
  puts "=== すべての差分 ==="
  diffs.each { |diff| print_job_diff(diff) }

  # 給与が上がった求人
  puts "\n=== 給与が上がった求人 ==="
  salary_increased = filter_by_salary_increased(diffs)
  puts "該当件数: #{salary_increased.length}件"
  salary_increased.each { |diff| print_job_diff(diff) }

  # スキルが追加された求人
  puts "\n=== スキルが追加された求人 ==="
  skills_added = filter_by_skills_added(diffs)
  puts "該当件数: #{skills_added.length}件"
  skills_added.each { |diff| print_job_diff(diff) }

  # 新規追加された求人
  puts "\n=== 新規追加された求人 ==="
  added_jobs = filter_by_change_type(diffs, 'added')
  puts "該当件数: #{added_jobs.length}件"
  added_jobs.each { |diff| print_job_diff(diff) }
end
