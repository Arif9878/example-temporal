package service

import (
	"log"
	"net/http"

	"github.com/Arif9878/example-temporal/activity"
	"github.com/Arif9878/example-temporal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Server() {
	// set up the worker
	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "user-namespace",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "cart", worker.Options{})

	// Get cart workflow
	w.RegisterWorkflow(workflow.GetCartWorkflow)
	w.RegisterActivity(activity.GetCart)

	// Set cart workflow
	w.RegisterWorkflow(workflow.SetCartWorkflow)
	w.RegisterActivity(activity.SetCart)

	mux := http.NewServeMux()
	mux.HandleFunc("/cart", CartGetHandler)     // curl -X GET http://localhost:5001/cart
	mux.HandleFunc("/cart/set", CartSetHandler) // curl -X POST http://localhost:5001/cart/set\?products\=1,2,3

	server := &http.Server{Addr: ":5001", Handler: mux}

	// start the worker and the web server
	go func() {
		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("unable to start Worker", err)
		}
	}()

	log.Fatal(server.ListenAndServe())
}
