package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kithix/go-tell/sender"
	"github.com/kithix/go-tell/sender/chat"
	"github.com/kithix/go-tell/store"
	"github.com/kithix/go-tell/store/mem"
	"github.com/kithix/go-tell/worker"
)

func main() {
	// store.Basic is the default behaviours for a storage system.
	jobStorage := store.Basic(
		// A storage system provides CRUD operations for behaviours
		memstorage.Open(),
	)

	// A handler takes an interface.
	// A handler itself is an interface see (worker.JobHandler)
	// allowing you to write your own business logic based on job information
	// where required.
	handler := chat.MakeHandler(sender.NewStdoutSender())

	// SMTP SERVER
	// handler := email.MakeHandler(smtp.New(smtp.Options{
	// 	Host:               "",
	// 	Port:               ,
	// 	InsecureSkipVerify: ,
	// 	Username:           "",
	// 	Password:           "",
	// }))

	// SLACK
	// handler := chat.MakeHandler(sender.WrapGoNotifyChat(slack.New(
	// 	slack.Options{
	// 		AuthToken: "",
	// 	},
	// )))

	_, err := worker.Open(jobStorage, handler, log.New(os.Stdout, "", 1))
	if err != nil {
		log.Fatal(err)
	}

	// Custom HTTP endpoint that creates a basic chat job
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Send a chat job specifically to no one.
		job, err := store.BuildChatJob("Hello world")
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		err = jobStorage.AddJob(job)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	})

	var port string
	flag.StringVar(&port, "p", "8080", "port of HTTP server")

	log.Println("Starting server on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
