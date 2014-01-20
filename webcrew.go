package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mahonia"
	"net/http"
	"regexp"
	"strings"
)

var (
	ptnUrl         = regexp.MustCompile(`<a href="(.+)" class="pic_link">`)
	ptnUrlmsg      = regexp.MustCompile(`<img src="(.+)" alt="(.+)" title="(.+)" width="200" height="140">`)
	ptnUrldesc     = regexp.MustCompile(`(.+)<a.*详细&gt;&gt;</a>`)
	ptnContentFlag = regexp.MustCompile(`<!-- 正文内容 begin -->(.+)<!-- 正文内容 end -->`)
	ptnBrTag       = regexp.MustCompile(`<br>`)
	ptnHTMLTag     = regexp.MustCompile(`(?s)</?.*?>`)
	ptnSpace       = regexp.MustCompile(`(^\s+)|( )`)
	ptnTab         = regexp.MustCompile(`\t|\n|\n\r`)
	ptnJsTag       = regexp.MustCompile(`<script.+</script>`)
	ptnCssTag      = regexp.MustCompile(`<style.+</style>`)
	ptnLinkTag     = regexp.MustCompile(`<p class='page'.+</p>`)
	ptnImgTag      = regexp.MustCompile(`<div class="img_wrapper"(.[^</div>]*)</div>`)
	ptnMete        = regexp.MustCompile(`<meta.*charset="(.*)".*>`)
)

func Get(url string) (content string, statusCode int) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		statusCode = -100
		return
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	statusCode = resp.StatusCode
	content = htmlToutf8(string(data))
	return
}

type IndexItem struct {
	url   string
	title string
	img   string
	//tag   string
	desc string
}

func htmlToutf8(content string) (res string) {
	mete := ptnMete.FindAllStringSubmatch(content, 1)
	charset := mete[0][1]
	if charset == "" {
		return
	}
	decoder := mahonia.NewDecoder("GBK")
	//encoder:= mahonia.NewEncoder("UTF-8")
	if charset == "utf-8" || charset == "UTF-8" || charset == "utf8" || charset == "utf-8" {
		res = content
	} else {
		content_g := decoder.ConvertString(content)
		res = content_g
	}
	return
}

func findIndex(content string) (index []IndexItem, err error) {
	//content = ptnTab.ReplaceAllString(content, "")
	matches1 := ptnUrl.FindAllStringSubmatch(content, 20)
	matches2 := ptnUrlmsg.FindAllStringSubmatch(content, 20)
	matches3 := ptnUrldesc.FindAllStringSubmatch(content, 20)
	if len(matches1) != len(matches2) || len(matches1) != len(matches3) {
		err = errors.New("get length of index error!")
		return
	}
	index = make([]IndexItem, len(matches1))
	for i := 0; i < len(matches1); i++ {
		index[i] = IndexItem{matches1[i][1], matches2[i][2], matches2[i][1], strings.Trim(matches3[i][1], "\t|\n|\n\r| ")}
	}
	return
}
func readContent(url string) (content string) {
	raw, statusCode := Get(url)
	if statusCode != 200 {
		fmt.Print("Fail to get the raw data from", url, "\n")
		return
	}
	raw = strings.Replace(raw, "\r\n", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	match := ptnContentFlag.FindStringSubmatch(raw)
	if match != nil {
		if len(match) >= 2 {
			content = match[1]
		}
	} else {
		return
	}

	content = ptnBrTag.ReplaceAllString(content, "\r\n")
	content = ptnJsTag.ReplaceAllString(content, "")
	content = ptnCssTag.ReplaceAllString(content, "")
	content = ptnImgTag.ReplaceAllString(content, "")
	content = ptnLinkTag.ReplaceAllString(content, "")
	content = ptnHTMLTag.ReplaceAllString(content, "")
	content = ptnSpace.ReplaceAllString(content, "")
	return
}

func main() {
	fmt.Println(`Get index ...`)
	s, statusCode := Get("http://travel.163.com/food/")
	if statusCode != 200 {
		return
	}
	index, _ := findIndex(s)
	fmt.Println(`Get contents and write to file ...`)
	//fmt.Println(index)
	for _, item := range index {
		fmt.Println("标题: ", item.title)
		fmt.Println("简述: ", item.desc)
		fmt.Println("链接: ", item.url)
		fmt.Println("图片: ", item.img)
		/*
			fmt.Printf("Get content %s from %s and write to file.\n", item.title, item.url)
			fileName := fmt.Sprintf("%s.html", item.title)
			content := readContent(item.url)
			ioutil.WriteFile(fileName, []byte(content), 0644)
			fmt.Printf("Finish writing to %s.\n", fileName)
			break
		*/
	}
}
