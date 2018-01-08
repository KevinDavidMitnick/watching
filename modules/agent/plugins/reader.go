// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugins

import (
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"os"
	"path/filepath"
	"strings"
)

func ListPlugins(relativePath string, params g.PluginParam) map[string]*Plugin {
	ret := make(map[string]*Plugin)
	if relativePath == "" {
		return ret
	}

	executeFile := filepath.Join(g.Config().Plugin.Dir, relativePath, params.ExecuteScript)
	executeFile = strings.TrimSpace(executeFile)

	if !file.IsExist(executeFile) || !file.IsFile(executeFile) {
		return ret
	}

	f, err := os.Open(executeFile)
	if err != nil {
		return ret
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return ret
	}

	cycle := params.ExecuteInterval
	param := strings.TrimSpace(params.ExecuteParam)
	plugin := &Plugin{FilePath: executeFile, MTime: fi.ModTime().Unix(), Cycle: cycle, Param: param}
	ret[executeFile] = plugin
	return ret
}
