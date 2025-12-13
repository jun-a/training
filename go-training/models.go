package main

import "encoding/xml"

// ============================================================
// Raw types - for reading XML data
// ============================================================

// JobRaw はXMLから読み込んだ生のデータを表す構造体
// NOTE: jobs_new.xmlとjobs_old.xmlにはstatusフィールドがあるが、
//       jobs.xmlにはstatusフィールドがないことに注意
type JobRaw struct {
	ID                  string              `xml:"id"`
	Title               string              `xml:"title"`
	Company             CompanyRaw          `xml:"company"`
	Description         string              `xml:"description"`
	SalaryRaw           string              `xml:"salary"`
	EmploymentType      string              `xml:"employment_type"`
	RequiredSkills      RequiredSkillsRaw   `xml:"required_skills"`
	PostedDate          string              `xml:"posted_date"`
	ApplicationDeadline string              `xml:"application_deadline"`
	IsRemoteRaw         string              `xml:"is_remote"`
	Status              string              `xml:"status"`
}

// CompanyRaw は企業情報の生データを表す構造体
type CompanyRaw struct {
	Name     string `xml:"name"`
	Location string `xml:"location"`
}

// RequiredSkillsRaw はXMLから読み込む必要スキルを表す構造体
type RequiredSkillsRaw struct {
	Skills []string `xml:"skill"`
}

// JobsRaw はJobRawのリストを表す構造体
type JobsRaw struct {
	XMLName xml.Name `xml:"jobs"`
	JobList []JobRaw `xml:"job"`
}

// ============================================================
// Normalized types - for Level 1 (main.go)
// These are used when XML has clean, structured salary data
// ============================================================

// JobLevel1 はレベル1用の求人情報を表す構造体（構造化された給与データ用）
type JobLevel1 struct {
	ID                  string              `xml:"id"`
	Title               string              `xml:"title"`
	Company             CompanyLevel1       `xml:"company"`
	Description         string              `xml:"description"`
	Salary              SalaryLevel1        `xml:"salary"`
	EmploymentType      string              `xml:"employment_type"`
	RequiredSkills      RequiredSkillsRaw   `xml:"required_skills"`
	PostedDate          string              `xml:"posted_date"`
	ApplicationDeadline string              `xml:"application_deadline"`
	IsRemote            bool                `xml:"is_remote"`
}

// CompanyLevel1 は企業情報を表す構造体
type CompanyLevel1 struct {
	Name     string `xml:"name"`
	Location string `xml:"location"`
}

// SalaryLevel1 は構造化された給与情報を表す構造体（jobs.xml用）
type SalaryLevel1 struct {
	Min      int    `xml:"min"`
	Max      int    `xml:"max"`
	Currency string `xml:"currency"`
}

// RequiredSkillsLevel1 は必要スキルを表す構造体
type RequiredSkillsLevel1 struct {
	Skills []string `xml:"skill"`
}

// JobsLevel1 はJobLevel1のリストを表す構造体
type JobsLevel1 struct {
	XMLName xml.Name     `xml:"jobs"`
	JobList []JobLevel1  `xml:"job"`
}

// ============================================================
// Normalized types - for Level 2+ (processing.go, diff.go, database.go)
// These are used after normalization and processing
// ============================================================

// Job は正規化後の求人情報を表す構造体
type Job struct {
	ID                  string
	Title               string
	Company             Company
	Description         string
	Salary              Salary
	EmploymentType      string
	RequiredSkills      []string
	PostedDate          string
	ApplicationDeadline string
	IsRemote            bool
	Status              string
}

// Company は企業情報を表す構造体
type Company struct {
	Name     string
	Location string
}

// Salary は給与情報を表す構造体
type Salary struct {
	Min      int    // 年収の最低額（円）
	Max      int    // 年収の最高額（円）
	Currency string // 通貨単位
	IsPublic bool   // 給与が公開されているか
}

// ============================================================
// Diff types - for Level 3 (diff.go)
// ============================================================

// JobDiff は新旧求人データの差分を表す構造体
type JobDiff struct {
	ID                 string
	ChangeType         string // "added", "updated", "deleted", "unchanged"
	TitleChanged       bool
	SalaryChanged      bool
	SalaryDiff         SalaryDiff
	SkillsAdded        []string
	SkillsRemoved      []string
	DeadlineChanged    bool
	OldDeadline        string
	NewDeadline        string
	StatusChanged      bool
	OldStatus          string
	NewStatus          string
	DescriptionChanged bool
}

// SalaryDiff は給与の差分を表す構造体
type SalaryDiff struct {
	OldMin int
	OldMax int
	NewMin int
	NewMax int
}
