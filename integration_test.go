package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"
	"reflect"
)

var binPath string

func TestMain(m *testing.M) {
	tmpFile, err := ioutil.TempFile("", "secretly")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	binPath = tmpFile.Name()
	if err := exec.Command("go", "build", "-o", binPath).Run(); err != nil {
		log.Fatal(err)
	}

	ex := m.Run()
	os.Remove(binPath) // defer doesn't run before os exit
	os.Exit(ex)
}

func Test_cliEnv(t *testing.T) {
	tests := []struct {
		name    string
		environ []string
		wantEnv []string
		wantErr bool
	}{
		{"totally empty", []string{}, []string{}, false},
		{"passes through", []string{"FOO_BAR=BAZ"}, []string{"FOO_BAR=BAZ"}, false},
		{"AWS error", []string{"SECRETLY_NAMESPACE=BAZ"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// we're going to override the PATH var, so look it up
			envPath, err := exec.LookPath("env")
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(binPath, envPath)
			cmd.Env = tt.environ
			out, err := cmd.Output()
			if err != nil {
				if !tt.wantErr {
					t.Fatal(err)
				}
				// errored as expected
				return
			}

			// always allocate array b/c reflect.DeepEqual treats empty and nil slices differently
			outputEnv := make([]string, 0)
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Trim(line, " ") != "" {
					outputEnv = append(outputEnv, strings.Trim(line, " "))
				}
			}
			sort.Strings(outputEnv)

			if !reflect.DeepEqual(outputEnv, tt.wantEnv) {
				t.Errorf("cli = got %v, want %v", outputEnv, tt.wantEnv)
			}
		})
	}
}
