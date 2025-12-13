package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// LoadJobsFromXMLAnswer はXMLファイルから求人データを読み込む（解答例）
func LoadJobsFromXMLAnswer(filename string) (*JobsRaw, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var jobs JobsRaw
	if err := xml.Unmarshal(data, &jobs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	return &jobs, nil
}

// TrimJobDataAnswer は求人データの不要な空白や改行を削除する（解答例）
func TrimJobDataAnswer(jobRaw JobRaw) JobRaw {
	// 複数の空白・改行を単一のスペースに置換する正規表現
	multiSpaceRegex := regexp.MustCompile(`\s+`)

	jobRaw.Title = strings.TrimSpace(jobRaw.Title)
	jobRaw.Company.Name = strings.TrimSpace(jobRaw.Company.Name)
	jobRaw.Company.Location = strings.TrimSpace(jobRaw.Company.Location)
	jobRaw.EmploymentType = strings.TrimSpace(jobRaw.EmploymentType)
	jobRaw.Status = strings.TrimSpace(jobRaw.Status)

	// 説明文は複数の空白・改行を単一のスペースに置換してからトリム
	jobRaw.Description = strings.TrimSpace(multiSpaceRegex.ReplaceAllString(jobRaw.Description, " "))

	// スキルの各要素をトリム
	for i := range jobRaw.RequiredSkills.Skills {
		jobRaw.RequiredSkills.Skills[i] = strings.TrimSpace(jobRaw.RequiredSkills.Skills[i])
	}

	return jobRaw
}

// NormalizeSalaryAnswer は様々な形式の給与データを正規化する（解答例）
func NormalizeSalaryAnswer(salaryStr string) Salary {
	salaryStr = strings.TrimSpace(salaryStr)

	// 空文字列または「非公開」
	if salaryStr == "" || strings.Contains(salaryStr, "非公開") {
		return Salary{
			Min:      0,
			Max:      0,
			Currency: "JPY",
			IsPublic: false,
		}
	}

	// カンマを削除
	salaryStr = strings.ReplaceAll(salaryStr, ",", "")
	// 全角数字を半角に変換（必要に応じて）
	salaryStr = strings.ReplaceAll(salaryStr, "〜", "～")

	// 月給の場合は12倍する
	isMonthly := strings.Contains(salaryStr, "月給")

	// 数値を抽出する正規表現（複数パターンに対応）
	numberRegex := regexp.MustCompile(`(\d+)`)
	matches := numberRegex.FindAllString(salaryStr, -1)

	if len(matches) == 0 {
		return Salary{
			Min:      0,
			Max:      0,
			Currency: "JPY",
			IsPublic: false,
		}
	}

	var min, max int

	// 数値が1つの場合
	if len(matches) == 1 {
		val, _ := strconv.Atoi(matches[0])

		// 「万円」が含まれている場合
		if strings.Contains(salaryStr, "万円") || strings.Contains(salaryStr, "万") {
			val = val * 10000
		}

		min = val
		max = val
	} else {
		// 数値が2つ以上の場合、最初と2番目を使用
		minVal, _ := strconv.Atoi(matches[0])
		maxVal, _ := strconv.Atoi(matches[1])

		// 「万円」が含まれている場合
		if strings.Contains(salaryStr, "万円") || strings.Contains(salaryStr, "万") {
			minVal = minVal * 10000
			maxVal = maxVal * 10000
		}

		min = minVal
		max = maxVal
	}

	// 月給の場合は12倍
	if isMonthly {
		min = min * 12
		max = max * 12
	}

	return Salary{
		Min:      min,
		Max:      max,
		Currency: "JPY",
		IsPublic: true,
	}
}

// NormalizeRemoteAnswer は様々な形式のリモートフラグを正規化する（解答例）
func NormalizeRemoteAnswer(remoteStr string) bool {
	remoteStr = strings.ToLower(strings.TrimSpace(remoteStr))

	switch remoteStr {
	case "true", "yes", "1", "remote":
		return true
	case "false", "no", "0", "":
		return false
	default:
		return false
	}
}

// NormalizeJobAnswer は生データを正規化して Job 構造体に変換する（解答例）
func NormalizeJobAnswer(jobRaw JobRaw) Job {
	// トリム処理
	jobRaw = TrimJobDataAnswer(jobRaw)

	// 給与の正規化
	salary := NormalizeSalaryAnswer(jobRaw.SalaryRaw)

	// リモートフラグの正規化
	isRemote := NormalizeRemoteAnswer(jobRaw.IsRemoteRaw)

	// Job 構造体を組み立てる
	return Job{
		ID:                  jobRaw.ID,
		Title:               jobRaw.Title,
		Company: Company{
			Name:     jobRaw.Company.Name,
			Location: jobRaw.Company.Location,
		},
		Description:         jobRaw.Description,
		Salary:              salary,
		EmploymentType:      jobRaw.EmploymentType,
		RequiredSkills:      jobRaw.RequiredSkills.Skills,
		PostedDate:          jobRaw.PostedDate,
		ApplicationDeadline: jobRaw.ApplicationDeadline,
		IsRemote:            isRemote,
		Status:              jobRaw.Status,
	}
}

// FilterBySkillAnswer は指定されたスキルを持つ求人をフィルタリングする（解答例）
func FilterBySkillAnswer(jobs []Job, skill string) []Job {
	var filtered []Job
	skillLower := strings.ToLower(skill)

	for _, job := range jobs {
		for _, s := range job.RequiredSkills {
			if strings.ToLower(s) == skillLower {
				filtered = append(filtered, job)
				break
			}
		}
	}
	return filtered
}

// FilterByRemoteAnswer はリモート可能な求人をフィルタリングする（解答例）
func FilterByRemoteAnswer(jobs []Job) []Job {
	var filtered []Job
	for _, job := range jobs {
		if job.IsRemote {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

// FilterBySalaryRangeAnswer は指定された給与範囲内の求人をフィルタリングする（解答例）
func FilterBySalaryRangeAnswer(jobs []Job, minSalary, maxSalary int) []Job {
	var filtered []Job
	for _, job := range jobs {
		// 給与が非公開の求人は除外
		if !job.Salary.IsPublic {
			continue
		}

		// 求人の最低給与がminSalary以上、かつ最高給与がmaxSalary以下
		if job.Salary.Min >= minSalary && job.Salary.Max <= maxSalary {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

// FilterByStatusAnswer はステータスで求人をフィルタリングする（解答例）
func FilterByStatusAnswer(jobs []Job, status string) []Job {
	var filtered []Job
	for _, job := range jobs {
		if job.Status == status {
			filtered = append(filtered, job)
		}
	}
	return filtered
}
