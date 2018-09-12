package main

import (
	"git.nevint.com/icall/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"os"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"strconv"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	session := tools.MgoConn()

	defer session.Close()

        router := gin.Default()

	router.GET("/",func(c *gin.Context) {
			c.String(http.StatusOK,"hello world")
	})

	router.POST("/play/",playHandlerFunc)

	router.GET("/selid/",selHandlerFunc)

        router.Run("0.0.0.0:8999")

}

type videos struct {
	Tex	string `json:"tex" binding:"required"`
	Number string `json:"number" binding:"required"`
	Spd string `json:"spd"`
	Pit string `json:"pit"`
	Vol string 	`json:"vol"`
	Per string 	`json:"per"`
	Timeout string `json:"timeout"`
}


type User struct {
        Teamname string
        Token string
        Tel string
}


func (this User) AuthLogin(token, tname string) bool {
	if this.Teamname == tname && this.Token == token {
		return true
	} else {
		return  false
	}
}


func (this User) SelTell() string {
	return this.Tel
}


func playHandlerFunc(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	tname := c.Request.Header.Get("TeamName")

	video := videos{
		Spd:"5",
		Pit:"5",
		Vol:"5",
		Per:"0",
		Timeout:"60",
	}

	err := c.BindJSON(&video)

	if err != nil {
		panic(err)
	}


	if len(video.Number) != 11 {
		c.String(401,"电话号码不足11位")
		return
	}


	session := tools.MgoConn()

	defer session.Close()

	m := session.DB("icall").C("author")

	var user User

	m.Find(bson.M{"token":token}).One(&user)

	islogin := user.AuthLogin(token,tname)

	if !islogin {
		c.String(http.StatusForbidden, "Access permission denied")
		return
	} else {

		rtime := time.Now().Unix()

		Mp3Handler, err :=  os.Open("/opt/video/example.mp3")

		if err != nil {
			panic(err.Error())
		}
		video_res, err := tools.GenVideo(video.Tex, video.Spd,video.Pit,video.Vol,video.Per)

		temp_mp3 := fmt.Sprintf("/opt/video/%d.mp3",rtime)

		TempHandler, err := os.OpenFile(temp_mp3, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

		if err != nil {
			panic(err)
		}

		mp3, _ := ioutil.ReadAll(Mp3Handler)
		if _, err := TempHandler.Write(mp3); err != nil {
			panic(err.Error())
		}

		if _, err := TempHandler.Write(video_res); err != nil {
			panic(err.Error())
		}


		defer Mp3Handler.Close()

		defer TempHandler.Close()

		call_number , _ := strconv.Atoi(user.SelTell())

		ishanghai := tools.IsShanghai(video.Number)

		cmd := exec.Command("fs_cli","-x  bgapi status ")
		output, err := cmd.Output()

		str_spl := strings.Split(string(output),":")

		result := str_spl[1]

		result1 := strings.TrimSpace(result)

		go tools.PlayCall(video.Number,temp_mp3,video.Timeout, result1,video.Tex,tname,ishanghai,call_number)

		c.JSON(http.StatusOK, gin.H{"jbid":result})
	}

}

func selHandlerFunc(c *gin.Context) {

	jbid := c.Query("jbid")

	if len(jbid) == 0 {
		c.String(http.StatusOK,"jbid not null")
		return
	}

	session := tools.MgoConn()

        defer session.Close()

        m := session.DB("icall").C("jbinfo")

        var mjob tools.MjobInfo

        m.Find(bson.M{"jbid":jbid}).One(&mjob)

	fmt.Println(mjob)

	if len(mjob.Body) == 0 {
		c.JSON(http.StatusOK,gin.H{"result":"jbid not found"})
		return
	} else {
		c.JSON(http.StatusOK,gin.H{"result":mjob.Body})
	}

}