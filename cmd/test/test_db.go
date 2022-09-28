package main

import (
	"log"

	"github.com/pdfkpb/gobo/pkg/patron"
)

func main() {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		log.Fatalf("failed to LoadPatronDB: %v", err)
	}

	err = patronDB.AddUser("test", 5000)
	if err != nil {
		switch err {
		case patron.ErrorUserAlreadyRegistered:
		default:
			log.Fatalf("failed to AddUser: %v", err)
		}
	}

	err = patronDB.CreateChallenge("test", "test1", 50)
	if err != nil {
		log.Fatalf("failed to CreateChallenge: %v", err)
	}

	escrow, err := patronDB.GetChallenge("test")
	if err != nil {
		log.Fatalf("failed to GetChallenge: %v", err)
	}

	log.Printf("Completed, escrow is %d\n", escrow)
}
