package scripting

import (
	"fmt"
	"github.com/google/skylark"
	"github.com/google/skylark/resolve"
)

func init() {
	resolve.AllowFloat = true
}

func exists(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var callable skylark.Callable
	if err := skylark.UnpackArgs("catch", args, kwargs, "callee", &callable); err != nil {
		return nil, err
	}
	if _, err := skylark.Call(thread, callable, nil, nil); err != nil {
		return skylark.String(err.Error()), nil
	}
	return skylark.None, nil
}

func RunScript(filename string, arguments skylark.StringDict) (skylark.StringDict, error) {
	thread := &skylark.Thread{
		Load: loader,
	}
	return skylark.ExecFile(thread, filename, nil, arguments)
}

func loader(thread *skylark.Thread, module string) (skylark.StringDict, error) {
	switch module {
	case HttpModule:
		return skylark.StringDict{
			"http": LoadHttp(),
		}, nil
	case TimeModule:
		return skylark.StringDict{
			"time": LoadTime(),
		}, nil
	case JsonModule:
		return skylark.StringDict{
			"json": LoadJson(),
		}, nil
	case FileModule:
		return skylark.StringDict{
			"file": LoadFile(),
		}, nil
	case Base64Module:
		return skylark.StringDict{
			"base64": LoadBase64(),
		}, nil
	}

	return nil, fmt.Errorf("invalid module")
}
