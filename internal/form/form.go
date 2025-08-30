// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/flamego/flamego"
	"github.com/wuhan005/govalid"
	"golang.org/x/text/language"

	"github.com/NekoWheel/NekoBox/internal/context"
)

func Bind(model interface{}) flamego.Handler {
	// Ensure not pointer.
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		panic("form: pointer can not be accepted as binding model")
	}

	return func(ctx context.Context) error {
		obj := reflect.New(reflect.TypeOf(model))
		r := ctx.Request().Request
		if r.Body != nil {
			defer func() { _ = r.Body.Close() }()
			if err := json.NewDecoder(r.Body).Decode(obj.Interface()); err != nil {
				return ctx.Error(http.StatusBadRequest, "Failed to parse form data")
			}
		}

		acceptLanguage := r.Header.Get("Accept-Language")
		languageTags, _, _ := language.ParseAcceptLanguage(acceptLanguage)
		languageTag := language.Chinese
		if len(languageTags) > 0 {
			languageTag = languageTags[0]
		}

		errors, ok := govalid.Check(obj.Interface(), languageTag)
		if !ok {
			var msg string
			if len(errors) > 0 {
				msg = errors[0].Error()
			}
			return ctx.Error(http.StatusBadRequest, "%s", msg)
		}

		// Validation passed.
		ctx.Map(obj.Elem().Interface())
		return nil
	}
}
