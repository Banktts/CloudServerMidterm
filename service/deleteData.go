package service

import (
	"net/http"
	"sync"
)

func DeleteMessage(uuid string, w *http.ResponseWriter) {
	
	Idx := getIdxFromMapsTable(uuid)
	if Idx == -1 {
		(*w).WriteHeader(http.StatusNotFound)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go insertMethodToUpdatesTable(Idx, true, &wg)
	wg.Wait()
	(*w).WriteHeader(http.StatusNoContent)
}
