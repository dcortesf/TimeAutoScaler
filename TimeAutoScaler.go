package main

import (
	"fmt"
	"time"
	"log"
	//"encoding/json"
  //"github.com/bitly/go-simplejson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TimeRule struct {
        Hour string
        DeploymentConfig string
				Scale int
}

func mytime() {

	session, err := mgo.Dial("192.168.99.100:27017")
	if err != nil {
					panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("rules").C("hours")
	err = c.Insert(&TimeRule{"12:34", "App1", 2},
					 &TimeRule{"13:30","App2",-1},&TimeRule{"12:34", "App3", -1})
	if err != nil {
					log.Fatal(err)
	}

	//result := TimeRule{}
	cResult :=[]TimeRule{}


	for {

		time.Sleep(1000 * time.Millisecond)
		currentTime := time.Now()
		fmt.Println(currentTime)
		tt := currentTime.String()
		fmt.Println("En String",tt)
		ttt := tt[:5]
		fmt.Println(ttt)

		//mm := currentTime[0:16]
		//fmt.Println(mm)
		err = c.Find(bson.M{"hour": "12:34"}).All(&cResult)
		if err != nil {
						log.Fatal(err)
		}

		fmt.Println(len(cResult))

		for _, item := range cResult {
        fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.DeploymentConfig, item.Scale)
    }

		//fmt.Println("Afected DCs:", result.DeploymentConfigs)
		//fmt.Println("Afected DCs:", cResult[0].DeploymentConfigs)
	}
}

func main() {

	//go mytime()
  mytime()
}
