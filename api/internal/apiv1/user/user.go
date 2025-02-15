package user

import (
	"encoding/json"

	"github.com/gotomicro/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-contrib/sessions"

	"github.com/shimohq/mogo/api/internal/invoker"
	"github.com/shimohq/mogo/api/pkg/component/core"
	"github.com/shimohq/mogo/api/pkg/model/db"
)

// Info get userinfo
func Info(c *core.Context) {
	session := sessions.Default(c.Context)
	user := session.Get("user")
	tmp, _ := json.Marshal(user)
	u := db.User{}
	_ = json.Unmarshal(tmp, &u)
	u.Password = ""
	c.JSONOK(u)
	return
}

type login struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login ...
func Login(c *core.Context) {
	var param login
	err := c.Bind(&param)
	if err != nil {
		c.JSONE(1, err.Error(), nil)
		return
	}
	conds := egorm.Conds{}
	conds["username"] = param.Username
	user, _ := db.UserInfoX(conds)
	// hash, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(hash))
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))
	if err != nil {
		c.JSONE(1, "account or password error", "")
		return
	}
	session := sessions.Default(c.Context)
	session.Set("user", user)
	_ = session.Save()
	c.JSONOK("")
	return
}

// Logout ..
func Logout(c *core.Context) {
	session := sessions.Default(c.Context)
	session.Delete("user")
	err := session.Save()
	if err != nil {
		c.JSONE(1, "logout fail", err.Error())
		return
	}
	c.JSONOK("succ")
	return
}

type password struct {
	Password    string `form:"password" binding:"required"`
	NewPassword string `form:"newPassword" binding:"required"`
	ConfirmNew  string `form:"confirmNew" binding:"required"`
}

func UpdatePassword(c *core.Context) {
	uid := cast.ToInt(c.Param("uid"))
	if uid == 0 {
		c.JSONE(1, "invalid parameter", nil)
		return
	}
	var param password
	err := c.Bind(&param)
	if err != nil {
		c.JSONE(1, err.Error(), nil)
		return
	}

	elog.Debug("UpdatePassword", elog.Any("uid", uid), elog.Any("param", param))

	if param.ConfirmNew != param.NewPassword {
		c.JSONE(1, "password not match", "")
		return
	}
	if len(param.NewPassword) < 5 || len(param.NewPassword) > 32 {
		c.JSONE(1, "password length should between 5 ~ 32", "")
		return
	}
	user, _ := db.UserInfo(uid)

	elog.Debug("UpdatePassword", elog.Any("user", user))

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))
	if err != nil {
		c.JSONE(1, "account or password error", "")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(param.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSONE(1, "account or password error", "")
		return
	}
	ups := make(map[string]interface{}, 0)
	ups["password"] = string(hash)
	err = db.UserUpdate(invoker.Db, uid, ups)
	if err != nil {
		c.JSONE(1, "password update error", err.Error())
		return
	}
	c.JSONOK("")
	return
}
