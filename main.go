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

type commandType string
const (
	command = commandType("command")
	customize = commandType("customize")
)

type commandContent string
const (
	turnOn = commandContent("turnOn")
	turnOff = commandContent("turnOff")
	brightnessUp = commandContent("brightnessUp")
	brightnessDown = commandContent("brightnessDown")
	warmer = commandContent("光色 赤+")
	cooler = commandContent("光色 青+")
)

type RequestBody struct {
	Command commandContent `json:"command"`
	Type    commandType `json:"commandType,omitempty"`
}

type RecieveBody struct {
	Command	commandContent `json:"command"`
	CommandType	commandType `json"commandType"`
	Interval	int	`json:"interval"`
	Times	*int	`json:"times"`
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

func periodic(interval int, times int, deviceId string, command *RequestBody){
	for i := 0; i < times; i++ {
		switchbot_post(deviceId, command)
		time.Sleep((time.Second * time.Duration(interval)))
	}
}

func main() {
	env_load()
	Lights := [2]string{os.Getenv("LIVINGLIGHT"), os.Getenv("PCLIGHT")}

	fmt.Println("Hello, World!")

	r := gin.Default()

	r.POST("", func(c *gin.Context) {
		var req RecieveBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Times == nil {
			defaultTimes := 1
			req.Times = &defaultTimes
		}

		fmt.Println(req)

		command := &RequestBody{
			Command: req.Command,
			Type: req.CommandType,
		}

		for _, v := range Lights {
			go periodic(req.Interval, *req.Times, v, command)
		}
		c.JSON(200, gin.H{"response": "OK!"})
	})

	r.Run()
}
