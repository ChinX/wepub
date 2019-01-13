package main

import (
	"fmt"
	"github.com/chinx/wepub/handler"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/chinx/wepub/model"
	"github.com/chinx/wepub/router"
	"github.com/chinx/wepub/setting"
	"github.com/go-session/redis"
	"github.com/go-session/session"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	opt, err := setting.LoadConfigFile("./cert/coupon_private.key", "./conf/windup.toml")
	if err != nil {
		log.Fatal(err)
	}

	err = model.InitORM("mysql",
		fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8",
			opt.Mysql.User, opt.Mysql.Password,
			opt.Mysql.Server, opt.Mysql.Port,
			opt.Mysql.Database))

	if err != nil {
		log.Fatal(err)
	}

	handler.StaticDir = opt.StaticDir
	session.InitManager(
		session.SetStore(redis.NewRedisStore(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", opt.Redis.Server, opt.Redis.Port),
			DB:       opt.Redis.Session,
			Password: opt.Redis.Password,
		})),
	)

	serveHandler, err := router.InitRouter()
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":"+strconv.Itoa(opt.HttpPort), serveHandler)
	if err != nil {
		log.Fatal(err)
	}
}
