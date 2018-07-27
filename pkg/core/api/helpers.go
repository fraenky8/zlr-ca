package api

import (
	"log"
	"strconv"
	"strings"
)

func ConvertIdsParam(sids string) (ids []int64) {

	sids = strings.TrimSpace(sids)
	if sids == "" {
		return []int64{}
	}

	for _, id := range strings.Split(sids, ",") {
		tid := strings.TrimSpace(id)

		// ignore empty ids
		if tid == "" {
			continue
		}

		id, err := strconv.ParseInt(tid, 10, 64)
		if err != nil {
			// ignore faulty ids - but give log message
			log.Printf("faulty id: %v", err)
			continue
		}

		ids = append(ids, id)
	}

	return ids
}
