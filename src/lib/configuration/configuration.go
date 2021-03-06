package configuration

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jalexanderII/MyEventsMicro/src/lib/persistence/dblayer"
)

var (
	DBTypeDefault       = dblayer.MONGODB
	DBConnectionDefault = "mongodb://127.0.0.1"
	DBNameDefault       = "MyEvents"
	RestfulEPDefault    = "localhost:8181"
	RestfulTLSEPDefault = "localhost:9191"
)

type ServiceConfig struct {
	Databasetype      dblayer.DBTYPE `json:"databasetype"`
	DBConnection      string         `json:"dbconnection"`
	DBName            string         `json:"dbname"`
	RestfulEndpoint   string         `json:"restfulapi_endpoint"`
	RestfulTLSEndPint string         `json:"restfulapi-tlsendpoint"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		DBNameDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}
