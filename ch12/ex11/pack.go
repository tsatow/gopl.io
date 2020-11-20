package ex11

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Pack(u *url.URL, data interface{}) error {
	v := reflect.ValueOf(data)
	params := make([]string, 0)

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)

		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}

		switch v.Kind() {
		case reflect.String:
			params = append(params, fmt.Sprintf("%s=%s", name, url.QueryEscape(v.String())))
		case reflect.Int:
			params = append(params, fmt.Sprintf("%s=%s", name, url.QueryEscape(strconv.FormatInt(v.Int(), 10))))
		case reflect.Bool:
			params = append(params, fmt.Sprintf("%s=%s", name, url.QueryEscape(strconv.FormatBool(v.Bool()))))
		default:
			return fmt.Errorf("unsupported kind %s", v.Type())
		}
	}
	query := strings.Join(params, "&")
	u.RawQuery = query
	return nil
}
