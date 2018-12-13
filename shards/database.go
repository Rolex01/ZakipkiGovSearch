package shards

import (
	"database/sql"
	"log"
	"strings"
)

const tableContracts string = `CREATE TABLE IF NOT EXISTS contracts_%s (
	row            SERIAL PRIMARY KEY,
	RegNum         text   UNIQUE,
	PurchaseNumber text,
	Published      timestamp,
	Signed         timestamp,
	Exec           timestamp,
	Budget         text,
	OKTMO          text,
	Suppliers      text,
	Customer       text,
	Placer         text,
	Price          decimal,
	Currency       text,
	Rate           decimal,
	Paid           decimal,
	Link           text,
	Status         text,
	Updated        timestamp,
	Version        integer,
	NMCK           decimal,
	Ntname         text,
	Lot            integer
);`

const indxContractsRegnum string = `CREATE INDEX c_%s_regnum ON contracts_%s(regnum);`
const indxContractsPnum string = `CREATE INDEX c_%s_regnum ON contracts_%s(purchasenumber);`
const indxContractsSigned string = `CREATE INDEX c_%s_signed ON contracts_%s(signed);`

const tableProducts string = `CREATE TABLE IF NOT EXISTS products_%s (
	row             SERIAL PRIMARY KEY,
	RegNum          text,
	PurchaseNumber  text,
	Number          text,
	Name            text,
	OKPD            text,
	OKPDInfo        text,
	Units           text,
	OKEI            text,
	Quantity        decimal,
	Price           decimal,
	PriceRu         decimal,
	Sum             decimal,
	SumRu           decimal,
	Rate            decimal,
	Currency        text,
	Published       timestamp,
	Signed          timestamp,
	Exec            timestamp,
	Budget          text,
	Customer        text,
	CustomerRegion  integer,
	CustomerINN     text,
	CustomerKPP     text,
	SupplierINN     text,
	SupplierKPP     text,
	SupplierAddress text,
	SupplierPhone   text,
	Supplier        text,
	ContractPrice   decimal,
	Paid            decimal,
	Status          text,
	PlacingWay      text,
	BudgetSource    text,
	Updated         timestamp,
	Version         integer,
	okpd1           integer,
	okpd2           integer
);`

const indxProductsRegnum string = `CREATE INDEX p_%s_regnum ON products_%s(regnum);`
const indxProductsPnum string = `CREATE INDEX p_%s_pnum ON products_%s(purchasenumber);`
const indxProductsSigned string = `CREATE INDEX p_%s_signed ON products_%s(signed);`

const tableContractsFiles = `CREATE TABLE IF NOT EXISTS contracts_files_%s (
	row         SERIAL PRIMARY KEY,
	regnum      text,
	description text,
	name        text,
	url         text,
	size        text,
	data        text,
	ignore      bool,
	downloaded  bool,
	published   timestamp
);`

const indxContractsFilesRegnum string = `CREATE INDEX cf_%s_regnum ON contracts_files_%s(regnum);`

const tableContractsProc = `CREATE TABLE IF NOT EXISTS contracts_proc_%s (
	row       SERIAL PRIMARY KEY,
	id        integer UNIQUE,
	regnum    text,
	published timestamp,
	version   integer,
	final     bool,
	stage     text,
	applied   bool,
	paid      decimal,
	shard     text
);`

const indxContractsProcRegnum string = `CREATE INDEX cp_%s_regnum ON contracts_proc_%s(regnum);`

const tableNotifications string = `CREATE TABLE IF NOT EXISTS notifications_%s (
	pnum           bigint PRIMARY KEY,
	purchaseNumber text   UNIQUE,
	objectInfo     text,
	url            text,
	org_regnum     text,
	org_inn        text,
	org_name       text,
	org_kpp        text,
	org_role       text,
	org_indx       bigint,
	org_region     int,
	type           text,
	published      timestamp,
	maxPrice       decimal,
	canceled       bool   default false,
	finished       bool   default false,
	Updated        timestamp,
	row            SERIAL
);`

const indxNotificationsPnum string = `CREATE INDEX n_%s_pnum ON notifications_%s(pnum);`
const indxNotificationsPublished string = `CREATE INDEX n_%s_pnum ON notifications_%s(published);`

const tableNotificationsLots string = `CREATE TABLE IF NOT EXISTS notifications_lots_%s (
	row            SERIAL PRIMARY KEY,
	lotnum         integer,
	regnum         text,
	pnum           bigint,
	maxPrice       decimal,
	currency       varchar(5),
	finance_source text,
	ctprice        decimal,
	canceled       bool
);`

const indxNotificationsLotsPnum string = `CREATE INDEX nl_%s_pnum ON notifications_lots_%s(pnum);`
const indxNotificationsLotsRegnum string = `CREATE INDEX nl_%s_regnum ON notifications_lots_%s(regnum);`

const tableNotificationsCustomers string = `CREATE TABLE IF NOT EXISTS notifications_customers_%s (
	row             SERIAL PRIMARY KEY,
	lotnum          integer,
	pnum            bigint,
	maxPrice        decimal,
	customer_name   text,
	customer_inn    text,
	customer_kpp    text,
	customer_indx   bigint,
	customer_regnum text,
	customer_region int,
	regnum          text,
	ctprice         decimal,
	ctobject        text
);`

const indxNotificationsCustomersPnum string = `CREATE INDEX nc_%s_pnum ON notifications_customers_%s(pnum);`
const indxNotificationsCustomersRegnum string = `CREATE INDEX nc_%s_regnum ON notifications_customers_%s(regnum);`
const indxNotificationsCustomersIndx string = `CREATE INDEX nc_%s_ci ON notifications_customers_%s(customer_indx);`

const tableNotificationsObjects string = `CREATE TABLE IF NOT EXISTS notifications_objects_%s (
	row       SERIAL PRIMARY KEY,
	pnum      bigint,
	lotnum    integer,
	customer  bigint,
	name      text,
	okpd      text,
	okei      text,
	quantity  decimal,
	countable bool   default false,
	price     decimal,
	sum       decimal
);`

const indxNotificationsObjectsPnum string = `CREATE INDEX nob_%s_pnum ON notifications_objects_%s(pnum);`
const indxNotificationsObjectsCustomer string = `CREATE INDEX nob_%s_customer ON notifications_objects_%s(customer);`

const tableNotificationsFiles = `CREATE TABLE IF NOT EXISTS notifications_files_%s (
	row            SERIAL PRIMARY KEY,
	purchaseNumber bigint,
	description    text,
	name           text,
	url            text,
	size           text,
	data           text,
	ignore         bool,
	downloaded     bool,
	fixorgs        bool
);`

const indxNotificationsFilesCustomer string = `CREATE INDEX nf_%s_pnum ON notifications_files_%s(purchasenumber);`
const indxNotificationsFilesNotNull string = `CREATE INDEX nf_%s_notnull_data ON notifications_files_%s(purchasenumber) where data is not NULL;`

const tableExecutions = `CREATE TABLE IF NOT EXISTS executions_%s (
	row      SERIAL PRIMARY KEY,
	regnum   text,
	name     text,
	number   text,
	date     timestamp,
	currency varchar(5),
	paid     decimal,
    UNIQUE (regnum, number, name, date)
);`

const indxExecutionsMain string = `CREATE INDEX executions_%s_pair ON executions_%s(regnum, number, name, date);`
const indxExecutionsRegnum string = `CREATE INDEX executions_%s_regnum ON executions_%s(regnum);`

const tableNsiBudget = `CREATE TABLE IF NOT EXISTS nsiBudget (
	row    SERIAL PRIMARY KEY,
	code   text   UNIQUE,
	parent text,
	name   text,
	actual bool
);`

const indxNsiBudget string = `CREATE INDEX nsbudget ON nsibudget(code);`

const tableNsiOKTMO = `CREATE TABLE IF NOT EXISTS nsiOKTMO (
	row      SERIAL PRIMARY KEY,
	code     text   UNIQUE,
	parent   text,
	section  text,
	fullName text,
	actual   bool
);`

const indxNsiOKTMO string = `CREATE INDEX nsoktmo ON nsioktmo(code);`

const tableNsiFO = `CREATE TABLE IF NOT EXISTS nsiFO (
	row       SERIAL PRIMARY KEY,
	shortname text,
	fullname  text,
	code      integer,
	region    text
);`

const tableNsiOKEI = `CREATE TABLE IF NOT EXISTS nsiOKEI (
	row                 SERIAL PRIMARY KEY,
	Code                text   UNIQUE,
	FullName            text,
	SectionCode         text,
	SectionName         text,
	GroupId             int,
	GroupName           text,
	LocalName           text,
	InternationalName   text,
	LocalSymbol         text,
	InternationalSymbol text,
	Actual              bool
);`

const indxNsiOKEI string = `CREATE INDEX nsokei ON nsiokei(code);`

const tableNsiOKPD = `CREATE TABLE IF NOT EXISTS nsiOKPD (
	row     SERIAL  PRIMARY KEY,
	prow    integer DEFAULT 0,
	code    text    UNIQUE,
	parent  text,
	name    text,
	comment text,
	actual  bool
);`

const indxNsiOKPD string = `CREATE INDEX nsokpd ON nsiokpd(code);`


const tableNsiOrg = `CREATE TABLE IF NOT EXISTS nsiOrg (
	row               SERIAL PRIMARY KEY,
	RegNum            text   UNIQUE,
	ShortName         text,
	FullName          text,
	OKVED             text,
	Url               text,
	INN               text,
	OGRN              text,
	KPP               text,
	OKTMO             text,
	PostalAddress     text,
	ContactPerson     text,
	SubordinationType text,
	Actual            bool
);`

const tableContracts223 string = `CREATE TABLE IF NOT EXISTS contracts223 (
	row                     serial primary key,
	ContractRegNumber       text   UNIQUE,
	URL                     text,
	CreateDateTime          timestamp,
	Customer                text,
	CustomerINN             text,
	CustomerKPP             text,
	CustomerName            text,
	Placer                  text,
	PlacerName              text,
	PlacerINN               text,
	PlacerKPP               text,
	PublicationDate         text,
	Status                  text,
	Version                 integer,
	ModificationDescription text,
	DigitalPurchase         bool,
	DigitalPurchaseCode     text,
	Provider                bool,
	ProviderCode            text,
	ChangeContract          bool,
	Attachments             text,
	Name                    text,
	ContractDate            timestamp,
	ApproveDate             timestamp,
	PurchaseNoticeNumber    text,
	PublicationDateTime     timestamp,
	PurchaseName            text,
	SubjectContract         text,
	PurchaseTypeInfoCode    text,
	PurchaseTypeInfoName    text,
	ResumeDate              timestamp,
	Supplier                text,
	SupplierName            text,
	SupplierINN             text,
	SupplierKPP             text,
	HasSubcontractor        bool,
	HasSubcontractorCode    text,
	SubcontractorsTotal     text,
	HasGoodInfo             bool,
	AdditionalInfo          text,
	Price                   decimal,
	ExchangeRate            decimal,
	RubPrice                decimal,
	Currency                text,
	StartExecutionDate      text,
	EndExecutionDate        text,
	ContractPositions       text,
	ContractChangeDocs      text
);`

const tablePurchase223 string = `CREATE TABLE IF NOT EXISTS notifications223 (
	row 					serial primary key,
	regnum					text unique,
	guid					text unique,
	name					text not null,
	customerINN				text not null,
	customerKPP				text not null,
	customerName			text not null,
	placerINN				text not null,
	placerKPP				text not null,
	placerName			    text not null,
	methodCode				text,
	methodName				text,
	created					timestamp,
	published				timestamp,
	modified				timestamp,
	joint					bool,
	emergency				bool,
	status					text,
	version					integer

);`
const tablePurchaseLots223 string = `CREATE TABLE IF NOT EXISTS lots223 (
	row						serial primary key,
	regnum					text,
	guid					text unique,
	subject					text not null,
	currency				text not null,
	initSum					decimal default -1,
	joint					bool,
	cancelled				bool,
	version					integer
);`
const tablePurchaseLotItems223 string = `CREATE TABLE IF NOT EXISTS lotItems223 (
	row						serial primary key,
	regnum					text not null,
	glot					text not null,
	guid					text unique not null,
	okpd2Code				text,
	okpd2Desc				text,
	okved2Code				text,
	okved2Desc				text,
	okeiCode				text,
	okeiDesc				text,
	qty						decimal,
	version					integer
);`
const indxNsiOrgRegnum string = `CREATE INDEX nsorg_regnum ON nsiorg(regnum);`
const indxNsiOrgINN string = `CREATE INDEX nsorg_inn ON nsiorg(inn);`

const tableNsiSuppliers = `CREATE TABLE IF NOT EXISTS nsiSuppliers (
	row     SERIAL PRIMARY KEY,
	name    text,
	INN     text,
	KPP     text,
	Address text,
	Phone   text
);`

const indxNsiSuppliersINN string = `CREATE INDEX nssuppliers_inn ON nsisuppliers(inn);`
const indxNsiSuppliersKPP string = `CREATE INDEX nssuppliers_kpp ON nsisuppliers(kpp);`

// InitSchema creates tables structure if the schema is empty
func InitSchema(db *sql.DB) error {
	for _, q := range generateQueries() {
		if err := runQuery(q, db); err != nil {
			log.Fatalf("Can't init schema correctly!\nThe error:%v\nThe query:%s", err, q)
		}
	}
	return nil
}

func InitIndexes(db *sql.DB) error {
	for _, q := range generateIndexesQueries() {
		if err := runQuery(q, db); err != nil {
			log.Fatalf("Can't init schema correctly!\nThe error:%v\nThe query:%s", err, q)
		}
	}
	return nil
}

func generateIndexesQueries() (queries []string) {
	for _, scheme := range singleIndexes {
		queries = append(queries, scheme)
	}
	for _, scheme := range rangeIndexes {
		queries = append(queries, strings.Join(GenRangeByPattern(scheme), ";"))
	}

	return
}

func generateQueries() (queries []string) {
	for _, scheme := range signleTables {
		queries = append(queries, scheme)
	}
	for _, scheme := range rangeTables {
		queries = append(queries, strings.Join(GenRangeByPattern(scheme), ";"))
	}

	return
}

func runQuery(q string, db *sql.DB) error {
	_, err := db.Exec(q)
	return err
}

var rangeTables = []string{
	tableContracts,
	tableContractsFiles,
	tableProducts,
	tableContractsProc,
	tableNotifications,
	tableNotificationsLots,
	tableNotificationsCustomers,
	tableNotificationsFiles,
	tableNotificationsObjects,
	tableExecutions }

var signleTables = []string{
	tableNsiBudget,
	tableNsiOKEI,
	tableNsiFO,
	tableNsiOKPD,
	tableNsiOKTMO,
	tableNsiOrg,
	tableContracts223,
	tablePurchase223,
	tablePurchaseLots223,
	tablePurchaseLotItems223,
	tableNsiSuppliers }

var singleIndexes = []string{
	indxNsiOKPD,
	indxNsiBudget,
	indxNsiOrgRegnum,
	indxNsiOrgINN,
	indxNsiOKEI,
	indxNsiOKTMO,
	indxNsiSuppliersINN,
	indxNsiSuppliersKPP }

var rangeIndexes = []string{
	indxContractsPnum,
	indxContractsSigned,
	indxProductsRegnum,
	indxProductsPnum,
	indxContractsFilesRegnum,
	indxContractsProcRegnum,
	indxNotificationsPublished,
	indxNotificationsLotsPnum,
	indxNotificationsLotsRegnum,
	indxNotificationsCustomersPnum,
	indxNotificationsCustomersRegnum,
	indxNotificationsCustomersIndx,
	indxNotificationsObjectsPnum,
	indxNotificationsObjectsCustomer,
	indxNotificationsFilesCustomer,
	indxNotificationsFilesNotNull,
	indxExecutionsMain,
	indxExecutionsRegnum }
