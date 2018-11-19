package module

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/chinx/mohist/random"
	"net/http"

	"github.com/chinx/wepub/model"
	"github.com/go-session/session"
)

const (
	StatusLogout int = iota
	StatusBinding
	StatusLogin

	PermissionAdmin int = 7

	idKey         = "openid"
	permissionKey = "permission"
)

type AdminLogin struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type Session struct {
	store session.Store
	w     http.ResponseWriter
	r     *http.Request
}

func NewSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return nil, err
	}
	return &Session{store: store, w: w, r: r}, nil
}

func (s *Session) IsAdmin() bool {
	key, ok := s.store.Get(permissionKey)
	if !ok {
		return false
	}
	return int(key.(float64)) >= PermissionAdmin
}

func (s *Session) SetAdminSession(data *AdminLogin) int {
	admin := &model.Admin{User: data.User}
	if ok := model.Get(admin); !ok {
		return s.Destroy()
	}

	sha := sha512.New()
	sha.Write([]byte(data.Password + admin.Salt))
	if hex.EncodeToString(sha.Sum(nil)) != admin.Password {
		return s.Destroy()
	}

	s.store.Set(idKey, data.User)
	s.store.Set(permissionKey, PermissionAdmin)
	err := s.store.Save()
	if err != nil {
		return s.Destroy()
	}
	return StatusLogin
}

func (s *Session) Destroy() int {
	session.Destroy(context.Background(), s.w, s.r)
	return StatusLogout
}

func (s *Session) UserID() (string, int) {
	userID, ok := s.store.Get(idKey)
	if !ok {
		return "", s.Destroy()
	}

	return userID.(string), StatusLogin
}

func CreateUser(data *AdminLogin) error {
	admin := &model.Admin{User: data.User}
	salt, err := random.NewStr(8)
	if err != nil {
		return  err
	}

	sha := sha512.New()
	sha.Write([]byte(data.Password + salt))
	admin.Password = hex.EncodeToString(sha.Sum(nil))
	admin.Salt = salt
	if ! model.Insert(admin){
		return errors.New("创建用户失败")
	}
	return nil
}
