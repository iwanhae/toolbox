package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
)

var (
	connected  int32 = 0
	terminated int32 = 0
)

func main() {
	s := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := map[string]string{}
			for k, v := range r.Header {
				result[fmt.Sprintf("header_%s", k)] = strings.Join(v, ", ")
			}
			result["method"] = r.Method
			result["url"] = r.URL.String()
			result["proto"] = r.Proto
			result["remote_addr"] = r.RemoteAddr
			result["host"] = r.Host
			result["stat_connected"] = fmt.Sprint(connected)
			result["stat_terminated"] = fmt.Sprint(terminated)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}),
		ConnState: func(c net.Conn, cs http.ConnState) {
			log.Printf("conn state: %s / connected: %d / terminated: %d", cs, connected, terminated)
			switch cs {
			case http.StateNew:
				atomic.AddInt32(&connected, 1)
			case http.StateClosed:
				atomic.AddInt32(&terminated, 1)
				atomic.AddInt32(&connected, -1)
			case http.StateHijacked:
				atomic.AddInt32(&terminated, 1)
				atomic.AddInt32(&connected, -1)
			case http.StateActive:
				// do nothing
			}
		}}
	log.Printf("Listening on %s", s.Addr)
	s.ListenAndServe()
}
