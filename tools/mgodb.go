package tools

import "gopkg.in/mgo.v2"

type MjobInfo struct {
	Body string	`json:"body"`
	Time string	`json:"time"`
	Status int	`json:"status"`
	Jbid string	`json:"jbid"`
	Context string `json:"context"`
	Phone string `json:"phone"`
	TeamName string `json:"team_name"`
}



func MgoConn() *mgo.Session {
	session, err := mgo.Dial("mongodb://xxxx:xxxx@127.0.0.1:27019")

	if err != nil {
		panic(err.Error())
	}

	return session

}


