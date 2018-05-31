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

package g

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"unsafe"

	"github.com/toolkits/file"
)

type GlobalConfig struct {
	RrdPath       string `json:"rrdPath"`
	DataPath      string `json:"dataPath"`
	HttpAddr      string `json:"httpAddr"`
	RaftAddr	  string `json:"raftAddr"`
	JoinAddr      string `json:"joinAddr"`
	LogPath       string `json:"logPath"`
}

var (
	ConfigFile string
	ptr        unsafe.Pointer
)

func Config() *GlobalConfig {
	return (*GlobalConfig)(atomic.LoadPointer(&ptr))
}

func ParseConfig(cfg string) {
	logFile := file.MustOpenLogFile("logs/opsultra-rrdlite.log")
	logger := log.New(logFile, "[config] ", log.LstdFlags|log.Lshortfile)
	if cfg == "" {
		logger.Fatalln("config file not specified: use -c $filename")
	}

	if !file.IsExist(cfg) {
		logger.Fatalln("config file specified not found:", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		logger.Fatalln("read config file", cfg, "error:", err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		logger.Fatalln("parse config file", cfg, "error:", err.Error())
	}

	// set config
	atomic.StorePointer(&ptr, unsafe.Pointer(&c))

	logger.Println("g.ParseConfig ok, file", cfg)
}
