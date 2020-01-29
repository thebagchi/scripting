package scripting

import (
	b64 "encoding/base64"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
)

const Base64Module = "base64.sky"

func LoadBase64() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(
		skylarkstruct.Default,
		skylark.StringDict{
			"encode": skylark.NewBuiltin("encode", encode),
		},
	)
}

func encode(
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
	data := b64.StdEncoding.EncodeToString([]byte(AsString(content)))
	return skylark.String(data), nil
}
