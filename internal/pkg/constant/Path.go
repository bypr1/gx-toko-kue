package constant

const PATH_EXCEL = "excels/"
const PATH_IMAGE = "images/"
const PATH_PDF = "pdfs/"

const PATH_TESTING = "testings/"
const PATH_TRANSACTION = "transactions/"
const PATH_CAKE = "cakes/"

func PathPDFTesting() string {
	return PATH_PDF + PATH_TESTING
}

func PathImageTesting() string {
	return PATH_IMAGE + PATH_TESTING
}

func PathExcelTesting() string {
	return PATH_EXCEL + PATH_TESTING
}

func PathExcelTransaction() string {
	return PATH_EXCEL + PATH_TRANSACTION
}

func PathImageCake() string {
	return PATH_IMAGE + PATH_CAKE
}
