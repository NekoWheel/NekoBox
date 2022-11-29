// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"embed"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flamego/template"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

var _ template.FileSystem = (*FileSystem)(nil)

type file struct {
	name string
	data []byte
	ext  string
}

func (f *file) Name() string          { return f.name }
func (f *file) Data() ([]byte, error) { return f.data, nil }
func (f *file) Ext() string           { return f.ext }

type FileSystem struct {
	files []template.File
}

func (fs *FileSystem) Files() []template.File { return fs.files }

// getExt returns the extension of given name, prefixed with the dot (".").
func getExt(name string) string {
	i := strings.Index(name, ".")
	if i == -1 {
		return ""
	}
	return name[i:]
}

func Minify(embedFS embed.FS, dir string, allowedExtensions []string) (*FileSystem, error) {
	var beforeSize, afterSize int64

	var files []template.File
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile(`^(application|text)/(x-)?(java|ecma)script$`), js.Minify)

	err := fs.WalkDir(embedFS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relpath, err := filepath.Rel(dir, path)
		if err != nil {
			return errors.Wrap(err, "get relative path")
		}

		ext := getExt(relpath)
		for _, allowed := range allowedExtensions {
			if ext != allowed {
				continue
			}

			data, err := embedFS.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "read")
			}
			beforeSize += int64(len(data))

			var buf bytes.Buffer
			switch filepath.Ext(path) {
			case ".html":
				if err := html.Minify(m, &buf, bytes.NewReader(data), nil); err != nil {
					return errors.Wrap(err, "minify html")
				}
			default:
				buf = *bytes.NewBuffer(data)
			}
			afterSize += int64(buf.Len())

			name := filepath.ToSlash(relpath[:len(relpath)-len(ext)])
			files = append(files,
				&file{
					name: name,
					data: buf.Bytes(),
					ext:  ext,
				},
			)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	logrus.WithFields(logrus.Fields{
		"before_size": beforeSize,
		"after_size":  afterSize,
		"ratio":       float64(afterSize) / float64(beforeSize),
	}).Info("Minify succeeded")

	return &FileSystem{
		files: files,
	}, nil
}
