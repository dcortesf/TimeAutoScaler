package main

import (
	"fmt"
	"time"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type TimeRule struct {
        Id string
        DeploymentConfig string
				Scale int
}

func mainLoop() {

	fmt.Println("[INFO] - Initialize autoscaler")

	host := os.Getenv("AUTOSCALER_DB_HOST")
	port:= os.Getenv("AUTOSCALER_DB_PORT")

	if(host==""){
		fmt.Println("[INFO] - Cannot find host and port ... stablished 192.168.99.100:27017 by default")
		host="192.168.99.100"
		port="27017"
	}

	session, err := mgo.Dial(host+":"+port)
	if err != nil {
					panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)


// Example Data.
	c := session.DB("rules").C("hours")
	cc := session.DB("rules").C("week")

	/*	err = c.Insert(&TimeRule{"12:34", "App1", 2},
					 &TimeRule{"13:30","App2",-1},&TimeRule{"16:28", "App3", -1})
				 if err != nil {
						log.Fatal(err)
					}

		err = cc.Insert(&TimeRule{"12:34", "WeekApp1", 2},
						 &TimeRule{"13:30","WeekApp2",-1},&TimeRule{"16:28", "WeekApp3", -1})
		if err != nil {
						log.Fatal(err)
		}*/

	//result := TimeRule{}
	cResult :=[]TimeRule{}

	for {

		time.Sleep(1000 * 60 * time.Millisecond)
		currentTime := time.Now()
		_CurrentTime := currentTime.String()
		hourNow := _CurrentTime[11:16]
		dayMonth := _CurrentTime[5:11]


		fmt.Println("Complete trace: ",_CurrentTime)
		fmt.Println("Week day : ", currentTime.Weekday().String())
		fmt.Println("Hour: ",hourNow)
		fmt.Println("DayMonth: ",dayMonth)

		err = c.Find(bson.M{"id": hourNow}).All(&cResult)
		if err != nil {
						log.Fatal(err)
		}

		fmt.Println(len(cResult))

		for _, item := range cResult {
        fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.DeploymentConfig, item.Scale)
    }

		////////
		err = cc.Find(bson.M{"id": hourNow}).All(&cResult)
		if err != nil {
				log.Fatal(err)
		}

		fmt.Println(len(cResult))

		for _, item := range cResult {
			fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.DeploymentConfig, item.Scale)
		}
	}
}

func main() {
   mainLoop()
}
