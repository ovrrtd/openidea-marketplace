package main

import "log"

func main() {
	err := migrate()

	if err != nil {
		log.Fatal(err)
	}
}
