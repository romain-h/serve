package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var SocketScript = `
<script type="text/javascript">
	(function connect() {
		var socket = new WebSocket("ws://localhost:%d/websocket");

		socket.onclose = function(event) {
			console.log("WS close");
			socket = null;

			// Try reconnect
			setTimeout(function() {
				connect();
			}, 5000)
		}

		socket.onmessage = function(event) {
			const ev = JSON.parse(event.data)
			if (ev.type && ev.type === 'reload') {
				socket.close(1000, "Reloading page");
				location.reload(true);
			}
		}
	})()
</script>
`

type wsEvent struct {
	Type string `json:"type"`
}

func getSocketHandler(buildDone chan bool) http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			for {
				select {
				case <-buildDone:
					if err := conn.WriteJSON(wsEvent{Type: "reload"}); err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}
}
