package scripting

import (
	json2 "encoding/json"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
)

const JsonModule = "json.sky"

func LoadJson() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(
		skylarkstruct.Default,
		skylark.StringDict{
			"json":   skylark.NewBuiltin("json", json),
			"object": skylark.NewBuiltin("object", object),
		},
	)
}

func json(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var value skylark.Value
	if err := skylark.UnpackArgs("json", args, kwargs, "value", &value); err != nil {
		return skylark.None, err
	}
	native, err := Unmarshal(value)
	if nil != err {
		return skylark.None, err
	}
	bytes, err := json2.Marshal(native)
	if nil != err {
		return skylark.None, err
	}
	return skylark.String(bytes), nil
}

func object(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var content skylark.String
	err := skylark.UnpackArgs("add", args, kwargs, "content", &content)
	if nil != err {
		return skylark.None, err
	}
	var value interface{}
	err = json2.Unmarshal([]byte(AsString(content)), &value)
	if nil != err {
		return skylark.None, err
	}
	return Marshal(value)
}
