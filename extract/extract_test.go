package extract

import (
	"fmt"
	"testing"
)

func TestExtract(t *testing.T) {
	res := Extract(Content{
		URL: "http://www.zhihu.com",
		Dir: "http://www.zhihu.com",
	})
	for _, c := range res {
		t.Log(c)
	}
}
func TestGetDirURL(t *testing.T) {
	urls := []string{"http:/www.zhihu.com/roundtable",
		"http:/www.zhihu.com/explore",
		"http:/www.zhihu.com/app",
		"http:/www.zhihu.com/contact",
		"https://app.mokahr.com/apply/zhihu",
		"http:/www.zhihu.com/org/signup",
		"https://tsm.miit.gov.cn/dxxzsp/",
		"http://www.beian.miit.gov.cn",
		"http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=11010802020088",
		"https://pic3.zhimg.com/80/v2-d0289dc0a46fc5b15b3363ffa78cf6c7.png"}
	for _, u := range urls {
		d := GetDirURL(u)
		fmt.Println("---------------------")
		fmt.Println("Origin:", u)
		fmt.Println("Dir:", d)
	}
}
