// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

const (
	fontWidth  = 11
	fontHeight = 18
)

var font = map[string][]byte{
	"0": {
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
	},
	"1": {
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	},
	"2": {
		0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0,
		0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	},
	"3": {
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0,
	},
	"4": {
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0,
		1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
	},
	"5": {
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0,
	},
	"6": {
		0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0,
		1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1,
		0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
	},
	"7": {
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0,
	},
	"8": {
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
	},
	"9": {
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0,
		1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1,
		0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0,
	},
}
