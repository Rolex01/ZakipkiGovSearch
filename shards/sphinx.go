package shards

import "strings"

const sphinxCfg string = `
searchd {
	listen = localhost:9312
	listen = localhost:9306:mysql41

	# log file, searchd run info is logged here
	# optional, default is 'searchd.log'
	log = /var/log/sphinxsearch/searchd.log

	# query log file, all search queries are logged here
	# optional, default is empty (do not log queries)
	query_log = /var/log/sphinxsearch/query.log

	# client read timeout, seconds
	# optional, default is 5
	read_timeout = 5

	# request timeout, seconds
	# optional, default is 5 minutes
	client_timeout = 300

	# maximum amount of children to fork (concurrent searches to run)
	# optional, default is 0 (unlimited)
	max_children = 30

	# PID file, searchd process ID file name
	# mandatory
	pid_file = /var/run/sphinxsearch/searchd.pid

	# max amount of matches the daemon ever keeps in RAM, per-index
	# WARNING, THERE'S ALSO PER-QUERY LIMIT, SEE SetLimits() API CALL
	# default is 1000 (just like Google) deprecated
	max_matches = 100000

	# seamless rotate, prevents rotate stalls if precaching huge datasets
	# optional, default is 1
	seamless_rotate = 1

	# whether to forcibly preopen all indexes on startup
	# optional, default is 1 (preopen everything)
	preopen_indexes = 1

	# whether to unlink .old index copies on succesful rotation.
	# optional, default is 1 (do unlink)
	unlink_old = 1

	# attribute updates periodic flush timeout, seconds
	# updates will be automatically dumped to disk this frequently
	# optional, default is 0 (disable periodic flush)
	#
	# attr_flush_period = 900

	# instance-wide ondisk_dict defaults (per-index value take precedence)
	# optional, default is 0 (precache all dictionaries in RAM)
	#
	# ondisk_dict_default = 1

	# MVA updates pool size
	# shared between all instances of searchd, disables attr flushes!
	# optional, default size is 1M
	mva_updates_pool = 1M

	# max allowed network packet size
	# limits both query packets from clients, and responses from agents
	# optional, default size is 8M
	max_packet_size = 12M

	# max allowed per-query filter count
	# optional, default is 256
	max_filters = 1024

	# max allowed per-filter values count
	# optional, default is 4096
	max_filter_values = 10000000
}


indexer {
	# memory limit, in bytes, kiloytes (16384K) or megabytes (256M)
	# optional, default is 32M, max is 2047M, recommended is 256M to 1024M
	mem_limit = 2047M

	# maximum IO calls per second (for I/O throttling)
	# optional, default is 0 (unlimited)
	#
	# max_iops = 40


	# maximum IO call size, bytes (for I/O throttling)
	# optional, default is 0 (unlimited)
	#
	# max_iosize = 1048576

	# maximum xmlpipe2 field length, bytes
	# optional, default is 2M
	#
	# max_xmlpipe2_field = 4M

	# write buffer size, bytes
	# several (currently up to 4) buffers will be allocated
	# write buffers are allocated in addition to mem_limit
	# optional, default is 1M
	#
	# write_buffer = 1M

	# maximum file field adaptive buffer size
	# optional, default is 8M, minimum is 1M
	#
	# max_file_field_buffer = 32M
}

source base_source {
	type     = pgsql
	sql_host = localhost
	sql_port = 5432
	sql_user = postgres
	sql_pass = lk147963
	sql_db   = $db
}

index base_index {
	docinfo     = extern
	dict        = keywords
	mlock       = 0
	morphology  = stem_ru stem_en
	enable_star = 1
}
`

// GenerateSphinxConfig generate sphinx indexes & sources. includeSysCfg - add || not system configuration to result string
func GenerateSphinxConfig(includeSysCfg bool, prefix string) string {

	var cfgData []string

	if includeSysCfg {
		cfgData = append(cfgData, prx(sphinxCfg, prefix))
	}

	for _, i := range singleSphinxIndexes {
		cfgData = append(cfgData, prx(i, prefix))
	}

	for _, i := range rangeSphinxIndexes {
		cfgData = append(cfgData, strings.Join(GenRangeByPattern(prx(i, prefix)), "\n"))
	}

	for _, i := range rangeSphinxIndexesIgnoreFirst {
		cfgData = append(cfgData, strings.Join(GenRangeByPattern(prx(i, prefix))[1:], "\n"))
	}

	return strings.Join(cfgData, "\n")
}

func prx(s, prefix string) string {
	s = strings.Replace(s, "$p", prefix, -1)
	s = strings.Replace(s, "$db", prefix, -1)
	return s
}

const notifications string = `

##
# Notifications %s
##

source $p_s_notifications_%s : base_source {
	sql_query          = \
		SELECT n.pnum::bigint AS id, n.pnum::bigint AS pnum, n.objectInfo AS info, n.org_regnum::bigint AS org, \
			n.org_region AS org_region, \
			EXTRACT(EPOCH FROM n.published AT TIME ZONE 'MSK') AS published, \
			n.maxPrice AS gMaxPrice, \
			n.canceled AS canceled, \
			n.finished AS finished, \
			o.name AS object, \
			code.row AS okpd, \
			code2.row AS okpd2 \
		FROM notifications_%s AS n \
		LEFT JOIN notifications_objects_%s AS o ON n.pnum = o.pnum \
		LEFT JOIN nsiokpd1 AS code ON o.okpd = code.code \
		LEFT JOIN nsiokpd2 AS code2 ON o.okpd = code2.code;
	sql_attr_timestamp = published
	sql_attr_float     = gMaxPrice
	sql_field_string   = object
	sql_field_string   = info
	sql_attr_uint      = okpd
	sql_attr_uint      = okpd2
	sql_attr_uint      = org_region
	sql_attr_uint      = org
	sql_attr_bigint    = pnum
}

source $p_s_ncustomers_%s : base_source {
	sql_query          = \
		SELECT n.pnum::bigint AS id, n.pnum::bigint AS pnum, n.objectInfo AS info, n.org_regnum::bigint AS org, \
			n.org_region AS org_region, \
			EXTRACT(EPOCH FROM n.published AT TIME ZONE 'MSK') AS published, \
			c.maxPrice AS gMaxPrice, \
			c.customer_regnum::bigint AS customer \
		FROM notifications_customers_%s AS c, notifications_%s AS n \
		WHERE c.pnum = n.pnum AND c.customer_regnum IS NOT NULL \
		  AND c.customer_regnum <> '' AND n.org_regnum IS NOT NULL;
	sql_attr_timestamp = published
	sql_attr_float     = gMaxPrice
	sql_field_string   = info
	sql_attr_uint      = customer
	sql_attr_uint      = org_region
	sql_attr_uint      = org
	sql_attr_bigint    = pnum
}


index $p_notifications_%s : base_index {
	source        = $p_s_notifications_%s
	path          = /var/lib/sphinxsearch/data/$p_notifications_%s
	min_word_len  = 1
	min_infix_len = 2
	# charset_type = utf-8
	# minimum indexed word length
	# default is 1 (index everything)
	# enable_star = 1 #deprecated
}

index $p_ncustomers_%s : base_index {
	source        = $p_s_ncustomers_%s
	path          = /var/lib/sphinxsearch/data/$p_ncustomers_%s
	min_word_len  = 1
	min_infix_len = 2
	# charset_type = utf-8
	# minimum indexed word length
	# default is 1 (index everything)
	# enable_star = 1 #deprecated
}
`

const products string = `

##
# Products %s
##

source $p_s_products_%s : base_source {
	sql_query_pre      = SET NAMES 'UTF8';
	sql_query          = \
		SELECT p.row, p.name AS Name, p.price AS Price, p.sum AS Sum, p.regnum AS RegNum, p.purchasenumber AS PurchaseNumber, \
			p.regnum::bigint AS rnum, \
			(CASE WHEN p.purchasenumber~E'^\\d+$' THEN p.purchasenumber::bigint ELSE 0 END) AS pnum, \
			p.customer AS Customer, \
			p.customerinn AS CINN, \
			p.customerinn::bigint AS CINNn, \
			p.customerregion AS Region, \
			p.customerkpp AS CKPP, \
			p.customerkpp::bigint AS CKPPn, \
			p.supplier AS Supplier, \
			p.supplierinn AS SINN, \
			EXTRACT(EPOCH FROM timestamptz(p.Signed) AT TIME ZONE 'MSK') AS Signed, \
			EXTRACT(EPOCH FROM timestamptz(p.exec) AT TIME ZONE 'MSK') AS Exec, \
			p.Paid AS Paid, \
			n.regnum::bigint AS orgcode, \
			p.okpd1 AS okpdCode, \
			p.okpd2 AS okpdCode2 \
		FROM products_%s AS p, nsiorg AS n \
		WHERE p.customerinn = n.inn AND p.customerkpp = n.kpp AND n.actual = true;
	sql_field_string   = Name
	sql_field_string   = RegNum
	sql_field_string   = PurchaseNumber
	sql_field_string   = SINN
	sql_attr_timestamp = Signed
	sql_attr_timestamp = Exec
	sql_attr_float     = Sum
	sql_attr_float     = Price
	sql_attr_float     = Paid
	sql_attr_bigint    = pnum
	sql_attr_bigint    = rnum
	sql_attr_uint      = Region
	sql_attr_uint      = CINNn
	sql_attr_uint      = CKPPn
	sql_attr_uint      = orgcode
	sql_attr_uint      = okpdCode
	sql_attr_uint      = okpdCode2
}

index $p_products_%s : base_index {
	source        = $p_s_products_%s
	path          = /var/lib/sphinxsearch/data/$p_products_%s
	min_word_len  = 2
	min_infix_len = 2
}
`

const fz223 = `
source $p_s_notifcations223 : base_source {
	sql_query_pre      = SET NAMES 'UTF8';
	sql_query          = \
		SELECT n.row,     \
			   n.regnum,     \
			   n.name,	    \
			   n.customerINN::bigint as cINN,	 \
			   n.customerKPP::bigint as cKPP,	  \
			   n.placerINN::bigint as pINN,	   \
			   n.placerKPP::bigint as pKPP,       \
			   substring(n.customerKPP,1,2)::integer as region, \
			   n.created,			 \
			   n.published,		  \
			   l.subject as lotname,			   \
			   l.initSum as sum
		FROM notifications223 as n, lots223 as l	 \
		WHERE n.regnum = l.regnum;

		sql_field_string = regnum
		sql_field_string = name
		sql_field_string = lotname
		sql_attr_uint = cINN
		sql_attr_uint = cKPP
		sql_attr_uint = pINN
		sql_attr_uint = pKPP
		sql_attr_uint = region
		sql_attr_timestamp = created
		sql_attr_timestamp = published
		sql_attr_float = sum
}
index $p_s_notifications223 : base_index {
	source        = $p_s_notifications223
	path          = /var/lib/sphinxsearch/data/$p_s_notifications223
	min_word_len  = 2
	min_infix_len = 2
}
`

const spz string = `

##
# SPZ
##

source $p_s_spz : base_source {
	sql_query        = \
		SELECT row, fullname, shortname, inn, substring(kpp, 1, 2)::int, regnum AS region \
		FROM nsiorg \
		WHERE actual = true;
	sql_field_string = fullname
	sql_field_string = shortname
	sql_field_string = inn
	sql_field_string = regnum
}

index $p_nsiorg : base_index {
	source        = $p_s_spz
	path          = /var/lib/sphinxsearch/data/$p_spz
	min_word_len  = 1
	min_infix_len = 2
	# charset_type = utf-8
	# minimum indexed word length
	# default is 1 (index everything)
	# enable_star = 1 #deprecated
}
`

const okpd string = `

##
# OKPD
##

source $p_s_okpd : base_source {
	sql_query        = \
		SELECT row, code, name \
		FROM nsiokpd1 \
		WHERE actual = true;
	sql_field_string = code
	sql_field_string = name
	sql_attr_uint    = row
}

index $p_nsiokpd : base_index {
	source        = $p_s_okpd
	path          = /var/lib/sphinxsearch/data/$p_okpd
	min_word_len  = 1
	min_infix_len = 2
	# charset_type = utf-8
	# minimum indexed word length
	# default is 1 (index everything)
	# enable_star = 1 #deprecated
}
`

const okpd2 string = `

##
# OKPD2
##

source $p_s_okpd2 : base_source {
	sql_query        = \
		SELECT row, code, name \
		FROM nsiokpd2 \
		WHERE actual = true;
	sql_field_string = code
	sql_field_string = name
	sql_attr_uint    = row
}

index $p_nsiokpd2 : base_index {
	source        = $p_s_okpd2
	path          = /var/lib/sphinxsearch/data/$p_okpd2
	min_word_len  = 1
	min_infix_len = 2
	# charset_type = utf-8
	# minimum indexed word length
	# default is 1 (index everything)
	# enable_star = 1 #deprecated
}
`

const okpd1 string = `

##
# OKPD1
##

source $p_s_okpd1 : base_source {
	sql_query        = \
		SELECT row, code, name \
		FROM nsiokpd1 \
		WHERE actual = true;
	sql_field_string = code
	sql_field_string = name
	sql_attr_uint    = row
}

index $p_nsiokpd1 : base_index {
	source        = $p_s_okpd1
	path          = /var/lib/sphinxsearch/data/$p_okpd1
	min_word_len  = 1
	min_infix_len = 2
	# charset_type = utf-8
	# minimum indexed word length
	# default is 1 (index everything)
	# enable_star = 1 #deprecated
}
`

const ndocs = `

##
# Notifications files %s
##

source $p_s_ndocs_%s : base_source {
	sql_query          = \
		SELECT f.purchasenumber::bigint AS id, f.data AS data, \
			EXTRACT(EPOCH FROM n.published AT TIME ZONE 'MSK') AS published \
		FROM notifications_files_%s AS f, notifications_%s AS n \
		WHERE f.data <> '' AND f.data IS NOT NULL \
		  AND n.pnum::bigint = f.purchasenumber::bigint;
	sql_attr_timestamp = published
	# sql_field_string = data
	sql_attr_bigint    = id
}

index $p_ndocs_%s : base_index {
	source        = $p_s_ndocs_%s
	path          = /var/lib/sphinxsearch/data/$p_ndocs_%s
	min_word_len  = 1
	min_infix_len = 2
	# charset_type = utf-8
	# minimum indexed word length
	# default is 1 (index everything)
	# enable_star = 1 #deprecated
}
`

const suppliers = `

##
# Suppliers
##

source $p_s_nsisuppliers : base_source {
	sql_query = \
		SELECT row, name, inn \
		FROM nsisuppliers;
	sql_field_string = name
	sql_field_string = inn
}

index $p_nsisuppliers : base_index {
	source        = $p_s_nsisuppliers
	path          = /var/lib/sphinxsearch/data/$p_nsisuppliers
	min_word_len  = 1
	min_infix_len = 2
}
`

var rangeSphinxIndexes = []string{products}
var rangeSphinxIndexesIgnoreFirst = []string{notifications, ndocs}
var singleSphinxIndexes = []string{spz, okpd, okpd1, okpd2, suppliers}
