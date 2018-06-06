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
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
)

// TODO 草草的写了一个db连接池,优化下
var (
	dbLock    sync.RWMutex
	dbConnMap map[string]*sql.DB
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = makeDbConn()
	if DB == nil || err != nil {
		log.Panicln("g.InitDB failed., get db conn fail:", err)
	} else {
		dbConnMap = make(map[string]*sql.DB)
		log.Println("g.InitDB ok")
	}
}

func GetDbConn(connName string) (c *sql.DB, e error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	log.Println("start to get db conn:", connName)
	var err error
	var dbConn *sql.DB
	dbConn = dbConnMap[connName]
	if dbConn != nil {
		err = dbConn.Ping()
		if err != nil {
			closeDbConn(dbConn)
			delete(dbConnMap, connName)
		}
	}
	if dbConn == nil || err != nil {
		dbConn, err = makeDbConn()
		if dbConn == nil || err != nil {
			closeDbConn(dbConn)
			return nil, err
		}
		dbConnMap[connName] = dbConn
	}
	err = dbConn.Ping()
	if err != nil {
		log.Println("get db conn is not nil,ping failed.")
		closeDbConn(dbConn)
		delete(dbConnMap, connName)
		return nil, err
	}
	log.Println("success to get db conn.")

	return dbConn, err
}

// 创建一个新的mysql连接
func makeDbConn() (conn *sql.DB, err error) {
	times := 1.0
	for {
		log.Printf("in make db connting,times:%f", times)
		conn, err = sql.Open("mysql", Config().DB.Dsn)
		if conn == nil || err != nil {
			log.Printf("get db conn fail: %s,try again,times:%f", err.Error(), times)
			if times >= 60 {
				log.Println("too may times try get db conn,give up.....")
				return conn, err
			} else {
				times += 1
			}
			time.Sleep(time.Duration(times) * time.Second)
		} else {
			conn.SetMaxIdleConns(Config().DB.MaxIdle)
			err = conn.Ping()
			return conn, err
		}
	}
}

func closeDbConn(conn *sql.DB) {
	if conn != nil {
		conn.Close()
	}
}
