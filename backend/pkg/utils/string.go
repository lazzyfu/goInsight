package utils

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// 允许的命名
var NamePattern = "[a-zA-Z][a-zA-Z0-9_]*"

const KeyJoinChar = "+"

// IsSubKey
func IsSubKey(keyA, keyB string) bool {
	var short, long string
	if len(keyA) < len(keyB) {
		short, long = keyA, keyB
	} else {
		short, long = keyB, keyA
	}
	shortSet := strings.Split(short, KeyJoinChar)
	longSet := strings.Split(long, KeyJoinChar)

	for i := 0; i < len(shortSet); i++ {
		if shortSet[i] != longSet[i] {
			return false
		}
	}
	return true
}

// IsByteContain
func IsByteContain(items []byte, item byte) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// HasPrefix 不区分大小写
func HasPrefix(s, prefix string, caseSensitive bool) bool {
	if caseSensitive {
		//  区分大小写
		return strings.HasPrefix(s, prefix)
	}
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

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

// 查重
func IsRepeat(items []string) (bool, []string) {
	repeat := make(map[string]bool)
	var repeatValue []string
	for _, item := range items {
		itemLower := strings.ToLower(item)
		if !repeat[itemLower] {
			repeat[itemLower] = true
		} else {
			repeatValue = append(repeatValue, item)
		}
	}
	if len(repeatValue) > 0 {
		return true, repeatValue
	}
	return false, repeatValue
}

// 获取最大值
func MaxInt(items []int) int {
	maxVal := items[0]
	for i := 1; i < len(items); i++ {
		//从第二个 元素开始循环比较，如果发现有更大的，则交换
		if maxVal < items[i] {
			maxVal = items[i]
		}
	}
	return maxVal
}

func RemoveElements(arr []string, toRemove []string) []string {
	// 创建一个map，用于快速查找要删除的元素
	removeMap := make(map[string]bool)
	for _, val := range toRemove {
		removeMap[val] = true
	}

	// 创建一个新的切片，用于存储不包含要删除元素的值
	result := []string{}

	// 遍历原始切片，将不在要删除元素中的值添加到新的切片中
	for _, val := range arr {
		if !removeMap[val] {
			result = append(result, val)
		}
	}

	return result
}

// 判断字符串是否匹配正则
func IsMatchPattern(pattern string, str string) bool {
	regCom := regexp.MustCompile(pattern)
	indices := regCom.FindAllStringIndex(str, -1)
	for _, indice := range indices {
		start, end := indice[0], indice[1]
		if unicode.IsDigit(rune(str[0])) {
			// 不能以数字开头
			return false
		} else {
			// 发现异常的字符
			if start != 0 || end != len(str) {
				return false
			}
		}
	}
	return true
}

func GenerateSimpleRandomString(length int) string {
	const charset = "0123456789ABCDEF"
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result)
}

func GenerateRandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()-_=+[]{};:,.<>?`~"
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result)
}
