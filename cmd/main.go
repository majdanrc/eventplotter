package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	//_ "github.com/denisenkom/go-mssqldb"

	"github.com/araddon/dateparse"

	"github.com/majdanrc/eventplotter/config"
	"github.com/majdanrc/eventplotter/events"
	"github.com/majdanrc/eventplotter/plotter"
	"github.com/majdanrc/eventplotter/providers"
	"github.com/majdanrc/eventplotter/streamer"
)

func main() {
	var (
		configPath  = flag.String("c", "e:\\plotterconf.json", "config file")
		env         = flag.String("e", "dev", "env")
		streamsFile = flag.String("s", "e:\\samplestreams.json", "streams file")
	)
	flag.Parse()

	streamSource, _ := streamer.ReadStreamConfig(*streamsFile)
	fmt.Println(streamSource)

	conf, err := dbconfig.Read(*env, *configPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf.Server)

	userName := streamSource.Parameters["userName"]
	errorDateParsed, err := dateparse.ParseAny(streamSource.AnalysisDate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(errorDateParsed)

	user, _ := providers.ProvideBasic(conf, "select idSupplier, idUser from [User] with(NOLOCK) where username='"+userName+"'")
	var typedUsers []events.BasicEvent

	for _, item := range user {
		if s, ok := item.(events.BasicEvent); ok {
			typedUsers = append(typedUsers, s)
		}
	}

	if len(typedUsers) == 0 || typedUsers[0].Values["idSupplier"] == "" || typedUsers[0].Values["idSupplier"] == "NULL" {
		log.Fatal("user problem")
	}

	streamSource.Parameters["@SupplierId"] = typedUsers[0].Values["idSupplier"]

	provider := new(providers.MsSqlProvider)
	provider.Config = conf

	var allEvents []events.Event

	for _, item := range streamSource.Streams {
		item.Query = strings.Replace(item.Query, "@SupplierId", streamSource.Parameters["@SupplierId"], -1)

		providedEvents, _ := provider.ProvideEvents(item)
		allEvents = append(allEvents, providedEvents...)
	}

	fmt.Println(len(allEvents))

	plotter := new(plotter.Plotter)
	plotter.Plot(allEvents, errorDateParsed, typedUsers[0].Values["idSupplier"])
}
