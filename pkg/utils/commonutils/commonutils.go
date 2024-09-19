package commonutils

import (
	"fmt"
	"reflect"
)

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 如果字段是可导出的（首字母大写）
		if field.PkgPath == "" {
			// 处理特殊类型
			switch fieldValue.Kind() {
			case reflect.Struct:
				// 递归处理嵌套结构体
				result[field.Name] = StructToMap(fieldValue.Interface())
			case reflect.Ptr:
				// 处理指针
				if !fieldValue.IsNil() {
					result[field.Name] = StructToMap(fieldValue.Elem().Interface())
				} else {
					result[field.Name] = nil
				}
			default:
				// 对于其他类型，直接获取其值
				result[field.Name] = fieldValue.Interface()
			}
		}
	}

	return result
}

func StructToStringMap(obj interface{}) map[string]string {
	result := make(map[string]string)
	v := reflect.ValueOf(obj)

	// 如果传入的是指针，获取其指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 确保我们处理的是结构体
	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 只处理可导出的字段
		if field.PkgPath == "" {
			switch fieldValue.Kind() {
			case reflect.Struct:
				nestedMap := StructToStringMap(fieldValue.Interface())
				for nestedKey, nestedValue := range nestedMap {
					result[field.Name+"."+nestedKey] = nestedValue
				}
			case reflect.Ptr:
				if !fieldValue.IsNil() {
					nestedMap := StructToStringMap(fieldValue.Elem().Interface())
					for nestedKey, nestedValue := range nestedMap {
						result[field.Name+"."+nestedKey] = nestedValue
					}
				} else {
					result[field.Name] = "nil"
				}
			case reflect.Slice, reflect.Array:
				result[field.Name] = fmt.Sprintf("%v", fieldValue.Interface())
			default:
				result[field.Name] = fmt.Sprintf("%v", fieldValue.Interface())
			}
		}
	}

	return result
}
