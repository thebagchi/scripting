package scripting

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
	"strconv"
	"time"
)

const TimeModule = "time.sky"

func LoadTime() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(
		skylarkstruct.Default,
		skylark.StringDict{
			"timestamp": skylark.NewBuiltin("timestamp", timestamp),
			"formatted": skylark.NewBuiltin("formatted", formatted),
			"format":    skylark.NewBuiltin("format", format),
			"parse":     skylark.NewBuiltin("format", parse),
		},
	)
}

func timestamp(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	now := time.Now().Unix()
	return skylark.MakeInt64(now), nil
}

func formatted(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	return skylark.String(now), nil
}

func format(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var format, timestamp skylark.String
	err := skylark.UnpackArgs("format", args, kwargs, "format", &format, "timestamp", &timestamp)
	if nil != err {
		return skylark.None, err
	}
	seconds, err := strconv.ParseInt(AsString(timestamp), 10, 64)
	if err != nil {
		return skylark.None, err
	}
	value := time.Unix(seconds, 0).Format(AsString(format))
	return skylark.String(value), nil
}

func parse(thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var format, dateTime skylark.String
	err := skylark.UnpackArgs("parse", args, kwargs, "format", &format, "dateTime", &dateTime)
	if nil != err {
		return skylark.None, err
	}
	timestamp, err := time.Parse(AsString(format), AsString(dateTime))
	if nil != err {
		glog.Error("Error: ", err)
		return skylark.String(fmt.Sprintf("%d", timestamp.Unix())), nil
	}
	return skylark.String(fmt.Sprintf("%d", time.Now().Unix())), nil
}
