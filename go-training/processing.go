package main

import (
	"fmt"
)

// processing.go - データ正規化・トリム処理の実装
//
// ベストプラクティス:
//
// 1. 正規表現のパッケージレベルでのコンパイル
//    関数内で regexp.MustCompile() を呼び出すと、関数が呼ばれるたびに
//    正規表現がコンパイルされてしまい、パフォーマンスが低下します。
//    
//    悪い例（関数内でコンパイル）:
//      func NormalizeSalary(s string) {
//          re := regexp.MustCompile(`\d+`)  // 毎回コンパイルされる
//          matches := re.FindAllString(s, -1)
//      }
//    
//    良い例（パッケージレベルでコンパイル）:
//      var numberRegex = regexp.MustCompile(`\d+`)
//      func NormalizeSalary(s string) {
//          matches := numberRegex.FindAllString(s, -1)  // 再利用
//      }
//
// 2. エラーハンドリング
//    strconv.Atoi() などの変換関数は必ずエラーを返します。
//    不正な入力データに対して適切にエラーハンドリングすることが重要です。
//
//    例:
//      if num, err := strconv.Atoi(str); err != nil {
//          // エラー処理: デフォルト値を使うか、エラーをログに記録
//          return 0
//      } else {
//          return num
//      }

// LoadJobsFromXML はXMLファイルから求人データを読み込む
func LoadJobsFromXML(filename string) (*JobsRaw, error) {
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

// TrimJobData は求人データの不要な空白や改行を削除する
func TrimJobData(jobRaw JobRaw) JobRaw {
	// TODO: ここに文字列のトリム処理を実装してください
	// 以下のフィールドから前後の空白・改行を削除する必要があります:
	// - Title
	// - Company.Name
	// - Company.Location
	// - Description（複数行の空白や改行も整理）
	// - EmploymentType
	// - RequiredSkills.Skills の各要素
	// - Status
	//
	// ヒント:
	// - strings.TrimSpace() で前後の空白を削除
	// - regexp を使って複数の空白・改行を単一のスペースに置換
	// 
	// 必要なインポート:
	// import (
	//     "regexp"
	//     "strings"
	// )

	return jobRaw
}

// NormalizeSalary は様々な形式の給与データを正規化する
func NormalizeSalary(salaryStr string) Salary {
	// TODO: ここに給与データの正規化処理を実装してください
	//
	// 対応すべき形式:
	// 1. "年収 6,000,000円 〜 10,000,000円" -> Min: 6000000, Max: 10000000
	// 2. "500万円〜800万円" -> Min: 5000000, Max: 8000000
	// 3. "7000000" -> Min: 7000000, Max: 7000000
	// 4. "月給 550,000円 〜 900,000円" -> Min: 6600000, Max: 10800000 (月給×12)
	// 5. "800〜1500万円" -> Min: 8000000, Max: 15000000
	// 6. "応相談（想定年収700〜1300万円）" -> Min: 7000000, Max: 13000000
	// 7. "¥4,500,000 - ¥7,000,000" -> Min: 4500000, Max: 7000000
	// 8. "年俸制 8,000,000円〜12,000,000円" -> Min: 8000000, Max: 12000000
	// 9. "" (空文字列) -> IsPublic: false
	// 10. "非公開" -> IsPublic: false
	//
	// ヒント:
	// - strings.TrimSpace() で前後の空白を削除
	// - strings.Contains() で特定の文字列が含まれているか確認
	// - regexp.MustCompile() で正規表現パターンを作成
	// - regexp.FindAllStringSubmatch() で数値を抽出
	// - strconv.Atoi() で文字列を数値に変換（エラーハンドリングを忘れずに）
	// - "万円"が含まれている場合は10000を掛ける
	// - "月給"が含まれている場合は12を掛ける
	// 
	// NOTE: ハイフン「-」とダッシュ「〜」の両方を考慮すること
	// NOTE: strconv.Atoi()のエラーは適切に処理すること（不正な文字列の場合）
	// 
	// 必要なインポート:
	// import (
	//     "regexp"
	//     "strconv"
	//     "strings"
	// )

	return Salary{
		Min:      0,
		Max:      0,
		Currency: "JPY",
		IsPublic: false,
	}
}

// NormalizeRemote は様々な形式のリモートフラグを正規化する
func NormalizeRemote(remoteStr string) bool {
	// TODO: ここにリモートフラグの正規化処理を実装してください
	//
	// 対応すべき形式:
	// - "true", "TRUE", "True" -> true
	// - "false", "FALSE", "False" -> false
	// - "yes", "YES", "Yes" -> true
	// - "no", "NO", "No" -> false
	// - "1" -> true
	// - "0" -> false
	// - "remote", "Remote", "REMOTE" -> true
	// - その他 -> false
	//
	// ヒント:
	// - strings.TrimSpace() で前後の空白を削除
	// - strings.ToLower() で小文字に変換して比較
	// 
	// 必要なインポート:
	// import "strings"

	return false
}

// NormalizeJob は生データを正規化して Job 構造体に変換する
func NormalizeJob(jobRaw JobRaw) Job {
	// TODO: ここに正規化処理を実装してください
	//
	// 手順:
	// 1. TrimJobData() でトリム処理
	// 2. NormalizeSalary() で給与を正規化
	// 3. NormalizeRemote() でリモートフラグを正規化
	// 4. Job 構造体を組み立てる
	// 
	// NOTE: RequiredSkills.Skills を Job.RequiredSkills にマッピングすることを忘れずに

	return Job{}
}

// NormalizeJobs は複数の求人データを正規化する
func NormalizeJobs(jobsRaw *JobsRaw) []Job {
	// TODO: ここに複数の求人データを正規化するコードを実装してください
	// ヒント: NormalizeJob() を繰り返し呼び出す

	var jobs []Job
	for _, jobRaw := range jobsRaw.JobList {
		jobs = append(jobs, NormalizeJob(jobRaw))
	}
	return jobs
}

// FilterBySkill は指定されたスキルを持つ求人をフィルタリングする
func FilterBySkill(jobs []Job, skill string) []Job {
	// TODO: ここに指定されたスキルを持つ求人をフィルタリングするコードを実装してください
	// 注意: スキル名の比較時は大文字小文字を無視する（strings.EqualFoldを使う）
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
	// ヒント:
	// - 給与が非公開の求人は除外する
	// - 求人の最低給与がminSalary以上、最高給与がmaxSalary以下のものを抽出
	return nil
}

// FilterByStatus はステータスで求人をフィルタリングする
func FilterByStatus(jobs []Job, status string) []Job {
	// TODO: ここに指定されたステータスの求人をフィルタリングするコードを実装してください
	// status: "active" または "closed"
	return nil
}

// PrintJob は求人情報を整形して出力する
func PrintJob(job Job) {
	fmt.Printf("=====================================\n")
	fmt.Printf("求人ID: %s\n", job.ID)
	fmt.Printf("タイトル: %s\n", job.Title)
	fmt.Printf("企業名: %s\n", job.Company.Name)
	fmt.Printf("勤務地: %s\n", job.Company.Location)

	if job.Salary.IsPublic {
		fmt.Printf("給与: %d万円 〜 %d万円\n", job.Salary.Min/10000, job.Salary.Max/10000)
	} else {
		fmt.Printf("給与: 非公開\n")
	}

	fmt.Printf("雇用形態: %s\n", job.EmploymentType)
	fmt.Printf("必要スキル: %v\n", job.RequiredSkills)
	fmt.Printf("リモート: %v\n", job.IsRemote)
	fmt.Printf("ステータス: %s\n", job.Status)
	fmt.Printf("掲載日: %s\n", job.PostedDate)
	fmt.Printf("応募締切: %s\n", job.ApplicationDeadline)
	fmt.Printf("説明: %s\n", job.Description)
	fmt.Printf("=====================================\n\n")
}
