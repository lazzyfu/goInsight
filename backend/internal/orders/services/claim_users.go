package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/pkg/utils"
	"gorm.io/datatypes"
)

func normalizeClaimUsers(users []string) ([]string, error) {
	if len(users) == 0 {
		return nil, fmt.Errorf("可领取人不能为空")
	}

	seen := make(map[string]struct{}, len(users))
	result := make([]string, 0, len(users))
	for _, user := range users {
		u := strings.TrimSpace(user)
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		result = append(result, u)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("可领取人不能为空")
	}
	return result, nil
}

func marshalClaimUsers(users []string) (datatypes.JSON, error) {
	normalized, err := normalizeClaimUsers(users)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(normalized)
	if err != nil {
		return nil, fmt.Errorf("序列化可领取人失败: %w", err)
	}
	return datatypes.JSON(data), nil
}

func parseClaimUsers(raw datatypes.JSON) ([]string, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("审批流未配置可领取人")
	}

	var users []string
	if err := json.Unmarshal(raw, &users); err != nil {
		return nil, fmt.Errorf("解析可领取人失败: %w", err)
	}

	normalized, err := normalizeClaimUsers(users)
	if err != nil {
		return nil, err
	}
	return normalized, nil
}

func canUserClaim(raw datatypes.JSON, username string) (bool, error) {
	claimUsers, err := parseClaimUsers(raw)
	if err != nil {
		return false, err
	}
	return utils.IsContain(claimUsers, username), nil
}
