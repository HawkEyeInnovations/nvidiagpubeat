/***********************************************************************
This module was automatically generated using the framework found below:
https://www.elastic.co/guide/en/beats/devguide/current/new-beat.html

Modifications to auto-generated code - Copyright 2018 eBay Inc.
Architect/Developer: Deepak Vasthimal

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
************************************************************************/

package beater

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/paths"

	"github.com/ebay/nvidiagpubeat/config"
	"github.com/ebay/nvidiagpubeat/nvidia"
)

type Nvidiagpubeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

//New Creates the Beat object
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Nvidiagpubeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

//Run Contains the main application loop that captures data and sends it to the defined output using the publisher
func (bt *Nvidiagpubeat) Run(b *beat.Beat) error {
	logp.Info("nvidiagpubeat is running. Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	if bt.config.AddHomePath {
		logp.Info("Adding %v to PATH", paths.Paths.Home)
		os.Setenv("PATH", paths.Paths.Home+";"+os.Getenv("PATH"))
	}
	ticker := time.NewTicker(bt.config.Period)

	// Construct an instance of Query and for each type in the array add it into a map
	query := nvidia.NewQuery()

	for _, t := range bt.config.System {
		query.System[t] = struct{}{}
	}

	for _, t := range bt.config.GPU {
		query.GPU[t] = struct{}{}
	}

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		metrics := nvidia.NewMetrics()
		events, err := metrics.Get(query)
		if err != nil {
			logp.Err("Event not generated, error: %s", err.Error())
		} else {
			logp.Debug("nvidiagpubeat", "Event generated, Attempting to publish to configured output.")
			for _, gpuevent := range events {
				if gpuevent != nil {
					event := beat.Event{
						Timestamp: time.Now(),
						Fields:    gpuevent,
					}
					bt.client.Publish(event)
				}
			}
		}
	}
}

//Stop Contains logic that is called when the Beat is signaled to stop
func (bt *Nvidiagpubeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
