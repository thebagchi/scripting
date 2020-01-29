package scripting

import (
	"bytes"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
	"os"
)

const FileModule = "file.sky"

func Append(values ...string) string {
	var buffer bytes.Buffer
	for _, value := range values {
		buffer.WriteString(value)
	}
	return buffer.String()
}

func AppendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

func LoadFile() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(
		skylarkstruct.Default,
		skylark.StringDict{
			"add": skylark.NewBuiltin("add", add),
		},
	)
}

func add(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var filename, content skylark.String
	err := skylark.UnpackArgs("add", args, kwargs, "filename", &filename, "content", &content)
	if nil != err {
		return skylark.None, err
	}
	AppendStringToFile(AsString(filename), Append(AsString(content), "\n"))
	return skylark.None, nil
}
