// Pipe - A small and beautiful blogging platform written in golang.
// Copyright (C) 2017-2018, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"net/http"
	"strings"
		"github.com/HelloWorldZQ/quintinblog/model"
	"github.com/HelloWorldZQ/quintinblog/service"
	"github.com/HelloWorldZQ/quintinblog/util"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	)

// DataModel represents data model.
type DataModel map[string]interface{}

const nilB3id = "H9oxzSym"

func logUrl(c *gin.Context) {
	logger.Debug(c.Request.URL)
	c.Next()
}
func fillUser(c *gin.Context) {
	inited := service.Init.Inited()
	if !inited && util.PathInit != c.Request.URL.Path {
		c.Redirect(http.StatusSeeOther, model.Conf.Server+util.PathInit)
		c.Abort()

		return
	}

	dataModel := &DataModel{}
	c.Set("dataModel", dataModel)
	session := util.GetSession(c)
	(*dataModel)["User"] = session
	if 0 != session.UID {
		c.Next()

		return
	}

	uaStr := c.Request.UserAgent()
	if isBot(uaStr) {
		logger.Tracef("Bot User-Agent [%s]", uaStr)
		c.Next()

		return
	}
	redirectURL := model.Conf.Server + c.Request.URL.Path
	if "/admin/" == c.Request.URL.Path { // https://github.com/HelloWorldZQ/quintinblog/issues/67
		redirectURL = model.Conf.Server + c.Request.URL.Path
	}
	if strings.HasPrefix(c.Request.URL.Path, util.PathBlogs) {
		name := c.Request.URL.Path[len(util.PathBlogs)+1:]
		name = strings.Split(name, "?")[0]
		name = strings.Split(name, "/")[0]
		if "" != name {
			user := service.User.GetUserByName(name)
			if nil != user {
				userBlog := service.User.GetOwnBlog(user.ID)
				blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, userBlog.ID)
				redirectURL = blogURLSetting.Value + strings.Split(c.Request.URL.Path, util.PathBlogs+"/"+name)[1]
				if "" != c.Request.URL.RawQuery {
					redirectURL += "?" + c.Request.URL.RawQuery
				}
			}
		}
	} else {
		if !strings.HasPrefix(redirectURL, model.Conf.Server) {
			redirectURL = model.Conf.Server + c.Request.URL.Path
		}
	}
	redirectURL = strings.TrimSpace(redirectURL)
	if "" == redirectURL {
		redirectURL = model.Conf.Server + c.Request.URL.Path
	}
	if redirectURL == model.Conf.Server+c.Request.URL.Path {
		c.Next()
	} else {
		c.Redirect(http.StatusSeeOther, redirectURL)
		c.Abort()
	}
}

func isBot(uaStr string) bool {
	var ua = user_agent.New(uaStr)

	return ua.Bot() || strings.Contains(uaStr, "HacPai")
}
