package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	defaultValidator = validator.New()
	defaultMu        sync.RWMutex
)

type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value,omitempty"`
}

func init() {
	defaultValidator.RegisterTagNameFunc(jsonTagName)
}

func V() *validator.Validate {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultValidator
}

func SetDefault(v *validator.Validate) {
	if v == nil {
		v = validator.New()
		v.RegisterTagNameFunc(jsonTagName)
	}
	defaultMu.Lock()
	defaultValidator = v
	defaultMu.Unlock()
}

func Struct(value any) error {
	return V().Struct(value)
}

func Bind(c *gin.Context, out any) error {
	if out == nil {
		return errors.New("bind target cannot be nil")
	}
	if err := c.ShouldBind(out); err != nil {
		return err
	}
	return Struct(out)
}

func Format(err error) []FieldError {
	if err == nil {
		return nil
	}
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []FieldError{{Field: "", Tag: err.Error()}}
	}
	fields := make([]FieldError, 0, len(validationErrors))
	for _, item := range validationErrors {
		field := item.Field()
		fields = append(fields, FieldError{
			Field: field,
			Tag:   item.Tag(),
			Value: item.Value(),
		})
	}
	return fields
}

func ErrorMessage(err error) string {
	fields := Format(err)
	if len(fields) == 0 {
		return ""
	}
	if fields[0].Field == "" {
		return fields[0].Tag
	}
	return fmt.Sprintf("%s validation failed on %s", fields[0].Field, fields[0].Tag)
}

func jsonTagName(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}
