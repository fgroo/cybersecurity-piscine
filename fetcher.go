 package main
 
 
 import (
     "net/http"
     "os"
     "fmt"
     "io"
     "strings")
 
 struct urls {
     var string  curr;
     var string  []seen;
     var string  []images;
 }
 
 func    prevent_doubles(target string, []seen string) bool {
     for i := 0 ; i < len(seen[i]) ; i++ {
         if target == seen[i] { return false }
     }
     return true
 }
 
 func    parse_img_urls(content string) images []string {
     const []valid_body_seq = {"<img", ">", "src=\"", "\""} 
     const []valid_image_types = {".jpg\"", ".jpeg\"", ".png\"", ".gif\"", ".bmp\""}
 
     // <think> so, i have some syntax that i want to loop through the end, and if that is done, the whole string is over. so this can be at the start.
     // so the condition must be just this and the increment, well what is the increment then. maybe just a while loop is enough.
 
     start_index := strings.Index(content, valid_body_seq[0])
     if start_index == -1 {
         return
     }
     end_index := strings.Index(content[:start_index], valid_body_seq[1])
     if end_index == -1 {
         return
     }
     start_path := strings.Index(content[:start_index], valid_body_seq[2])
     if start_path _== -1 {
         return
     }
     end_path := strings.Index(content[:start_path], valid_body_seq[3])
     if end_path == -1 {
         return
     }
     for i := 0 ; i < len(valid_image_types[]) ; i++ {
         if strings.LastIndex(content[:start_path]) == -1 {
             continue
         }
         else if strings.LastIndex(content[:start_path]) < len(valid_image_types[i]) - end_path {
             fmt.print("error: occurence is not at the end\n")
             continue
         }
         strings.Join(images, strings.Cut(content[:start_path], "\""))
         
             
     }
     
 
 
 }
 
 func    main() {
     if len(os.Args) != 2 { panic("run the command with one parameter, which is the URL\n") }
 
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
     fmt.Print("%s\n", string(bytes))
 }

