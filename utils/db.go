package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func FieldToWhere(field reflect.Value, columnTag string, patternSearch bool) string {
	var fieldValue any
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return ""
		}
		fieldValue = field.Elem().Interface()
	} else {
		fieldValue = field.Interface()
	}
	switch v := fieldValue.(type) {
	case string:
		if v == "" {
			return ""
		}
		if patternSearch {
			return columnTag + " LIKE '%" + v + "%'"
		}
		return columnTag + " = '" + v + "'"
	case int, int64, float64:
		return columnTag + " = " + reflect.ValueOf(v).String()
	case bool:
		if v {
			return columnTag + " = true"
		}
		return columnTag + " = false"
	case []string:
		if len(v) == 0 {
			return ""
		}
		if patternSearch {
			likeClauses := make([]string, len(v))
			for i, val := range v {
				likeClauses[i] = columnTag + " LIKE '%" + val + "%'"
			}
			return "(" + strings.Join(likeClauses, " OR ") + ")"
		}
		return columnTag + " IN (" + strings.Join(v, ", ") + ")"
	case []int, []int64, []float64:
		if reflect.ValueOf(v).Len() == 0 {
			return ""
		}
		return columnTag + " IN (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), ","), "[]") + ")"
	default:
		return columnTag + " = '" + reflect.ValueOf(v).String() + "'"
	}
}

func ConvertToWhere(req any) string {
	reqValue := reflect.ValueOf(req)
	reqType := reflect.TypeOf(req)
	whereClauses := make([]string, 0)

	if reqValue.Kind() == reflect.Ptr {
		reqValue = reqValue.Elem()
		reqType = reqType.Elem()
	}
	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		typeField := reqType.Field(i)

		columnTag := typeField.Tag.Get("column")
		if columnTag == "" || columnTag == "-" {
			continue
		}
		where := FieldToWhere(field, columnTag, typeField.Tag.Get("pattern_search") == "true")
		if where != "" {
			whereClauses = append(whereClauses, where)
		}
	}
	if len(whereClauses) > 0 {
		return strings.Join(whereClauses, " AND ")
	}
	return ""
}
