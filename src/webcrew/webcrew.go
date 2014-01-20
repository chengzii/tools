package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mahonia"
	"mydb"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ptnUrl     = regexp.MustCompile(`<a href="(.+)" class="pic_link">`)
	ptnUrlmsg  = regexp.MustCompile(`<img src="(.+)" alt="(.+)" title="(.+)" width="200" height="140">`)
	ptnUrldesc = regexp.MustCompile(`(.+)<a.*详细&gt;&gt;</a>`)
	ptnTag     = regexp.MustCompile(`<div class="news_tag" >(.+)</div>`)
	ptnTaga    = regexp.MustCompile(`<a href=".*">(.+)`)

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
	ptnMete2       = regexp.MustCompile(`<meta.*charset=(.*)".*>`)

	sohuUrl     = regexp.MustCompile(`<h4><a href="(.+)" target="_blank">(.+)</a></h4>`)
	sohuUrlimg  = regexp.MustCompile(`<div class="l"><img src="(.+)" style="width:(.+)px" alt="" /></div>`)
	sohuUrldesc = regexp.MustCompile(`<p>(.*)<a href=".*" target="_blank">阅读全文&gt;&gt;</a></p>`)
	//ptnTag     = regexp.MustCompile(`<div class="news_tag" >(.+)</div>`)
	//ptnTaga    = regexp.MustCompile(`<a href=".*">(.+)`)
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
	url       string
	title     string
	img       string
	tag       string
	desc      string
	date_time int64
}

func htmlToutf8(content string) (res string) {
	mete := ptnMete.FindAllStringSubmatch(content, 1)
	if len(mete) < 1 {
		mete = ptnMete2.FindAllStringSubmatch(content, 1)
	}
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

func findIndex163(content string) (index []IndexItem, err error) {
	matches1 := ptnUrl.FindAllStringSubmatch(content, 20)
	matches2 := ptnUrlmsg.FindAllStringSubmatch(content, 20)
	matches3 := ptnUrldesc.FindAllStringSubmatch(content, 20)
	content = strings.Replace(content, "\r\n", "", -1)
	content = strings.Replace(content, "\n", "", -1)
	content = strings.Replace(content, "<div class=\"news_tag\"", "\n<div class=\"news_tag\"", -1)
	content = strings.Replace(content, "<span class=\"clear\">", "\n<span class=\"clear\">", -1)
	matches4 := ptnTag.FindAllStringSubmatch(content, 20)
	if len(matches1) != len(matches2) || len(matches1) != len(matches3) || len(matches1) != len(matches4) {
		err = errors.New("get length of index error!")
		return
	}
	taginfo := make([]string, len(matches4))
	for ti, tags := range matches4 {
		var info string
		tas := strings.Split(tags[1], "</a>")
		for i := 0; i < len(tas)-1; i++ {
			mta := ptnTaga.FindStringSubmatch(tas[i])
			if info == "" {
				info += mta[1]
			} else {
				info += " " + mta[1]
			}
		}
		taginfo[ti] = info
	}
	//fmt.Println(taginfo, len(taginfo))
	index = make([]IndexItem, len(matches1))
	//index = make([]IndexItem, 1)
	for i := 0; i < len(matches1); i++ {
		index[i] = IndexItem{matches1[i][1], matches2[i][2], matches2[i][1], taginfo[i], strings.Trim(matches3[i][1], "\t|\n|\n\r| "), getDateTime(matches2[i][1])}
		//break
	}
	return
}
func findIndexsohu(content string) (index []IndexItem, err error) {
	matches1 := sohuUrl.FindAllStringSubmatch(content, 7)
	matches2 := sohuUrlimg.FindAllStringSubmatch(content, 7)
	matches3 := sohuUrldesc.FindAllStringSubmatch(content, 7)
	if len(matches1) != len(matches2) || len(matches1) != len(matches3) {
		err = errors.New("get length of index error!")
		//fmt.Println(len(matches1), len(matches2), len(matches3))
		return
	}
	tag := "餐饮 餐处"
	//fmt.Println(taginfo, len(taginfo))
	index = make([]IndexItem, len(matches1))
	//index = make([]IndexItem, 1)
	for i := 0; i < len(matches1); i++ {
		index[i] = IndexItem{matches1[i][1], matches1[i][2], matches2[i][1], tag, strings.Trim(matches3[i][1], "\t|\n|\n\r| "), getDateTimeSohu(matches2[i][1])}
		//fmt.Println(matches2[i][1], getDateTimeSohu(matches2[i][1]))
		//break
	}
	return
}
func getDateTime(imgurl string) (unixtime int64) {
	arr := strings.Split(imgurl, "/")
	year, err := strconv.Atoi(arr[4])
	month, err := strconv.Atoi(arr[5])
	day, err := strconv.Atoi(arr[6])
	if len(arr) < 7 || err != nil {
		unixtime = time.Now().Unix()
	} else {
		unixtime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Unix()
		//time.Date(year, month, day, hour, min, sec, nsec, loc)
	}
	return
}
func getDateTimeSohu(imgurl string) (unixtime int64) {
	arr := strings.Split(imgurl, "/")
	if len(arr) < 4 {
		unixtime = time.Now().Unix()
	} else {
		timeflag := arr[3]

		year, err := strconv.Atoi(timeflag[0:4])
		month, err := strconv.Atoi(timeflag[4:6])
		day, err := strconv.Atoi(timeflag[6:8])
		//fmt.Println(year, month, day)
		if err != nil {
			unixtime = time.Now().Unix()
		} else {
			unixtime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Unix()
		}
		//time.Date(year, month, day, hour, min, sec, nsec, loc)
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
	fmt.Println(`Get 163 index ...`)
	s163, statusCode := Get("http://travel.163.com/food/")
	if statusCode != 200 {
		return
	}
	index, _ := findIndex163(s163)
	fmt.Println(`Get contents and write to file ...`)
	//fmt.Println(index)
	for i := len(index) - 1; i >= 0; i-- {
		item := index[i]
		if item.date_time > time.Now().Unix()-24*3600 {
			mydb.Insert_fd(item.title, item.desc, item.url, item.img, item.tag, item.date_time)
		}
		/*
			fmt.Println("标题: ", item.title)
			fmt.Println("简述: ", item.desc)
			fmt.Println("链接: ", item.url)
			fmt.Println("图片: ", item.img)
			fmt.Println("标签: ", item.tag)
			fmt.Println("日期: ", item.date_time)
		*/
	}
	fmt.Println(`Get sohu index ...`)
	ssohu, statusCode := Get("http://chihe.sohu.com/restaurant/")
	if statusCode != 200 {
		fmt.Println(statusCode)
		return
	}
	indexs, _ := findIndexsohu(ssohu)
	fmt.Println(`Get contents and write to file ...`)
	//fmt.Println(indexs)
	for j := len(indexs) - 1; j >= 0; j-- {
		item := indexs[j]
		if item.date_time > time.Now().Unix()-24*3600 {
			mydb.Insert_fd(item.title, item.desc, item.url, item.img, item.tag, item.date_time)
		}
		/*
			fmt.Println("标题: ", item.title)
			fmt.Println("简述: ", item.desc)
			fmt.Println("链接: ", item.url)
			fmt.Println("图片: ", item.img)
			fmt.Println("标签: ", item.tag)
			fmt.Println("日期: ", item.date_time)
		*/
	}
}
