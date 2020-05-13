package main

import (
	"fmt"
	"github.com/nlnwa/veidemann/dev/test/test-frontier/pkg"
	"github.com/nlnwa/veidemann/dev/test/test-frontier/pkg/db"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Infof("Starting Frontier tests")
	robots, err := pkg.NewRobotsEvaluatorMock("7100")
	if err != nil {
		log.Fatal(err)
	}
	dns, err := pkg.NewDnsResolverMock("7101")
	if err != nil {
		log.Fatal(err)
	}
	oos, err := pkg.NewOosHandlerMock("7102")
	if err != nil {
		log.Fatal(err)
	}

	frontierConn, err := pkg.NewFrontierConn("veidemann-frontier:7700")
	if err != nil {
		log.Fatal(err)
	}
	frontierClient := pkg.NewFrontierClient(frontierConn)

	log.Infof("Connecting DB")
	database := db.NewConnection("rethinkdb-proxy", 28015, "admin", "rethinkdb", "veidemann")
	database.Connect()
	log.Infof("DB connected")

	runTests(frontierClient, database)
	log.Infof("Frontier tests done")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Debugf("Got signal: %v", sig)
			panic("KILLED")
			frontierConn.Close()
			robots.Close()
			dns.Close()
			oos.Close()
			return
		}
	}()

}

func runTests(frontier *pkg.FrontierClient, database *db.DbConnection) {
	database.ResetDB()
	time.Sleep(2 * time.Second)
	go frontier.Harvester(database)

	done := make(chan interface{})
	time.Sleep(5 * time.Second)
	start := time.Now()
	RunJob(database, frontier, done)
	<-done
	fmt.Printf("Job finished in %v\n", time.Since(start))
	database.ResetDB()
	//time.Sleep(20*time.Second)
}
