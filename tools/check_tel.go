package tools

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"github.com/garyburd/redigo/redis"
)

func IsShanghai(phone string) bool {

	c , err := redis.Dial("tcp","127.0.0.1:6379")

	if err != nil {

		panic(err.Error())
	}

	defer  c.Close()

	//c.Do("AUTH","xxxx")
	isbool, err := redis.String(c.Do("GET",phone))

	if isbool == "0" {
		return true
	} else if isbool == "1" {
		return false
	}


	client := &http.Client{}

	url := "https://ip.cn/db.php?num=" + phone

	reqest, err := http.NewRequest("GET", url, nil)

	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Mobile Safari/537.36")

	reqest.Header.Add("upgrade-insecure-requests","1")

	if err != nil {
		panic(err)
	}

	response, _ := client.Do(reqest)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		panic(err)
	}

	tell := doc.Find("#result").Text()

	x, _ := regexp.Compile("上海")

	result := x.FindString(tell)


	if result == "上海" {
		c.Do("SET",phone,"0")
		return true
	}else {
		c.Do("SET",phone,"1")
		return false
	}
}
