package serde

import (
	"strconv"
	"strings"
)

func Deserialize(str string) (Command, error) {
	return deserialize(str)
}

func deserialize(str string) (Command, error) {
	if str == "" || len(str) < len(End)+1 || !strings.HasSuffix(str, End) {
		return Command{}, ErrInvalidString
	}

	switch str[0] {
	case '-':
		return errorToCmd(str[1:])

	case '+':
		return simpleStringToCmd(str[1:])

	case '$':
		return bulkStringToCmd(str[1:])

	case ':':
		return integerToCmd(str[1:])

	case '*':
		return arrayToCmd(str[1:])

	default:
		return Command{}, ErrInvalidCommandType
	}
}

func errorToCmd(str string) (Command, error) {
	var err error
	str, err = removeEnd(str)
	if err != nil {
		return Command{}, err
	}

	parts := strings.Split(str, " ")
	if len(parts) < 2 || len(parts[0]) < 2 {
		return Command{}, ErrInvalidCommandError
	}

	t := strings.ToUpper(parts[0])
	m := strings.TrimSpace(str[len(t):])

	ce := CommandError{
		Type:    t,
		Message: m,
	}

	cmd := Command{
		Type:  Error,
		Error: &ce,
	}

	return cmd, nil
}

func simpleStringToCmd(str string) (Command, error) {
	var err error
	str, err = removeEnd(str)
	if err != nil {
		return Command{}, err
	}

	cmd := Command{
		Type:  SimpleString,
		Value: str,
	}

	return cmd, nil
}

func bulkStringToCmd(str string) (Command, error) {
	var err error
	str, err = removeEnd(str)
	if err != nil {
		return Command{}, err
	}

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

	size, err := strconv.Atoi(str[:pos])
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

func integerToCmd(str string) (Command, error) {
	var err error
	str, err = removeEnd(str)
	if err != nil {
		return Command{}, err
	}

	_, err = strconv.Atoi(str)
	if err != nil {
		return Command{}, ErrInvalidString
	}

	cmd := Command{
		Type:  Integer,
		Value: str,
	}

	return cmd, nil
}

func arrayToCmd(str string) (Command, error) {
	size, str, err := arraySize(str)
	if err != nil {
		return Command{}, err
	}

	var items []string

	if size > 0 {
		items = strings.Split(str, End)
		if items[0] == "" {
			return Command{}, ErrInvalidString
		}

		if items[len(items)-1] == "" {
			items = items[:len(items)-1]
		}

		for i, item := range items {
			items[i] = item + End
		}
	}

	items, err = arrangeArrayItems(items)
	if err != nil {
		return Command{}, err
	}

	return arrayItemsToCmd(size, items, Command{Type: Array})
}

func arrayItemsToCmd(size int, items []string, cmd Command) (Command, error) {
	var (
		c   Command
		err error

	)

	if size != len(items) {
		return Command{}, ErrInvalidString
	}

	if size == 0 {
		return cmd, nil
	}

	item := items[0]

	if item[0] == '*' {
		c, err = arrayToCmd(item[1:])
	} else {
		c, err = deserialize(item)
	}
	if err != nil {
		return Command{}, err
	}

	ca := CommandArray{
		Type:  c.Type,
		Value: c.Value,
		Array: c.Array,
	}

	size--
	items = items[1:]

	cmd.Array = append(cmd.Array, ca)
	return arrayItemsToCmd(size, items, cmd)
}

func arraySize(str string) (int, string, error) {
	if str == "" {
		return 0, "", ErrInvalidString
	}

	var (
		size int
		err  error
	)

	pos := strings.Index(str, End)
	if pos == -1 {
		size, err = strconv.Atoi(str)
		str = ""
	} else {
		size, err = strconv.Atoi(str[:pos])
		str = str[pos+len(End):]
	}
	if err != nil {
		return 0, "", ErrInvalidString
	}

	if size == 0 && str != "" || size > 0 && str == "" {
		return 0, "", ErrInvalidString
	}

	return size, str, nil
}

func arrangeArrayItems(items []string) ([]string, error) {
	copy := []string{}
	pos := 0
	for i, item := range items {
		if i != pos {
			continue
		}

		switch item[0] {
		case '$':
			var builder strings.Builder
			builder.WriteString(item)
			builder.WriteString(items[i+1])

			copy = append(copy, builder.String())
			pos++

		case '*':
			if i == len(items)-1 {
				return nil, ErrInvalidString
			}

			size, _, err := arraySize(item[1:] + items[i+1])
			if err != nil {
				return nil, err
			}

			var builder strings.Builder
			builder.WriteString(item)
			for j := 1; j <= size; j++ {
				if j > len(items) {
					return nil, ErrInvalidString
				}

				builder.WriteString(items[i+j])
			}

			copy = append(copy, builder.String())
			pos += size

		default:
			copy = append(copy, item)
		}

		pos++
	}

	return copy, nil
}

func removeEnd(str string) (string, error) {
	if len(str) < len(End) {
		return "", ErrInvalidString
	}

	str = str[:len(str)-len(End)]
	return str, nil
}
