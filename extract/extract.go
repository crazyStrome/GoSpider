package extract

import (
	"log"
	"net/http"
	"path"
	"strings"

	"golang.org/x/net/html"
)

// Content contents the full path url and the base path
type Content struct {
	URL string
	Dir string
}

// Extract return a set of Content for a content input
func Extract(c Content) (res []Content, err error) {
	log.Println("Proceding", c.URL)
	resp, err := http.Get(c.URL)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	res = []Content{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if strings.ToLower(attr.Key) == "href" || strings.ToLower(attr.Key) == "src" {
					if strings.HasPrefix(attr.Val, "http://") || strings.HasPrefix(attr.Val, "https://") {
						res = append(res, Content{
							URL: attr.Val,
							Dir: GetDirURL(attr.Val),
						})
					} else if strings.HasPrefix(attr.Val, "/") {
						res = append(res, Content{
							URL: path.Join(c.Dir, attr.Val),
							Dir: c.Dir,
						})
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
		return
	}
	f(doc)
	log.Println("Proceded", c.URL)
	return
}

// GetDirURL return the base url of input url
// for example input http://www.baidu.com/haha.js
// return http://www.baidu.com
func GetDirURL(full string) string {
	dir := path.Dir(full)
	// dir="http:" || dir="https:"
	if len(dir) == 5 || len(dir) == 6 {
		return full
	}
	return dir
}
