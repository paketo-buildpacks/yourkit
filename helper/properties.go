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

package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
)

type Properties struct {
	Logger bard.Logger
}

func (p Properties) Execute() (map[string]string, error) {
	if _, ok := os.LookupEnv("BPL_YOURKIT_ENABLED"); !ok {
		return nil, nil
	}

	agentPath, ok := os.LookupEnv("BPI_YOURKIT_AGENT_PATH")
	if !ok {
		return nil, fmt.Errorf("$BPI_YOURKIT_AGENT_PATH must be set")
	}

	snapshotPath, ok := os.LookupEnv("BPI_YOURKIT_SNAPSHOT_PATH")
	if !ok {
		return nil, fmt.Errorf("$BPI_YOURKIT_SNAPSHOT_PATH must be set")
	}

	logPath, ok := os.LookupEnv("BPI_YOURKIT_LOG_PATH")
	if !ok {
		return nil, fmt.Errorf("$BPI_YOURKIT_LOG_PATH must be set")
	}

	port := "10001"
	if s, ok := os.LookupEnv("BPL_YOURKIT_PORT"); ok {
		port = s
	}

	sessionName := os.Getenv("BPL_YOURKIT_SESSION_NAME")

	s := "YourKit session"
	if sessionName != "" {
		s = fmt.Sprintf("%s %s", s, sessionName)
	}
	s = fmt.Sprintf("%s started on port %s", s, port)
	p.Logger.Info(s)

	var values []string
	if s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); ok {
		values = append(values, s)
	}

	values = append(values, fmt.Sprintf(
		"-agentpath:%s=dir=%s,logdir=%s,port=%s,sessionname=%s", agentPath, snapshotPath, logPath, port, sessionName))

	return map[string]string{"JAVA_TOOL_OPTIONS": strings.Join(values, " ")}, nil
}
