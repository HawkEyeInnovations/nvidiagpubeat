//******************************************************************
//Copyright 2018 eBay Inc.
//Architect/Developer: Deepak Vasthimal

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at

// https://www.apache.org/licenses/LICENSE-2.0

//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//******************************************************************

package nvidia

import (
	"github.com/elastic/beats/libbeat/common"
)

//Metrics implements one flavour of GPUMetrics interface.
type Metrics struct {
}

//NewMetrics returns instance of Metrics
func NewMetrics() Metrics {
	return Metrics{}
}

//Get return a slice of GPU metrics
func (m Metrics) Get(query Query) ([]common.MapStr, error) {

	gpuUtilization := NewUtilization()
	gpuUtilizationCmd := gpuUtilization.command()
	events, err := gpuUtilization.run(gpuUtilizationCmd, query, NewLocal())
	return events, err
}
