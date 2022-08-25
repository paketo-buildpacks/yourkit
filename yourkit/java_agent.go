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

package yourkit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type JavaAgent struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaAgent(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (JavaAgent, libcnb.BOMEntry) {
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return JavaAgent{LayerContributor: contributor}, entry
}

func (j JavaAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Expanding to %s", layer.Path)

		if err := crush.ExtractZip(artifact, layer.Path, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand YourKit\n%w", err)
		}

		layer.LaunchEnvironment.Default("BPI_YOURKIT_AGENT_PATH", filepath.Join(layer.Path, "bin", "linux-x86-64", "libyjpagent.so"))
		layer.LaunchEnvironment.Default("BPI_YOURKIT_SNAPSHOT_PATH", filepath.Join(layer.Path, "snapshots"))
		layer.LaunchEnvironment.Default("BPI_YOURKIT_LOG_PATH", filepath.Join(layer.Path, "logs"))

		return layer, nil
	})
}

func (j JavaAgent) Name() string {
	return j.LayerContributor.LayerName()
}
