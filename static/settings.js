var 
	mapWidth		= 620;
	mapHeight		= 360;

	shadowWidth		= 2;
	shadowOpacity		= 0.3;
	shadowColor		= "black";
	shadowX			= 1;
	shadowY			= 2;

	iPhoneLink		= true,

	isNewWindow		= false,

	borderColor		= "#ffffff",
	borderColorOver		= "#ffffff",

	nameColor		= "#ffffff",
	nameFontSize		= "11px",
	nameFontWeight		= "bold",
	nameStroke = false,

	overDelay		= 300,

	map_data = {
    st1: {
   	 	id: 1,
		name: "Адыгея",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st2: {
   	 	id: 2,
		name: "Алтай",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st3: {
   	 	id: 3,
		name: "Башкортостан",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st4: {
   	 	id: 4,
		name: "Бурятия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st5: {
   	 	id: 5,
		name: "Дагестан",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st6: {
   	 	id: 6,
		name: "Северная Осетия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st7: {
   	 	id: 7,
		name: "Кабардино-Балкария",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st8: {
   	 	id: 8,
		name: "Калмыкия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st9: {
   	 	id: 9,
		name: "Карачаево-Черкесия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st10: {
   	 	id: 10,
		name: "Карелия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st11: {
   	 	id: 11,
		name: "Коми",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st12: {
   	 	id: 12,
		name: "Марий Эл",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st13: {
   	 	id: 13,
		name: "Мордовия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st14: {
   	 	id: 14,
		name: "Саха (Якутия)",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st15: {
   	 	id: 15,
		name: "Ингушетия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st16: {
   	 	id: 16,
		name: "Татарстан",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st17: {
   	 	id: 17,
		name: "Тыва",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st18: {
   	 	id: 18,
		name: "Удмуртия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st19: {
   	 	id: 19,
		name: "Хакасия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st20: {
   	 	id: 20,
		name: "Чечня",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st21: {
   	 	id: 21,
		name: "Чувашия",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st22: {
   	 	id: 22,
		name: "Алтайский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st23: {
   	 	id: 23,
		name: "Забайкальский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st24: {
   	 	id: 24,
		name: "Камчатский край",
		shortname: "",
		link: "",
		comment: "City name city name<br>City name city name",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st25: {
   	 	id: 25,
		name: "Краснодарский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st26: {
   	 	id: 26,
		name: "Красноярский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st27: {
   	 	id: 27,
		name: "Пермский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st28: {
   	 	id: 28,
		name: "Приморский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st29: {
   	 	id: 29,
		name: "Ставропольский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st30: {
   	 	id: 30,
		name: "Хабаровский край",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st31: {
   	 	id: 31,
		name: "Амурская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st32: {
   	 	id: 32,
		name: "Архангельская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st33: {
   	 	id: 33,
		name: "Астраханская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st34: {
   	 	id: 34,
		name: "Белгородская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st35: {
   	 	id: 35,
		name: "Брянская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st36: {
   	 	id: 36,
		name: "Владимирская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st37: {
   	 	id: 37,
		name: "Волгоградская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st38: {
   	 	id: 38,
		name: "Вологодская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st39: {
   	 	id: 39,
		name: "Воронежская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st40: {
   	 	id: 40,
		name: "Ивановская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st41: {
   	 	id: 41,
		name: "Иркутская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st42: {
   	 	id: 42,
		name: "Калининградская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st43: {
   	 	id: 43,
		name: "Калужская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st44: {
   	 	id: 44,
		name: "Кемеровская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st45: {
   	 	id: 45,
		name: "Кировская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st46: {
   	 	id: 46,
		name: "Костромская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st47: {
   	 	id: 47,
		name: "Курганская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st48: {
   	 	id: 48,
		name: "Курская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st49: {
   	 	id: 49,
		name: "Ленинградская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st50: {
   	 	id: 50,
		name: "Липецкая область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st51: {
   	 	id: 51,
		name: "Магаданская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st52: {
   	 	id: 52,
		name: "Московская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st53: {
   	 	id: 53,
		name: "Мурманская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st54: {
   	 	id: 54,
		name: "Нижегородская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st55: {
   	 	id: 55,
		name: "Новгородская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st56: {
   	 	id: 56,
		name: "Новосибирская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st57: {
   	 	id: 57,
		name: "Омская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st58: {
   	 	id: 58,
		name: "Оренбургская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st59: {
   	 	id: 59,
		name: "Орловская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st60: {
   	 	id: 60,
		name: "Пензенская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st61: {
   	 	id: 61,
		name: "Псковская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st62: {
   	 	id: 62,
		name: "Ростовская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st63: {
   	 	id: 63,
		name: "Рязанская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st64: {
   	 	id: 64,
		name: "Самарская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st65: {
   	 	id: 65,
		name: "Саратовская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st66: {
   	 	id: 66,
		name: "Сахалинская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st67: {
   	 	id: 67,
		name: "Свердловская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st68: {
   	 	id: 68,
		name: "Смоленская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st69: {
   	 	id: 69,
		name: "Тамбовская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st70: {
   	 	id: 70,
		name: "Тверская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st71: {
   	 	id: 71,
		name: "Томская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st72: {
   	 	id: 72,
		name: "Тульская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st73: {
   	 	id: 73,
		name: "Тюменская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st74: {
   	 	id: 74,
		name: "Ульяновская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st75: {
   	 	id: 75,
		name: "Челябинская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st76: {
   	 	id: 76,
		name: "Ярославская область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st77: {
   	 	id: 77,
		name: "Москва",
		shortname: "МСК",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st78: {
   	 	id: 78,
		name: "Санкт-Петербург",
		shortname: "СПБ",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st79: {
   	 	id: 79,
		name: "Еврейская автономная область",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st80: {
   	 	id: 80,
		name: "Ненецкий",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st81: {
   	 	id: 81,
		name: "Ханты-Мансийский",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st82: {
   	 	id: 82,
		name: "Чукотский",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st83: {
   	 	id: 83,
		name: "Ямало-Ненецкий",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st84: {
   	 	id: 84,
		name: "Республика Крым",
		shortname: "",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	},
    st85: {
   	 	id: 85,
		name: "Севастополь",
		shortname: "СЕВ",
		link: "",
		comment: "",
		image: "",
		color_map: "#A76DC8", 
		color_map_over: "#FF9840"
	}
};