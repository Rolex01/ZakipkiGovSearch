package cbr

import (
	cbrrate "bitbucket.org/company-one/cbrrate"
	"github.com/jasonlvhit/gocron"
	"time"
	"fmt"
)

func InitScheduler(rate *[]cbrrate.CurrencyRate) {
	UpdateCurrencyRate(rate)
	gocron.Every(1).Day().At("05:00").Do(UpdateCurrencyRate, rate)
	gocron.Start()
}

func UpdateCurrencyRate(rate *[]cbrrate.CurrencyRate) {
	currYear := time.Now().Year()

	*rate = cbrrate.LoadCBRRate("01.01.2014", fmt.Sprintf("31.12.%d",currYear))
}
