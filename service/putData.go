package service

import (
	"fmt"
	"net/http"
	"sync"
)

func UpdateMessage(data Message, w *http.ResponseWriter) {
	Idx := getIdxFromMapsTable(data.Uuid)
	if Idx == -1 {
		(*w).WriteHeader(http.StatusNotFound)
		return
	}
	// wait group
	var wg sync.WaitGroup
	wg.Add(2)
	go updateMessageDatasTable(data, Idx, &wg)
	go insertMethodToUpdatesTable(Idx, false, &wg)
	wg.Wait()
	(*w).WriteHeader(http.StatusNoContent)
}

func getIdxFromMapsTable(uuid string) int {
	db := connectSqlDB()
	defer db.Close()
	var idx int
	fmt.Println("uuid :", uuid)
	err := db.QueryRow("select idx from maps_table where uuid=? ", uuid).Scan(&idx)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	fmt.Println(idx)
	return idx
}

func updateMessageDatasTable(data Message, idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	db := connectSqlDB()
	defer db.Close()
	res, err := db.Exec("update datas_table set uuid= ?,message= ?, author= ?, likes= ? where idx=?", data.Uuid, data.Message, data.Author, data.Likes, idx)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res.RowsAffected())
}

func insertMethodToUpdatesTable(idx int, method bool, wg *sync.WaitGroup) {
	defer wg.Done()
	db := connectSqlDB()
	defer db.Close()
	res, err := db.Exec("INSERT INTO updates_table (idx,deleteMethod) VALUES (?,?) ", idx, method)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res.RowsAffected())
}
