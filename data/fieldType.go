package data

type FieldType int

const (
	AddressLine2 FieldType = iota
	AirportCode            // AIRPORT
	AirporContinent
	AirportCountry
	AirportElevation
	AirportGPSCode
	AirportLatitude
	AirportLongitude
	AirportMunicipality
	AirportName
	AirportRegionCode
	AnimalCommonName // ANIMAL
	AnimalScientificName
	AppBundleId // APP
	AppName
	AppVersion
	Avatar
	Base64ImageUrl
	BinomialDistribution
	BitcoinAddress
	Blank
	Boolean
	Buzzword
	CarMake
	CarModel
	CarModelYear
	CarVIN
	CatchPhrase
	CharacterSequence
	City
	Color
	CompanyName
	ConstructionHeavyEquipment // Construction
	ConstructionMaterial
	ConstructionRole
	ConstructionStandardCostCode
	ConstructionSubcontractCategory
	ConstructionTrade
	Country // Country
	CountryCode
	CreditCardNumber // Credit Card
	CreditCardType
	Currency // Currency
	CurrencyCode
	CustomList
	Datetime
	DepartmentCorporate // Departments
	DepartmentRetail
	DigitSequence
	DomainName
	DrugNameBrand // Drug
	DrugNameGeneric
	DummyImageUrl
	DUNSNumber
	EIN
	EmailAddress
	Encrypt
	EhtereumAddress
	FakeCompanyName
	FamilyNameChinese
	FDANDCCode
	FileName
	FirstName // Names
	FirstNameEuropean
	FirstNameFemale
	FirstNameMale
	Formula // using other columns
	Frequency
	FullName
	Gender
	GenderAbrev
	GenderBinary
	GenderFacebook
	GivenNameChinese
	GUID
	HexColor
	IBAN
	ICD10DiagnosisCode // ICD
	ICD10DxDescLong
	ICD10DxDescShort
	ICD10ProcDescLong
	ICD10ProcDescShort
	ICD10ProcedureCode
	ICD9DxDescDiagnosisCode
	ICD9DxDescLong
	ICD9DescShort
	ICD9ProcDescLong
	ICD9ProcDescShort
	ICDProcedureCode
	IpAddressV4 // ip address
	IpAddressV6
	IpAddressV6CIDR
	ISBN
	JobTitle
	Language
	LastName // Name
	Latitude // Coordinates
	Longitude
	MACAddress
	MD5
	MedicareBeneficiaryId
	MimeType
	MobileDeviceBrand // Mobilde Device
	MobileDeviceModel
	MobileDeviceOS
	MobileDeviceReleaseDate
	Money
	MongoDBObjectID
	MovieGenres
	MovieTitle
	NatoPhonetic
	NaughtyString
	NHSNumber
	Number
	Paragraphs
	Password
	Phone
	PlantCommonName // Plants
	PlantFamily
	PlantScientificName
	PostalCode // Address
	ProductGrocery
	Race
	RepeatingElement
	RowNumber
	Sentences
	SHA1
	SHA256
	ShirtSize
	ShortHexColor // Color
	Slogan
	SQLExpression // sql expression to use as value
	SSN
	State
	StockIndustry // Stock Market
	StockMarket
	StockMarketCap
	StockName
	StockSector
	StockSymbol
	StreetAddress // Address
	StreetName
	StreetNumber
	StreetSuffix
	Suffix   // Name
	Template // Concatenate different values
	Time
	TimeZone
	Title          // Name
	TopLevelDomain // .com .edu
	ULID
	University
	URL
	UserAgent
	Username
	Words
)
