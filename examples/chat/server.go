package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// chatServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type chatServer struct {
}

func (s chatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:   []string{"chat"},
		OriginPatterns: []string{"localhost:8007", "127.0.0.1:8007"},
	})
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	if c.Subprotocol() != "chat" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the chat subprotocol")
		return
	}

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = echo(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Printf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		return err
	}
	log.Printf("received: %v", v)

	return nil
}
