package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"net/http"
)

type Message struct {
	Uuid    string `json:"uuid"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Likes   uint32    `json:"likes"`
}

func CreateNewMessage(data Message, w *http.ResponseWriter) {
	testIdx := getIdxFromMapsTable(data.Uuid)
	if testIdx != -1 {
		(*w).WriteHeader(http.StatusConflict)
		return
	}

	Idx := getMessageLastIdx()
	var wg sync.WaitGroup
	wg.Add(2)
	go insertNewUUIDToMapsTable(data.Uuid, Idx+1, &wg)
	go insertNewMessageToDatasTable(data, &wg)
	wg.Wait()
	(*w).WriteHeader(http.StatusCreated)
}

func insertNewMessageToDatasTable(data Message, wg *sync.WaitGroup) {
	defer wg.Done()	
	db := connectSqlDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO datas_table (uuid,message,author,likes) VALUES (?,?,?,?) ")
	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}
	res, err2 := stmt.Exec(data.Uuid, data.Message, data.Author, data.Likes)
	if err2 != nil {
		panic(err2.Error())
	}
	fmt.Println(res.RowsAffected())
}

func insertNewUUIDToMapsTable(uuid string, idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	db := connectSqlDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO maps_table (uuid,idx) VALUES (?,?) ")
	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}
	res, err2 := stmt.Exec(uuid, idx)
	if err2 != nil {
		panic(err2.Error())
	}
	fmt.Println(res.RowsAffected())
}

func getMessageLastIdx() int {
	var Idx int
	db := connectSqlDB()
	defer db.Close()
	err1 := db.QueryRow("select idx from datas_table ORDER BY idx DESC LIMIT 1").Scan(&Idx)
	if err1 != nil {
		fmt.Println(err1)
	}
	return Idx
}
