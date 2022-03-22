package service

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type MessageWithId struct {
	Idx     int    `json:"idx"`
	Uuid    string `gorm:"size:36" json:"uuid"`
	Message string `gorm:"size:1024" json:"message"`
	Author  string `gorm:"size:64" json:"author"`
	Likes   uint32 `json:"likes"`
}

type SyncResponse struct {
	NewMessages	   map[int][]MessageWithId `json:"newMessages"`
	// NewMessages    []MessageWithId `json:"newMessages"`
	UpdateMessages []MessageWithId `json:"updateMessages"`
	DeleteListIdx  []int           `json:"deleteListIdx"`
	LastIdx        int             `json:"lastIdx"`
	LastQidx       int             `json:"lastQidx"`
}

func SyncMessages(idx int, qidx int) SyncResponse {
	// wait group
	var wg sync.WaitGroup
	// pointer variables
	// var newMessages, updateMessages []MessageWithId
	var updateMessages []MessageWithId
	var lastIdx, lastQidx int
	var deleteListIdx []int
	newMessages := make(map[int][]MessageWithId)
	// run in parallel
	wg.Add(4)
	go GetNewMessages(idx, newMessages, &lastIdx, &wg)
	go GetUpdateMessages(qidx, &updateMessages, &wg)
	go GetDeleteListIdx(qidx, &deleteListIdx, &wg)
	go GetLastQidx(&lastQidx, &wg)
	// return
	wg.Wait()
	return SyncResponse{newMessages, updateMessages, deleteListIdx, lastIdx, lastQidx}
}

// GetNewMessages config
const itemPerQuery = 5000
const maxQuery = 20

// Get new messages from datas_table
func GetNewMessages(idx int, newMessages map[int][]MessageWithId, lastIdx *int, wg *sync.WaitGroup) {
	defer wg.Done()
	// get last idx
	GetLastIdx(lastIdx)
	// create queries for concurrent query
	var queries []string
	cidx := idx
	for cidx+itemPerQuery < *lastIdx {
		queries = append(queries, "SELECT idx,uuid,message,likes,author FROM datas_table WHERE idx > "+strconv.Itoa(cidx)+" && idx <= "+strconv.Itoa(cidx+itemPerQuery))
		cidx = cidx + itemPerQuery
	}
	queries = append(queries, "SELECT idx,uuid,message,likes,author FROM datas_table WHERE idx > "+strconv.Itoa(cidx))
	// select all new messages
	var wg2 sync.WaitGroup
	bufferChan := make(chan int, maxQuery)
	wg2.Add(len(queries))
	for coroutineIdx, query := range queries {
		// add buffer chan to limit number of thread; will stall until have space
		bufferChan <- 1
		go GetNewMessagesFragment(query, coroutineIdx, newMessages, &wg2, bufferChan)
		time.Sleep(10 * time.Millisecond)

	}
	close(bufferChan)
	wg2.Wait()
	// concat all new messages
	// for coroutineIdx, _ := range queries {
	// 	*newMessages = append(*newMessages, resMap[coroutineIdx]...)
	// }
}

func GetNewMessagesFragment(query string, coroutineIdx int, resMap map[int][]MessageWithId, wg2 *sync.WaitGroup, bufferChan chan int) {
	fmt.Println("routine index", coroutineIdx)
	defer wg2.Done()
	db := connectSqlDB()
	defer db.Close()
	res, err1 := db.Query(query)
	if err1 != nil {
		panic(err1.Error())
	}
	var fragmentNewMessages []MessageWithId
	// convert all response to MessageWithId struct and store inside slice
	for res.Next() {
		var m MessageWithId
		err3 := res.Scan(&m.Idx, &m.Uuid, &m.Message, &m.Likes, &m.Author)
		if err3 != nil {
			panic(err3.Error())
		}
		fragmentNewMessages = append(fragmentNewMessages, m)
	}
	resMap[coroutineIdx] = fragmentNewMessages
	// release buffer chan
	<-bufferChan
}

// Get updated message from data_tables by look at updates_table
func GetUpdateMessages(qidx int, updateMessages *[]MessageWithId, wg *sync.WaitGroup) {
	defer wg.Done()
	// select all update message
	db := connectSqlDB()
	defer db.Close()
	stmt, err := db.Prepare("SELECT datas_table.idx,datas_table.uuid,datas_table.message,datas_table.author,datas_table.likes " +
		"FROM datas_table " +
		"INNER JOIN updates_table " +
		"ON datas_table.idx = updates_table.idx " +
		"WHERE updates_table.qidx > ? AND updates_table.deleteMethod = false")
	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}
	res, err2 := stmt.Query(qidx)
	if err2 != nil {
		panic(err2.Error())
	}
	// convert all response to MessageWithId struct and store inside slice
	for res.Next() {
		var m MessageWithId
		err3 := res.Scan(&m.Idx, &m.Uuid, &m.Message, &m.Author, &m.Likes)
		if err3 != nil {
			panic(err3.Error())
		}
		*updateMessages = append(*updateMessages, m)
	}
}

// Get list of delete idx from updates_table
func GetDeleteListIdx(qidx int, deleteListIdx *[]int, wg *sync.WaitGroup) {
	defer wg.Done()
	// select all idx that mark as delete
	db := connectSqlDB()
	defer db.Close()
	stmt, err := db.Prepare("SELECT idx FROM updates_table WHERE qidx > ? AND deleteMethod = true")
	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}
	res, err2 := stmt.Query(qidx)
	if err2 != nil {
		panic(err2.Error())
	}
	// convert all response to slice
	for res.Next() {
		var m int
		err3 := res.Scan(&m)
		if err3 != nil {
			panic(err3.Error())
		}
		*deleteListIdx = append(*deleteListIdx, m)
	}
}

// Get latest idx from datas_table
func GetLastIdx(lastIdx *int) {
	// select lastest qidx
	db := connectSqlDB()
	defer db.Close()
	res := db.QueryRow("SELECT idx FROM datas_table ORDER BY idx DESC LIMIT 1")
	*lastIdx = -1
	if res != nil {
		err1 := res.Scan(lastIdx)
		if err1 != nil {
			fmt.Println(err1)
			// panic(err1.Error())
		}
	}
}

// Get lastest qidx from updates_table
func GetLastQidx(lastQidx *int, wg *sync.WaitGroup) {
	defer wg.Done()
	// select lastest qidx
	db := connectSqlDB()
	defer db.Close()
	res := db.QueryRow("SELECT qidx FROM updates_table ORDER BY qidx DESC LIMIT 1")
	*lastQidx = -1
	if res != nil {
		err1 := res.Scan(lastQidx)
		if err1 != nil {
			fmt.Println(err1)
			// panic(err1.Error())
		}
	}
}
