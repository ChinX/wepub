package model

type Admin struct {
	User     string `xorm:"user varchar(40) notnull pk"`
	Password string `xorm:"password varchar(128) notnull"`
	Salt     string `xorm:"salt varchar(8) notnull"`
}

func (a *Admin) TableName() string {
	return "admin"
}

func init() {
	register(&Admin{})
}
