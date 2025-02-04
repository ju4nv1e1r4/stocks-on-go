package main

import (
	"log"
	"os"
	"stocknews/requests"
)

func main()  {
	aplication := requests.Start()
	erro := aplication.Run(os.Args)

	if erro != nil {
		log.Fatal(erro)
	}
}
