package eventservicehandler

import (
	"bookmyevents/lib"
	"bookmyevents/repository/dblayer/mongolayer"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type eventServiceHandler struct {
	dbhandler mongolayer.DatabaseHandler
}

func NewEventServiceHandler(databasehandler mongolayer.DatabaseHandler) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler: databasehandler,
	}
}

// ServeAPI serves events APIs
func ServeAPI(tlsendpoint, endpoint string, dbHandler mongolayer.DatabaseHandler) (chan error, chan error) {
	handler := NewEventServiceHandler(dbHandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{Search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	httpErrorChan := make(chan error)
	httpTLSErrorChan := make(chan error)

	go func() {
		httpErrorChan <- http.ListenAndServe(endpoint, r)
	}()

	go func() {
		httpTLSErrorChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.pem", r)
	}()

	return httpErrorChan, httpTLSErrorChan
}

func (eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{error : No search criteria found, you can either search by id lie /id/4 or name like /name/new_year_party}`)
		log.Println("No SearchCriteria found in the path")
	}

	searchKey, ok := vars["Search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{error : No search key found }`)
		log.Println("No seatrch key found in the path")
	}

	var event lib.Event
	var err error

	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchKey)
	case "id":
		id, err := hex.DecodeString(searchKey)
		if err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		fmt.Fprintf(w, "{error : %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json:charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		fmt.Fprintf(w, "{ error : Error occurred while trying to find all available vents %s }", err)
		return
	}
	w.Header().Set("Content-Type", "application/json:charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{ error: Error encoding events %s}", err)
	}

}

func (eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := lib.Event{}

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error decoding event data %s}", err)
		return
	}

	id, err := eh.dbhandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error adding an event %d %s}", id, err)
		return
	}

}
