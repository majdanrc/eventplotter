package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/majdanrc/eventplotter/events"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/araddon/dateparse"

	"github.com/majdanrc/eventplotter/config"
	"github.com/majdanrc/eventplotter/plotter"
	"github.com/majdanrc/eventplotter/providers"
)

func main() {
	var (
		configPath  = flag.String("c", "e:\\plotterconf.json", "config file")
		env         = flag.String("e", "dev", "env")
		streamsFile = flag.String("s", "", "streams file")
	)
	flag.Parse()
	_ = streamsFile

	conf, err := dbconfig.Read(*env, *configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf.Server)

	userName := "ErnestBarron41688@example.com"
	errodDate := "10/14/2017 03:06:19 PM"
	errorDateParsed, err := dateparse.ParseAny(errodDate)
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

	//typedUsers = append(typedUsers, events.BasicEvent{Values: map[string]string{"idSupplier": "1307"}})

	subsQuery := "select * from Subscription with(NOLOCK) where idSupplier=@SupplierId order by dateCreation asc"
	invoicesQuery := "select * from Invoice with(NOLOCK) where idSupplier=@SupplierId order by dateCreation asc"
	suppTypesQuery := "select * from SUPPLIERTYPE with(NOLOCK) where idSupplier=@SupplierId order by ACTIONDATE asc"

	subs, _ := providers.ProvideVertical(conf, strings.Replace(subsQuery, "@SupplierId", typedUsers[0].Values["idSupplier"], -1))
	invoices, _ := providers.ProvideProgressing(conf, strings.Replace(invoicesQuery, "@SupplierId", typedUsers[0].Values["idSupplier"], -1))
	suppTypes, _ := providers.ProvideBasic(conf, strings.Replace(suppTypesQuery, "@SupplierId", typedUsers[0].Values["idSupplier"], -1))

	var events []interface{}

	events = append(events, subs...)
	events = append(events, invoices...)
	events = append(events, suppTypes...)

	//_ = suppTypeEvents
	fmt.Println(len(events))

	plotter := new(plotter.Plotter)
	plotter.Plot(events, errorDateParsed, typedUsers[0].Values["idSupplier"])
}
