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
	"encoding/xml"
	"errors"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/elastic/beats/libbeat/common"
)

// Utilization implements one flavour of GPUCount interface.
type Utilization struct {
}

// NewUtilization returns instance of Utilization
func NewUtilization() Utilization {
	return Utilization{}
}

func (g Utilization) command() *exec.Cmd {
	return exec.Command("nvidia-smi", "-q", "-x")
}

// TrimmedInt allows for an Unmarshal which will crop off unit types
type TrimmedInt int64

// UnmarshalXML for the TrimmedInt type will trim off units from the string
// e.g. "48 MB" becomes "48"
func (t *TrimmedInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	strs := strings.Split(v, " ")

	if len(strs) == 0 {
		return errors.New("No values in string " + v)
	}

	i, err := strconv.ParseInt(strs[0], 10, 64)
	if err != nil {
		return err
	}

	*t = TrimmedInt(i)

	return nil
}

// PopulateArgs takes in a reflected type and adds the field and its value
// to the libbeat map if it has been configured to add it in the Map
func PopulateArgs(m Map, val reflect.Value, event *common.MapStr) {
	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		valueField := val.Field(i)
		typeField := valType.Field(i)

		tag := typeField.Tag.Get("json")

		if _, ok := m[tag]; ok {
			event.Put(tag, valueField.Interface())
		}
	}
}

//Run the nvidiasmi command to collect GPU metrics
//Parse output and return events.
func (g Utilization) run(cmd *exec.Cmd, query Query, action Action) ([]common.MapStr, error) {
	reader := action.start(cmd)

	decoder := xml.NewDecoder(reader)

	type Utilization struct {
		GPU     TrimmedInt `xml:"gpu_util" json:"gpu"`
		Memory  TrimmedInt `xml:"memory_util" json:"memory"`
		Encoder TrimmedInt `xml:"encoder_util" json:"encoder"`
		Decoder TrimmedInt `xml:"decoder_util" json:"decoder"`
	}

	type Temperature struct {
		GPU TrimmedInt `xml:"gpu_temp" json:"gpu"`
	}

	type Memory struct {
		Total TrimmedInt `xml:"total" json:"total"`
		Used  TrimmedInt `xml:"used" json:"used"`
		Free  TrimmedInt `xml:"free" json:"free"`
	}

	type GPU struct {
		Name        string      `xml:"product_name" json:"name"`
		Utilization Utilization `xml:"utilization" json:"utilization"`
		Temperature Temperature `xml:"temperature" json:"temperature"`
		Memory      Memory      `xml:"fb_memory_usage" json:"memory"`
	}

	type Data struct {
		XMLName       xml.Name `xml:"nvidia_smi_log"`
		CudaVersion   string   `xml:"cuda_version" json:"cuda_version"`
		DriverVersion string   `xml:"driver_version" json:"driver_version"`
		GPU           []GPU    `xml:"gpu"`
	}

	v := Data{}

	err := decoder.Decode(&v)
	if err != nil {
		return nil, err
	}

	gpuCount := len(v.GPU)

	events := make([]common.MapStr, gpuCount, gpuCount)

	for i, gpu := range v.GPU {
		event := common.MapStr{
			"gpu_index": i,
			"type":      "nvidiagpubeat",
		}

		PopulateArgs(query.System, reflect.ValueOf(&v).Elem(), &event)
		PopulateArgs(query.GPU, reflect.ValueOf(&gpu).Elem(), &event)

		events[i] = event
	}

	cmd.Wait()
	return events, nil
}
