package main

import (
	"fmt"
	"time"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	//flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type TimeRule struct {
        Id string
        DeploymentConfig string
				Scale int
}

func mainLoop() {

	fmt.Println("[INFO] - Initialize autoscaler")

	viper.BindEnv("GPCMONGO_SERVICE_HOST")
	viper.BindEnv("GPCMONGO_SERVICE_PORT")
	viper.BindEnv("GPCMONGO_USER")
	viper.BindEnv("GPCMONGO_PASS")

	var mongo_host string = "0.0.0.0"
	var mongo_port string = "27017"
	var mongo_user string = ""
	var mongo_pass string = ""

	mongo_host = viper.GetString("GPCMONGO_SERVICE_HOST")
	mongo_port = viper.GetString("GPCMONGO_SERVICE_PORT")
	mongo_user = viper.GetString("GPCMONGO_USER")
	mongo_pass = viper.GetString("GPCMONGO_PASS")


	if(mongo_host==""){
		fmt.Println("[INFO] - Cannot find host and port ... stablished 0.0.0.0:27017 by default")
		mongo_host = "0.0.0.0"
		mongo_port = "27017"
	}else{
		fmt.Println("[INFO] - Rule database in "+mongo_host+":"+mongo_port)
	}

	if(mongo_user=="" || mongo_pass==""){
		fmt.Println("[ERROR] - Cannot find OSE credential. You must provide them in GPCMONGO_USER and GPCMONGO_PASS environment vars")
		os.Exit(0)
	}

	fmt.Printf("[INFO] - Connecting to %s:%s, database gpc with user %s and %s\n",mongo_host,mongo_port,mongo_user,mongo_pass)

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongo_host+":"+mongo_port},
		Timeout:  60 * time.Second,
		Database: "gpc",
		Username: mongo_user,
		Password: mongo_pass,
	}

	//session, err := mgo.Dial("0.0.0.0:27017")
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
					panic(err)
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	fmt.Println("tengo la session")
	c := session.DB("gpc").C("timeRule")
	//c := session.DB("test").C("test")

	fmt.Println("[INFO] - Connection was complete successfully")

	cResult :=[]TimeRule{}

	for {
		loc, _ := time.LoadLocation("Europe/Madrid")
		time.Sleep(1000 * 60 * time.Millisecond)
		currentTime := time.Now()
		fmt.Println(currentTime.In(loc))
		_CurrentTime := currentTime.String()
		hourNow := _CurrentTime[11:16]
		day := _CurrentTime[8:10]
		month := _CurrentTime[5:7]

		expr1 := hourNow
		expr2 := day+" "+hourNow
		expr3 := day+"/"+month+" "+hourNow
		expr4 := currentTime.Weekday().String()[0:3]+" "+hourNow

		//Valid expressions

		fmt.Println("expr1: ",expr1)
		fmt.Println("expr2: ",expr2)
		fmt.Println("expr3: ",expr3)
		fmt.Println("expr4: ",expr4)
		fmt.Println(len(expr1))


		err = c.Find(bson.M{"timeExpression": expr1}).All(&cResult)
		if err != nil {
						log.Fatal(err)
		}

		fmt.Println(len(cResult))

		for _, item := range cResult {
        fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.DeploymentConfig, item.Scale)
    }

		////////
		/*err = cc.Find(bson.M{"id": hourNow}).All(&cResult)
		if err != nil {
				log.Fatal(err)
		}

		fmt.Println(len(cResult))

		for _, item := range cResult {
			fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.DeploymentConfig, item.Scale)
		}*/
	}
}

func main() {
   mainLoop()
}
