package main

import (
	"bufio"
	"fmt"
	"io"
	"mahonia"
	"net/url"
	"os"
	"strings"
	"urltocontent"
	//"strconv"
)

type Res struct {
	Keyword string
	City    string
	Freq    string
}

func main() {
	file := "lse2.log"
	re, err := getFileContent(file)
	if err != nil {
		return
	}
	uri := geturi("")
	for _, v := range re {
		vurl := v.Geturl(uri)
		fmt.Println(urltocontent.Geturl(vurl))
		//break
	}
}
func getFileContent(file string) (re []Res, err error) {
	fi, err := os.Open(file)
	defer fi.Close()
	if err != nil {
		return
	}
	r := bufio.NewReader(fi)
	decoder := mahonia.NewDecoder("GBK")
	for {
		linebyte, _, err := r.ReadLine()
		if err != nil || err == io.EOF {
			break
		}
		line := decoder.ConvertString(string(linebyte))
		lines := strings.Split(line, "\t")
		if len(lines) < 3 {
			continue
		}
		re = append(re, Res{lines[0], lines[1], lines[2]})
	}
	return
}
func geturi(key string) (uri string) {
	if key == "amap_mse" {
		uri = "http://10.13.4.25:2099/bin/sp?data_type=POI&query_type=TQUERY&city=010&keywords=1&category=&field=&page_num=10&page=1&center=&x=&y=&range=3000&id=&sstype=0&geoobj=&geotype=&sortrule=0&districtname=&templateid=&extid=&poiid=&userid=&eid=&custom=&custom_and=true&except=&result_filter=&sort_rule=0&sort_fields=&xy_cluster=false&show_group=&group_by=&group_by_parent=true&relative_score_filter=0.5&show_score=false&show_fields=&query_config=&auto_resize=true&server_ip=10.13.4.25&server_port=13333&qii_server_port=14001&qii=true&show_uuid=true&show_task_code=true&log_server_ip=127.0.0.1&log_server_port=13339&use_log=false&addr_poi_merge=true&app=fake_app&outfmt=pb&query_src=test&user_info=test"
	} else if key == "amap_lse2" {
		uri = "http://10.13.4.25:2099/bin/sp?data_type=POI&query_type=TQUERY&city=010&keywords=1&category=&field=&page_num=10&page=1&center=&x=&y=&range=3000&id=&sstype=0&geoobj=&geotype=&sortrule=0&districtname=&templateid=&extid=&poiid=&userid=&eid=&custom=&custom_and=true&except=&result_filter=&sort_rule=0&sort_fields=&xy_cluster=false&show_group=&group_by=&group_by_parent=true&relative_score_filter=0.5&show_score=false&show_fields=&query_config=&auto_resize=true&server_ip=10.13.4.25&server_port=13333&qii_server_port=14001&qii=true&show_uuid=true&show_task_code=true&log_server_ip=127.0.0.1&log_server_port=13339&use_log=false&addr_poi_merge=true&app=fake_app&outfmt=pb&query_src=test&user_info=test"
	} else if key == "snowman" {
		uri = "http://10.13.4.25:2099/bin/sp?data_type=POI&query_type=TQUERY&city=010&keywords=1&category=&field=&page_num=10&page=1&center=&x=&y=&range=3000&id=&sstype=0&geoobj=&geotype=&sortrule=0&districtname=&templateid=&extid=&poiid=&userid=&eid=&custom=&custom_and=true&except=&result_filter=&sort_rule=0&sort_fields=&xy_cluster=false&show_group=&group_by=&group_by_parent=true&relative_score_filter=0.5&show_score=false&show_fields=&query_config=&auto_resize=true&server_ip=10.13.4.25&server_port=13333&qii_server_port=14001&qii=true&show_uuid=true&show_task_code=true&log_server_ip=127.0.0.1&log_server_port=13339&use_log=false&addr_poi_merge=true&app=fake_app&outfmt=pb&query_src=test&user_info=test"
	} else if key == "chelianwang" {
		uri = "http://10.13.4.25:2099/bin/sp?data_type=POI&query_type=TQUERY&city=010&keywords=1&category=&field=&page_num=10&page=1&center=&x=&y=&range=3000&id=&sstype=0&geoobj=&geotype=&sortrule=0&districtname=&templateid=&extid=&poiid=&userid=&eid=&custom=&custom_and=true&except=&result_filter=&sort_rule=0&sort_fields=&xy_cluster=false&show_group=&group_by=&group_by_parent=true&relative_score_filter=0.5&show_score=false&show_fields=&query_config=&auto_resize=true&server_ip=10.13.4.25&server_port=13333&qii_server_port=14001&qii=true&show_uuid=true&show_task_code=true&log_server_ip=127.0.0.1&log_server_port=13339&use_log=false&addr_poi_merge=true&app=fake_app&outfmt=pb&query_src=test&user_info=test"
	} else {
		uri = "http://10.13.4.25:2099/bin/sp?data_type=POI&query_type=TQUERY&city=010&keywords=1&category=&field=&page_num=10&page=1&center=&x=&y=&range=3000&id=&sstype=0&geoobj=&geotype=&sortrule=0&districtname=&templateid=&extid=&poiid=&userid=&eid=&custom=&custom_and=true&except=&result_filter=&sort_rule=0&sort_fields=&xy_cluster=false&show_group=&group_by=&group_by_parent=true&relative_score_filter=0.5&show_score=false&show_fields=&query_config=&auto_resize=true&server_ip=10.13.4.25&server_port=13333&qii_server_port=14001&qii=true&show_uuid=true&show_task_code=true&log_server_ip=127.0.0.1&log_server_port=13339&use_log=false&addr_poi_merge=true&app=fake_app&outfmt=pb&query_src=test&user_info=test"
	}
	return
}
func (r Res) Geturl(uri string) (reurl string) {
	encoder := mahonia.NewEncoder("GBK")
	wd := encoder.ConvertString(r.Keyword)
	wd = url.QueryEscape(wd)
	kvs := strings.Split(uri, "&")
	for _, v := range kvs {
		//fmt.Println(v)
		if reurl != "" {
			reurl = reurl + "&"
		}
		vs := strings.SplitN(v, "=", 2)
		//fmt.Println(vs)
		if len(vs) == 2 {
			if strings.ToUpper(vs[0]) == "KEYWORDS" {
				reurl = reurl + vs[0] + "=" + wd
			} else if strings.ToUpper(vs[0]) == "CITY" {
				reurl = reurl + vs[0] + "=" + r.City
			} else {
				reurl = reurl + vs[0] + "=" + vs[1]
			}
		} else {
			reurl = reurl + v
		}
		//fmt.Println(reurl)
	}
	return
}

/*
func urlcode(inurl string) (reurl string){
        var varr []string
        for _, v := range strings.Split(inurl, "&") {
                varr = strings.Split(v, "=")
                if varr[0]=="keywords" {
                        reurl+=varr[0]+"="+url.QueryEscape(varr[1])+"&"
                }else{
                        reurl+= v + "&"
                }
        }
        return
}*/
