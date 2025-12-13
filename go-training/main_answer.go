package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// LoadJobsFromXMLAnswer はXMLファイルから求人データを読み込む（解答例）
func LoadJobsFromXMLAnswer(filename string) (*Jobs, error) {
	// ファイルを開く
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// ファイルの内容を読み込む
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// XMLをパース
	var jobs Jobs
	if err := xml.Unmarshal(data, &jobs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	return &jobs, nil
}

// FilterBySkillAnswer は指定されたスキルを持つ求人をフィルタリングする（解答例）
func FilterBySkillAnswer(jobs []Job, skill string) []Job {
	var filtered []Job
	for _, job := range jobs {
		for _, s := range job.RequiredSkills.Skills {
			if s == skill {
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
		// 求人の最低給与がminSalary以上、かつ最高給与がmaxSalary以下
		if job.Salary.Min >= minSalary && job.Salary.Max <= maxSalary {
			filtered = append(filtered, job)
		}
	}
	return filtered
}
