package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var tables = make([]xorm.TableName, 0, 16)
var engine *xorm.Engine

func InitORM(driver, source string) error {
	eng, err := xorm.NewEngine(driver, source)
	if err != nil {
		return err
	}

	if err = eng.Ping(); err != nil {
		return err
	}

	eng.SetMapper(core.GonicMapper{})
	engine = eng

	syncTables()
	return nil
}

func syncTables() {
	for _, value := range tables {
		engine.Sync2(value)
	}
}

func register(tab xorm.TableName) {
	if isExist(tab.TableName()) {
		fmt.Printf("table %s is already exists", tab.TableName())
	}
	tables = append(tables, tab)
}

func isExist(name string) bool {
	for _, value := range tables {
		if value.TableName() == name {
			return true
		}
	}
	return false
}

func NewSession() *xorm.Session {
	return engine.NewSession()
}
