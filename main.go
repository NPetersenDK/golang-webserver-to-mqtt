package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	http.HandleFunc("/post", handlePost)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var message Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Received message: %s\n", string(body))

	publishMQTT(body)
	w.WriteHeader(http.StatusOK)
}

func publishMQTT(jsonBytes []byte) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883") // Replace with your MQTT broker address
	opts.SetClientID("mqtt-golang")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Failed to connect to MQTT broker:", token.Error())
		return
	}
	// defer client.Disconnect(250)

	topic := "test-golang/topic" // Replace with the desired topic to publish the JSON data

	token := client.Publish(topic, 0, false, jsonBytes)
	token.Wait()

	if token.Error() != nil {
		log.Println("Failed to publish MQTT message:", token.Error())
		return
	}

	log.Println("Published MQTT message successfully")
	client.Disconnect(50)
}
