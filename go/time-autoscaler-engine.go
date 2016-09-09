package main

import (
	"crypto/tls"
	"encoding/json"
	"bytes"
	"fmt"
	"time"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	//flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

//Data structure which contains the
type TimeRule struct {
	Region          string            `json:"region" yaml:"region"`
	Project         string            `json:"project" yaml:"project"`
	Dc              string            `json:"dc" yaml:"dc"` //alias deploymentconfig
	TimeExpression  string            `json:"timeExpression" yaml:"timeExpression"`
	Alias						string						`json:"alias" yaml:"alias"`
	Instances       int               `json:"instances" yaml:"instances"`
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

		//Check expressions in database

		//1.
		err = c.Find(bson.M{"timeExpression": expr1}).All(&cResult)
		if err != nil {
						log.Fatal(err)
		}

		for _, item := range cResult {
        fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.Dc, item.Instances)
				scale(item)
    }

		//2.
		err = c.Find(bson.M{"timeExpression": expr2}).All(&cResult)
		if err != nil {
						log.Fatal(err)
		}

		for _, item := range cResult {
        fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.Dc, item.Instances)
				scale(item)
    }

		//3.
		err = c.Find(bson.M{"timeExpression": expr3}).All(&cResult)
		if err != nil {
						log.Fatal(err)
		}

		for _, item := range cResult {
        fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.Dc, item.Instances)
				scale(item)
    }

		//4.
		err = c.Find(bson.M{"timeExpression": expr4}).All(&cResult)
		if err != nil {
						log.Fatal(err)
		}

		for _, item := range cResult {
        fmt.Printf(" DeploymentConfig: %s - scale: %d \n", item.Dc, item.Instances)
				scale(item)
    }


	}
}

func main() {
   mainLoop()
}

func scale(data TimeRule) bool{
		fmt.Printf("Scaling DeploymentConfig %s in namespace %s in region %s", data.Dc, data.Project, data.Region)

		ioJsonData := new(bytes.Buffer)
    json.NewEncoder(ioJsonData).Encode(data)


		transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	scaleUrl := "https://gpc-api-service-globalpaas-dev.appls.boaw.paas.gsnetcloud.corp/api/ose/scale/"+data.Region+"/"+data.Project+"/"+data.Dc+"/"+strconv.Itoa(data.Instances)
	fmt.Println(scaleUrl)

	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", scaleUrl, ioJsonData)

	req.Header.Add("Authorization", "eyJ0b2tlbiI6eyJjYXAxIjoiVzZuQVpqRnVkQlJfQ3hmN1BRTzV5a0pRVmpfLTl4QXdnOVZwR2lDdzN2SSIsImN0bzIiOiJzRkJpWHJSSlF1U2xTY2xHaFhWV1NYbXJuY05kQXVOaWdTWE9EV3NlYzBnIiwiYm9hdyI6ImxRN1NhUEZPMEtISjBnajRGd1Atak9aMWxaWVFBNFl5WS1yekZ4S2JoSnMiLCJib2FlIjoiTGhweEdHN2lFYnBvUG52cjRobVZMOW9wMlEtLVJOb2lSYXdBYmhKaGIxSSJ9fQ==")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Unknow error scaling dc %s, in namespace %s in region %s\n",data.Dc,data.Project,data.Region)
		fmt.Println(err)
	} else {
		if res.StatusCode == 200 {
			fmt.Println("Time-rule has been added successfully")
			return true
		}else{
			switch res.StatusCode {
	    	case 403:
					fmt.Printf("You're not allowed in project %s\n",data.Project)
	        return false
				case 404:
					fmt.Printf("No rules are found for selected DC %s\n",data.Project)
	        return false
				default:
					fmt.Println("Error: %d\n",res.StatusCode)
					return false
	    	}
		}

	}

	return true
}
