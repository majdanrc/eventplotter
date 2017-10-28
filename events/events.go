package events

import "time"
import "fmt"

type EventCollection []interface{}

type BasicEvent struct {
	On     time.Time
	Values map[string]string
}

func (e BasicEvent) KeyEvents() []time.Time {
	var res []time.Time

	res = append(res, convertToDate(e.Values["ACTIONDATE"]))

	return res
}

func (e BasicEvent) Description() []string {
	var res []string

	item := fmt.Sprintf("%s: %s -> %s", "STATUS", e.Values["OLDSTATUS"], e.Values["NEWSTATUS"])
	res = append(res, item)

	return res
}

type VerticalEvent struct {
	On         time.Time
	Values     map[string]string
	DateValues []string
	InfoValues []string
}

func (e VerticalEvent) KeyEvents() []time.Time {
	var res []time.Time

	for _, k := range e.DateValues {
		res = append(res, convertToDate(e.Values[k]))
	}

	return res
}

type ProgressingEvent struct {
	On     time.Time
	Values map[string]string
}

func (e ProgressingEvent) KeyEvents() []time.Time {
	var res []time.Time

	res = append(res, convertToDate(e.Values["dateCreation"]))
	res = append(res, convertToDate(e.Values["issueDate"]))
	res = append(res, convertToDate(e.Values["paymentDate"]))

	return res
}

func (e ProgressingEvent) Description() []string {
	var res []string

	item := fmt.Sprintf("%s: %s", "invoiceNumber", e.Values["invoiceNumber"])
	res = append(res, item)

	return res
}

func convertToDate(date string) time.Time {
	res, _ := time.Parse("2006-01-02 15:04:05.999", date)
	return res
}
