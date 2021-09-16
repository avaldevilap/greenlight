package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	*r = Runtime(i)
	return nil
}

func (runtime Runtime) GormDataType() string {
	return "int32"
}

// func (runtime *Runtime) Scan(value interface{}) error {
// 	str := fmt.Sprintf("%v", value)
// 	fmt.Println(str)
// 	unquotedJSONValue, err := strconv.Unquote(string(str))
// 	if err != nil {
// 		return ErrInvalidRuntimeFormat
// 	}
// 	parts := strings.Split(unquotedJSONValue, " ")
// 	if len(parts) != 2 || parts[1] != "mins" {
// 		return ErrInvalidRuntimeFormat
// 	}
// 	i, err := strconv.ParseInt(parts[0], 10, 32)
// 	if err != nil {
// 		return ErrInvalidRuntimeFormat
// 	}
// 	*runtime = Runtime(i)
// 	return nil
// }
// func (runtime Runtime) Value() (driver.Value, error) {
// 	return int64(runtime), nil
// }
