package main

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Test_chatServer tests the chatServer by sending it 5 different messages
// and ensuring the responses all match.
func Test_chatServer(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(chatServer{})
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, s.URL, &websocket.DialOptions{
		Subprotocols: []string{"chat"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for i := 0; i < 5; i++ {
		err = wsjson.Write(ctx, c, map[string]int{
			"i": i,
		})
		if err != nil {
			t.Fatal(err)
		}

		v := map[string]int{}
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v["i"])
		if v["i"] != i {
			t.Fatalf("expected %v but got %v", i, v)
		}
	}

	c.Close(websocket.StatusNormalClosure, "")
}
