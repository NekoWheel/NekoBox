// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package censor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveTrustedURL(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		want string
	}{
		{name: "empty string", text: "", want: ""},
		{name: "normal", text: "你喜欢看什么书？", want: "你喜欢看什么书？"},
		{name: "xiaohongshu", text: `59 茄酱发布了一篇小红书笔记，快来看吧！ 😆 jk14SBf2G75FTW2 😆 http://xhslink.com/a/WYaQ6Tjn86Hab，复制本条信息，打开【小红书】App查看精彩内容！`, want: `59 茄酱发布了一篇小红书笔记，快来看吧！ 😆 jk14SBf2G75FTW2 😆 /a/WYaQ6Tjn86Hab，复制本条信息，打开【小红书】App查看精彩内容！`},
		{name: "weibo", text: `TK 发微博啦~ https://weibo.com/1401527553/5144129806795225 真好看！`, want: `TK 发微博啦~ /1401527553/5144129806795225 真好看！`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := RemoveTrustedURL(tc.text)
			assert.Equal(t, tc.want, got)
		})
	}
}
