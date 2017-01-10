package main

import (
	"log"
	"time"
	"os"
	"syscall"
	"os/signal"
)

const DefaultSyncFrequency = 10 * time.Second

type ServiceSyncerWithUaa interface {
	SyncUsers(currentUsers []UaaUser, uaaClient UaaClient) error
	SyncOrganizations(currentOrganizations []UaaOrganization, uaaClient UaaClient) error
}

func shutdownChan() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func mainLoop(frequency time.Duration, uaaClient UaaClient, syncers ... ServiceSyncerWithUaa) {
	tick := time.Tick(frequency)
	doneChan := shutdownChan()
	log.Println("Starting main loop")
	for {
		singleAllSync(uaaClient, syncers)
		select {
		case <- tick:
		case s := <- doneChan:
			log.Printf("Got: %v, ending main loop", s)
			close(doneChan)
			return
		}
	}
}

func singleAllSync(uaaClient UaaClient, syncers []ServiceSyncerWithUaa) {
	err := uaaClient.EnsureCredentials()
	if err != nil {
		log.Println("Problem ensuring UAA credentials", err)
		return
	}
	uaaUsers, err := uaaClient.GetUsers()
	if err != nil {
		log.Println("ERROR: couldn't get users from UAA", err)
		return
	}
	log.Printf("Single sync loop with %d users from UAA", len(uaaUsers))
	for _, syncer := range syncers {
		go syncer.SyncUsers(uaaUsers, uaaClient)
	}
}

func handlePotentialSetupError(err error, msg string) {
	if err != nil {
		log.Fatalf(msg + " %v", err)
	}
}

func main() {

	var err error
	frequency := DefaultSyncFrequency
	if f := os.Getenv("SYNC_FREQUENCY"); f != "" {
		frequency, err = time.ParseDuration(f)
		handlePotentialSetupError(err, "Coudln't parse sync frequency.")
	}

	uaaClient, err := NewUaaClientFromEnv()
	handlePotentialSetupError(err, "Error setting up UAA client.")

	grafanaSync, err := NewGrafanaSyncOperatorFromEnv()
	handlePotentialSetupError(err, "Couldn't create grafana syncer.")

	log.Println("Initialization completed")
	mainLoop(frequency, uaaClient, grafanaSync)
	log.Println("Going down")
}
