package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	err := run()
	if err != nil {
		log.Fatal().Stack().Caller().Err(err).Send()
	}
}

func run() error {
	router := mux.NewRouter()
	router.Handle("/get/{id}", get()).Methods(http.MethodGet)
	router.Handle("/put/{id}", put()).Methods(http.MethodPut)
	router.Handle("/post/{id}", post()).Methods(http.MethodPost)
	router.Handle("/delete/{id}", delete()).Methods(http.MethodDelete)
	router.Handle("/head/{id}", head()).Methods(http.MethodHead)

	srv := http.Server{
		Addr:         "localhost:3000",
		Handler:      cors.Default().Handler(router),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	log.Info().Msgf("Server running on %v", srv.Addr)
	return srv.ListenAndServe()
}

func close(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r := recover(); r != nil {
		var err error
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("unknown error")
		}
		respond(w, http.StatusInternalServerError, "text/html", err.Error())
	}
}

type response struct {
	ID       string        `json:"id" xml:"id"`
	String   string        `json:"string" xml:"string"`
	Integer  int           `json:"integer" xml:"integer"`
	Float    float64       `json:"float" xml:"float"`
	Date     *time.Time    `json:"date" xml:"date"`
	Duration time.Duration `json:"duration" xml:"duration"`
	Boolean  bool          `json:"boolean" xml:"boolean"`

	Nested struct {
		String   string        `json:"string" xml:"string"`
		Integer  int           `json:"integer" xml:"integer"`
		Float    float64       `json:"float" xml:"float"`
		Date     *time.Time    `json:"date" xml:"date"`
		Duration time.Duration `json:"duration" xml:"duration"`
		Boolean  bool          `json:"boolean" xml:"boolean"`
	} `json:"nested" xml:"nested"`
}

func get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer close(w, r)
		vars := mux.Vars(r)
		id := vars["id"]

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "get").Str("id", id).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		if id == "" {
			respond(w, http.StatusNotFound, "text/html", http.StatusText(http.StatusNotFound))
			return
		}
		now := time.Now()
		respond(w, http.StatusOK, r.Header.Get("Accept"), response{
			ID:       id,
			String:   "asdf safasdf asdf",
			Integer:  rand.Int(),
			Float:    rand.Float64(),
			Date:     &now,
			Duration: time.Duration(rand.Int()),
			Boolean:  rand.Intn(1) > 0,
			Nested: struct {
				String   string        `json:"string" xml:"string"`
				Integer  int           `json:"integer" xml:"integer"`
				Float    float64       `json:"float" xml:"float"`
				Date     *time.Time    `json:"date" xml:"date"`
				Duration time.Duration `json:"duration" xml:"duration"`
				Boolean  bool          `json:"boolean" xml:"boolean"`
			}{
				String:   "asdf safasdf asdf",
				Integer:  rand.Int(),
				Float:    rand.Float64(),
				Date:     &now,
				Duration: time.Duration(rand.Int()),
				Boolean:  rand.Intn(1) > 0,
			},
		})
	}
}

func head() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer close(w, r)

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "head").Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		respond(w, http.StatusOK, "text/html", "head")
	}
}

func post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer close(w, r)

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "post").Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		response := response{}
		accept := r.Header.Get("Accept")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respond(w, http.StatusUnprocessableEntity, "text/html", err.Error())
		}

		if err := unmarshal(accept, body, &response); err != nil {
			respond(w, http.StatusUnprocessableEntity, "text/html", err.Error())
		}

		contentType := r.Header.Get("Content-Type")
		respond(w, http.StatusCreated, contentType, response)
	}
}

func put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer close(w, r)
		vars := mux.Vars(r)
		id := vars["id"]

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "put").Str("id", id).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		if id == "" {
			respond(w, http.StatusNotFound, "text/html", http.StatusText(http.StatusNotFound))
			return
		}

		response := response{}
		accept := r.Header.Get("Accept")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respond(w, http.StatusUnprocessableEntity, "text/html", err.Error())
		}

		if err := unmarshal(accept, body, &response); err != nil {
			respond(w, http.StatusUnprocessableEntity, "text/html", err.Error())
		}

		contentType := r.Header.Get("Content-Type")
		respond(w, http.StatusCreated, contentType, response)
	}
}

func patch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer close(w, r)
		vars := mux.Vars(r)
		id := vars["id"]

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "put").Str("id", id).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		if id == "" {
			respond(w, http.StatusNotFound, "text/html", http.StatusText(http.StatusNotFound))
			return
		}

		response := response{}
		accept := r.Header.Get("Accept")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respond(w, http.StatusUnprocessableEntity, "text/html", err.Error())
		}

		if err := unmarshal(accept, body, &response); err != nil {
			respond(w, http.StatusUnprocessableEntity, "text/html", err.Error())
		}

		contentType := r.Header.Get("Content-Type")
		respond(w, http.StatusCreated, contentType, response)
	}
}

func delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer close(w, r)
		vars := mux.Vars(r)
		id := vars["id"]

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "put").Str("id", id).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		if id == "" {
			respond(w, http.StatusNotFound, "text/html", http.StatusText(http.StatusNotFound))
			return
		}
		respond(w, http.StatusNoContent, "text/html", nil)
	}
}

func respond(w http.ResponseWriter, code int, contentType string, payload interface{}) {
	body, _ := marshal(contentType, payload)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write(body)
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	// 	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	// 	time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
}

func unmarshal(accept string, body []byte, v interface{}) error {
	switch accept {
	case "application/json":
		return json.Unmarshal(body, v)
	case "application/xml":
		return xml.Unmarshal(body, v)
	default:
		return json.Unmarshal(body, v)
	}
}

func marshal(contentType string, response interface{}) ([]byte, error) {
	switch contentType {
	case "application/json":
		return json.Marshal(response)
	case "application/xml":
		return xml.Marshal(response)
	default:
		return json.Marshal(response)
	}
}
