package main

import (
	"fmt"
)

// LoadJobsFromXMLLevel1 はXMLファイルから求人データを読み込む（レベル1用）
// jobs.xmlのように給与が構造化されたデータを読み込む
func LoadJobsFromXMLLevel1(filename string) (*JobsLevel1, error) {
	// TODO: ここにXMLファイルを読み込むコードを実装してください
	// ヒント:
	// 1. os.Open()でファイルを開く
	// 2. defer file.Close()でファイルを閉じる
	// 3. io.ReadAll()でファイルの内容を読み込む
	// 4. xml.Unmarshal()でXMLをパース
	// 
	// 必要なインポート:
	// import (
	//     "encoding/xml"
	//     "io"
	//     "os"
	// )
	return nil, fmt.Errorf("not implemented")
}

// FilterBySkillLevel1 は指定されたスキルを持つ求人をフィルタリングする（レベル1用）
func FilterBySkillLevel1(jobs []JobLevel1, skill string) []JobLevel1 {
	// TODO: ここに指定されたスキルを持つ求人をフィルタリングするコードを実装してください
	// HINT: strings.EqualFold() を使うと大文字小文字を無視した比較ができます
	return nil
}

// FilterByRemoteLevel1 はリモート可能な求人をフィルタリングする（レベル1用）
func FilterByRemoteLevel1(jobs []JobLevel1) []JobLevel1 {
	// TODO: ここにリモート可能な求人をフィルタリングするコードを実装してください
	return nil
}

// FilterBySalaryRangeLevel1 は指定された給与範囲内の求人をフィルタリングする（レベル1用）
func FilterBySalaryRangeLevel1(jobs []JobLevel1, minSalary, maxSalary int) []JobLevel1 {
	// TODO: ここに指定された給与範囲内の求人をフィルタリングするコードを実装してください
	// ヒント: 求人の最低給与がminSalary以上、最高給与がmaxSalary以下のものを抽出
	return nil
}

// PrintJobLevel1 は求人情報を整形して出力する（レベル1用）
func PrintJobLevel1(job JobLevel1) {
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
	jobs, err := LoadJobsFromXMLLevel1("jobs.xml")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		return
	}

	fmt.Printf("読み込んだ求人数: %d件\n\n", len(jobs.JobList))

	// すべての求人を表示
	fmt.Println("=== すべての求人 ===")
	for _, job := range jobs.JobList {
		PrintJobLevel1(job)
	}

	// Goスキルを持つ求人をフィルタリング
	fmt.Println("\n=== Goスキルを持つ求人 ===")
	goJobs := FilterBySkillLevel1(jobs.JobList, "Go")
	for _, job := range goJobs {
		PrintJobLevel1(job)
	}

	// リモート可能な求人をフィルタリング
	fmt.Println("\n=== リモート可能な求人 ===")
	remoteJobs := FilterByRemoteLevel1(jobs.JobList)
	for _, job := range remoteJobs {
		PrintJobLevel1(job)
	}

	// 給与範囲でフィルタリング（600万円〜1000万円）
	fmt.Println("\n=== 給与600万円〜1000万円の求人 ===")
	salaryJobs := FilterBySalaryRangeLevel1(jobs.JobList, 6000000, 10000000)
	for _, job := range salaryJobs {
		PrintJobLevel1(job)
	}
}
