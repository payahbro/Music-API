package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Duration int64

func (d Duration) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d seconds", d)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

var ErrInvalidDurationFormat = errors.New("invalid duration format")

func (d *Duration) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidDurationFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	// ex: 123 seconds
	if len(parts) != 2 || parts[1] != "seconds" {
		return ErrInvalidDurationFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return ErrInvalidDurationFormat
	}

	*d = Duration(i)

	return nil
}
