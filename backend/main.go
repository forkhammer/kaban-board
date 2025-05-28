package main

import "main/cmd/web"

func main() {
	app := web.NewApplication()
	app.Run()
}
