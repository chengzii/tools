package urltocontent

import (
	"encoding/xml"
	"io/ioutil"
	"mahonia"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	//"log"
)

type Lbsresult struct {
	XMLName    xml.Name   `xml:"LbsResult"`
	LQII       Lqii       `xml:"LQII"`
	KeywordRes Keywordres `xml:"KeywordRes"`
	LocRes     LocRes     `xml:"LocRes"`
}
type Lqii struct {
	LQIIVersion string     `xml:"LQIIVersion"`
	LQIIPraser  string     `xml:"LQIIPraser"`
	QUERYTYPE   string     `xml:"QUERYTYPE"`
	Searchtime  string     `xml:"searchtime"`
	ServerAddr  Serveraddr `xml:"ServerAddr"`
}
type Serveraddr struct {
	QiiAddr   string `xml:"QiiAddr"`
	LseAddr   string `xml:"LseAddr"`
	GeoAddr   string `xml:"GeoAddr"`
	RegeoAddr string `xml:"RegeoAddr"`
}
type Keywordres struct {
	Searchtype string `xml:"searchtype"`
	Keyword    string `xml:"keyword"`
	Searchtime string `xml:"searchtime"`
	List       List   `xml:"list"`
}
type LocRes struct {
	Loctype      string       `xml:"loctype"`
	Keyword      string       `xml:"keyword"`
	Searchtime   string       `xml:"searchtime"`
	Locationlist Locationlist `xml:"locationlist"`
}
type List struct {
	Searchresult Searchresult `xml:"searchresult"`
}
type Locationlist struct {
	Searchresult Searchresultloc `xml:"searchresult"`
}
type Searchresult struct {
	SEngineVersion  string             `xml:"SEngineVersion"`
	Count           int                `xml:"count"`
	Searchtime      int                `xml:"searchtime"`
	Pinyin          Pinyin             `xml:"pinyin"`
	Bounds          string             `xml:"bounds"`
	Citysuggestions Citysuggestionlist `xml:"citysuggestions"`
	Query_info      Queryinfo          `xml:"query_info"`
	List            Poilist            `xml:"list"`
}
type Searchresultloc struct {
	Count        int     `xml:"count"`
	Spellcorrect string  `xml:"spellcorrect"`
	List         Poilist `xml:"list"`
}
type Pinyin struct {
	List Type `xml:"list"`
}
type Type struct {
	List string `xml:"type,attr"`
}
type Citysuggestionlist struct {
	Citysuggestionlisttype string `xml:"type,attr"`
	Citysuggestion         []Poi  `xml:"citysuggestion"`
}
type Queryinfo struct {
	Geoword bool `xml:"geoword"`
}
type Poilist struct {
	Poilisttype string `xml:"type,attr"`
	Addr_poi    []Poi  `xml:"addr_poi"`
	Busline     []Poi  `xml:"busline"`
	Poi         []Poi  `xml:"poi"`
}
type Poi struct {
	Pguid      string `xml:"pguid"`
	Name       string `xml:"name"`
	Adcode     string `xml:"adcode"`
	Citycode   string `xml:"citycode"`
	Address    string `xml:"address"`
	Typecode   string `xml:"typecode"`
	Ptype      string `xml:"type"`
	X          string `xml:"x"`
	Y          string `xml:"y"`
	Gridcode   string `xml:"gridcode"`
	Poiweight  string `xml:"poiweight"`
	Distance   string `xml:"distance"`
	Last_match Match  `xml:"last_match"`
}
type Match struct {
	Match_string string `xml:"match_string"`
}

func Url2string(url string) (content string) {
	var recontent string
	pois, err, count := Geturl(url)
	if err == nil {
		num := len(pois)
		if num == 0 {
			decoder := mahonia.NewDecoder("GBK")
			encoder := mahonia.NewEncoder("UTF-8")
			url = decoder.ConvertString(url)
			url = encoder.ConvertString(url)
			recontent += "poicount:" + strconv.Itoa(count) + "\t" + url + "\t"
		} else {
			recontent += "poicount:" + strconv.Itoa(count) + "\t"
		}
		for i := 0; i < num; i++ {
			poi, ok := pois[i]
			if ok {
				recontent += strconv.Itoa(i+1) + "=>"
				if poi.Name != "" {
					recontent += "name:" + poi.Name
				} else {
					recontent += " "
					break
				}
				if poi.Pguid != "" {
					recontent += "pguid:" + poi.Pguid
				}
				if poi.Address != "" {
					recontent += "address:" + poi.Address
				}
				if poi.Typecode != "" {
					recontent += "typecode:" + poi.Typecode
				}
				if poi.Last_match.Match_string != "" {
					recontent += "match_string:" + poi.Last_match.Match_string
				}
				recontent += ";"
			}
		}
	} else {
		recontent += "error"
	}
	return recontent
}
func Geturl(url string) (repoi map[int]Poi, reerr error, count int) {
	resp, err := http.Get(url)
	poi := make(map[int]Poi)
	decoder := mahonia.NewDecoder("gbk")
	encoder := mahonia.NewEncoder("UTF-8")
	if err == nil {
		re, err2 := ioutil.ReadAll(resp.Body)
		if err2 == nil {
			html := decoder.ConvertString(string(re))
			uhtml := encoder.ConvertString(html)
			content := strings.Replace(uhtml, "GBK", "utf-8", 1)
			poi, err, count = xmlana(content)
		} else {
			err = err2
		}
	}
	return poi, err, count
}
func xmlana(content string) (repoi map[int]Poi, reerr error, num int) {
	var result Lbsresult
	var searchresult Searchresult
	var count int
	poi := make(map[int]Poi)
	index := 0
	err1 := xml.Unmarshal([]byte(content), &result)
	err2 := xml.Unmarshal([]byte(content), &searchresult)
	if err1 == nil {
		count = result.KeywordRes.List.Searchresult.Count
		if count == 0 {
			count = result.LocRes.Locationlist.Searchresult.Count
		}
		for k, v := range result.KeywordRes.List.Searchresult.Citysuggestions.Citysuggestion {
			v.Pguid = "SUGGESTION"
			poi[k] = v
		}
		for _, v := range result.KeywordRes.List.Searchresult.List.Addr_poi {
			poi[index] = v
			index++
		}
		for _, v := range result.KeywordRes.List.Searchresult.List.Busline {
			poi[index] = v
			index++
		}
		for _, v := range result.KeywordRes.List.Searchresult.List.Poi {
			poi[index] = v
			index++
		}
		for k, v := range result.LocRes.Locationlist.Searchresult.List.Poi {
			v.Pguid = "LOCATION00"
			poi[k] = v
		}
	} else if err2 == nil {
		err1 = err2
		count = searchresult.Count
		for _, v := range searchresult.List.Busline {
			poi[index] = v
			index++
		}
		for _, v := range searchresult.List.Poi {
			poi[index] = v
			index++
		}
	}
	return poi, err1, count
}
func Getpoicolumn(poi Poi, column string) (flag bool, result string) {
	flag = false
	s := reflect.ValueOf(&poi).Elem()
	typeOfj := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if strings.ToUpper(column) == strings.ToUpper(typeOfj.Field(i).Name) {
			if a, ok := f.Interface().(string); ok {
				result = a
				flag = true
			}
		}
	}
	return
}
