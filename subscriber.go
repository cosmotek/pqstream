package pqstream

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"time"

	"github.com/lib/pq"
)

type Subscription struct {
	C      <-chan Event
	cancel func()
}

func (e Subscription) Cancel() {
	e.cancel()
}

func NewSubscription(ctx context.Context, postgresConnectionConfig string, topics ...string) (Subscription, error) {
	listener := pq.NewListener(postgresConnectionConfig, time.Millisecond*1000, time.Second*10, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	})

	err := listener.Listen("eventstream")
	if err != nil {
		return Subscription{}, err
	}

	ctx, cancelfn := context.WithCancel(ctx)
	notify := make(chan Event)
	go func(l *pq.Listener) {
		for {
			select {
			case n := <-l.Notify:
				log.Println("Received event", n)

				var event Event
				err := json.Unmarshal([]byte(n.Extra), &event)
				if err != nil {
					log.Println("Error processing JSON: ", err)
					continue
				}

				for _, topic := range topics {
					matches, err := filepath.Match(topic, event.Topic)
					if err != nil {
						log.Println("Error matching topic: ", err)
						continue
					}

					log.Println(topic, event.Topic, matches)

					if matches {
						notify <- event
					}
				}

				break
			case <-time.After(90 * time.Second):
				log.Println("Received no events for 90 seconds, checking connection")
				go l.Ping()
				break
			case <-ctx.Done():
				log.Println("Exiting gracefully")
				listener.Close()
				break
			}
		}
	}(listener)

	return Subscription{
		C:      notify,
		cancel: cancelfn,
	}, nil
}
