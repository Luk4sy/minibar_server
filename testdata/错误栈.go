package main

import (
	"errors"
	"fmt"
	e "github.com/pkg/errors"
	"runtime"
)

func main() {
	err := errors.New("错误")
	SetError1(err)
	SetError2(err)
}

func SetError1(err error) {
	var msg = make([]byte, 1024)
	n := runtime.Stack(msg, false) // 把错误信息扫描到最近的栈
	fmt.Println(string(msg[:n]))
}

func SetError2(err error) {
	msg := e.WithStack(err)
	fmt.Printf("%+v\n", msg)
}
