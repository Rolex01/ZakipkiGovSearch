package documentFields

var NotificationFields []Field = []Field{
	Field{"FixPurchaseNumber", "Номер извещения о проведении торгов", 0},
	Field{"ObjectInfo", "Наименование аукциона", 0},
	Field{"PublishedXlsDate", "Дата размещения", 0},
	Field{"PriceExcel", "НМЦК, руб.", 0},
	Field{"Org", "Организация, осуществляющая закупку", 0},
	Field{"OrgINN", "ИНН организации, осуществляющей заказ", 0},
	Field{"Customer", "Заказчик", 0},
	Field{"CustomerINN", "ИНН заказчика", 0},
	Field{"GetRegnum", "Номер реестровой записи", 1},
	Field{"FormatCtPrice", "Сумма контракта, руб.", 1},
	//Field{"CtPName", "Наименование контракта (for Яна with <3)", 90},
}

var Notification223Fields []Field = []Field{
	Field{"FixPurchaseNumber", "Номер извещения о проведении торгов", 0},
	Field{"Name", "Наименование аукциона", 0},
	Field{"Lotname", "Наименование лота", 0},
	Field{"PublishedXlsDate", "Дата размещения", 0},
	Field{"PriceExcel", "Цена лота, руб.", 0},
	Field{"TotalPriceExcel", "НМЦК, руб.", 0},
	Field{"Placer", "Организация, осуществляющая закупку", 0},
	Field{"PlacerINN", "ИНН организации, осуществляющей закупку", 0},
	Field{"Customer", "Заказчик", 0},
	Field{"CustomerINN", "ИНН заказчика", 0},
	Field{"Region", "Регион", 0},
	Field{"MethodName", "Способ определения поставщика", 0},
	Field{"PlaceName", "Площадка", 0},
}
