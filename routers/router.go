// Copyright 2013 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package routers implemented controller methods of beego.
package routers

import (
	"net/url"
	"strings"

	"github.com/astaxie/beego"
)

var langTypes []*langType // Languages are supported.

// langType represents a language type.
type langType struct {
	Lang, Name string
}

func init() {
	// Initialized language type list.
	langs := strings.Split(beego.AppConfig.String("langs"), "|")
	names := strings.Split(beego.AppConfig.String("langNames"), "|")
	langTypes = make([]*langType, 0, len(langs))
	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Lang: v,
			Name: names[i],
		})
	}
}

// setLangVer sets site language version.
func setLangVer(ctx *beego.Context, input url.Values, data map[interface{}]interface{}) langType {
	// 1. Check URL arguments.
	lang := input.Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		ck, err := ctx.Request.Cookie("lang")
		if err == nil {
			lang = ck.Value
		}
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := ctx.Request.Header.Get("Accept-Language")
		if len(al) > 2 {
			al = al[:2] // Only compare first two letters.
			for _, v := range langTypes {
				if al == v.Lang {
					lang = al
					break
				}
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en"
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	ctx.SetCookie("lang", curLang.Lang+";path=/", 0)

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	// Set language properties.
	data["Lang"] = curLang.Lang
	data["CurLang"] = curLang.Name
	data["RestLangs"] = restLangs

	return curLang
}
