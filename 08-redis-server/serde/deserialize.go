package serde

import (
	"strconv"
	"strings"
)

func Deserialize(str string) (Command, error) {
	if str == "" || len(str) < 5 || !strings.HasSuffix(str, End) {
		return Command{}, ErrInvalidString
	}

	str = str[:len(str)-4] // remove the end (\\r\\n)
	switch str[0] {
	case '-':
		return errorToCmd(str[1:])

	case '+':
		return simpleStringToCmd(str[1:])

	case '$':
		return bulkStringToCmd(str[1:])

	default:
		return Command{}, ErrInvalidCommandType
	}
}

func errorToCmd(str string) (Command, error) {
	parts := strings.Split(str, " ")
	if len(parts) < 2 || len(parts[0]) < 2 {
		return Command{}, ErrInvalidCommandError
	}

	t := strings.ToUpper(parts[0])
	m := strings.TrimSpace(str[len(t):])

	err := CommandError{
		Type:    t,
		Message: m,
	}

	cmd := Command{
		Type:  Error,
		Error: &err,
	}

	return cmd, nil
}

func simpleStringToCmd(str string) (Command, error) {
	cmd := Command{
		Type:  SimpleString,
		Value: str,
	}

	return cmd, nil
}

func bulkStringToCmd(str string) (Command, error) {
	if str == "0"+End {
		cmd := Command{
			Type:  BulkString,
			Value: "",
		}

		return cmd, nil
	}

	pos := strings.Index(str, End)
	if pos == -1 {
		return Command{}, ErrInvalidString
	}

	size, err := strconv.Atoi(str[:pos-1])
	if err != nil {
		return Command{}, ErrInvalidString
	}

	str = str[pos+len(End):]
	if len(str) != size {
		return Command{}, ErrInvalidString
	}

	cmd := Command{
		Type:  BulkString,
		Value: str,
	}

	return cmd, nil
}
