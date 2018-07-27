package api

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ConvertIdsParam(sids string) (ids []int64, err error) {

	sids = strings.TrimSpace(sids)
	if sids == "" {
		return []int64{}, fmt.Errorf("no id(s) provided")
	}

	for _, id := range strings.Split(sids, ",") {
		tid := strings.TrimSpace(id)

		// ignore empty ids
		if tid == "" {
			continue
		}

		id, parseErr := strconv.ParseInt(tid, 10, 64)
		if parseErr != nil {
			log.Printf("faulty id: %v", parseErr)
			err = parseErr
			continue
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []int64{}, fmt.Errorf("no valid id(s) provided")
	}

	if err != nil {
		return []int64{}, fmt.Errorf("at least one invalid id detected: %v", err)
	}

	return ids, nil
}
