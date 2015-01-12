package main

import (
	"fmt"
	"gumbo"
	"time"
)

func main() {
	st := time.Now()
	s := "<h1 class='xx'>Hello, World!</h1>"
	gb := gumbo.NewGumboParse(s)
	defer gb.Destory()
	//n, x := gb.GetNodeByTag(GUMBO_TAG_BODY, nil)
	//fmt.Println(x)
	v, x := gb.GetNodeByTagAndAttr(gumbo.GUMBO_TAG_H1, "class", "xx", nil)
	//v, x := gb.GetNodeByTag(GUMBO_TAG_H1, n)
	fmt.Println(x)
	//fmt.Println(gb.GetText(v))
	fmt.Println(gb.GetAttribute(v, "class"))
	fmt.Println(time.Now().Sub(st).Seconds())
}
