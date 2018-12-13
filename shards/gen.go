package shards

import (
	"fmt"
	"strings"
)

// GetByRegnum returns shard id from regnum record
func GetByRegnum(r string) (shId string) {
	/* select substring(regnum,1,1) as sh_id,substring(regnum,2,10) as inn, customer,substring(regnum,12,2) as Year, substring(regnum,14) from contracts_12_2015; */
	shId = r[:1]
	return
}

// GetByPnum returns shard id from purchasenumber record
func GetByPnum(r string) (shId string) {
	return fmt.Sprintf("%019s", r)[1:2]
}

var shRange = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

// GetByRangeByPattern inserts shIds from shRange into defined pattern
func GenRangeByPattern(p string) (shSlice []string) {
	for _, shId := range shRange {
		shSlice = append(shSlice, strings.Replace(p, "%s", shId, -1))
	}
	return
}
