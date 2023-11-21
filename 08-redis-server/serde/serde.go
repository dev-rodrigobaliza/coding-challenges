package serde

import "errors"

type CommandType int

const (
	Unknown CommandType = iota
	SimpleString
	BulkString
	Integer
	Array
	Error

	End = "\\r\\n"
)

var (
	ErrInvalidCommandType  = errors.New("invalid command type")
	ErrInvalidCommandValue = errors.New("invalid command value")
	ErrInvalidCommandError = errors.New("invalid command error")
	ErrInvalidArray        = errors.New("invalid array")
	ErrInvalidString       = errors.New("invalid string to deserialize")
)

type CommandArray struct {
	Type  CommandType
	Value string
	Array []CommandArray
}

type CommandError struct {
	Type    string
	Message string
}

type Command struct {
	Type  CommandType
	Value string
	Array []CommandArray
	Error *CommandError
}
