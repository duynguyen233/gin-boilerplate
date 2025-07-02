package utils

import (
	"reflect"
)

// ModelDump converts a struct to a map[string]any, excluding unset fields if excludeUnset is true.
func ModelDump(req any, excludeUnset bool) map[string]any {
	reqValue := reflect.ValueOf(req)
	reqType := reflect.TypeOf(req)

	if reqValue.Kind() == reflect.Ptr && reqValue.IsValid() {
		reqValue = reqValue.Elem()
		reqType = reqType.Elem()
	}

	result := make(map[string]any)

	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		typeField := reqType.Field(i)

		jsonTag := typeField.Tag.Get("json")
		if columnTag := typeField.Tag.Get("column"); columnTag != "" {
			jsonTag = columnTag
		}
		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		// Check if it's schema.Optional[T]
		if method := field.MethodByName("IsOptional"); method.IsValid() && excludeUnset {
			if field.MethodByName("IsNullDefined").Call(nil)[0].Bool() {
				result[jsonTag] = nil
			} else if field.MethodByName("HasValue").Call(nil)[0].Bool() {
				result[jsonTag] = field.FieldByName("Value").Elem().Interface()
			}
			continue
		}

		result[jsonTag] = field.Interface()
	}

	return result
}
