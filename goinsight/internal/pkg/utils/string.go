package utils

import "strings"

// IsContain等值比较，忽略大小写
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if strings.EqualFold(eachItem, item) {
			return true
		}
	}
	return false
}

func ErrsJoin(str string, err []error) string {
	if len(err) < 1 {
		return ""
	}
	result := ""
	for i, v := range err {
		if v == nil {
			continue
		}

		if i == 0 {
			result += v.Error()
			continue
		}
		result += v.Error() + str
	}
	return result
}

// 去重
func RemoveDuplicate(s []string) []string {
	result := []string{}
	temp := map[string]struct{}{}
	for _, item := range s {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
