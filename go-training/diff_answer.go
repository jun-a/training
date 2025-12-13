package main

// CompareJobsAnswer は新旧の求人データを比較して差分を抽出する（解答例）
func CompareJobsAnswer(oldJobs, newJobs []Job) []JobDiff {
	var diffs []JobDiff

	// 新旧の求人をIDでマップ化
	oldJobMap := make(map[string]Job)
	for _, job := range oldJobs {
		oldJobMap[job.ID] = job
	}

	newJobMap := make(map[string]Job)
	for _, job := range newJobs {
		newJobMap[job.ID] = job
	}

	// 新しい求人データをループ
	for _, newJob := range newJobs {
		oldJob, exists := oldJobMap[newJob.ID]

		if !exists {
			// 新規追加
			diffs = append(diffs, JobDiff{
				ID:         newJob.ID,
				ChangeType: "added",
			})
		} else {
			// 既存の求人 - 差分をチェック
			diff := JobDiff{
				ID:         newJob.ID,
				ChangeType: "unchanged",
			}

			// タイトルの変更
			if oldJob.Title != newJob.Title {
				diff.TitleChanged = true
				diff.ChangeType = "updated"
			}

			// 給与の変更
			if oldJob.Salary.Min != newJob.Salary.Min || oldJob.Salary.Max != newJob.Salary.Max {
				diff.SalaryChanged = true
				diff.SalaryDiff = SalaryDiff{
					OldMin: oldJob.Salary.Min,
					OldMax: oldJob.Salary.Max,
					NewMin: newJob.Salary.Min,
					NewMax: newJob.Salary.Max,
				}
				diff.ChangeType = "updated"
			}

			// スキルの差分
			skillsAdded, skillsRemoved := FindSkillsDiffAnswer(oldJob.RequiredSkills, newJob.RequiredSkills)
			if len(skillsAdded) > 0 || len(skillsRemoved) > 0 {
				diff.SkillsAdded = skillsAdded
				diff.SkillsRemoved = skillsRemoved
				diff.ChangeType = "updated"
			}

			// 締切日の変更
			if oldJob.ApplicationDeadline != newJob.ApplicationDeadline {
				diff.DeadlineChanged = true
				diff.OldDeadline = oldJob.ApplicationDeadline
				diff.NewDeadline = newJob.ApplicationDeadline
				diff.ChangeType = "updated"
			}

			// ステータスの変更
			if oldJob.Status != newJob.Status {
				diff.StatusChanged = true
				diff.OldStatus = oldJob.Status
				diff.NewStatus = newJob.Status
				diff.ChangeType = "updated"
			}

			// 説明文の変更
			if oldJob.Description != newJob.Description {
				diff.DescriptionChanged = true
				diff.ChangeType = "updated"
			}

			diffs = append(diffs, diff)
		}
	}

	// 削除された求人を検索
	for _, oldJob := range oldJobs {
		if _, exists := newJobMap[oldJob.ID]; !exists {
			diffs = append(diffs, JobDiff{
				ID:         oldJob.ID,
				ChangeType: "deleted",
			})
		}
	}

	return diffs
}

// FindSkillsDiffAnswer はスキルの差分を検出する（解答例）
func FindSkillsDiffAnswer(oldSkills, newSkills []string) (added []string, removed []string) {
	// oldSkillsをマップ化
	oldSkillMap := make(map[string]bool)
	for _, skill := range oldSkills {
		oldSkillMap[skill] = true
	}

	// newSkillsをマップ化
	newSkillMap := make(map[string]bool)
	for _, skill := range newSkills {
		newSkillMap[skill] = true
	}

	// 追加されたスキルを検索
	for _, skill := range newSkills {
		if !oldSkillMap[skill] {
			added = append(added, skill)
		}
	}

	// 削除されたスキルを検索
	for _, skill := range oldSkills {
		if !newSkillMap[skill] {
			removed = append(removed, skill)
		}
	}

	return added, removed
}

// FilterByChangeTypeAnswer は指定された変更タイプの差分をフィルタリングする（解答例）
func FilterByChangeTypeAnswer(diffs []JobDiff, changeType string) []JobDiff {
	var filtered []JobDiff
	for _, diff := range diffs {
		if diff.ChangeType == changeType {
			filtered = append(filtered, diff)
		}
	}
	return filtered
}

// FilterBySalaryIncreasedAnswer は給与が上がった求人をフィルタリングする（解答例）
func FilterBySalaryIncreasedAnswer(diffs []JobDiff) []JobDiff {
	var filtered []JobDiff
	for _, diff := range diffs {
		if diff.ChangeType == "updated" && diff.SalaryChanged {
			if diff.SalaryDiff.NewMin > diff.SalaryDiff.OldMin ||
				diff.SalaryDiff.NewMax > diff.SalaryDiff.OldMax {
				filtered = append(filtered, diff)
			}
		}
	}
	return filtered
}

// FilterBySkillsAddedAnswer はスキルが追加された求人をフィルタリングする（解答例）
func FilterBySkillsAddedAnswer(diffs []JobDiff) []JobDiff {
	var filtered []JobDiff
	for _, diff := range diffs {
		if len(diff.SkillsAdded) > 0 {
			filtered = append(filtered, diff)
		}
	}
	return filtered
}

// FilterByStatusClosedAnswer はステータスがクローズに変更された求人をフィルタリングする（解答例）
func FilterByStatusClosedAnswer(diffs []JobDiff) []JobDiff {
	var filtered []JobDiff
	for _, diff := range diffs {
		if diff.StatusChanged && diff.NewStatus == "closed" {
			filtered = append(filtered, diff)
		}
	}
	return filtered
}
