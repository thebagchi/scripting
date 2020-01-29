package scripting

import (
	"fmt"
	"github.com/google/skylark"
	"strconv"
)

func AsString(input skylark.Value) string {
	value, err := strconv.Unquote(input.String())
	if nil != err {
		return ""
	}
	return value
}

func Unmarshal(v skylark.Value) (interface{}, error) {
	switch v.Type() {
	case "NoneType":
		return nil, nil
	case "bool":
		return v.Truth() == skylark.True, nil
	case "int":
		return skylark.AsInt32(v)
	case "float":
		if float, ok := skylark.AsFloat(v); ok {
			return float, nil
		} else {
			return nil, fmt.Errorf("couldn't parse float")
		}
	case "string":
		return strconv.Unquote(v.String())
	case "dict":
		if dict, ok := v.(*skylark.Dict); ok {
			var values = map[string]interface{}{}
			for _, key := range dict.Keys() {
				value, _, err := dict.Get(key)
				if err != nil {
					return nil, err
				}
				temp, err := Unmarshal(value)
				if err != nil {
					return nil, err
				}
				values[AsString(key)] = temp
			}
			return values, nil
		} else {
			return nil, fmt.Errorf("error parsing dict. invalid type: %v", v)
		}
	case "list":
		if list, ok := v.(*skylark.List); ok {
			var element skylark.Value
			var iterator = list.Iterate()
			var value = make([]interface{}, 0)
			for iterator.Next(&element) {
				temp, err := Unmarshal(element)
				if err != nil {
					return nil, err
				}
				value = append(value, temp)
			}
			iterator.Done()
			return value, nil
		} else {
			return nil, fmt.Errorf("error parsing list. invalid type: %v", v)
		}
	case "tuple":
		if tuple, ok := v.(skylark.Tuple); ok {
			var element skylark.Value
			var iterator = tuple.Iterate()
			var value = make([]interface{}, 0)
			for iterator.Next(&element) {
				temp, err := Unmarshal(element)
				if err != nil {
					return nil, err
				}
				value = append(value, temp)
			}
			iterator.Done()
			return value, nil
		} else {
			return nil, fmt.Errorf("error parsing dict. invalid type: %v", v)
		}
	case "set":
		return nil, fmt.Errorf("sets aren't yet supported")
	default:
		return nil, fmt.Errorf("unrecognized skylark type: %s", v.Type())
	}
}

func Marshal(v interface{}) (skylark.Value, error) {
	switch x := v.(type) {
	case nil:
		return skylark.None, nil
	case bool:
		return skylark.Bool(x), nil
	case string:
		return skylark.String(x), nil
	case int:
		return skylark.MakeInt(x), nil
	case float64:
		return skylark.Float(x), nil
	case []interface{}:
		var elements = make([]skylark.Value, 0)
		for _, value := range x {
			element, err := Marshal(value)
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
		}
		return skylark.NewList(elements), nil
	case map[string]interface{}:
		dict := &skylark.Dict{}
		for key, value := range x {
			element, err := Marshal(value)
			if err != nil {
				return nil, err
			}
			if err = dict.Set(skylark.String(key), element); err != nil {
				return nil, err
			}
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("unknown type %T", v)
	}
}
