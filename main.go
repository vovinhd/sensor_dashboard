package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseUrl       string       `yaml:"base_url"`
	MqttUrl       string       `yaml:"mqtt_url"`
	MqttDevices   []MqttDevice `yaml:"mqtt_devices"`
	DatasourceDir string       `yaml:"datasource_dir"`
	ServeFromFS   bool         `yaml:"serve_from_fs"`
}

type MqttDevice struct {
	Name  string `yaml:"name"`
	Topic string `yaml:"topic"`
	Type  string `yaml:"type"`
	Tag   string `yaml:"tag"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, all config coming from system env")
	}
	//
	//dashboardHost := os.Getenv("SENSOR_DASHBOARD_WEB_BASE_URL")
	//mqttServer := os.Getenv("SENSOR_DASHBOARD_MQTT_SERVER")
	//mqttTopicsJson := os.Getenv("SENSOR_DASHBOARD_MQTT_TOPICS")
	//datasourceDir := os.Getenv("SENSOR_DASHBOARD_DATASOURCE_DIR")

	configPath := os.Getenv("SENSOR_DASHBOARD_CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}
	f, err := os.Open(configPath)

	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Panic(err)
		}
	}(f)

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	queries := InitDB(config.DatasourceDir)
	_, err = InitMqtt(config.MqttUrl, config.MqttDevices, *queries)
	if err != nil {
		panic(err)
	}

	handler, err := Api(*queries, config)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server on " + config.BaseUrl)
	log.Fatal(http.ListenAndServe(config.BaseUrl, handler))

	//client.Disconnect(250)

}
