package providers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// used only for provider
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/majdanrc/eventplotter/config"
	"github.com/majdanrc/eventplotter/events"
	"github.com/majdanrc/eventplotter/streamer"
)

type MsSqlProvider struct {
	Config dbconfig.Config
}

func (p *MsSqlProvider) ProvideEvents(streamer streamer.Stream) ([]events.Event, error) {
	var res []events.Event

	dsn := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;database=%s",
		p.Config.Server, p.Config.Login, p.Config.Password, p.Config.Database)

	db, err := sql.Open("mssql", dsn)
	if err != nil {
		return res, fmt.Errorf("cannot connect: %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		return res, fmt.Errorf("cannot connect: %s", err.Error())
	}
	defer db.Close()

	params := make(map[string]string)

	eventValues, _ := exec(db, streamer.Query, params)

	for _, v := range eventValues {

		switch streamer.Type {
		case 2:
			ev := events.BasicEvent{On: time.Now().UTC(), Values: v}
			for _, dv := range streamer.DateValues {
				ev.DateValues = append(ev.DateValues, convertToDate(v[dv]))
			}
			res = append(res, ev)
		case 3:
			ev := events.VerticalEvent{On: time.Now().UTC(), Values: v}
			for _, dv := range streamer.DateValues {
				ev.DateValues = append(ev.DateValues, convertToDate(v[dv]))
			}
			res = append(res, ev)
		case 4:
			ev := events.ProgressingEvent{On: time.Now().UTC(), Values: v}
			for _, dv := range streamer.DateValues {
				ev.DateValues = append(ev.DateValues, convertToDate(v[dv]))
			}
			res = append(res, ev)
		}
	}

	return res, nil
}

func ProvideBasic(conf dbconfig.Config, query string) ([]events.Event, error) {
	var res []events.Event

	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", conf.Server, conf.Login, conf.Password, conf.Database)

	db, err := sql.Open("mssql", dsn)
	if err != nil {
		return res, fmt.Errorf("cannot connect: %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		return res, fmt.Errorf("cannot connect: %s", err.Error())
	}
	defer db.Close()

	params := make(map[string]string)

	suppType, _ := exec(db, query, params)

	for _, v := range suppType {
		ev := events.BasicEvent{On: time.Now().UTC(), Values: v}

		res = append(res, ev)
	}

	return res, nil
}

func exec(db *sql.DB, cmd string, params map[string]string) ([]map[string]string, error) {
	for k, v := range params {
		cmd = strings.Replace(cmd, k, v, -1)
	}

	fmt.Println(cmd)

	rows, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if cols == nil {
		return nil, nil
	}
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	var result []map[string]string

	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}

		rowValues := make(map[string]string)

		for i := 0; i < len(vals); i++ {
			rowValues[cols[i]] = parseValue(vals[i].(*interface{}))
		}

		result = append(result, rowValues)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return result, nil
}

func parseValue(pval *interface{}) string {
	var res string

	switch v := (*pval).(type) {
	case nil:
		res = fmt.Sprint("NULL")
	case bool:
		if v {
			res = fmt.Sprint("1")
		} else {
			res = fmt.Sprint("0")
		}
	case []byte:
		res = fmt.Sprint(string(v))
	case time.Time:
		res = fmt.Sprint(v.Format("2006-01-02 15:04:05.999"))
	default:
		res = fmt.Sprint(v)
	}

	return res
}

func convertToDate(date string) time.Time {
	res, _ := time.Parse("2006-01-02 15:04:05.999", date)
	return res
}
