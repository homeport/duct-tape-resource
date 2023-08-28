// Copyright © 2022 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dtr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var defaultCommand = "/bin/bash"

type Config struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Source struct {
	ID string `json:"id"`

	Check Custom `json:"check"`
	In    Custom `json:"in"`
	Out   Custom `json:"out"`
}

type Version map[string]string

type Custom struct {
	Env map[string]string `json:"env"`

	Before string `json:"before"`
	Run    string `json:"run"`
}

type CheckResult []Version

type InOutResult struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func LoadConfig(in io.Reader) (*Config, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// In case no explicit ID is defined, used generic term "ref"
	if config.Source.ID == "" {
		config.Source.ID = "ref"
	}

	return &config, nil
}

func command(ctx context.Context, envs map[string]string, run string, out io.Writer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, defaultCommand, "-c", run)
	cmd.Stdin = nil
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	cmd.Env = cmd.Environ()
	for key, value := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	return cmd
}

func execute(entry Custom) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if entry.Before != "" {
		before := command(ctx, entry.Env, entry.Before, os.Stderr)
		if err := before.Run(); err != nil {
			return nil, fmt.Errorf("failure while running before command: %w", err)
		}
	}

	var outStream bytes.Buffer
	// User does not always need to define Run
	// e.g. user just implements Check but not In to trigger jobs
	// In this scenario a user would still use get on the resource. Then In would be executed.
	// For this scenario we would still like to allow no-op In runs
	// Therefore we can skip command execution if nothing is defined
	if entry.Run != "" {
		run := command(ctx, entry.Env, entry.Run, &outStream)
		if err := run.Run(); err != nil {
			return nil, fmt.Errorf("failure while running run command: %w", err)
		}
	}

	var result []string
	for _, entry := range strings.Split(outStream.String(), "\n") {
		if entry != "" {
			result = append(result, strings.Trim(entry, " "))
		}
	}

	return result, nil
}

func metadata(list []string) []Metadata {
	var metadata []Metadata
	for _, entry := range list {
		split := strings.SplitN(entry, " ", 2)
		metadata = append(metadata, Metadata{
			Name:  split[0],
			Value: split[1]},
		)
	}

	return metadata
}
