package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Arif9878/example-temporal/messages"
	"github.com/Arif9878/example-temporal/workflow"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func CartGetHandler(w http.ResponseWriter, r *http.Request) {
	// create a new temporal client
	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "user-namespace",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// get the cart
	we, err := c.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		ID:        "GetCartWorkflow_" + uuid.New().String(),
		TaskQueue: "cart",
	}, workflow.GetCartWorkflow)
	if err != nil {
		http.Error(w, "unable to start workflow", http.StatusInternalServerError)
		return
	}

	result := &[]messages.Product{}
	if err := we.Get(r.Context(), &result); err != nil {
		http.Error(w, "unable to get workflow result", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "unable to write response", http.StatusInternalServerError)
		return
	}
}

func CartSetHandler(w http.ResponseWriter, r *http.Request) {
	productString := r.URL.Query().Get("products")
	// create a new temporal client
	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "user-namespace",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// split productString
	stringArr := strings.Split(productString, ",")
	data := messages.Cart{Products: stringArr}

	we, err := c.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		ID:        "SetCartWorkflow_" + uuid.New().String(),
		TaskQueue: "cart",
	}, workflow.SetCartWorkflow, &data)
	if err != nil {
		http.Error(w, "unable to start workflow", http.StatusInternalServerError)
		return
	}

	if err := we.Get(r.Context(), nil); err != nil {
		http.Error(w, "unable to get workflow result", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal("{status: 'ok'}")
	if err != nil {
		http.Error(w, "unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "unable to write response", http.StatusInternalServerError)
		return
	}
}
