package main


import (
	"net/http"
	"os"
	"fmt"
	"io"
	"strings")

type vars struct {
    curr   string
    seen   []string
    images []string
}

func	prevent_doubles(seen []string, target string) bool {
	for i, _ := range seen {
		if target == seen[i] { return true }
	}
	return false
}

// there are many other tags like: <input type="image" src>, <link rel="icon">, <img srcset>, <img data-src>

func	process_hrefs(seen string[]) {
	for i, _ := range seen {
		main
	}
}

func	parse_img_urls_for_img_src(content string) string {
	valid_body_seq := []string{"<img", ">", "src=\"", "\""} 
	valid_image_types := []string{".jpg\"", ".jpeg\"", ".png\"", ".gif\"", ".bmp\""}

	start_index := strings.Index(content, valid_body_seq[0])
	if start_index == -1 {
		return ""
	}
	end_index := strings.Index(content[:start_index], valid_body_seq[1])
	if end_index == -1 {
		return ""
	}
	start_path := strings.Index(content[:start_index], valid_body_seq[2])
	if start_path == -1 {
		return ""
	}
	end_path := strings.Index(content[:start_path], valid_body_seq[3])
	if end_path == -1 {
		return ""
	}

	path_str := content[start_path:end_path]
	for i, _ := range valid_image_types {
		if strings.LastIndex(path_str, valid_image_types[i]) == -1 {
			continue
		} else if strings.LastIndex(path_str, valid_image_types[i]) < len(path_str) - len(valid_image_types[i]) {
			fmt.Print("error: occurence is not at the end\n")
			continue
		}
		return path_str
	}
	return ""
}

// href=" , 

func	find_child_urls(content string, depth int) []string {
	
	list := []string

	for depth ; depth-- {
		if find_line := strings.Index(content, "<a") ; find_line == -1 {
			return list
		}
		if find_line := strings.Index(content[find_line:], "<\\a>") ; find_line == -1 {
			fmt.Print("error: parsing failed due to missing <\\a>\n")
			return list
		}
		if find_line := strings.Index(content[find_line:], "href=\"") ; find_line == -1 {
			continue
		}
		start_path := find_line
		if find_line := strings.Index(content[find_line:], "\"") ; find_line == -1 {
			fmt.Print("error: parsing failed due to missing \" in href=\"\n")
			return list
		}
		url := content[start_path:find_line]
		if prevent_doubles(list.seen, url) == false { continue }
		append(list.seen, url)
		content = content[find_line:]
	}
	return list
}



func	parser(content string, depth int) (urls vars, err error) {

	var	urls vars

	

}

func	

func	main() {

	if len(os.Args) == 1 { fmt.Errorf("run the command at least with one parameter, which is the URL\n", "arguments missing") } // return fmt.Errorf("missing required parameter: %s", "user_id")
	if len(os.Args) == 2 and len(os.Args[1]) <  len("https://") { fmt.Errorf("if you only use one parameter, then that one has to be the URL\n", "URL missing") }
	if len(os.Args) == 3 and len(os.Args[1]) <  len("https://") { fmt.Errorf("if you only use one parameter, then that one has to be the URL\n", "URL missing") }
	resp, err := http.Get(os.Args[1]) // https://de.wikipedia.org/wiki/6-7
	if err != nil {
		fmt.Print("error: %d - not available\n", err) // handle error
		return 
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("error: ReadAll")
		return
	}
	content := string(bytes)
	fmt.Print("%s\n", content)
	if new_image := parse_img_urls_for_img_src(content); new_image != "" {
		urls.images = append(urls.images, new_image)
	}
}