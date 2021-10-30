package pqstream

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

type EventstreamHandler struct {
	PostgreSQLConnectionString string
}

func (e EventstreamHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	topicsStr := req.URL.Query().Get("topics")
	topics := strings.Split(topicsStr, ",")

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	listener := pq.NewListener(e.PostgreSQLConnectionString, time.Millisecond*1000, time.Second*10, reportProblem)

	err = listener.Listen("eventstream")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	subscription, err := NewSubscription(req.Context(), e.PostgreSQLConnectionString, topics...)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer subscription.Cancel()

	for {
		select {
		case event := <-subscription.C:
			err = conn.WriteJSON(event)
			if err != nil {
				log.Println("Write:", err)
				break
			}
		}
	}
}
