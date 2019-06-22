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
	"io"
	"os/exec"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
)

//GPUCount provides interface to get gpu count command and run it.
type GPUCount interface {
	command() *exec.Cmd
	run(cmd *exec.Cmd, env string, action Action) (int, error)
}

//Count implements one flavour of GPUCount interface.
type Count struct {
}

//NewCount returns instance of Count
func newCount() Count {
	return Count{}
}

func (g Count) command() *exec.Cmd {
	return exec.Command("nvidia-smi", "--list-gpus")
}

func (g Count) run(cmd *exec.Cmd, env string, action Action) (int, error) {
	if env == "test" {
		return 4, nil
	}
	reader := action.start(cmd)
	gpuCount := 0
	for {
		line, err := reader.ReadString('\n')
		
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}
		logp.Info("Line: %v", line)
		if len(line) != 0 && !strings.HasPrefix(line, "GPU ") {
			return -1, errors.New("Unable to query GPUs")
		}

		
		gpuCount++
	}
	return gpuCount, nil
}
