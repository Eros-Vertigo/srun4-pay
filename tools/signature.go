package tools

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// GenerateSignature 按key排序生成签名
func GenerateSignature(data interface{}) string {
	// 将结构体字段转换为键值对
	values := make(map[string]string)
	v := reflect.ValueOf(data)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		value := fmt.Sprintf("%v", field.Interface())
		values[t.Field(i).Name] = value
	}
	// 对 key 排序
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接键值对为字符串
	var query strings.Builder
	for _, key := range keys {
		value := values[key]
		query.WriteString(fmt.Sprintf("%s=%s", key, value))
	}
	return strings.TrimRight(query.String(), "&")
}

func Md5sum(str string) [md5.Size]byte {
	return md5.Sum([]byte(str))
}
