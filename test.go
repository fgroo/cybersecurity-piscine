package main


import (
	"fmt"
	"strings"
	"os"
	"io"
	"net/http"
)

func	main() {

	resp, err := http.Get("https://go.dev/doc/tutorial/create-module")
	if err != nil {
		panic("wjfoisg")
	}
	bytes, err := io.ReadAll(resp.Body)
	str := string(bytes)
	tag_begin := strings.Index(str, "<img")
	tag_end := strings.Index(str[tag_begin:], ">")
	tag := tag_begin + tag_end
	fmt.Println(tag_begin)
	fmt.Println(tag_end)
	fmt.Println(tag)
	// fmt.Println(str[tag_begin:tag_begin+tag_end + 1])
	newstr := str[tag_begin:tag_begin+tag_end + 1]
	fmt.Println(newstr)
	// fmt.Print(str)
	if err := os.WriteFile("debug.html", bytes, 0644); err != nil {
		panic(err)
	}
}