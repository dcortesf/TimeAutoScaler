package main

import (
	"fmt"
	"time"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TimeRule struct {
        Hour string
				ToHour string
        DeploymentConfig string
				scale int
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
	err = c.Insert(&TimeRule{"12:34","13:50", "MyDeploymentConfig",2},
					 &TimeRule{"13:30","15:00", "aaa",3})
	if err != nil {
					log.Fatal(err)
	}

	result := TimeRule{}


	for {
		time.Sleep(1000 * time.Millisecond)
		currentTime := time.Now()
		fmt.Println(currentTime)
		err = c.Find(bson.M{"hour": "12:34"}).One(&result)
		if err != nil {
						log.Fatal(err)
		}

		fmt.Println("Rule:", result.DeploymentConfig)
	}
}

func main() {

	//go mytime()
  mytime()
}
