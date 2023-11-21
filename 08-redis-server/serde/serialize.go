package serde

import (
	"fmt"
	"strconv"
	"strings"
)

func Serialize(cmd Command) (string, error) {
	return serialize(cmd.Type, cmd.Value, cmd.Array, cmd.Error)
}

func serialize(ct CommandType, value string, array []CommandArray, err *CommandError) (string, error) {
	switch ct {
	case Error:
		return cmdToError(err)

	case SimpleString:
		return cmdToSimpleString(value)

	case BulkString:
		return cmdToBulkString(value)

	case Integer:
		return cmdToInteger(value)

	case Array:
		return cmdToArray(array)

	default:
		return "", ErrInvalidCommandType
	}
}

func cmdToError(err *CommandError) (string, error) {
	if err == nil || err.Type == "" || err.Message == "" {
		return "", ErrInvalidCommandError
	}

	err.Type = strings.ToUpper(err.Type)
	str := fmt.Sprintf("-%s %s%s", err.Type, err.Message, End)
	return str, nil
}

func cmdToSimpleString(value string) (string, error) {
	if value == "" {
		return fmt.Sprintf("+%s", End), nil
	}

	str := fmt.Sprintf("+%s%s", value, End)
	return str, nil
}

func cmdToBulkString(value string) (string, error) {
	if value == "" {
		return fmt.Sprintf("$0%s%s", End, End), nil
	}

	str := fmt.Sprintf("$%d%s%s%s", len(value), End, value, End)
	return str, nil
}

func cmdToInteger(value string) (string, error) {
	if _, err := strconv.ParseInt(value, 10, 64); err != nil {
		return "", ErrInvalidCommandValue
	}

	str := fmt.Sprintf(":%s%s", value, End)
	return str, nil
}

func cmdToArray(array []CommandArray) (string, error) {
	if len(array) == 0 {
		return fmt.Sprintf("*0%s", End), nil
	}

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("*%d%s", len(array), End))

	for _, v := range array {
		str, err := serialize(v.Type, v.Value, v.Array, nil)
		if err != nil {
			return "", err
		}

		s.WriteString(str)
	}

	return s.String(), nil
}
