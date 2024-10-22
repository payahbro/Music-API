package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

// Paginatings create pagination
type Paginatings struct {
	Page     int64
	PageSize int64
}

func (p Paginatings) Limit() int {
	return int(p.PageSize)
}

func (p Paginatings) Offset() int {
	return (int(p.Page) - 1) * int(p.PageSize)
}

func (p Paginatings) Validate(validate *validator.Validate) error {
	var err error

	err = validate.Var(p.Page, "min=1")
	if err != nil {
		return errors.New("page must be greater than 0")
	}

	err = validate.Var(p.PageSize, "min=1")
	if err != nil {
		return errors.New("page size must be greater than 0")
	}

	err = validate.Var(p.PageSize, "max=50")
	if err != nil {
		return errors.New("page size has maximum of 50")
	}

	return nil
}

// Sortings sort the given column from query parameter
type Sortings struct {
	Sorts         string
	SafeSortLists []string
}

func (s Sortings) SortName() string {
	for _, v := range s.SafeSortLists {
		if s.Sorts == v {
			return strings.TrimPrefix(s.Sorts, "-")
		}
	}
	panic("unsafe sort")
}

func (s Sortings) SortDirection() string {
	if strings.HasPrefix(s.Sorts, "-") {
		return "DESC"
	}
	return "ASC"
}

// Background goroutines for sent email
func Background(fn func()) {
	// use goroutines to immediately return JSON Response without waiting email to be sent
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()

		fn()
	}()
}

// ReadStrQuery read query parameter return as String
func ReadStrQuery(e echo.Context, key string, def string) string {
	var s string
	s = e.QueryParam(key)

	if s == "" {
		return def
	}

	return s
}

// ReadIntQuery read query parameter return as Integer
func ReadIntQuery(e echo.Context, key string, def int64) int64 {
	var s = e.QueryParam(key)
	if s == "" {
		return def
	}

	var i, err = strconv.Atoi(s)
	if err != nil {
		return def
	}

	return int64(i)
}

// ReadCSVQuery read query parameter return as CSV
func ReadCSVQuery(e echo.Context, key string, def []string) []string {
	var csv = e.QueryParam(key)
	if csv == "" {
		return def
	}

	return strings.Split(csv, ",")

}

// ReadIdParam for read ID in path parameter
func ReadIdParam(e echo.Context) (int64, error) {
	params := e.Param("tracksId")
	id, err := strconv.ParseInt(params, 10, 64)

	if err != nil || id == 0 {
		return 0, err
	}

	return id, nil
}

// ReadJSON for read request body
func ReadJSON(e echo.Context, dst any) error {
	err := e.Bind(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
			// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
			// for syntax errors in the JSON. So we check for this using errors.Is() and
			// return a generic error message. There is an open issue regarding this at
			// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
			// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
			// JSON value is the wrong type for the target destination. If the error relates
			// to a specific field, then we include that in our error message to make it
			// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
			// An io.EOF error will be returned by Decode() if the request body is empty. We
			// check for this with errors.Is() and return a plain-english error message
			// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
			// A json.InvalidUnmarshalError error will be returned if we pass something
			// that is not a non-nil pointer to Decode(). We catch this and panic,
			// rather than returning an error to our handler. At the end of this chapter
			// we'll talk about panicking versus returning errors, and discuss why it's an
			// appropriate thing to do in this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
			// For anything else
		default:
			return err
		}
	}
	return nil
}
