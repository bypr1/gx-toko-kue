package core

import (
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func SetDateRange(parameters url.Values) (time.Time, time.Time) {
	var fromDate, toDate time.Time

	now := time.Now()

	if fromDateReq := parameters.Get("fromDate"); len(fromDateReq) > 0 {
		fromDate, _ = time.Parse("02/01/2006 15:04:05", fromDateReq+" 00:00:00")
	} else {
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)
	}

	if toDateReq := parameters.Get("toDate"); len(toDateReq) > 0 {
		toDate, _ = time.Parse("02/01/2006", toDateReq)
	} else {
		toDate = now
	}

	return fromDate, toDate
}

func StrPadLeft(original string, padLength int, padChar rune) string {
	if len(original) >= padLength {
		return original
	}
	padding := strings.Repeat(string(padChar), padLength-len(original))
	return padding + original
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
