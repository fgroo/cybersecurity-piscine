package main

import (
	"fmt"
	exif "github.com/dsoprea/go-exif/v3"
	"os"
	//"io"
	//"path"
	//"path/filepath"
)

func	readexif(files ...string) {
	//valid_image_types := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}



}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("run the command at least with one parameter, which is the URL (arguments missing)")
		return
	}
	files := os.Args[1:]
	//var paths []*os.File
	var	rawexif []byte
	var flatexif []exif.ExifTag
	for i, _ := range files {
		n, err := os.ReadFile(files[i])
		if err != nil {
			fmt.Printf("error: %v - %d - %s\n", err, n, files[i])
			return
		}
		//defer n.Close()
		//fmt.Println(string(string(n)))
		rawexif, err = exif.SearchAndExtractExif(n)
		if err != nil {
			fmt.Printf("error: %v - %d - %s\n", err, n, files[i])
			return
		}
		flatexif, _, err = exif.GetFlatExifData(rawexif, nil)
		if err != nil {
			fmt.Printf("error: %v - %d - %s\n", err, n, files[i])
			return
		}
		//paths = append(paths, n)
	}
	//for _, r 
	
	//fmt.Println(paths)
	for _, tag := range flatexif {
	fmt.Println(tag.String())
	}
}
