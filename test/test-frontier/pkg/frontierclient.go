package pkg

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	frontierV1 "github.com/nlnwa/veidemann-api-go/frontier/v1"
	"github.com/nlnwa/veidemann/dev/test/test-frontier/pkg/db"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

type FrontierClient struct {
	conn *FrontierConn
}

func NewFrontierClient(conn *FrontierConn) *FrontierClient {
	return &FrontierClient{
		conn: conn,
	}
}

func (f *FrontierClient) CrawlSeed(jobExecutionId string, job, seed *configV1.ConfigObject) (crawlExecutionId string, err error) {
	ceid, err := f.conn.Client().CrawlSeed(context.Background(), &frontierV1.CrawlSeedRequest{
		JobExecutionId: jobExecutionId,
		Job:            job,
		Seed:           seed,
	})
	if err == nil {
		return ceid.Id, nil
	}
	return "", err
}

//func (f *FrontierClient) RunJob(database *db.DbConnection, done chan interface{}) {
//	jobConfig, err := database.GetConfig(&configV1.ConfigRef{
//		Kind: configV1.Kind_crawlJob,
//		Id:   "e46863ae-d076-46ca-8be3-8a8ef72e709e",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	jeid, err := database.CreateJob(&frontierV1.JobExecutionStatus{
//		State: frontierV1.JobExecutionStatus_CREATED,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Created job: %v\n", jeid)
//	go func() {
//		fmt.Printf("JOB DONE %v\n", database.WaitForJobCompletion(jeid))
//		close(done)
//	}()
//
//	for i := 0; i < 10; i++ {
//		seed, err := database.CreateSeed(jeid, fmt.Sprintf("http://foo%02d.bar/", i))
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("Created seed: %v\n", seed.Id)
//
//		eid, err := f.conn.Client().CrawlSeed(context.Background(), &frontierV1.CrawlSeedRequest{
//			JobExecutionId: jeid,
//			Job:            jobConfig,
//			Seed:           seed,
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		log.Infof("JEID: %v, EID: %v", jeid, eid.Id)
//	}
//	//for i := 0; i < 5; i++ {
//	//	seed, err := database.CreateSeed(jeid, fmt.Sprintf("http://foo%02d.bar/test01", i))
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//	fmt.Printf("Created seed: %v\n", seed.Id)
//	//
//	//	eid, err := f.conn.Client().CrawlSeed(context.Background(), &frontierV1.CrawlSeedRequest{
//	//		JobExecutionId: jeid,
//	//		Job:            jobConfig,
//	//		Seed:           seed,
//	//	})
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//
//	//	log.Infof("JEID: %v, EID: %v", jeid, eid.Id)
//	//}
//}

func (f *FrontierClient) Harvester(database *db.DbConnection) {
	for {
		stream, err := f.conn.Client().GetNextPage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		waitc := make(chan struct{}) // Closed when session is completed/failed
		waitn := make(chan struct{}) // Closed when initial response from Frontier is received
		var harvestSpec *frontierV1.PageHarvestSpec
		go func() {
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					// read done.
					close(waitc)
					return
				}
				if err != nil {
					log.Errorf("Failed to receive Frontier response: %v", err)
					close(waitc)
					return
				}
				harvestSpec = in
				close(waitn)
			}
		}()
		if err := stream.Send(&frontierV1.PageHarvest{Msg: &frontierV1.PageHarvest_RequestNextPage{RequestNextPage: true}}); err != nil {
			log.Errorf("Failed to send request for new page to Frontier: %v", err)
			_ = stream.CloseSend()
			return
		}
		select {
		case <-waitn:
		}

		log.WithField("uri", harvestSpec.QueuedUri.Uri).Infof("Starting fetch")

		//if renderResult.Error != nil {
		//	if err := stream.Send(&frontierV1.PageHarvest{Msg: &frontierV1.PageHarvest_Error{Error: renderResult.Error}}); err != nil {
		//		log.Errorf("Failed to send error response to Frontier: %v", err)
		//		_ = stream.CloseSend()
		//		return
		//	}
		//} else {
		if err := stream.Send(&frontierV1.PageHarvest{
			Msg: &frontierV1.PageHarvest_Metrics_{
				Metrics: &frontierV1.PageHarvest_Metrics{
					UriCount:        1,
					BytesDownloaded: 1,
				},
			},
		}); err != nil {
			log.Errorf("Failed to send metrics to Frontier: %v", err)
			_ = stream.CloseSend()
			return
		}

		//for _, outlink := range renderResult.Outlinks {
		for i := 0; i < 2; i++ {
			u := harvestSpec.QueuedUri.Uri
			u = u[:strings.LastIndex(u, "/")]
			outlink := &frontierV1.QueuedUri{
				ExecutionId:         harvestSpec.QueuedUri.ExecutionId,
				DiscoveredTimeStamp: ptypes.TimestampNow(),
				Uri:                 fmt.Sprintf("%s/test%02d", u, i),
				DiscoveryPath:       harvestSpec.QueuedUri.DiscoveryPath + "L",
				Referrer:            harvestSpec.QueuedUri.Uri,
				PageFetchTimeMs:     1,
				JobExecutionId:      harvestSpec.QueuedUri.JobExecutionId,
			}
			if err := stream.Send(&frontierV1.PageHarvest{Msg: &frontierV1.PageHarvest_Outlink{Outlink: outlink}}); err != nil {
				log.Errorf("Failed to send outlink to Frontier: %v", err)
			}
		}
		//}

		_ = stream.CloseSend()
		<-waitc
	}
}
