package main

//import "log"
import (
	"fmt"
	"strings"
	"html/template"
)

// Region describes Region item
type region struct {
	title      string
	value      string
	isSelected bool
}
type district struct {
	title   string
	regions []region
}

func (d *district) generateHTML(selected []string) (html string) {
	html = fmt.Sprintf("<optgroup label=\"%s\">", d.title)
	//log.Println("Regions genHTML", d, selected)
	for _, region := range d.regions {
		for _, selectedItem := range selected {
			//log.Println("Looking for ", selectedItem, " in ", region.value, ":", strings.Compare(region.value, selectedItem))
			if region.isSelected = strings.Compare(region.value, selectedItem) == 0; region.isSelected {
				//log.Println("Regions genHTML found!", region, selected)
				html += fmt.Sprintf("<option value=\"%s\" selected>%s</option>", region.value, region.title)
			} else {
				html += fmt.Sprintf("<option value=\"%s\">%s</option>", region.value, region.title)
			}
		}
		if len(selected) == 0 {
			html += fmt.Sprintf("<option value=\"%s\">%s</option>", region.value, region.title)
		}
	}
	html += "</optgroup>"
	return
}
func (p ContractResult) GenFilterRegionContracts() template.HTML {
	//log.Println("GenFilterRegion For contracts", p.FilterRegion)
	return genFilterRegion(p.FilterRegion)
}
func (p NotificationResult) GenFilterRegionNotifications() template.HTML {
	//log.Println("GenFilterRegion For contracts", p.FilterRegion)
	return genFilterRegion(p.FilterRegion)
}
func (p Notification223Result) GenFilterRegionNotifications() template.HTML {
	//log.Println("GenFilterRegion For contracts", p.FilterRegion)
	return genFilterRegion(p.FilterRegion)
}
func genFilterRegion(regions []string) (st template.HTML) {
	//var districts []district
	//s := "<select id=\"region-set\" style='width: 85%;outline: none;border: 1px solid #cccccc;height: 26px;background: #ffffff url(\"http://online.monitoring-crm.ru/admin/template/css/../images/in-bg.png\") center top repeat-x;' name=\"filterRegion\" multiple=\"multiple\" >"
	s := ""//"<select id=\"region-set\" name=\"filterRegion\" multiple=\"multiple\" >"

	tmpD := new(district)
	tmpD.title = "Дальневосточный федеральный округ"
	tmpR := new(region)

	tmpR.value = "14"
	tmpR.title = "Саха/Якутия"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "25"
	tmpR.title = "Приморский край"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "27"
	tmpR.title = "Хабаровский край"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "28"
	tmpR.title = "Амурская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "41"
	tmpR.title = "Камчатский край"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "49"
	tmpR.title = "Магаданская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "65"
	tmpR.title = "Сахалинская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "79"
	tmpR.title = "Еврейская область"
	tmpD.regions = append(tmpD.regions, *tmpR)

	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	tmpD.title = "Приволжский федеральный округ"

	tmpR.value = "2"
	tmpR.title = "Республика Башкортостан"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "12"
	tmpR.title = "Республика Марий Эл"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "13"
	tmpR.title = "Республика Мордовия"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "16"
	tmpR.title = "Республика Татарстан"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "18"
	tmpR.title = "Удмуpтская Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "21"
	tmpR.title = "Чувашия"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "43"
	tmpR.title = "Кировская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "52"
	tmpR.title = "Нижегоpодская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "56"
	tmpR.title = "Оренбургская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "58"
	tmpR.title = "Пензенская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "59"
	tmpR.title = "Пермский край"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "63"
	tmpR.title = "Самарская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "64"
	tmpR.title = "Саратовская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "73"
	tmpR.title = "Ульяновская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	//	districts = append(districts, *tmpD)
	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	tmpD.title = "Северо-Западный федеральный округ"

	tmpR.value = "10"
	tmpR.title = "Карелия Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "11"
	tmpR.title = "Республика Коми"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "29"
	tmpR.title = "Аpхангельская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "35"
	tmpR.title = "Вологодская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "39"
	tmpR.title = "Калининградская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "47"
	tmpR.title = "Ленинградская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "51"
	tmpR.title = "Мурманская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "53"
	tmpR.title = "Новгородская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "60"
	tmpR.title = "Псковская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "78"
	tmpR.title = "г. Санкт-Петербург"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "83"
	tmpR.title = "Ненецкий АО"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	//	districts = append(districts, *tmpD)
	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	tmpD.title = "Северо-Кавказский федеральный округ"

	tmpR.value = "5"
	tmpR.title = "Республика Дагестан"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "6"
	tmpR.title = "Республика Ингушетия"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "7"
	tmpR.title = "Кабаpдино-Балкаpская Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "9"
	tmpR.title = "Карачаево-Черкесская Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "15"
	tmpR.title = "Республика Северная Осетия - Алания"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "20"
	tmpR.title = "Чеченская Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "26"
	tmpR.title = "Ставропольский край"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	//	districts = append(districts, *tmpD)
	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	tmpD.title = "Сибирский федеральный округ"

	tmpR.value = "3"
	tmpR.title = "Бурятия Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "4"
	tmpR.title = "Республика Алтай"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "17"
	tmpR.title = "Тыва Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "19"
	tmpR.title = "Хакасия Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "22"
	tmpR.title = "Алтайский кpай"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "24"
	tmpR.title = "Красноярский край"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "38"
	tmpR.title = "Иркутская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "42"
	tmpR.title = "Кемеровская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "54"
	tmpR.title = "Новосибирская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "55"
	tmpR.title = "Омская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "70"
	tmpR.title = "Томская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "75"
	tmpR.title = "Забайкальский край"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "85"
	tmpR.title = "Иркутская область Усть-Ордынский"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	//	districts = append(districts, *tmpD)
	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	tmpD.title = "Уральский федеральный округ"

	tmpR.value = "45"
	tmpR.title = "Курганская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "66"
	tmpR.title = "Свеpдловская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "72"
	tmpR.title = "Тюменская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "74"
	tmpR.title = "Челябинская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "86"
	tmpR.title = "Ханты-Мансийский автономный округ - Югра"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "89"
	tmpR.title = "Ямало-Ненецкий АО"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	//	districts = append(districts, *tmpD)
	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	tmpD.title = "Центральный федеральный округ"

	tmpR.value = "31"
	tmpR.title = "Белгородская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "32"
	tmpR.title = "Брянская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "33"
	tmpR.title = "Владимиpская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "36"
	tmpR.title = "Воронежская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "37"
	tmpR.title = "Ивановская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "40"
	tmpR.title = "Калужская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "44"
	tmpR.title = "Костpомская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "46"
	tmpR.title = "Курская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "48"
	tmpR.title = "Липецкая область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "50"
	tmpR.title = "Московская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "57"
	tmpR.title = "Орловская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "62"
	tmpR.title = "Рязанская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "67"
	tmpR.title = "Смоленская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "68"
	tmpR.title = "Тамбовская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "69"
	tmpR.title = "Тверская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "71"
	tmpR.title = "Тульская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "76"
	tmpR.title = "Ярославская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "77"
	tmpR.title = "г. Москва"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	//	districts = append(districts, *tmpD)
	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	tmpD.title = "Южный федеральный округ"

	tmpR.value = "1"
	tmpR.title = "Адыгея Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "8"
	tmpR.title = "Калмыкия Республика"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "23"
	tmpR.title = "Кpаснодаpский кpай"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "30"
	tmpR.title = "Астраханская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "34"
	tmpR.title = "Волгогpадская область"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "61"
	tmpR.title = "Ростовская область"
	tmpD.regions = append(tmpD.regions, *tmpR)

	tmpR.value = "91"
	tmpR.title = "Республика Крым"
	tmpD.regions = append(tmpD.regions, *tmpR)
	tmpR = new(region)

	tmpR.value = "92"
	tmpR.title = "г. Севастополь"
	tmpD.regions = append(tmpD.regions, *tmpR)

	s += tmpD.generateHTML(regions)
	tmpD.regions = nil

	//s += "</select>"
	st=template.HTML(s)
	return
}

/*<select id="region-set" style='width: 85%;outline: none;border: 1px solid #cccccc;height: 26px;background: #ffffff url("http://online.monitoring-crm.ru/admin/template/css/../images/in-bg.png") center top repeat-x;' name="filterRegion" multiple="multiple" >





</select>
*/
