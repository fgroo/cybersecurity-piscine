package main

import (
	"fmt"
	exif "github.com/dsoprea/go-exif/v3"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("run the command at least with one parameter, which is the URL (arguments missing)")
		return
	}
	files := os.Args[1:]
	var	rawexif []byte
	var flatexif []exif.ExifTag
	for i, _ := range files {
		n, err := os.ReadFile(files[i])
		if err != nil {
			fmt.Printf("error: %v - %s\n", err, files[i])
			return
		}
		rawexif, err = exif.SearchAndExtractExif(n)
		if err != nil {
			fmt.Printf("error: %v - %s\n\n", err, files[i])
			continue
		}
		flatexif, _, err = exif.GetFlatExifData(rawexif, nil)
		if err != nil {
			fmt.Printf("error: %v - %s\n", err, files[i])
		}
		var creationdate string
		fmt.Printf("FILE%d: %s\n", i + 1, files[i])
		for _, tag := range flatexif {
			fmt.Println(tag.String())
			switch tag.TagName {
			case "DateTimeOriginal":
				creationdate = tag.Formatted
			case "DateTime":
				if creationdate == "" {
					creationdate = tag.Formatted
				}
			}
		}
		if creationdate != "" {
			fmt.Printf("Creation date: %s\n\n", creationdate)
		}
	}
}
