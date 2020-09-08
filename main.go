package main

import "pkt_activity/src/service"

func main() {
	srv := service.NewActivityService()
	srv.Run()
}
