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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/homeport/duct-tape-resource/internal/dtr"
)

var _ = Describe("Check", func() {
	Context("invalid configuration", func() {
		It("should fail if no run command is configured", func() {
			_, err := Check(feed(Config{}))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("valid configuration", func() {
		It("should just fail nicely if the provided command returns with non-zero exit code", func() {
			_, err := Check(feed(Config{
				Source: Source{
					Check: Custom{Run: "false"},
				},
			}))

			Expect(err).To(HaveOccurred())
		})

		It("should return provided version when no line is return by run", func() {
			versionIn := Version{"ref": "barfoo"}
			result, err := Check(feed(Config{
				Source: Source{
					Check: Custom{Run: "true"},
				},
				Version: versionIn,
			}))

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(CheckResult{versionIn}))
		})

		It("should return two versions when two lines are returned", func() {
			result, err := Check(feed(Config{
				Source: Source{
					Check: Custom{Run: "echo foobar && echo barfoo"},
				},
			}))

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(CheckResult{
				Version{"ref": "foobar"},
				Version{"ref": "barfoo"},
			}))
		})

		It("should run before commands before returning a version", func() {
			result, err := Check(feed(Config{
				Source: Source{
					Check: Custom{
						Before: "echo foobar >/tmp/version",
						Run:    "cat /tmp/version",
					},
				},
			}))

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(CheckResult{
				Version{"ref": "foobar"},
			}))
		})

		It("should just fail nicely if the provided before command return with non-zero exit code", func() {
			_, err := Check(feed(Config{
				Source: Source{
					Check: Custom{
						Before: "false",
						Run:    "true",
					},
				},
			}))

			Expect(err).To(HaveOccurred())
		})

		It("should set environment variables that can be used in run", func() {
			result, err := Check(feed(Config{
				Source: Source{
					Check: Custom{
						Env: map[string]string{
							"FOO": "foo",
							"BAR": "bar",
						},
						Run: "echo ${FOO}${BAR}",
					},
				},
			}))

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(CheckResult{
				Version{"ref": "foobar"},
			}))
		})
	})
})
