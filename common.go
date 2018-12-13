package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"
)

func StrToDate(d string) time.Time {
	s := strings.Split(d, ".")
	loc, _ := time.LoadLocation("Europe/Moscow")
	day, _ := strconv.ParseInt(s[0], 10, 32)
	year, _ := strconv.ParseInt(s[2], 10, 32)
	month, _ := strconv.ParseInt(s[1], 10, 32)
	return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, loc)
}

func filterDateTimestamp(s string) uint64 {

	dt := strings.Split(s, ".")
	if len(dt) < 3 {
		return 0
	}
	loc, _ := time.LoadLocation("Europe/Moscow")
	d, _ := strconv.Atoi(dt[0])
	m, _ := strconv.Atoi(dt[1])
	y, _ := strconv.Atoi(dt[2])
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, loc)
	//t = t.Add(time.Hour * -24)
	// log.Println("s | ", t)
	return uint64(t.Unix())
}

func filterToDateTimestamp(s string) uint64 {

	dt := strings.Split(s, ".")
	if len(dt) < 3 {
		return 0
	}
	loc, _ := time.LoadLocation("Europe/Moscow")
	d, _ := strconv.Atoi(dt[0])
	m, _ := strconv.Atoi(dt[1])
	y, _ := strconv.Atoi(dt[2])
	t0 := time.Date(y, time.Month(m), d, 0, 0, 0, 0, loc)
	t := t0.Add(time.Hour * 24)
	// log.Println("s | ", t)
	return uint64(t.Unix())
}

func getOrgsInfo(regnums []uint64, db *sql.DB) (inns, kpps []uint64) {
	for _, v := range regnums {
		i, k := getOrgInfoByRegnum(v, db)
		inns = append(inns, i)
		kpps = append(kpps, k)
	}
	return
}

func getOrgInfoByRegnum(regnum uint64, db *sql.DB) (inn, kpp uint64) {
	err := db.QueryRow("select inn::bigint,kpp::bigint from nsiorg where regnum::bigint=$1", regnum).Scan(&inn, &kpp)
	if err != nil {
		log.Println("Can't resolve inn & kpp by regnum:", regnum)
	}
	return
}

func hackMKDC(arr []uint64) []uint64 {
	hacked := make([]uint64, len(arr))
	for i, v := range arr {
		if v == 7112000059 {
			v = 2817032763
		}
		hacked[i] = v
	}
	return hacked
}
