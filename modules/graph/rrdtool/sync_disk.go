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

package rrdtool

import (
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
	"log"
	"time"
)

const (
	_ = iota
	IO_TASK_M_READ
	IO_TASK_M_WRITE
	IO_TASK_M_FLUSH
	IO_TASK_M_FETCH
)

type io_task_t struct {
	method int
	args   interface{}
	done   chan error
}

var (
	Out_done_chan chan int
	io_task_chan  chan *io_task_t
)

func init() {
	Out_done_chan = make(chan int, 1)
	// chan len up to 1600,allow 1600 metric sync disk the same time.
	io_task_chan = make(chan *io_task_t, 1600)
}

func syncDisk() {
	time.Sleep(time.Second * g.CACHE_DELAY)
	ticker := time.NewTicker(time.Millisecond * g.FLUSH_DISK_STEP)
	defer ticker.Stop()
	var idx int = 0

	for {
		select {
		case <-ticker.C:
			idx = idx % store.GraphItems.Size
			FlushRRD(idx, false)
			idx += 1
		case <-Out_done_chan:
			log.Println("cron recv sigout and exit...")
			return
		}
	}
}

func ioWorker() {
	var err error
	for {
		select {
		case task := <-io_task_chan:
			if task.method == IO_TASK_M_FLUSH {
				if args, ok := task.args.(*flushfile_t); ok {
					task.done <- flushrrd(args.filename, args.items)
				}
			} else if task.method == IO_TASK_M_FETCH {
				if args, ok := task.args.(*fetch_t); ok {
					args.data, err = fetch(args.filename, args.cf, args.start, args.end, args.step)
					task.done <- err
				}
			}
		}
	}
}
