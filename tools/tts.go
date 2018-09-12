package tools

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

const (
	API_KEY = "xxxx"
	SECRET_KEY = "xxxxxx"
	TOKEN_URL = "xxxxxx"
	VIDEO_URL = "xxxxx"
)


func GenVideo(text,spd,pit,vol,per string) (response []uint8, err error) {
	var rst map[string]interface{}

	rq , err := http.Get(TOKEN_URL + "?grant_type=client_credentials&client_id=" + API_KEY + "&client_secret=" +  SECRET_KEY)

	if err != nil {
		panic(err)
	}

	defer rq.Body.Close()

	resp, _ := ioutil.ReadAll(rq.Body)

	json.Unmarshal(resp, &rst)
	token := rst["access_token"]

	parma := fmt.Sprintf("lan=zh&ctp=1&cuid=123456&tok=%s&per=%s&spd=%s&pit=%s&vol=%s&tex=%s",token,per,spd,pit,vol,text)

	content_resp ,err := http.Post(VIDEO_URL, "application/x-www-form-urlencoded", strings.NewReader(parma))


	if err != nil {
		panic(err)
	}

	
	response, err = ioutil.ReadAll(content_resp.Body)


	return response, err
}

