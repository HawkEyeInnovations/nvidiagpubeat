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
	"errors"
	"strconv"
	"strings"
	"os/exec"
	"encoding/xml"

	"github.com/elastic/beats/libbeat/common"
)

//GPUUtilization provides interface to utilization metrics and state of GPU.
type GPUUtilization interface {
	command(env string) *exec.Cmd
	run(cmd *exec.Cmd, gpuCount int, query string, action Action) ([]common.MapStr, error)
}

//Utilization implements one flavour of GPUCount interface.
type Utilization struct {
}

//newUtilization returns instance of Utilization
func newUtilization() Utilization {
	return Utilization{}
}

func (g Utilization) command(env string, query string) *exec.Cmd {
	if env == "test" {
		return exec.Command("localnvidiasmi")
	}
	return exec.Command("nvidia-smi", "-q", "-x")
}

type TrimmedInt int64

func(t *TrimmedInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	strs := strings.Split(v, " ")

	if len(strs) == 0 {
		return errors.New("No values in string " + v)
	}

	i, err := strconv.ParseInt(strs[0], 10, 64);
	if err != nil {
		return err
	}

	*t = TrimmedInt(i);

	return nil
}

//Run the nvidiasmi command to collect GPU metrics
//Parse output and return events.
func (g Utilization) run(cmd *exec.Cmd, gpuCount int, query string, action Action) ([]common.MapStr, error) {
	reader := action.start(cmd)
	events := make([]common.MapStr, gpuCount, 2*gpuCount)

	decoder := xml.NewDecoder(reader)

	type Utilization struct {
		GPU TrimmedInt `xml:"gpu_util" json:"gpu"`
		Memory TrimmedInt `xml:"memory_util" json:"memory"`
		Encoder TrimmedInt `xml:"encoder_util" json:"encoder"`
		Decoder TrimmedInt `xml:"decoder_util" json:"decoder"`
	}

	type Temperature struct {
		GPU TrimmedInt `xml:"gpu_temp" json:"gpu"`
	}

	type Memory struct {
		Total TrimmedInt `xml:"total" json:"total"`
		Used TrimmedInt `xml:"used" json:"used"`
		Free TrimmedInt `xml:"free" json:"free"`
	}

	type GPU struct {
		Name string `xml:"product_name"`
		Utilization Utilization `xml:"utilization"`
		Temperature Temperature `xml:"temperature"`
		Memory Memory `xml:"fb_memory_usage"`
	}

	type Data struct {
		XMLName xml.Name `xml:"nvidia_smi_log"`
		DriverVersion string `xml:"driver_version"`
		GPU []GPU `xml:"gpu"`
	}

	v := Data{}

	err := decoder.Decode(&v)
	if err != nil {
		return nil, err
	}

	for i, gpu := range v.GPU {
		event := common.MapStr{
			"gpu_index": i,
			"type": "nvidiagpubeat",
			"driver_version": v.DriverVersion,
			"name" : gpu.Name,
			"utilization": gpu.Utilization,
			"temperature" : gpu.Temperature,
			"memory" : gpu.Memory,
		}

		// event.Put("name", gpu.Name)
		// event.Put("utilziation", gpu.Utilization)
		// util := common.MapStr{
		// 	"gpu" : gpu.Utilization.GPU,
		// 	"memory" : gpu.Utilization.Memory,
		// 	"encoder" : gpu.Utilization.Encoder,
		// 	"decoder" : gpu.Utilization.Decoder,
		// }

		// event.Put("utilization", util)
		events[i] = event
	}

	cmd.Wait()
	return events, nil
}
