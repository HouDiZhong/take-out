package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/global"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

type User struct {
	Code string `json:"code"`
}
type WeCahtResp struct {
	Openid string `json:"openid"`
	Errmsg string `json:"errmsg"`
}
type UserController struct {
	service service.UserService
}

func NewUserController(uService service.UserService) *UserController {
	return &UserController{service: uService}
}

func (u *UserController) GetWechatOpenId(code string) (string, error) {
	/* wechat := global.Config.Wechat
	const WeCahtUrl = "https://api.weixin.qq.com/sns/jscode2session&grant_type=authorization_code"
	resp, err := http.Get(fmt.Sprintf("%s&appid=%s&secret=%s&js_code=%s", WeCahtUrl, wechat.AppId, wechat.Secret, code))

	if err != nil {
		global.Log.Error("微信请求uid接口报错", "Error", err.Error())
		return "", err
	}

	defer resp.Body.Close()

	var weCahtResp = new(WeCahtResp)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, weCahtResp)

	return weCahtResp.Openid, nil */
	// 没有小程序相关信息，暂时将code作为openid返回
	return code, nil
}

func (u *UserController) Login(c *gin.Context) {
	var user User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	openId, err := u.GetWechatOpenId(user.Code)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
		return
	}

	token, err := u.service.Login(openId)
	if err != nil {
		global.Log.Error("用户登录失败", "Error", err.Error())
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: "登录失败"})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: token, Msg: "登录成功"})
}

func (u *UserController) Logout(c *gin.Context) {
	err := u.service.Logout(c)
	if err != nil {
		global.Log.Error("用户登出失败", "Error", err.Error())
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: "登出失败"})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "登出成功"})
}
