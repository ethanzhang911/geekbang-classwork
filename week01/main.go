package main

import (
	"fmt"
	"github.com/ethanzhang911/geekbang-classwork/dao"
)

func main() {
	err := dao.Query()
	fmt.Println(err.Error())

}
