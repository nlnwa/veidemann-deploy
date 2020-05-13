/*
 * Copyright 2020 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package db

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	frontierV1 "github.com/nlnwa/veidemann-api-go/frontier/v1"
	log "github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/encoding"
)

//type DbConnection interface {
//	Connect() error
//	GetConfig(ref *configV1.ConfigRef) (*configV1.ConfigObject, error)
//	GetConfigsForSelector(kind configV1.Kind, label *configV1.Label) ([]*configV1.ConfigObject, error)
//	WriteCrawlLog(crawlLog *frontierV1.CrawlLog) error
//	WritePageLog(pageLog *frontierV1.PageLog) error
//}

// connection holds the connections for ContentWriter and Veidemann database
type DbConnection struct {
	dbConnectOpts r.ConnectOpts
	dbSession     r.QueryExecutor
}

// NewConnection creates a new connection object
func NewConnection(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string) *DbConnection {
	c := &DbConnection{
		dbConnectOpts: r.ConnectOpts{
			Address:    fmt.Sprintf("%s:%d", dbHost, dbPort),
			Username:   dbUser,
			Password:   dbPassword,
			Database:   dbName,
			NumRetries: 10,
		},
	}
	return c
}

// connect establishes connections
func (c *DbConnection) Connect() error {
	// Set up database connection
	dbSession, err := r.Connect(c.dbConnectOpts)
	if err != nil {
		log.Errorf("fail to connect to database: %v", err)
		return err
	}
	c.dbSession = dbSession

	log.Infof("Recorder Proxy is using DB at: %s", c.dbConnectOpts.Address)
	return nil
}

func (c *DbConnection) GetConfig(ref *configV1.ConfigRef) (*configV1.ConfigObject, error) {
	res, err := r.Table("config").Get(ref.Id).Run(c.dbSession)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var result configV1.ConfigObject
	err = res.One(&result)

	if err != nil {
		return nil, fmt.Errorf("DB error: %w", err)
	}
	return &result, nil
}

func (c *DbConnection) GetConfigsForSelector(kind configV1.Kind, label *configV1.Label) ([]*configV1.ConfigObject, error) {
	res, err := r.Table("config").GetAllByIndex("label", r.Expr([]string{label.Key, label.Value})).Run(c.dbSession)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var r []configV1.ConfigObject
	err = res.All(&r)
	if err != nil {
		return nil, fmt.Errorf("DB error: %w", err)
	}

	var result []*configV1.ConfigObject
	for _, co := range r {
		if co.Kind == kind {
			result = append(result, &co)
		}
	}

	return result, nil
}

func (c *DbConnection) WriteCrawlLog(crawlLog *frontierV1.CrawlLog) error {
	_, err := r.Table("crawl_log").Insert(crawlLog).RunWrite(c.dbSession)
	return err
}

func (c *DbConnection) WritePageLog(pageLog *frontierV1.PageLog) error {
	_, err := r.Table("page_log").Insert(pageLog).RunWrite(c.dbSession)
	return err
}

func (c *DbConnection) WaitForJobCompletion(jobId string) error {
	resp, err := r.Table("job_executions").Get(jobId).Changes(r.ChangesOpts{IncludeInitial: true}).Map(r.Row.Field("new_val")).Run(c.dbSession)
	if err != nil {
		return err
	}

	var job frontierV1.JobExecutionStatus
	for {
		if ok := resp.Next(&job); ok {
			log.Debugf("JOB: %v\n", job.State)
			if job.State == frontierV1.JobExecutionStatus_FINISHED {
				return nil
			}
		} else {
			return resp.Err()
		}
	}
}

func (c *DbConnection) Prune(table string) {
	var res []interface{}

	if resp, err := r.Table(table).Delete().Run(c.dbSession); err == nil {
		if err = resp.All(&res); err != nil {
			log.Errorf("error getting delete status from %s: %v", table, err)
		}
		log.Infof("Delete '%s': %v\n", table, res)
	} else {
		log.Errorf("error deleting data from %s: %v", table, err)
	}
}

func (c *DbConnection) ResetDB() {
	c.Prune("config_seeds")
	c.Prune("config_crawl_entities")
	c.Prune("uri_queue")
	c.Prune("executions")
	c.Prune("job_executions")
}

func (c *DbConnection) CreateJob(job *frontierV1.JobExecutionStatus) (string, error) {
	resp, err := r.Table("job_executions").Insert(job).RunWrite(c.dbSession)

	return resp.GeneratedKeys[0], err
}

func (c *DbConnection) CreateSeed(jobid, uri string) (*configV1.ConfigObject, error) {
	entity := &configV1.ConfigObject{
		ApiVersion: "v1",
		Kind:       configV1.Kind_crawlEntity,
		Meta: &configV1.Meta{
			Name:  uri,
			Label: []*configV1.Label{&configV1.Label{Key: "foo", Value: "bar"}},
		},
	}
	e, err := r.Table("config_crawl_entities").Insert(entity).RunWrite(c.dbSession)
	if err != nil {
		return nil, err
	}

	seed := &configV1.ConfigObject{
		ApiVersion: "v1",
		Kind:       configV1.Kind_seed,
		Meta: &configV1.Meta{
			Name:    uri,
			Label:   []*configV1.Label{&configV1.Label{Key: "rompa", Value: "bar"}},
			Created: ptypes.TimestampNow(),
		},
		Spec: &configV1.ConfigObject_Seed{
			Seed: &configV1.Seed{
				EntityRef: &configV1.ConfigRef{Kind: configV1.Kind_crawlEntity, Id: e.GeneratedKeys[0]},
				JobRef:    []*configV1.ConfigRef{&configV1.ConfigRef{Kind: configV1.Kind_crawlJob, Id: jobid}},
			},
		},
	}
	s, err := r.Table("config_seeds").Insert(seed, r.InsertOpts{ReturnChanges: true}).RunWrite(c.dbSession)
	if err != nil {
		return nil, err
	}

	data := s.Changes[0].NewValue
	var dest configV1.ConfigObject
	err = encoding.Decode(&dest, data)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}
