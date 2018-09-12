package tools

import (
	"github.com/0x19/goesl"
	"strings"
	"fmt"
	"regexp"
	"time"
	"encoding/json"
)


const  (
	fshost = "127.0.0.1"
	fsport = 8021
	fspass = "ClueCon"
	fstimeout = 5
)


func PlayCall(phone, video, timeout, jbid , tex, tname string, isshanghai bool, call_number int) {
	c , err := goesl.NewClient(fshost,fsport,fspass,fstimeout)
	if err != nil {
		goesl.Error("Error while creating new client: %s", err)
	}

	defer c.Close()

	go c.Handle()


	session := MgoConn()

	defer session.Close()


	c.Send("events json BACKGROUND_JOB")


	if isshanghai {
		p := fmt.Sprintf("originate {origination_uuid=%s,ignore_early_media=true,originate_timeout=%s,origination_caller_id_number=%d}sofia/gateway/aok/%s &playback(%s)",jbid,timeout,call_number,phone,video)
		err := c.BgApi(p)

		if err != nil {
			panic(err)
		}
	} else {
		 p := fmt.Sprintf("originate {origination_uuid=%s,ignore_early_media=true,originate_timeout=%s,origination_caller_id_number=%d}sofia/gateway/aok/0%s &playback(%s)",jbid,timeout,call_number,phone,video)
		err := c.BgApi(p)

		if err != nil {
			panic(err)
		}
	}



	for {

		msg, err := c.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				goesl.Error("Error while reading  message: %s", err)
			}
		}


		body := msg.Body



		if len(body) != 0 {
			r ,_:= regexp.Compile("OK")
			switch {
			case r.FindString(string(body)) == "OK" :
				timeStr := time.Now().Format("2006-01-02 15:04:05")
				body := string(body)

				res := MjobInfo{
					Body: body,
					Context: tex,
					TeamName:tname,
					Time: timeStr,
					Phone:phone,
					Status: 0,
					Jbid: jbid,
				}

				json.Marshal(res)
				d := session.DB("icall").C("jbinfo")
				d.Insert(res)


			default:
				timeStr := time.Now().Format("2006-01-02 15:04:05")
				body := string(body)
				res := MjobInfo{
					Body: body,
					Context: tex,
					TeamName:tname,
					Time: timeStr,
					Phone:phone,
					Status: 1,
					Jbid: jbid,
				}

				json.Marshal(res)
				d := session.DB("icall").C("jbinfo")
				d.Insert(res)

			}
			goesl.Debug("body is ",string(body))
		}

		goesl.Debug("Got new message: %s", msg)

	}
}


