package function

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	cannedResponse = CannedResponse{
		Status: http.StatusOK,
	}
)

type CannedResponse struct {
	Status int    `json:"status"`
	Body   string `json:"body"`

	Delay Duration `json:"delay"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.URL.Path, r.Method)

	if r.URL.Path == "/get" && r.Method == http.MethodGet {
		defer r.Body.Close()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		bytes, err := json.Marshal(cannedResponse)
		if err != nil {
			log.Printf("/set failed %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(bytes)

		return
	}

	if r.URL.Path == "/set" && r.Method == http.MethodPost {
		defer r.Body.Close()

		buffer, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(buffer, &cannedResponse)
		if err != nil {
			log.Printf("/set failed %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if cannedResponse.Delay.Duration > time.Millisecond*0 {
		time.Sleep(cannedResponse.Delay.Duration)
	}

	w.WriteHeader(cannedResponse.Status)

	if len(cannedResponse.Body) > 0 {
		w.Write([]byte(cannedResponse.Body))
	}
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
