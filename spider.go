package main


import (
	"net/http"
	"os"
	"fmt"
	"io"
	"strings"
	"strconv"
	"errors"
	"net/url"
	"path"	// FIXED: last path element of a URL -> the file name
	"path/filepath"	// FIXED: build the local destination path
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
	
	// FIXED: lazy-load pages (img.com does this) leave src="..." as a tiny blur-up
	// placeholder and put the real file in data-src="...". Prefer data-src, keep src as fallback.
	src_key := "data-src=\""
	src_index := strings.Index(tag, src_key)
	if src_index == -1 {
		src_key = " src=\""	// OLD: src_index := strings.Index(tag, " src=\"")
		src_index = strings.Index(tag, src_key)
	}
	if src_index == -1 {
		return "", snipped_content
	}
	src_start := start_index + src_index + len(src_key)
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

// FIXED: a tag must be "<a" + something that ENDS the tag name, or "<article>",
// "<aside>", "<area>", "<abbr>" match too. Mirrors parse_img_urls_for_img_src:
// ok == false means there is no <a> left, which is the end of the page, not an error.
func	parse_href_for_a_tag(content string) (path_str string, snipped_content string, ok bool) {
	start_index := strings.Index(content, "<a")
	if start_index == -1 {
		return "", content, false
	}
	tag_end := strings.Index(content[start_index:], ">")
	if tag_end == -1 {
		return "", content, false
	}
	tag := content[start_index : start_index + tag_end + 1]
	snipped_content = content[start_index + tag_end + 1:]

	if c := tag[2] ; c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != '/' && c != '>' {
		return "", snipped_content, true	// false match on <article> etc.
	}
	href_index := strings.Index(tag, "href=\"")
	if href_index == -1 {
		return "", snipped_content, true	// <a name="top"> without href
	}
	href_start := href_index + len("href=\"")
	end_index := strings.Index(tag[href_start:], "\"")
	if end_index == -1 {
		return "", snipped_content, true
	}
	return tag[href_start : href_start + end_index], snipped_content, true
}

func	find_hrefs_urls(content string, depth int) (seen []string, err error) {

	// OLD:
	// 	var find_line int
	// 	for depth > 0 {
	// 		if find_line = strings.Index(content, "<a") ; find_line == -1 { return seen, fmt.Errorf("<a") }
	// 		if find_line := strings.Index(content[find_line:], "<\\a>") ; find_line == -1 {
	// 			fmt.Print("error: parsing failed due to missing <\\a>\n"); return seen, fmt.Errorf("<\\a>") }
	// 		if find_line := strings.Index(content[find_line:], "href=\"") ; find_line == -1 { continue }
	// 		start_path := find_line
	// 		if find_line := strings.Index(content[find_line:], "\"") ; find_line == -1 {
	// 			fmt.Print("error: parsing failed due to missing \" in href=\"\n"); return seen, fmt.Errorf("\"") }
	// 		url := content[start_path:find_line]
	// 		if prevent_doubles(seen, url) == false { continue }
	// 		seen = append(seen, url); content = content[find_line:]; depth--
	// 	}
	// 	return seen, nil
	// FIXED (all of the above):
	//   - "<\\a>" is the literal string "<\a>", never present in HTML (the closing tag is "</a>"),
	//     so the very first link made the whole function return an error and 0 urls.
	//   - "find_line :=" declares a NEW variable inside the if, so the outer find_line never moved:
	//     the offsets returned by content[find_line:] were then used as absolute indices.
	//   - prevent_doubles(...) == false { continue } skipped the NEW urls and kept the duplicates.
	//   - no more links left is normal, so it returns nil, not an error.
	// FIXED: depth == -1 means "collect every link on the page" so that the filtering
	// can happen before the limit. Old: "for depth > 0 { ... depth-- }" capped the RAW
	// list, so the social links in the nav ate the budget first:
	//   -rl 5 on img.com/news/our-news collected /, linkedin, twitter, instagram, facebook
	//   and then filter_links threw 4 of them away -> only 1 page was ever opened.
	for depth == -1 || depth > 0 {
		href := ""
		for {
			var ok bool
			href, content, ok = parse_href_for_a_tag(content)
			if !ok {
				return seen, nil
			}
			if href != "" && prevent_doubles(seen, href) == false {
				break
			}
		}
		seen = append(seen, href)
		if depth > 0 {
			depth--
		}
	}
	return seen, nil
}

// FIXED: a real page mixes three kinds of href: "/portfolio" (relative), "#" (scroll anchor)
// and "https://twitter.com/IMG" (leaves the site). Only same-host ones can be fetched,
// otherwise -r crawls off to linkedin/instagram instead of staying on the target site.
func	filter_links(base string, seen []string) (out []string) {
	base_url, base_err := url.Parse(base)
	if base_err != nil {
		return nil
	}
	for _, href := range seen {
		href = strings.TrimSpace(href)
		if href == "" || strings.HasPrefix(href, "#") {
			continue
		}
		target, err := url.Parse(href)
		if err != nil {
			continue
		}
		if target.Scheme != "" && target.Scheme != "http" && target.Scheme != "https" {
			continue	// mailto:, tel:, javascript:, ...
		}
		full := base_url.ResolveReference(target)
		if full.Host != base_url.Host {
			continue
		}
		full.Fragment = ""
		if full.Path == "" {
			full.Path = "/"
		}
		if prevent_doubles(out, full.String()) == false {
			out = append(out, full.String())
		}
	}
	return out
}

func	fetch_url(url string) (content string) {
	resp, err := http.Get(url) // https://de.wikipedia.org/wiki/6-7
	if err != nil {
		// OLD: fmt.Errorf(...) — the error value was thrown away, so nothing was ever printed
		fmt.Printf("error: %v - not available\n", err)	// FIXED: Printf shows it; %d on an error is wrong too
		return ""
	}
	defer resp.Body.Close()	// FIXED: without this every link leaks the connection of its body
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: ReadAll %v\n", err)	// FIXED: same as above
		return ""
	}
	content = string(bytes)
	return // :D
}

func	executer(recursive bool, depth int, path string) {

	var	urls vars
	page := os.Args[len(os.Args)-1]
	content := fetch_url(page)
	orig := content


	image := ""
	for content != "" {
		image, content = parse_img_urls_for_img_src(content)
		if image == "" {
			continue
		}
		if (prevent_doubles(urls.images, image) == false) {	urls.images = append(urls.images, image) }
	}
	var err error
	if recursive == true || depth > 0 {
		if urls.seen, err = find_hrefs_urls(orig, depth) ; err != nil {
			return
		}
		urls.seen = filter_links(page, urls.seen)	// FIXED: relative hrefs, bare "#" and off-site links
		for _, link := range urls.seen {
			content = fetch_url(link)

			for content != "" {
				image, content = parse_img_urls_for_img_src(content)
				if image == "" {
					continue
				}
				if (prevent_doubles(urls.images, image) == false) {	urls.images = append(urls.images, image) }
			}
		}
	}
	for i, _  := range urls.images {
		fmt.Println(urls.images[i])
	}
	err = os.Mkdir(path, 0755)
	if errors.Is(err, os.ErrExist) {
		fmt.Println("Folder already exists")
	} else if err != nil {
		fmt.Println("Something went wrong:", err)
		// OPEN: return might be better to prevent a loop of failed written images on the disk
	}
	// TODO: download images and put into folder
	// FIXED: strings.Trim(link, "https://") is NOT a prefix trim. The 2nd argument is a SET OF
	// CHARACTERS: it cuts any of h,t,p,s,: and / off the front AND the back. The scheme happens
	// to go away, but "https://s.x.com/..." also eats the first host letter, and "?width=800"
	// stays on. What is left is still a whole path ("dyl347hiwv3ct.cloudfront.net/app/uploads/
	// 2026/11/PR-Horizontal.jpeg"), so os.Create fails with "no such file or directory" because
	// those folders do not exist under ./data/. The name I want is the LAST path element:
	// url.Path drops ?query and #fragment, path.Base cuts everything before the last "/".
	// OLD (kept for reference):
	// 	img, err := os.Create(path + strings.Trim(link, "https://"))
	// 	...
	// 	defer resp.Body.Close()	// inside the for loop
	// 	defer img.Close()		// inside the for loop, deferred a second time at the end
	for _, link := range urls.images {
		dest := image_file_name(link)
		if dest == "" {
			continue	// no file name in the URL at all, nothing to store
		}
		// filepath.Join so it also works when -p was given without a trailing "/"
		if err := download_image(filepath.Join(path, dest), link) ; err != nil {
			fmt.Printf("error: %v - %s\n", err, link)
			continue
		}
	}
}

// FIXED: the last path element of a URL, ready to be a file name.
func	image_file_name(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	// "path" and not "filepath": URL paths always use "/", filepath would also treat
	// "\" as a separator on Windows.
	name := path.Base(u.Path)
	if name == "" || name == "." || name == ".." || name == "/" {
		return ""
	}
	return name
}

// FIXED: a function of its own so the two defers run at the END of THIS call, once per image.
// The old version had "defer resp.Body.Close()" inside the for loop in executer, and defers
// only run when the ENCLOSING function returns: with 50 images nothing was closed until the
// whole program ended, and it runs out of file descriptors / connections before that.
func	download_image(dest string, link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()	// FIXED: once is enough; the old code deferred img.Close() twice

	_, err = io.Copy(file, resp.Body)
	return err
}


func	parse_options() (recursive bool, depth int, path string) {
	default_path := "./data/"
	var	err error

	if len(os.Args) == 2 {
		return false, 0, default_path
	}

	if os.Args[1] == "-r" {
		return true, 5, default_path
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
		return true, 5, os.Args[2]
	}
	if os.Args[1] == "-rpl" || os.Args[1] == "-prl" || os.Args[1] == "-plr" {
		if depth, err = strconv.Atoi(os.Args[3]) ; err != nil {
			return true, 5, os.Args[2]
		}
		return true, depth, os.Args[2]
	}
	if os.Args[1] == "-lpr" || os.Args[1] == "-rlp" {
		if depth, err = strconv.Atoi(os.Args[2]) ; err != nil {
			return true, 5, os.Args[3]
		}
		return true, depth, os.Args[3]
	}
	panic("parse_options: bad programming skills")
}

func	main() {
	if len(os.Args) == 1 {
		// FIXED: fmt.Errorf only BUILDS an error and the result was thrown away, so these
		// messages never reached the terminal. Println actually writes them out.
		fmt.Println("run the command at least with one parameter, which is the URL (arguments missing)")
		return
	} // return fmt.Errorf("missing required parameter: %s", "user_id")
	if len(os.Args) == 2 && len(os.Args[1]) <  len("http://x.x") {
		fmt.Println("if you only use one parameter, then that one has to be the URL (URL missing)")
		return
	} // there might be a way better option to check for valid parameters
	if len(os.Args) == 3 && len(os.Args[1]) >  len("-rlp") {
		fmt.Println("the first parameters must be the option [-rlp], with the URL as the last parameter (wrong format)")
		return
	}
	if len(os.Args) == 3 && os.Args[1] != "-r" && os.Args[1] != "-p" && os.Args[1] != "-rl" && os.Args[1] != "-lr" && os.Args[1] != "-rlp" && os.Args[1] != "-prl" && os.Args[1] != "-rpl" && os.Args[1] != "-lpr" && os.Args[1] != "-pr"  && os.Args[1] != "-rp" {
		fmt.Println("invalid parameter, use [-rlp] (invalid parameter)")
		return
	}
	executer(parse_options())
}