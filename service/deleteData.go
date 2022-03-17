package service

func DeleteMessage(uuid string) {
	
	Idx := getIdxFromMapsTable(uuid)

	go insertMethodToUpdatesTable(Idx, true)
}
