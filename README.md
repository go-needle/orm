<!-- markdownlint-disable MD033 MD041 -->
<div align="center">

# 🪡orm

<!-- prettier-ignore-start -->
<!-- markdownlint-disable-next-line MD036 -->
A simple orm framework for Golang
<!-- prettier-ignore-end -->

<img src="https://img.shields.io/badge/golang-1.11+-blue" alt="golang">
</div>

## introduction
This is a simple orm framework for Golang that supports Sqlite.

## installing
Select the version to install

`go get github.com/go-needle/orm@version`

If you have already get , you may need to update to the latest version

`go get -u github.com/go-needle/orm`


## quickly start
```golang
package main

import (
	"fmt"
	"github.com/go-needle/orm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `orm:"constraint:PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func main() {
	engine, err := orm.NewEngine("sqlite3", "g.db")
	if err != nil {
		panic(err)
	}
	s := engine.NewSession().Model(&User{}).Debug()
	if s.HasTable() {
		_ = s.DropTable()
	}
	_ = s.CreateTable()
	n, _ := s.Insert(user1, user2, user3)
	fmt.Println(n)
	var users []User
	_ = s.Where("Age >= ?", 20).Find(&users)
	fmt.Println(users)
}
```
