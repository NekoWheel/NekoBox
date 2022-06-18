// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

import (
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

// File is the configuration object.
var File *ini.File

func Init() error {
	var err error
	File, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, "conf/app.ini")
	if err != nil {
		return errors.Wrap(err, "parse 'conf/app.ini'")
	}

	if err := File.Section("app").MapTo(&App); err != nil {
		return errors.Wrap(err, "map 'server'")
	}

	if err := File.Section("server").MapTo(&Server); err != nil {
		return errors.Wrap(err, "map 'server'")
	}

	if err := File.Section("database").MapTo(&Database); err != nil {
		return errors.Wrap(err, "map 'database'")
	}

	if err := File.Section("recaptcha").MapTo(&Recaptcha); err != nil {
		return errors.Wrap(err, "map 'recaptcha'")
	}

	if err := File.Section("upload").MapTo(&Upload); err != nil {
		return errors.Wrap(err, "map 'upload'")
	}

	if err := File.Section("mail").MapTo(&Mail); err != nil {
		return errors.Wrap(err, "map 'mail'")
	}

	return nil
}
