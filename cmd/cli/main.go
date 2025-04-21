package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/jankaczmarski/http-server"
)

const dbFileName = "game.db.json"

func main() {
	store, closeStore, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	blindAlerter := poker.BlindAlerterFunc(poker.Alerter)
	game := poker.NewTexasHoldem(blindAlerter, store)
	if err != nil {
		log.Fatal(err)
	}
	defer closeStore()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	poker.NewCLI(os.Stdin, os.Stdout, game).PlayPoker()
}
