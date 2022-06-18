// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"reflect"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/unknwon/com"
	"github.com/wuhan005/govalid"

	"github.com/NekoWheel/NekoBox/internal/context"
)

type ErrorCategory string

const (
	ErrorCategoryDeserialization ErrorCategory = "deserialization"
	ErrorCategoryValidation      ErrorCategory = "validation"
)

type Error struct {
	Category ErrorCategory
	Error    error
}

func Bind(model interface{}) flamego.Handler {
	// Ensure not pointer.
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		panic("form: pointer can not be accepted as binding model")
	}

	return func(c context.Context, data template.Data) {
		obj := reflect.New(reflect.TypeOf(model))
		defer func() { c.Map(obj.Elem().Interface()) }()

		r := c.Request()
		if err := r.ParseForm(); err != nil {
			c.Map(Error{Category: ErrorCategoryDeserialization, Error: err})
			return
		}

		// Bind the form data to the given struct.
		typ := reflect.TypeOf(obj.Interface())
		val := reflect.ValueOf(obj.Interface())
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			val = val.Elem()
		}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldName := typ.Field(i).Tag.Get("form")
			if fieldName == "" {
				fieldName = com.ToSnakeCase(field.Name)
			}
			val.Field(i).Set(reflect.ValueOf(r.Form.Get(fieldName)))
		}

		errors, ok := govalid.Check(obj.Interface())
		if !ok {
			Assign(obj.Interface(), data)

			c.SetError(errors[0])
			c.Map(Error{Category: ErrorCategoryValidation, Error: errors[0]})
			return
		}
	}
}

// Assign assigns form values back to the template data.
func Assign(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := com.ToSnakeCase(field.Name)

		data[fieldName] = val.Field(i).Interface()
	}
}
