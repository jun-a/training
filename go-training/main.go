package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// Job は求人情報を表す構造体
type Job struct {
	ID                  string          `xml:"id"`
	Title               string          `xml:"title"`
	Company             Company         `xml:"company"`
	Description         string          `xml:"description"`
	Salary              Salary          `xml:"salary"`
	EmploymentType      string          `xml:"employment_type"`
	RequiredSkills      RequiredSkills  `xml:"required_skills"`
	PostedDate          string          `xml:"posted_date"`
	ApplicationDeadline string          `xml:"application_deadline"`
	IsRemote            bool            `xml:"is_remote"`
}

// Company は企業情報を表す構造体
type Company struct {
	Name     string `xml:"name"`
	Location string `xml:"location"`
}

// Salary は給与情報を表す構造体
type Salary struct {
	Min      int    `xml:"min"`
	Max      int    `xml:"max"`
	Currency string `xml:"currency"`
}

// RequiredSkills は必要スキルを表す構造体
type RequiredSkills struct {
	Skills []string `xml:"skill"`
}

// Jobs はJobのリストを表す構造体
type Jobs struct {
	XMLName xml.Name `xml:"jobs"`
	JobList []Job    `xml:"job"`
}

// LoadJobsFromXML はXMLファイルから求人データを読み込む
func LoadJobsFromXML(filename string) (*Jobs, error) {
	// TODO: ここにXMLファイルを読み込むコードを実装してください
	// ヒント:
	// 1. os.Open()でファイルを開く
	// 2. defer file.Close()でファイルを閉じる
	// 3. io.ReadAll()でファイルの内容を読み込む
	// 4. xml.Unmarshal()でXMLをパース
	return nil, fmt.Errorf("not implemented")
}

// FilterBySkill は指定されたスキルを持つ求人をフィルタリングする
func FilterBySkill(jobs []Job, skill string) []Job {
	// TODO: ここに指定されたスキルを持つ求人をフィルタリングするコードを実装してください
	return nil
}

// FilterByRemote はリモート可能な求人をフィルタリングする
func FilterByRemote(jobs []Job) []Job {
	// TODO: ここにリモート可能な求人をフィルタリングするコードを実装してください
	return nil
}

// FilterBySalaryRange は指定された給与範囲内の求人をフィルタリングする
func FilterBySalaryRange(jobs []Job, minSalary, maxSalary int) []Job {
	// TODO: ここに指定された給与範囲内の求人をフィルタリングするコードを実装してください
	// ヒント: 求人の最低給与がminSalary以上、最高給与がmaxSalary以下のものを抽出
	return nil
}

// PrintJob は求人情報を整形して出力する
func PrintJob(job Job) {
	fmt.Printf("=====================================\n")
	fmt.Printf("求人ID: %s\n", job.ID)
	fmt.Printf("タイトル: %s\n", job.Title)
	fmt.Printf("企業名: %s\n", job.Company.Name)
	fmt.Printf("勤務地: %s\n", job.Company.Location)
	fmt.Printf("給与: %d万円 〜 %d万円\n", job.Salary.Min/10000, job.Salary.Max/10000)
	fmt.Printf("雇用形態: %s\n", job.EmploymentType)
	fmt.Printf("必要スキル: %v\n", job.RequiredSkills.Skills)
	fmt.Printf("リモート: %v\n", job.IsRemote)
	fmt.Printf("掲載日: %s\n", job.PostedDate)
	fmt.Printf("応募締切: %s\n", job.ApplicationDeadline)
	fmt.Printf("説明: %s\n", job.Description)
	fmt.Printf("=====================================\n\n")
}

func main() {
	// XMLファイルから求人データを読み込む
	jobs, err := LoadJobsFromXML("jobs.xml")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		return
	}

	fmt.Printf("読み込んだ求人数: %d件\n\n", len(jobs.JobList))

	// すべての求人を表示
	fmt.Println("=== すべての求人 ===")
	for _, job := range jobs.JobList {
		PrintJob(job)
	}

	// Goスキルを持つ求人をフィルタリング
	fmt.Println("\n=== Goスキルを持つ求人 ===")
	goJobs := FilterBySkill(jobs.JobList, "Go")
	for _, job := range goJobs {
		PrintJob(job)
	}

	// リモート可能な求人をフィルタリング
	fmt.Println("\n=== リモート可能な求人 ===")
	remoteJobs := FilterByRemote(jobs.JobList)
	for _, job := range remoteJobs {
		PrintJob(job)
	}

	// 給与範囲でフィルタリング（600万円〜1000万円）
	fmt.Println("\n=== 給与600万円〜1000万円の求人 ===")
	salaryJobs := FilterBySalaryRange(jobs.JobList, 6000000, 10000000)
	for _, job := range salaryJobs {
		PrintJob(job)
	}
}
