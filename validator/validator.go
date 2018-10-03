package validator

import (
	"fmt"
	"github.com/gamelife1314/go-note/models"
	"regexp"
	"unicode/utf8"
)

func Length(field, input string, min, max int) (bool, string) {
	length := utf8.RuneCountInString(input)
	if min > length || length > max {
		return false, fmt.Sprintf("%s 长度必须介于 %d 和 %d 之间", field, min, max)
	}

	return true, ""
}

func Regexp(field, input, pattern string) (pass bool, msg string) {

	if match, err := regexp.MatchString(pattern, input); err != nil || match == false {
		return false, fmt.Sprintf("%s 不满足正则表达式：%s", field, pattern)
	}

	return true, ""
}

func Unique(field, table, value string) (pass bool, msg string) {
	var count uint
	models.Database.Table(table).Where(fmt.Sprintf("%s = ?", field), value).Count(&count)
	if count > 0 {
		return false, fmt.Sprintf("%s 的值：%s 已经存在了", field, value)
	}
	return true, ""
}

func Exists(field, table, value string) (pass bool, msg string) {
	var count uint
	models.Database.Table(table).Where(map[string]interface{}{field: value}).Count(&count)
	if count == 0 {
		return false, fmt.Sprintf("%s 的值：%s 不存在", field, value)
	}
	return true, ""
}
