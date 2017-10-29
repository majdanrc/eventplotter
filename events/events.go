package events

import "time"
import "fmt"

type Event interface {
	KeyEvents() []time.Time
	Description() []string
}

type BasicEvent struct {
	On         time.Time
	Values     map[string]string
	DateValues []time.Time
}

func (e BasicEvent) KeyEvents() []time.Time {
	return e.DateValues
}

func (e BasicEvent) Description() []string {
	var res []string

	item := fmt.Sprintf("%s: %s -> %s", "STATUS_TEST", e.Values["OLDSTATUS"], e.Values["NEWSTATUS"])
	res = append(res, item)

	return res
}

type VerticalEvent struct {
	On         time.Time
	Values     map[string]string
	DateValues []time.Time
	InfoValues []string
}

func (e VerticalEvent) KeyEvents() []time.Time {
	return e.DateValues
}

func (e VerticalEvent) Description() []string {
	var res []string

	item := fmt.Sprintf("vertical")
	res = append(res, item)

	return res
}

type ProgressingEvent struct {
	On         time.Time
	Values     map[string]string
	DateValues []time.Time
}

func (e ProgressingEvent) KeyEvents() []time.Time {
	return e.DateValues
}

func (e ProgressingEvent) Description() []string {
	var res []string

	item := fmt.Sprintf("%s: %s", "invoiceNumber", e.Values["invoiceNumber"])
	res = append(res, item)

	return res
}
