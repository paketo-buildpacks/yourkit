// Copyright 2018-2026 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/paketo-buildpacks/libdependency/upstream"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

const (
	id   = "yourkit-java"
	name = "YourKit Java Agent"
	purl = "pkg:generic/yourkit-java-profiler"

	download     = "https://www.yourkit.com/java/profiler/download"
	downloadBase = "https://download.yourkit.com"
)

var (
	// Current download page format: "Version: 2026.9.149"
	versionPattern = regexp.MustCompile(`Version:\s*(\d+)\.(\d+)\.(\d+)`)
	// Older download page format: "Version: 2026.3 ... Build: #176"
	legacyVersionPattern = regexp.MustCompile(`(?s)Version:\s*(\d+)\.(\d+).*?Build:\s*#(\d+)`)
)

// assetArch maps a buildpack architecture to the architecture suffix used in
// the YourKit download file name.
var assetArch = map[string]string{
	"amd64": "x64",
	"arm64": "arm64",
}

func main() {
	retrieve.NewMetadataWithPlatforms(id, getAllVersions, generateMetadata)
}

func getAllVersions() (versionology.VersionFetcherArray, error) {
	resp, err := http.Get(download)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch %s\n%w", download, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch %s: status code %d", download, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s\n%w", download, err)
	}

	matches := versionPattern.FindStringSubmatch(string(body))
	if matches == nil {
		matches = legacyVersionPattern.FindStringSubmatch(string(body))
	}
	if matches == nil {
		return nil, fmt.Errorf("unable to find a version on %s", download)
	}

	versionString := fmt.Sprintf("%s.%s.%s", matches[1], matches[2], matches[3])
	version, err := semver.NewVersion(versionString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse version %s\n%w", versionString, err)
	}

	return versionology.VersionFetcherArray{versionology.NewSimpleVersionFetcher(version)}, nil
}

func generateMetadata(versionFetcher versionology.VersionFetcher, platform retrieve.Platform) ([]versionology.Dependency, error) {
	versionString := versionFetcher.Version().String()

	arch, ok := assetArch[platform.Arch]
	if !ok {
		fmt.Printf("Skipping %s: unsupported architecture %s\n", versionString, platform.Arch)
		return nil, nil
	}

	uri := fmt.Sprintf("%s/yjp/%[2]s/YourKit-Java-Profiler-%[2]s-%[3]s.zip", downloadBase, versionString, arch)
	checksum, err := upstream.GetSHA256OfRemoteFile(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to checksum %s\n%w", uri, err)
	}

	dependency := cargo.ConfigMetadataDependency{
		Arch:     platform.Arch,
		Checksum: fmt.Sprintf("sha256:%s", checksum),
		CPE:      fmt.Sprintf("cpe:2.3:a:yourkit:java_profiler:%s:*:*:*:*:*:*:*", versionString),
		ID:       id,
		Licenses: []interface{}{
			map[string]string{
				"uri": "https://www.yourkit.com/purchase/eula.pdf",
			},
		},
		Name:    name,
		OS:      platform.OS,
		PURL:    fmt.Sprintf("%s@%s?arch=%s", purl, versionString, platform.Arch),
		SHA256:  checksum,
		Stacks:  []string{"*"},
		URI:     uri,
		Version: versionString,
	}

	d, err := versionology.NewDependency(dependency, fmt.Sprintf("%s-%s", platform.OS, platform.Arch))
	if err != nil {
		return nil, err
	}

	return []versionology.Dependency{d}, nil
}
