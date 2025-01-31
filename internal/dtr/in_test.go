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

package dtr_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/homeport/duct-tape-resource/internal/dtr"
)

var _ = Describe("In", func() {
	Context("valid configuration", func() {
		It("should return the given version", func() {
			version := Version{"ref": "foobar"}
			result, err := In(feed(Config{
				Source: Source{
					In: Custom{Run: "true"},
				},
				Version: version,
			}))

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Version).To(Equal(version))
		})

		It("should return the given version and metadata if output was available", func() {
			version := Version{"ref": "foobar"}
			result, err := In(feed(Config{
				Source: Source{
					In: Custom{Run: "echo foo bar && echo bar foo"},
				},
				Version: version,
			}))

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Version).To(Equal(version))
			Expect(result.Metadata).To(Equal([]Metadata{
				{Name: "foo", Value: "bar"},
				{Name: "bar", Value: "foo"},
			}))
		})

		It("should have access to the OS arguments", func() {
			version := Version{"ref": "foobar"}
			result, err := In(
				feed(Config{
					Source: Source{
						In: Custom{Run: "echo foo $1"},
					},
					Version: version,
				}),
				"bar",
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Version).To(Equal(version))
			Expect(result.Metadata).To(Equal([]Metadata{
				{Name: "foo", Value: "bar"},
			}))
		})

		It("should store the input version in a file", func() {
			withTempDir(func(dir string) {
				version := Version{"ref": "foobar"}
				result, err := In(
					feed(Config{
						Source:  Source{In: Custom{Run: "true"}},
						Version: version,
					}),
					dir,
				)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Version).To(Equal(version))

				refFile := filepath.Join(dir, "ref")
				Expect(refFile).To(BeAnExistingFile())

				data, err := os.ReadFile(refFile)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal([]byte("foobar")))
			})
		})
	})

	Context("empty/no-op configuration", func() {
		It("should not fail", func() {
			_, err := In(feed(Config{}))

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return given version", func() {
			version := Version{"ref": "foobar"}
			result, err := In(feed(Config{
				Version: version,
			}))

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(InOutResult{
				Version: version,
			}))
		})
	})
})
