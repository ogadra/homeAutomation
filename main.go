package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const DefaultEndpoint = "https://api.switch-bot.com"

type RequestBody struct {
	Command string `json:"command"`
	Type    string `json:"commandType,omitempty"`
}

func env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func switchbot_post(deviceId string, command *RequestBody) {
	reqBody, _ := json.Marshal(*command)
	client := &http.Client{}
	fmt.Println(string(reqBody))
	req, _ := http.NewRequest("POST", DefaultEndpoint+"/v1.0/devices/"+deviceId+"/commands", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", os.Getenv("SBTOKEN"))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func periodic(sec int, times int, deviceId string, command *RequestBody){
	for i := 0; i < times; i++ {
		switchbot_post(deviceId, command)
		time.Sleep((time.Second * time.Duration(sec)))
	}
}

func main() {
	env_load()
	Lights := [2]string{os.Getenv("LIVINGLIGHT"), os.Getenv("PCLIGHT")}

	fmt.Println("Hello, World!")

	r := gin.Default()

	r.GET("", func(c *gin.Context) {
		fmt.Println(c.ClientIP())
		fmt.Println(c.GetHeader("X-Real-IP"))
		command := &RequestBody{
			Command: "brightnessDown",
//			Type: "customize",
		}
		for _, v := range Lights {
			go periodic(10, 20, v, command)
		}
		c.JSON(200, gin.H{"response": "OK!"})

	})

	r.Run()
}
