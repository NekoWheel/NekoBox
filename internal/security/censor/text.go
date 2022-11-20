// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package censor

import (
	"context"
	"encoding/json"
	"fmt"
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
	ForbiddenTypeHarmful     ForbiddenType = "harmful"
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
		ForbiddenTypeHarmful:     "不良场景",
	}[f]
}

type TextCensorResponse struct {
	SourceName    string          `json:"source_name"`
	Pass          bool            `json:"pass"`
	ForbiddenType ForbiddenType   `json:"forbidden_type"`
	Hint          string          `json:"hint"`
	Confidence    float64         `json:"confidence"`
	RawResponse   json.RawMessage `json:"raw_response"`
}

func (r *TextCensorResponse) ToJSON() []byte {
	jsonBytes, _ := json.Marshal(r)
	return jsonBytes
}

func (r *TextCensorResponse) ErrorMessage() string {
	errorMessage := "内容安全检查不通过"
	if r.ForbiddenType != "" {
		errorMessage += fmt.Sprintf(" [%s]", r.ForbiddenType.String())
	}
	if r.Hint != "" {
		errorMessage += fmt.Sprintf("，相关标签：%q", r.Hint)
	}
	return errorMessage
}

type TextCensor interface {
	Censor(ctx context.Context, text string) (*TextCensorResponse, error)
	String() string
}
