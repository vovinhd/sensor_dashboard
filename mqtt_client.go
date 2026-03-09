package main

import (
	"context"
	"encoding/json"
	"log"
	"sensor_dashboard/db"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttHumiditySensorMessage struct {
	Time       IsoTime            `json:"Time"`
	SensorData HumiditySensorData `json:"AM2301"`
	TempUnit   string             `json:"TempUnit"`
}

type MqttPowerSensorMessage struct {
	Time              IsoTime           `json:"Time"`
	SensorData        PowerSensorData   `json:"ENERGY"`
	SensorTemperature SensorTemperature `json:"ESP32"`
	TempUnit          string            `json:"TempUnit"`
}

type MqttSwitchStateSensorMessage struct {
	PowerState string `json:"POWER"`
}

type HumiditySensorData struct {
	Temperature float64 `json:"Temperature"`
	Humidity    float64 `json:"Humidity"`
	DewPoint    float64 `json:"DewPoint"`
}
type PowerSensorData struct {
	TotalStartTime IsoTime `json:"TotalStartTime"`
	Total          float64 `json:"Total"`
	Yesterday      float64 `json:"Yesterday"`
	Today          float64 `json:"Today"`
	Period         float64 `json:"Period"`
	Power          float64 `json:"Power"`
	ApparentPower  float64 `json:"ApparentPower"`
	ReactivePower  float64 `json:"ReactivePower"`
	Factor         float64 `json:"Factor"`
	Voltage        float64 `json:"Voltage"`
	Current        float64 `json:"Current"`
}

type SensorTemperature struct {
	Temperature float64 `json:"Temperature"`
}

type IsoTime struct {
	time.Time
}

const sensorTimeLayout = "2006-01-02T15:04:05"

func (ct *IsoTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(sensorTimeLayout, s)
	return
}

var wildcardHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("got wildcard message on topic %s with payload: %s\n", msg.Topic(), msg.Payload())
}

func messageHandler(queries db.Queries, sensorName string, sensorType string) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("handling message for %s! topic: %s, payload: %s\n", sensorName, msg.Topic(), msg.Payload())

		switch strings.ToLower(sensorType) {
		case "humidity":
			var message MqttHumiditySensorMessage
			err := json.Unmarshal(msg.Payload(), &message)
			if err != nil {
				log.Printf("error unmarshalling mqtt message: %s\n", err)
				return
			}

			_, err = queries.CreateHumidityLog(context.Background(), db.CreateHumidityLogParams{
				Sensor:      sensorName,
				Time:        message.Time.Time,
				Temperature: message.SensorData.Temperature,
				Humidity:    message.SensorData.Humidity,
				DewPoint:    message.SensorData.DewPoint,
			})

			if err != nil {
				log.Printf("error creating log: %s\n", err)
			}
		case "power":
			var message MqttPowerSensorMessage
			err := json.Unmarshal(msg.Payload(), &message)
			if err != nil {
				log.Printf("error unmarshalling mqtt message: %s\n", err)
			}

			_, err = queries.CreatePowerLog(context.Background(), db.CreatePowerLogParams{
				Time:              message.Time.Time,
				Sensor:            sensorName,
				TotalStartTime:    message.SensorData.TotalStartTime.Time,
				Total:             message.SensorData.Total,
				Yesterday:         message.SensorData.Yesterday,
				Today:             message.SensorData.Today,
				Period:            message.SensorData.Period,
				Power:             message.SensorData.Power,
				ApparentPower:     message.SensorData.ApparentPower,
				ReactivePower:     message.SensorData.ReactivePower,
				Factor:            message.SensorData.Factor,
				Voltage:           message.SensorData.Voltage,
				Current:           message.SensorData.Current,
				SensorTemperature: message.SensorTemperature.Temperature,
			})
			if err != nil {
				log.Printf("error creating log: %s\n", err)
			}
		case "switch":
			var message MqttSwitchStateSensorMessage
			err := json.Unmarshal(msg.Payload(), &message)
			if err != nil {
				log.Printf("error unmarshalling mqtt message: %s\n", err)
				return
			}

			_, err = queries.CreateSwitchStateLog(context.Background(), db.CreateSwitchStateLogParams{
				Time:        time.Now(),
				Sensor:      sensorName,
				SwitchState: message.PowerState,
			})

			if err != nil {
				log.Printf("error creating log: %s\n", err)
			}
		}
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Printf("Connected to broker!")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("connection lost: %v", err)
}

func createDeviceTags(queries db.Queries, sensors []MqttDevice) error {

	err := queries.TruncateDeviceTag(context.Background())

	if err != nil {
		return err
	}
	for _, sensor := range sensors {
		_, err = queries.CreateDeviceTag(context.Background(), db.CreateDeviceTagParams{
			DeviceName: sensor.Name,
			Tag:        sensor.Tag,
			DeviceType: sensor.Type,
		})
		if err != nil {
			return err
		}
	}

	return nil

}

func InitMqtt(brokerUrl string, clientId string, sensors []MqttDevice, queries db.Queries) (mqtt.Client, error) {

	err := createDeviceTags(queries, sensors)

	if err != nil {
		return nil, err
	}

	log.Println("Listening for sensors:", sensors)

	opts := mqtt.NewClientOptions().AddBroker(brokerUrl)
	opts.SetClientID(clientId)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetDefaultPublishHandler(wildcardHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, sensor := range sensors {
		sub(client, sensor, queries)
	}
	//subAll(client)

	return client, nil
}

func subAll(client mqtt.Client) {
	topic := "#"
	token := client.Subscribe(topic, 1, wildcardHandler)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func sub(client mqtt.Client, sensor MqttDevice, queries db.Queries) {
	topic := sensor.Topic
	token := client.Subscribe(topic, 1, messageHandler(queries, sensor.Name, sensor.Type))
	token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}
	log.Printf("Subscribed to topic: %s\n", topic)
}
