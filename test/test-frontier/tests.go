package main

import (
	"fmt"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	frontierV1 "github.com/nlnwa/veidemann-api-go/frontier/v1"
	"github.com/nlnwa/veidemann/dev/test/test-frontier/pkg"
	"github.com/nlnwa/veidemann/dev/test/test-frontier/pkg/db"
	log "github.com/sirupsen/logrus"
)

func RunJob(database *db.DbConnection, f *pkg.FrontierClient, done chan interface{}) {
	jobConfig, err := database.GetConfig(&configV1.ConfigRef{
		Kind: configV1.Kind_crawlJob,
		Id:   "e46863ae-d076-46ca-8be3-8a8ef72e709e",
	})
	if err != nil {
		log.Fatal(err)
	}

	jeid, err := database.CreateJob(&frontierV1.JobExecutionStatus{
		State: frontierV1.JobExecutionStatus_CREATED,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created job: %v\n", jeid)
	go func() {
		fmt.Printf("JOB DONE %v\n", database.WaitForJobCompletion(jeid))
		close(done)
	}()

	for i := 0; i < 10; i++ {
		seed, err := database.CreateSeed(jeid, fmt.Sprintf("http://foo%02d.bar/", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created seed: %v\n", seed.Id)

		eid, err := f.CrawlSeed(jeid, jobConfig, seed)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("JEID: %v, EID: %v", jeid, eid)
	}
	//for i := 0; i < 5; i++ {
	//	seed, err := database.CreateSeed(jeid, fmt.Sprintf("http://foo%02d.bar/test01", i))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Created seed: %v\n", seed.Id)
	//
	//	eid, err := f.conn.Client().CrawlSeed(context.Background(), &frontierV1.CrawlSeedRequest{
	//		JobExecutionId: jeid,
	//		Job:            jobConfig,
	//		Seed:           seed,
	//	})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	log.Infof("JEID: %v, EID: %v", jeid, eid.Id)
	//}
}
