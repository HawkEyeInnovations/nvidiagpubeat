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
	"os"
	"testing"
)

func TestCommand(t *testing.T) {
	util := NewUtilization()
	cmd := util.command()

	if len(cmd.Args) != 3 {
		t.Errorf("Expected %d, Actual %d", 3, len(cmd.Args))
	}

	if cmd.Args[0] != "nvidia-smi" {
		t.Errorf("Expected %s, Actual %s", "nvidia-smi", cmd.Args[0])
	}

	if cmd.Args[1] != "-q" {
		t.Errorf("Expected %s, Actual %s", "--q", cmd.Args[1])
	}

	if cmd.Args[2] != "-x" {
		t.Errorf("Expected %s, Actual %s", "-x", cmd.Args[2])
	}
}

func TestRunTestEnv(t *testing.T) {
	t.Log("Running test env")
	util := NewUtilization()
	cmd := util.command()
	_, err := util.run(cmd, MockQuery(), MockSingle{})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestRunProdEnv(t *testing.T) {
	util := NewUtilization()
	cmd := util.command()
	os.Setenv("PATH", ".")
	output, err := util.run(cmd, MockQuery(), NewLocal())

	if err != nil {
		t.Errorf("%v", err)
	}

	if output == nil {
		t.Errorf("Output cannot be nil")
	}

	for _, o := range output {
		if o == nil {
			t.Errorf("output cannot be nil.")
		}
	}
}

func TestEventContainsTypeField(t *testing.T) {
	util := NewUtilization()
	cmd := util.command()
	t.Logf("Command: %s", cmd.Path)

	output, _ := util.run(cmd, MockQuery(), MockSingle{})
	t.Logf("Nr. of Events: %d", len(output))

	for _, o := range output {
		if o["type"] != "nvidiagpubeat" {
			t.Errorf("event does not contain 'type' field equal to 'nvidiagpubeat'")
		}
	}
}

func TestEventCountCorrectSingle(t *testing.T) {
	util := NewUtilization()
	cmd := util.command()
	t.Logf("Command: %s", cmd.Path)

	output, _ := util.run(cmd, MockQuery(), MockSingle{})
	t.Logf("Nr. of Events: %d", len(output))

	if count := len(output); count != 1 {
		t.Errorf("Expected %d Actual %d", 1, count)
	}
}

func TestEventCountCorrectDual(t *testing.T) {
	util := NewUtilization()
	cmd := util.command()
	t.Logf("Command: %s", cmd.Path)

	output, _ := util.run(cmd, MockQuery(), MockDual{})
	t.Logf("Nr. of Events: %d", len(output))

	if count := len(output); count != 2 {
		t.Errorf("Expected %d Actual %d", 2, count)
	}
}
