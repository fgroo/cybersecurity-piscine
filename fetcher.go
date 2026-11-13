package main


import (
	"net/http"
	"os"
	"fmt"
	"io"
	"strings"
	"strconv"
	"errors"
)

type vars struct {
	seen	[]string
	images	[]string
}

func	prevent_doubles(seen []string, target string) bool {
	for i, _ := range seen {
		if target == seen[i] { return true }
	}
	return false
}

// there are many other tags like: <input type="image" src>, <link rel="icon">, <img srcset>, <img data-src>

// func	process_hrefs(seen string[]) {
// 	for i, _ := range seen {
// 		main
// 	}
// }

func	parse_img_urls_for_img_src(content string) (path_str string, snipped_content string) {
	valid_image_types := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	start_index := strings.Index(content, "<img")
	if start_index == -1 {
		return "", ""
	}
	tag_end := strings.Index(content[start_index:], ">")
	if tag_end == -1 {
		return "", ""
	}
	tag := content[start_index : start_index + tag_end + 1]
	snipped_content = content[start_index + tag_end + 1:]
	
	src_index := strings.Index(tag, " src=\"")
	if src_index == -1 {
		return "", snipped_content
	}
	src_start := start_index + src_index + len(" src=\"")
	end_index := strings.Index(content[src_start:], "\"")
	if end_index == -1 {
		return "", ""
	}
	path_str = content[src_start : src_start+end_index]

	// INFO: a src may carry ?query or #fragment; only the file name decides the type
	file_name := path_str
	if cut := strings.IndexAny(file_name, "?#"); cut != -1 {
		file_name = file_name[:cut]
	}
	for _, value := range valid_image_types {
		if strings.HasSuffix(strings.ToLower(file_name), value) {
			return path_str, content[src_start + end_index + 1:]
		}
	}
	return "", snipped_content
}

// href=" , 

func	find_hrefs_urls(content string, depth int) (seen []string, err error) {

	var find_line int

	for depth > 0 {
		if find_line = strings.Index(content, "<a") ; find_line == -1 {
			return seen, fmt.Errorf("<a")
		}
		if find_line := strings.Index(content[find_line:], "<\\a>") ; find_line == -1 {
			fmt.Print("error: parsing failed due to missing <\\a>\n")
			return seen, fmt.Errorf("<\\a>")
		}
		if find_line := strings.Index(content[find_line:], "href=\"") ; find_line == -1 {
			continue
		}
		start_path := find_line
		if find_line := strings.Index(content[find_line:], "\"") ; find_line == -1 {
			fmt.Print("error: parsing failed due to missing \" in href=\"\n")
			return seen, fmt.Errorf("\"")
		}
		url := content[start_path:find_line]
		if prevent_doubles(seen, url) == false { continue }
		seen = append(seen, url)
		content = content[find_line:]
		depth--
	}
	return seen, nil
}

func	fetch_url(url string) (content string) {
	resp, err := http.Get(url) // https://de.wikipedia.org/wiki/6-7
	if err != nil {
		fmt.Errorf("error: %d - not available\n", err) // handle error
		return ""
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("error: ReadAll", err)
		return ""
	}
	content = string(bytes)
	return // :D
}

func	executer(recursive bool, depth int, path string) {

	var	urls vars
	content := fetch_url(os.Args[len(os.Args)-1])
	content = string(content)

	if recursive == false || depth == 0 {
		image := ""
		for content != "" {
			image, content = parse_img_urls_for_img_src(content)
			if image == "" {
				continue
			}
			urls.images = append(urls.images, image)
		}

	} else {
		if urls.seen, err = find_hrefs_urls(content, depth) ; err != nil {
			return
		for i, cont
		for content != "" {
			image, content = parse_img_urls_for_img_src(content)
			if image == "" {
				continue
			}
			urls.images = append(urls.images, image)
		}
	}
	err := os.Mkdir(path, 0755)
	if errors.Is(err, os.ErrExist) {
		fmt.Println("Folder already exists")
	} else if err != nil {
		fmt.Println("Something went wrong:", err)
	}
	// TODO: download images and put into folder
}


func	parse_options() (recursive bool, depth int, path string) {
	default_path := "./data/"
	var	err error

	if len(os.Args) == 2 {
		return false, 0, default_path
	}

	if os.Args[1] == "-r" {
		return true, 0, default_path
	}
	if os.Args[1] == "-rl" || os.Args[1] == "-lr" {
		if depth, err = strconv.Atoi(os.Args[2]) ; err != nil {
			return true, 5, default_path
		}
		return true, depth, default_path
	}
	if os.Args[1] == "-p" {
		return false, 0, os.Args[2]
	}
	if os.Args[1] == "-rp" || os.Args[1] == "-pr" {
		return true, 0, os.Args[2]
	}
	if os.Args[1] == "-rpl" || os.Args[1] == "-prl" || os.Args[1] == "-plr" {
		if depth, err = strconv.Atoi(os.Args[3]) ; err != nil {
			return true, 5, os.Args[2]
		}
		return true, depth, os.Args[2]
	}
	if os.Args[1] == "-lpr" || os.Args[1] == "-rp" || os.Args[1] == "-pr" {
		if depth, err = strconv.Atoi(os.Args[2]) ; err != nil {
			return true, 5, os.Args[3]
		}
		return true, depth, os.Args[3]
	}
	panic("parse_options: bad programming skills")
}

func	main() {
	if len(os.Args) == 1 {
		fmt.Errorf("run the command at least with one parameter, which is the URL\n", "arguments missing")
		return
	} // return fmt.Errorf("missing required parameter: %s", "user_id")
	if len(os.Args) == 2 && len(os.Args[1]) <  len("http://x.x") {
		fmt.Errorf("if you only use one parameter, then that one has to be the URL\n", "URL missing")
		return
	} // there might be a way better option to check for valid parameters
	if len(os.Args) == 3 && len(os.Args[1]) >  len("-rlp") {
		fmt.Errorf("the first parameters must be the option [-rlp], with the URL as the last parameter\n", "wrong format")
		return
	}
	if len(os.Args) == 3 && os.Args[1] != "-r" && os.Args[1] != "-p" && os.Args[1] != "-rl" && os.Args[1] != "-lr" && os.Args[1] != "-rlp" && os.Args[1] != "-prl" && os.Args[1] != "-rpl" && os.Args[1] != "-lpr" && os.Args[1] != "-pr"  && os.Args[1] != "-rp" {
		fmt.Errorf("invalid parameter, use [-rlp]\n", "invalid parameter")
		return
	}
	executer(parse_options())
	// fmt.Println("%s", content)
	// if new_image := parse_img_urls_for_img_src(content); new_image != "" {
	// 	urls.images = append(urls.images, new_image)
	// }
}