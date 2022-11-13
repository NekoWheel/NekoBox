// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package censor

import (
	"encoding/json"
)

type ForbiddenType string

const (
	ForbiddenTypeSpam        ForbiddenType = "spam"
	ForbiddenTypeAd          ForbiddenType = "ad"
	ForbiddenTypePolitics    ForbiddenType = "politics"
	ForbiddenTypeTerrorism   ForbiddenType = "terrorism"
	ForbiddenTypeAbuse       ForbiddenType = "abuse"
	ForbiddenTypePorn        ForbiddenType = "porn"
	ForbiddenTypeFlood       ForbiddenType = "flood"
	ForbiddenTypeContraband  ForbiddenType = "contraband"
	ForbiddenTypeMeaningless ForbiddenType = "meaningless"
)

type TextCensorResponse struct {
	SourceName    string          `json:"source_name"`
	Pass          bool            `json:"pass"`
	ForbiddenType ForbiddenType   `json:"forbidden_type"`
	Hint          string          `json:"hint"`
	Confidence    float64         `json:"confidence"`
	Metadata      json.RawMessage `json:"metadata"`
}

type TextCensor interface {
	Censor(text string) (*TextCensorResponse, error)
}
