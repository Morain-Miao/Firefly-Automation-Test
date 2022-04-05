package task_worker

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, TransferMoneyTaskQueue, worker.Options{})
	w.RegisterWorkflow(TransferMoney)
	w.RegisterActivity(Withdraw)
	w.RegisterActivity(Deposit)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

// @@@SNIPEND
