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
	"io"
	"os"
	"path/filepath"
)

func In(in io.Reader, args ...string) (InOutResult, error) {
	config, err := LoadConfig(in)
	if err != nil {
		return InOutResult{}, err
	}

	if len(args) > 0 {
		// First argument should be a directory, but just to be on the safe side
		// let's check that it's a directory and only then put the version details
		// into a file in the input directory
		var dir = args[0]
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			name := filepath.Join(dir, config.Source.ID)
			data := []byte(config.Version[config.Source.ID])
			if err := os.WriteFile(name, data, os.FileMode(0644)); err != nil {
				return InOutResult{}, err
			}
		}
	}

	output, err := execute(config.Source.In, args...)
	if err != nil {
		return InOutResult{}, err
	}

	return InOutResult{
		Version:  config.Version,
		Metadata: metadata(output),
	}, nil
}
