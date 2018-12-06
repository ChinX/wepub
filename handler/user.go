package handler

import (
	"encoding/hex"
	"errors"
	"github.com/chinx/wepub/crypts"
	"net/http"
	"strings"
	"time"

	"github.com/chinx/wepub/api"
	"github.com/chinx/wepub/module"
)

var baseFormat = "2006-01-02 15:04:05"

func UserLogout(w http.ResponseWriter, r *http.Request) {
	result := &api.CommonResult{Status: module.StatusLogout}
	userData, err := module.NewSession(w, r)
	if err != nil {
		result.Message = "未登录"
		reply(w, http.StatusUnauthorized, result, err)
		return
	}
	result.Status = userData.Destroy()
	reply(w, http.StatusCreated, result, err)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	admin := &module.AdminLogin{}
	result := &api.CommonResult{Status: module.StatusLogout}
	err := readBody(r.Body, admin)
	if err != nil || admin.User == "" || admin.Password == "" {
		result.Message = "账号或密码不能为空"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	userData, err := module.NewSession(w, r)
	if err != nil {
		result.Message = "创建登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return
	}

	result.Status = userData.SetAdminSession(admin)
	if result.Status != module.StatusLogin {
		result.Message = "登录失败，请检测账户信息，稍后尝试"
		reply(w, http.StatusUnauthorized, result, nil)
		return
	}

	reply(w, http.StatusCreated, result, nil)
}

func RegistryHandler(w http.ResponseWriter, r *http.Request) {
	admin := &module.AdminLogin{}
	result := &api.CommonResult{Status: module.StatusLogout}
	err := readBody(r.Body, admin)
	if err != nil || admin.User == "" || admin.Password == "" || admin.Code == ""{
		result.Message = "账号、密码或邀请码不能为空"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	if ok, err := checkCode(admin.Code); !ok{
		result.Message = "邀请码已失效"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	err = module.CreateUser(admin)
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusInternalServerError, result, err)
		return
	}

	result.Status = module.StatusBinding
	reply(w, http.StatusCreated, result, nil)
}

func checkAdmin(w http.ResponseWriter, r *http.Request) *api.CommonResult {
	return  &api.CommonResult{Status: module.StatusLogin}
	result := &api.CommonResult{Status: module.StatusLogout}
	userData, err := module.NewSession(w, r)
	if err != nil {
		result.Message = "获取登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return result
	}

	if !userData.IsAdmin() {
		result.Message = "无权限"
		reply(w, http.StatusForbidden, result, err)
		return result
	}

	result.UserID, result.Status = userData.UserID()
	if result.Status == module.StatusLogout {
		result.Message = "登录状态已失效，请重新登录"
		reply(w, http.StatusUnauthorized, result, nil)
		return result
	}
	return result
}

func checkCode(code string) (bool, error) {
	data, err := hex.DecodeString(code)
	if err != nil {
		return false, err
	}

	rawData, err := crypts.AesDecrypt(data)
	if err != nil {
		return false, err
	}

	keys := strings.Split(string(rawData),"&")
	if len(keys) < 2{
		return false, errors.New("邀请码格式错误")
	}

	t, err := time.Parse(baseFormat, keys[1])
	if err != nil {
		return false, err
	}

	return string(keys[0]) == "windup" && t.After(time.Now()), nil
}
