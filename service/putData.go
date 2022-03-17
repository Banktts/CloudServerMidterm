package service

import "fmt"

func UpdateMessage(data Message) {
	Idx := getIdxFromMapsTable(data.Uuid)
	go updateMessageDatasTable(data, Idx)
	go insertMethodToUpdatesTable(Idx, false)
}

func getIdxFromMapsTable(uuid string) int {
	db := connectSqlDB()
	var idx int
	fmt.Println("uuid :", uuid)
	err := db.QueryRow("select idx from maps_table where uuid=?", uuid).Scan(&idx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(idx)
	return idx
}

func updateMessageDatasTable(data Message, idx int) {
	db := connectSqlDB()
	res, err := db.Exec("update datas_table set uuid= ?,message= ?, author= ?, likes= ? where idx=?", data.Uuid, data.Message, data.Author, data.Likes, idx)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(res.RowsAffected())
}

func insertMethodToUpdatesTable(idx int, method bool) {
	db := connectSqlDB()
	res, err := db.Exec("INSERT INTO updates_table (idx,deleteMethod) VALUES (?,?) ", idx, method)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(res.RowsAffected())
}
