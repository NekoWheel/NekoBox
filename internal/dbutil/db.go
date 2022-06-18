// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"database/sql"

	"gorm.io/gorm"
)

type Transactor interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
}
