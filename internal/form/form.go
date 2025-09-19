// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"mime"
	"mime/multipart"
	"reflect"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/unknwon/com"
	"github.com/wuhan005/govalid"

	"github.com/wuhan005/NekoBox/internal/context"
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

		contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		switch contentType {
		case "application/x-www-form-urlencoded", "":
			if err := r.ParseForm(); err != nil {
				c.Map(Error{Category: ErrorCategoryDeserialization, Error: err})
				return
			}
		case "multipart/form-data":
			if err := r.ParseMultipartForm(10 * 1 << 20); err != nil { // 10 MiB
				c.Map(Error{Category: ErrorCategoryDeserialization, Error: err})
				return
			}
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

			if r.MultipartForm != nil {
				fhType := reflect.TypeOf((*multipart.FileHeader)(nil))

				value, ok := r.MultipartForm.Value[fieldName]
				if ok {
					if val.Field(i).Kind() == reflect.Slice && val.Field(i).Type().Elem() == fhType {
						continue
					}

					// FIXME: We don't implement the slice type yet.
					val.Field(i).Set(reflect.ValueOf(value[0]))
				} else {
					value, ok := r.MultipartForm.File[fieldName]
					if ok {
						numElems := len(value)

						if val.Field(i).Kind() == reflect.Slice && numElems > 0 && val.Field(i).Type().Elem() == fhType {
							slice := reflect.MakeSlice(val.Field(i).Type(), numElems, numElems)
							for i := 0; i < numElems; i++ {
								slice.Index(i).Set(reflect.ValueOf(value[i]))
							}
							val.Field(i).Set(slice)
						} else if val.Field(i).Type() == fhType {
							val.Field(i).Set(reflect.ValueOf(value[0]))
						}
					}
				}
			} else {
				val.Field(i).Set(reflect.ValueOf(r.Form.Get(fieldName)))
			}
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
