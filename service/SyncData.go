package service

import (
	"sync"
)

type MessageWithId struct {
	Idx		int	   `json:"idx"`
	Uuid    string `json:"uuid"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Likes   int    `json:"likes"`
}

type SyncResponse struct {
	NewMessages		[]MessageWithId `json:"newMessages"`
	UpdateMessages 	[]MessageWithId `json:"updateMessages"`
	DeleteListIdx 	[]int 			`json:"deleteListIdx"`
	LastIdx 		int 			`json:"lastIdx"`
	LastQidx 		int 			`json:"lastQidx"`
}

func SyncMessages(idx int, qidx int) SyncResponse {
	// wait group
	var wg sync.WaitGroup
	// pointer variables
	var newMessages, updateMessages []MessageWithId
	var lastIdx, lastQidx int
	var deleteListIdx []int
	// run in parallel
	wg.Add(4)
	go GetNewMessages(idx, &newMessages, &lastIdx, &wg)
	go GetUpdateMessages(qidx, &updateMessages, &wg)
	go GetDeleteListIdx(qidx, &deleteListIdx, &wg)
	go GetLastQidx(&lastQidx, &wg)
	// return
	wg.Wait()
	return SyncResponse{newMessages, updateMessages, deleteListIdx, lastIdx, lastQidx}
}

// Get new messages from datas_table
func GetNewMessages(idx int, newMessages *[]MessageWithId, lastIdx *int, wg *sync.WaitGroup) {
	defer wg.Done()
	// select all new messages
	db := connectSqlDB()
	defer db.Close()
	stmt, err := db.Prepare("SELECT * FROM datas_table WHERE idx > ? ")
	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}
	res, err2 := stmt.Query(idx)
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
		*newMessages = append(*newMessages, m)
	}
	if len(*newMessages) > 0 {
		*lastIdx = (*newMessages)[len(*newMessages)-1].Idx
	} else {
		*lastIdx = idx
	}
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

// Get lastest qidx from updates_table
func GetLastQidx(lastQidx *int, wg *sync.WaitGroup) {
	defer wg.Done()
	// select lastest qidx
	db := connectSqlDB()
	defer db.Close()
	err1 := db.QueryRow("SELECT qidx FROM updates_table ORDER BY qidx DESC LIMIT 1").Scan(lastQidx)
	if err1 != nil {
		panic(err1.Error())
	}
}