package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	testCases := []testCase{
		{
			name:      "std out echo",
			input:     args("-p", "testdata/project", "run", "hello_stdout"),
			stdoutput: "Hello stdout\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "std err echo",
			input:     args("-p", "testdata/project", "run", "hello_stderr"),
			stdoutput: "",
			erroutput: "Hello stderr\n",
			err:       nil,
		},
		{
			name:      "exit 0",
			input:     args("-p", "testdata/project", "run", "exit_0"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "exit 1",
			input:     args("-p", "testdata/project", "run", "exit_1"),
			stdoutput: ``,
			erroutput: "Error: exit code 4 - Failed executing script `exit_1`: shell script `exit 1`\nExit code: 1\n",
			err: errors.New(
				"exit code 4 - Failed executing script `exit_1`: shell script `exit 1`\nExit code: 1",
			),
		},
		{
			name:      "project with absolute path",
			input:     args("-p", filepath.Join(pwd, "testdata/project"), "run", "hello_stdout"),
			stdoutput: "Hello stdout\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "project without shuttle.yaml",
			input:     args("-p", "testdata/base", "run", "hello_stdout"),
			stdoutput: "",
			erroutput: "Error: shuttle run is not available in this context. To use shuttle run you need to be in a project with a shuttle.yaml file\n",
			err: errors.New(
				"shuttle run is not available in this context. To use shuttle run you need to be in a project with a shuttle.yaml file",
			),
		},
		{
			name:      "script fails when required argument is missing",
			input:     args("-p", "testdata/project", "run", "required_arg"),
			stdoutput: "",
			erroutput: `Error: required flag(s) "foo" not set
`,
			err: errors.New(`required flag(s) "foo" not set`),
		},
		{
			name:      "script succeeds with required argument",
			input:     args("-p", "testdata/project", "run", "required_arg", "foo=bar"),
			stdoutput: "bar\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "script succeeds with required argument with hash in value",
			input:     args("-p", "testdata/project", "run", "required_arg", "foo=bar="),
			stdoutput: "bar=\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "script succeeds with required argument missing and validation disabled",
			input:     args("-p", "testdata/project", "run", "--validate=false", "required_arg"),
			stdoutput: "\n",
			erroutput: "",
			err:       nil,
		},
		{
			name: "script fails when validation is disabled and argument is not in valid format",
			input: args(
				"-p",
				"testdata/project",
				"run",
				"--validate=false",
				"required_arg",
				"foo",
			),
			stdoutput: "\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "script fails on unkown argument",
			input:     args("-p", "testdata/project", "run", "required_arg", "foo=bar", "--a b"),
			stdoutput: "",
			erroutput: `Error: unknown flag: --a b
`,
			err: errors.New(`unknown flag: --a b`),
		},
		{
			name:      "branched git plan",
			input:     args("-p", "testdata/project-git-branched", "run", "say"),
			stdoutput: "something clever\n",
			erroutput: "Cloning plan https://github.com/kjuulh/shuttle-example-go-plan.git\n",
			err:       nil,
		},
		{
			name:      "git plan",
			input:     args("-p", "testdata/project-git", "run", "say"),
			stdoutput: "something masterly\n",
			erroutput: "Cloning plan https://github.com/kjuulh/shuttle-example-go-plan.git\n",
			initErr:   errors.New("something"),
		},
		{
			name:      "tagged git plan",
			input:     args("-p", "testdata/project-git", "--plan", "#tagged", "run", "say"),
			stdoutput: "something tagged\n",
			erroutput: "\x1b[032;1mOverload git plan branch/tag/sha with tagged\x1b[0m\nCloning plan https://github.com/kjuulh/shuttle-example-go-plan.git\n",
			err:       nil,
		},
		{
			name: "Local project",
			input: args(
				"--project",
				"./testdata/project-local/service",
				"--plan",
				"./testdata/project-local/plan",
				"run",
				"hello-plan",
			),
			stdoutput: "Hello from plan\n",
			erroutput: "Using overloaded plan ./testdata/project-local/plan\n",
			err:       nil,
		},
		// FIXME: This case actually hits a bug as we do not support fetching specific commits
		// {
		// 	name:      "sha git plan",
		// 	input:     args("-p", "testdata/project-git", "--plan", "#df4630118c7dfb594b4de903621681e677534638", "run", "say"),
		// 	stdoutput: "\x1b[032;1mOverload git plan branch/tag/sha with 2b52c21\x1b[0m\nCloning plan https://github.com/kjuulh/shuttle-example-go-plan.git\nsomething minor\n",
		// 	erroutput: "",
		// 	err:       nil,
		// },
	}
	executeTestCases(t, testCases)

	testContainsCases := []testCase{
		{
			name: "Local project fail",
			input: args(
				"--project",
				"./testdata/project-local/service",
				"--plan",
				"./testdata/wrong-project-local/plan",
				"run",
				"hello-plan",
			),
			stdoutput: "",
			erroutput: "shuttle/cmd/testdata/wrong-project-local/plan: no such file or directory",
			err: errors.New(
				"shuttle/cmd/testdata/wrong-project-local/plan: no such file or directory",
			),
			initErr: errors.New(
				`failed to copy plan to .shuttle/plan, make sure the upstream plan exists`,
			),
		},
	}
	executeTestContainsCases(t, testContainsCases)
}
