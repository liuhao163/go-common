package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidateError struct {
	fieldName   string
	fieldReason []string
	fieldValue  interface{}
}

func (v ValidateError) Error() string {
	return fmt.Sprintf("verify failed.filed=%v,val=%v", v.fieldName, v.fieldValue)
}

/*
  true,nil 校验成功
  false 失败，ValidateError是失败的字段以及原因
*/
func Validate(input interface{}) (bool, []ValidateError) {
	typeOf := reflect.TypeOf(input)
	valueOf := reflect.ValueOf(input)

	ret := true
	errors := make([]ValidateError, 0)
	for i := 0; i < typeOf.NumField(); i++ {
		f := typeOf.Field(i)
		v := valueOf.Field(i)

		reason := make([]string, 0)
		if _, ok := f.Tag.Lookup("verify-nonempty"); ok {
			ret = ret && verifyRequirement(v)
			if !ret {
				reason = append(reason, "verify-nonempty")
			}
		}
		if r, ok := f.Tag.Lookup("verify-range"); ok {
			r := strings.Split(r, "-")
			ret = ret && len(r) > 0 && verifyRange(v, r)
			if !ret {
				reason = append(reason, "verify-range")
			}
		}
		if len(reason) > 0 {
			errors = append(errors, ValidateError{fieldName: f.Name, fieldReason: reason})
		}
	}
	return ret, errors
}

func verifyRequirement(v reflect.Value) bool {
	ret := false
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		ret = v.Len() > 0
		break
	case reflect.String:
		ret = v.Len() > 0
		break
	case reflect.Ptr:
		ret = !v.IsNil()
		break
	default:
	}
	return ret
}

func verifyRange(v reflect.Value, r []string) bool {
	//get val
	var val float64
	var err error
	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		val = float64(v.Int())
		break
	case reflect.Float32, reflect.Float64:
		val = v.Float()
		break
	case reflect.String:
		val, err = strconv.ParseFloat(v.String(), 10)
		if err != nil {
			return false
		}
		break
	case reflect.Slice, reflect.Array:
		val = float64(v.Len())
		break
	default:
		return false
	}

	//校验
	ret := true
	if len(r[0]) > 0 {
		min, err := strconv.ParseFloat(r[0], 10)
		ret = err == nil && val >= min
		if !ret {
			return ret
		}
	}

	if len(r) > 1 {
		if len(r[1]) > 0 {
			max, err := strconv.ParseFloat(r[1], 10)
			ret = ret && err == nil && val <= max
		}
		if !ret {
			return ret
		}
	}

	return ret
}
