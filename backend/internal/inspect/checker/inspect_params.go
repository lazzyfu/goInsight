package checker

import (
	"fmt"
	"strconv"
	"strings"
)

// parseInspectParamTypedValue 解析从 InsightGlobalInspectParams / InsightInstanceInspectParams 读出的 Value。
// rawType 仅允许：number / boolean / string（与 initializeInspectParams 的 Type 一致）。
// string 类型目前只有两种“写法”：
// 1) TABLE_SUPPORT_CHARSET: "utf8,utf8_general_ci;utf8mb4,utf8mb4_general_ci" -> []map[string]string
// 2) TABLE_SUPPORT_ENGINE（以及同样 CSV 写法的 INNODB_ROW_FORMAT）: "InnoDB,MyISAM" -> []string
func parseInspectParamTypedValue(key, rawValue, rawType string) (any, error) {
	value := strings.TrimSpace(rawValue)
	typeLower := strings.ToLower(strings.TrimSpace(rawType))

	switch key {
	case "TABLE_SUPPORT_ENGINE", "INNODB_ROW_FORMAT":
		if typeLower != "" && typeLower != "string" {
			return nil, fmt.Errorf("invalid type for %s: %s", key, rawType)
		}
		return splitCSV(value), nil
	case "TABLE_SUPPORT_CHARSET":
		if typeLower != "" && typeLower != "string" {
			return nil, fmt.Errorf("invalid type for %s: %s", key, rawType)
		}
		return parseTableSupportCharsetValue(value)
	}

	switch typeLower {
	case "number":
		if value == "" {
			return 0, nil
		}
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return n, nil
	case "boolean":
		v := strings.ToLower(value)
		switch v {
		case "true", "1":
			return true, nil
		case "false", "0", "":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid boolean: %s", value)
		}
	case "string", "":
		return value, nil
	default:
		return value, nil
	}
}

// splitCSV 将逗号分隔字符串拆成 []string（自动 trim，忽略空项）。
func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	items := strings.Split(value, ",")
	res := make([]string, 0, len(items))
	for _, it := range items {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		res = append(res, it)
	}
	return res
}

// 兼容 initializeInspectParams 默认格式："utf8,utf8_general_ci;utf8mb4,utf8mb4_general_ci"
func parseTableSupportCharsetValue(value string) ([]map[string]string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	entries := splitBy(value, ";")
	res := make([]map[string]string, 0, len(entries))
	for _, e := range entries {
		parts := splitBy(e, ",")
		if len(parts) == 0 {
			continue
		}
		m := map[string]string{"charset": parts[0]}
		if len(parts) > 1 {
			m["recommend"] = parts[1]
		} else {
			m["recommend"] = ""
		}
		res = append(res, m)
	}
	return res, nil
}

func splitBy(value, sep string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	items := strings.Split(value, sep)
	res := make([]string, 0, len(items))
	for _, it := range items {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		res = append(res, it)
	}
	return res
}
