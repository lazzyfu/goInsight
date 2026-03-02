package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"
	"gorm.io/datatypes"
)

type approvalDefinitionStage struct {
	Type      string   `json:"type"`
	Stage     int      `json:"stage"`
	StageName string   `json:"stage_name"`
	Approvers []string `json:"approvers"`
}

func normalizeAndValidateApprovalDefinition(raw datatypes.JSON) (datatypes.JSON, []string, error) {
	if len(raw) == 0 {
		return nil, nil, fmt.Errorf("审批流定义不能为空")
	}

	var stages []approvalDefinitionStage
	if err := json.Unmarshal(raw, &stages); err != nil {
		return nil, nil, fmt.Errorf("审批流定义格式不正确: %w", err)
	}
	if len(stages) == 0 {
		return nil, nil, fmt.Errorf("审批流定义不能为空")
	}

	stageSet := make(map[int]struct{}, len(stages))
	userSet := make(map[string]struct{})

	for idx := range stages {
		stage := &stages[idx]

		if stage.Stage <= 0 {
			return nil, nil, fmt.Errorf("第%d个阶段编号必须大于0", idx+1)
		}
		if _, ok := stageSet[stage.Stage]; ok {
			return nil, nil, fmt.Errorf("审批阶段编号%d重复", stage.Stage)
		}
		stageSet[stage.Stage] = struct{}{}

		stage.StageName = strings.TrimSpace(stage.StageName)
		if stage.StageName == "" {
			return nil, nil, fmt.Errorf("第%d个阶段名称不能为空", idx+1)
		}

		stage.Type = strings.ToUpper(strings.TrimSpace(stage.Type))
		if stage.Type != "AND" && stage.Type != "OR" {
			return nil, nil, fmt.Errorf("第%d个阶段审批类型无效，仅支持AND或OR", idx+1)
		}

		if len(stage.Approvers) == 0 {
			return nil, nil, fmt.Errorf("第%d个阶段审批人不能为空", idx+1)
		}
		approverSet := make(map[string]struct{}, len(stage.Approvers))
		approvers := make([]string, 0, len(stage.Approvers))
		for _, approver := range stage.Approvers {
			approver = strings.TrimSpace(approver)
			if approver == "" {
				continue
			}
			if _, ok := approverSet[approver]; ok {
				continue
			}
			approverSet[approver] = struct{}{}
			approvers = append(approvers, approver)
			userSet[approver] = struct{}{}
		}
		if len(approvers) == 0 {
			return nil, nil, fmt.Errorf("第%d个阶段审批人不能为空", idx+1)
		}
		stage.Approvers = approvers
	}

	for stage := 1; stage <= len(stages); stage++ {
		if _, ok := stageSet[stage]; !ok {
			return nil, nil, fmt.Errorf("审批阶段必须从1开始连续编号")
		}
	}

	sort.Slice(stages, func(i, j int) bool {
		return stages[i].Stage < stages[j].Stage
	})

	normalized, err := json.Marshal(stages)
	if err != nil {
		return nil, nil, fmt.Errorf("审批流定义序列化失败: %w", err)
	}

	users := make([]string, 0, len(userSet))
	for user := range userSet {
		users = append(users, user)
	}
	sort.Strings(users)

	return datatypes.JSON(normalized), users, nil
}

func validateActiveUsersExist(users []string) error {
	if len(users) == 0 {
		return nil
	}

	userSet := make(map[string]struct{}, len(users))
	normalizedUsers := make([]string, 0, len(users))
	for _, user := range users {
		user = strings.TrimSpace(user)
		if user == "" {
			continue
		}
		if _, ok := userSet[user]; ok {
			continue
		}
		userSet[user] = struct{}{}
		normalizedUsers = append(normalizedUsers, user)
	}
	if len(normalizedUsers) == 0 {
		return fmt.Errorf("审批流用户不能为空")
	}

	var existingUsers []string
	if err := global.App.DB.Table("insight_users").
		Where("username in ? and is_active=1", normalizedUsers).
		Pluck("username", &existingUsers).Error; err != nil {
		return err
	}

	exists := make(map[string]struct{}, len(existingUsers))
	for _, user := range existingUsers {
		exists[user] = struct{}{}
	}

	var missing []string
	for _, user := range normalizedUsers {
		if _, ok := exists[user]; !ok {
			missing = append(missing, user)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("用户不存在或未启用: %s", strings.Join(missing, ","))
	}

	return nil
}
