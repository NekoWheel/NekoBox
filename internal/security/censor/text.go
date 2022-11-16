// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package censor

import (
	"bytes"
	"context"
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

func (f ForbiddenType) String() string {
	return map[ForbiddenType]string{
		ForbiddenTypeSpam:        "含垃圾信息",
		ForbiddenTypeAd:          "广告",
		ForbiddenTypePolitics:    "涉政",
		ForbiddenTypeTerrorism:   "暴恐",
		ForbiddenTypeAbuse:       "辱骂",
		ForbiddenTypePorn:        "色情",
		ForbiddenTypeFlood:       "灌水",
		ForbiddenTypeContraband:  "违禁",
		ForbiddenTypeMeaningless: "无意义",
	}[f]
}

type TextCensorResponse struct {
	SourceName    string          `json:"source_name"`
	Pass          bool            `json:"pass"`
	ForbiddenType ForbiddenType   `json:"forbidden_type"`
	Hint          string          `json:"hint"`
	Confidence    float64         `json:"confidence"`
	Metadata      json.RawMessage `json:"metadata"`
}

func (r *TextCensorResponse) ToJSON() []byte {
	jsonBytes, _ := json.Marshal(r)
	return jsonBytes
}

func CheckTextCensorResponseValid(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}

	if bytes.EqualFold(raw, []byte("null")) {
		return false
	}

	var r TextCensorResponse
	if err := json.Unmarshal(raw, &r); err != nil {
		return false
	}
	return r.SourceName != ""
}

type TextCensor interface {
	Censor(ctx context.Context, text string) (*TextCensorResponse, error)
}
