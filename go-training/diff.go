package main

import (
	"fmt"
)

// CompareJobs は新旧の求人データを比較して差分を抽出する
func CompareJobs(oldJobs, newJobs []Job) []JobDiff {
	// TODO: ここに新旧データの比較処理を実装してください
	//
	// 手順:
	// 1. 新旧の求人をIDでマップ化する
	// 2. 新しい求人データをループし、以下を判定:
	//    - 旧データに存在しない -> "added"
	//    - 旧データに存在し、内容が異なる -> "updated"
	//    - 旧データに存在し、内容が同じ -> "unchanged"
	// 3. 旧データにあるが新データにないものを検索 -> "deleted"
	// 4. 更新された求人については、何が変更されたかを詳細に記録
	//    - タイトルの変更
	//    - 給与の変更
	//    - スキルの追加・削除
	//    - 締切日の変更
	//    - ステータスの変更
	//    - 説明文の変更
	//
	// ヒント:
	// - map[string]Job を使って ID をキーにした検索を高速化
	// - スキルの差分は、追加されたスキルと削除されたスキルを別々に抽出

	return nil
}

// FindSkillsDiff はスキルの差分を検出する
func FindSkillsDiff(oldSkills, newSkills []string) (added []string, removed []string) {
	// TODO: ここにスキルの差分を検出するコードを実装してください
	//
	// 手順:
	// 1. oldSkillsにあってnewSkillsにないもの -> removed
	// 2. newSkillsにあってoldSkillsにないもの -> added
	//
	// ヒント:
	// - map[string]bool を使って存在チェックを高速化

	return nil, nil
}

// FilterByChangeType は指定された変更タイプの差分をフィルタリングする
func FilterByChangeType(diffs []JobDiff, changeType string) []JobDiff {
	// TODO: ここに指定された変更タイプの差分をフィルタリングするコードを実装してください
	// changeType: "added", "updated", "deleted", "unchanged"

	return nil
}

// FilterBySalaryIncreased は給与が上がった求人をフィルタリングする
func FilterBySalaryIncreased(diffs []JobDiff) []JobDiff {
	// TODO: ここに給与が上がった求人をフィルタリングするコードを実装してください
	// ヒント:
	// - ChangeType が "updated" かつ SalaryChanged が true
	// - NewMin > OldMin または NewMax > OldMax

	return nil
}

// FilterBySkillsAdded はスキルが追加された求人をフィルタリングする
func FilterBySkillsAdded(diffs []JobDiff) []JobDiff {
	// TODO: ここにスキルが追加された求人をフィルタリングするコードを実装してください
	// ヒント: len(SkillsAdded) > 0

	return nil
}

// FilterByStatusClosed はステータスがクローズに変更された求人をフィルタリングする
func FilterByStatusClosed(diffs []JobDiff) []JobDiff {
	// TODO: ここにステータスがクローズに変更された求人をフィルタリングするコードを実装してください
	// ヒント: StatusChanged が true かつ NewStatus が "closed"

	return nil
}

// PrintJobDiff は差分情報を整形して出力する
func PrintJobDiff(diff JobDiff) {
	fmt.Printf("=====================================\n")
	fmt.Printf("求人ID: %s\n", diff.ID)
	fmt.Printf("変更タイプ: %s\n", diff.ChangeType)

	if diff.ChangeType == "updated" {
		if diff.TitleChanged {
			fmt.Printf("  [変更] タイトルが更新されました\n")
		}

		if diff.SalaryChanged {
			fmt.Printf("  [変更] 給与: %d万円〜%d万円 -> %d万円〜%d万円\n",
				diff.SalaryDiff.OldMin/10000,
				diff.SalaryDiff.OldMax/10000,
				diff.SalaryDiff.NewMin/10000,
				diff.SalaryDiff.NewMax/10000,
			)
		}

		if len(diff.SkillsAdded) > 0 {
			fmt.Printf("  [追加] スキル: %v\n", diff.SkillsAdded)
		}

		if len(diff.SkillsRemoved) > 0 {
			fmt.Printf("  [削除] スキル: %v\n", diff.SkillsRemoved)
		}

		if diff.DeadlineChanged {
			fmt.Printf("  [変更] 締切日: %s -> %s\n", diff.OldDeadline, diff.NewDeadline)
		}

		if diff.StatusChanged {
			fmt.Printf("  [変更] ステータス: %s -> %s\n", diff.OldStatus, diff.NewStatus)
		}

		if diff.DescriptionChanged {
			fmt.Printf("  [変更] 説明文が更新されました\n")
		}
	}

	fmt.Printf("=====================================\n\n")
}

// PrintDiffSummary は差分のサマリーを出力する
func PrintDiffSummary(diffs []JobDiff) {
	// TODO: ここに差分のサマリーを出力するコードを実装してください
	//
	// 出力内容:
	// - 新規追加された求人の数
	// - 更新された求人の数
	// - 削除された求人の数
	// - 変更なしの求人の数
	//
	// ヒント: ChangeType をカウントする

	added := 0
	updated := 0
	deleted := 0
	unchanged := 0

	for _, diff := range diffs {
		switch diff.ChangeType {
		case "added":
			added++
		case "updated":
			updated++
		case "deleted":
			deleted++
		case "unchanged":
			unchanged++
		}
	}

	fmt.Printf("=== 差分サマリー ===\n")
	fmt.Printf("新規追加: %d件\n", added)
	fmt.Printf("更新: %d件\n", updated)
	fmt.Printf("削除: %d件\n", deleted)
	fmt.Printf("変更なし: %d件\n", unchanged)
	fmt.Printf("合計: %d件\n\n", added+updated+deleted+unchanged)
}
