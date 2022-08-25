/*
 * Copyright 2018-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper_test

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/yourkit/v4/helper"
)

func testProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		p = helper.Properties{}
	)

	it("returns if $BPL_YOURKIT_ENABLED is not set", func() {
		Expect(p.Execute()).To(BeNil())
	})

	context("$BPL_YOURKIT_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv("BPL_YOURKIT_ENABLED", "")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BPL_YOURKIT_ENABLED")).To(Succeed())
		})

		it("returns error if $BPI_YOURKIT_AGENT_PATH is not set", func() {
			_, err := p.Execute()
			Expect(err).To(MatchError("$BPI_YOURKIT_AGENT_PATH must be set"))
		})

		context("$BPI_YOURKIT_AGENT_PATH", func() {
			it.Before(func() {
				Expect(os.Setenv("BPI_YOURKIT_AGENT_PATH", "test-agent-path")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BPI_YOURKIT_AGENT_PATH")).To(Succeed())
			})

			it("returns error of $BPI_YOURKIT_SNAPSHOT_PATH is not set", func() {
				_, err := p.Execute()
				Expect(err).To(MatchError("$BPI_YOURKIT_SNAPSHOT_PATH must be set"))
			})

			context("$BPI_YOURKIT_SNAPSHOT_PATH", func() {
				it.Before(func() {
					Expect(os.Setenv("BPI_YOURKIT_SNAPSHOT_PATH", "test-snapshot-path")).To(Succeed())
				})

				it.After(func() {
					Expect(os.Unsetenv("BPI_YOURKIT_SNAPSHOT_PATH")).To(Succeed())
				})

				it("returns error of $BPI_YOURKIT_LOG_PATH is not set", func() {
					_, err := p.Execute()
					Expect(err).To(MatchError("$BPI_YOURKIT_LOG_PATH must be set"))
				})

				context("$BPI_YOURKIT_LOG_PATH", func() {
					it.Before(func() {
						Expect(os.Setenv("BPI_YOURKIT_LOG_PATH", "test-log-path")).To(Succeed())
					})

					it.After(func() {
						Expect(os.Unsetenv("BPI_YOURKIT_LOG_PATH")).To(Succeed())
					})

					it("contributes configuration", func() {
						Expect(p.Execute()).To(Equal(map[string]string{
							"JAVA_TOOL_OPTIONS": "-agentpath:test-agent-path=dir=test-snapshot-path,logdir=test-log-path,port=10001,sessionname=",
						}))
					})

					context("$BPL_YOURKIT_PORT", func() {
						it.Before(func() {
							Expect(os.Setenv("BPL_YOURKIT_PORT", "10002")).To(Succeed())
						})

						it.After(func() {
							Expect(os.Unsetenv("BPL_YOURKIT_PORT")).To(Succeed())
						})

						it("contributes port configuration from $BPL_YOURKIT_PORT", func() {
							Expect(p.Execute()).To(Equal(map[string]string{
								"JAVA_TOOL_OPTIONS": "-agentpath:test-agent-path=dir=test-snapshot-path,logdir=test-log-path,port=10002,sessionname=",
							}))
						})
					})

					context("$BPL_YOURKIT_SESSION_NAME", func() {
						it.Before(func() {
							Expect(os.Setenv("BPL_YOURKIT_SESSION_NAME", "test-session-name")).To(Succeed())
						})

						it.After(func() {
							Expect(os.Unsetenv("BPL_YOURKIT_SESSION_NAME")).To(Succeed())
						})

						it("contributes session name configuration from $BPL_YOURKIT_SESSION_NAME", func() {
							Expect(p.Execute()).To(Equal(map[string]string{
								"JAVA_TOOL_OPTIONS": "-agentpath:test-agent-path=dir=test-snapshot-path,logdir=test-log-path,port=10001,sessionname=test-session-name",
							}))
						})
					})

					context("$JAVA_TOOL_OPTIONS", func() {
						it.Before(func() {
							Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
						})

						it.After(func() {
							Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
						})

						it("contributes configuration appended to existing $JAVA_TOOL_OPTIONS", func() {
							Expect(p.Execute()).To(Equal(map[string]string{
								"JAVA_TOOL_OPTIONS": "test-java-tool-options -agentpath:test-agent-path=dir=test-snapshot-path,logdir=test-log-path,port=10001,sessionname=",
							}))
						})
					})

				})
			})
		})
	})
}
