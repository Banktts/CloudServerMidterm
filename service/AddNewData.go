package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Uuid    string `json:"uuid"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Likes   int    `json:"likes"`
}

func CreateNewMessage(data Message) {

	Idx := getMessageLastIdx()
	go insertNewUUIDToMapsTable(data.Uuid, Idx+1)
	go insertNewMessageToDatasTable(data)

}

func insertNewMessageToDatasTable(data Message) {
	db := connectSqlDB()
	stmt, err := db.Prepare("INSERT INTO datas_table (uuid,message,author,likes) VALUES (?,?,?,?) ")
	if err != nil {
		panic(err.Error())
	}
	res, err2 := stmt.Exec(data.Uuid, data.Message, data.Author, data.Likes)
	if err2 != nil {
		panic(err2.Error())
	}
	fmt.Println(res.RowsAffected())
}

func insertNewUUIDToMapsTable(uuid string, idx int) {
	db := connectSqlDB()
	stmt, err := db.Prepare("INSERT INTO maps_table (uuid,idx) VALUES (?,?) ")
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
	err1 := db.QueryRow("select idx from datas_table ORDER BY idx DESC LIMIT 1").Scan(&Idx)
	if err1 != nil {
		fmt.Println(err1)
	}
	return Idx
}
