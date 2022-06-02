package main

import (
	"flag"
	"fmt"

	"github.com/jalexanderII/MyEventsMicro/eventsservice/rest"
	"github.com/jalexanderII/MyEventsMicro/lib/configuration"
	"github.com/jalexanderII/MyEventsMicro/lib/persistence/dblayer"
)

func main() {

	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	// extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath)

	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	// RESTful API start
	rest.ServeAPI(config.RestfulEndpoint, dbhandler)
}
