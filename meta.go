package spider

const (
	ContentTypeACC    = ".acc"
	ContentTypeABW    = ".abw"
	ContentTypeARC    = ".arc"
	ContentTypeAVI    = ".avi"
	ContentTypeAZW    = ".azw"
	ContentTypeBIN    = ".bin"
	ContentTypeBMP    = ".bmp"
	ContentTypeBZ     = ".bz"
	ContentTypeBZ2    = ".bz2"
	ContentTypeCSH    = ".csh"
	ContentTypeCSS    = ".css"
	ContentTypeCSV    = ".csv"
	ContentTypeDOC    = ".doc"
	ContentTypeDOCX   = ".docx"
	ContentTypeEOT    = ".eot"
	ContentTypeEPUB   = ".epub"
	ContentTypeGIF    = ".gif"
	ContentTypeHTM    = ".htm"
	ContentTypeHTML   = ".html"
	ContentTypeICO    = ".ico"
	ContentTypeICS    = ".ics"
	ContentTypeJAR    = ".jar"
	ContentTypeJPEG   = ".jpeg"
	ContentTypeJPG    = ".jpg"
	ContentTypeJS     = ".js"
	ContentTypeJSON   = ".json"
	ContentTypeJSONLD = ".jsonld"
	ContentTypeMID    = ".mid"
	ContentTypeMIDI   = ".midi"
	ContentTypeMJS    = ".mjs"
	ContentTypeMP3    = ".mp3"
	ContentTypeMPEG   = ".mpeg"
	ContentTypeMPKG   = ".mpkg"
	ContentTypeODP    = ".odp"
	ContentTypeODS    = ".ods"
	ContentTypeODT    = ".odt"
	ContentTypeOGA    = ".oga"
	ContentTypeOGV    = ".ogv"
	ContentTypeOGX    = ".ogx"
	ContentTypeOTF    = ".otf"
	ContentTypePNG    = ".png"
	ContentTypePDF    = ".pdf"
	ContentTypePPT    = ".ppt"
	ContentTypePPTX   = ".pptx"
	ContentTypeRAR    = ".rar"
	ContentTypeRTF    = ".rtf"
	ContentTypeSH     = ".sh"
	ContentTypeSVG    = ".svg"
	ContentTypeSWF    = ".swf"
	ContentTypeTAR    = ".tar"
	ContentTypeTIF    = ".tif"
	ContentTypeTIFF   = ".tiff"
	ContentTypeTTF    = ".ttf"
	ContentTypeTXT    = ".txt"
	ContentTypeVSD    = ".vsd"
	ContentTypeWAV    = ".wav"
	ContentTypeWEBA   = ".weba"
	ContentTypeWEBM   = ".webm"
	ContentTypeWEBP   = ".webp"
	ContentTypeWOFF   = ".woff"
	ContentTypeWOFF2  = ".woff2"
	ContentTypeXHTML  = ".xhtml"
	ContentTypeXLS    = ".xls"
	ContentTypeXLSX   = ".xlsx"
	ContentTypeXML    = ".xml"
	ContentTypeXUL    = ".xul"
	ContentTypeZIP    = ".zip"
	ContentType3GP    = ".3GP"
	ContentType3G2    = ".3G2"
	ContentType7Z     = ".7Z"
)

var ContentTypes = map[string]string{
	"audio/aac":                    ContentTypeACC,
	"application/x-abiwor":         ContentTypeABW,
	"application/x-freearc":        ContentTypeARC,
	"video/x-msvideo":              ContentTypeAVI,
	"application/vnd.amazon.ebook": ContentTypeAZW,
	"application/octet-stream":     ContentTypeBIN,
	"image/bmp":                    ContentTypeBMP,
	"application/x-bzip":           ContentTypeBZ,
	"application/x-bzip2":          ContentTypeBZ2,
	"application/x-csh":            ContentTypeCSH,
	"text/css":                     ContentTypeCSS,
	"text/csv":                     ContentTypeCSV,
	"application/msword":           ContentTypeDOC,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ContentTypeDOCX,

	"application/vnd.ms-fontobject":       ContentTypeEOT,
	"application/epub+zip":                ContentTypeEPUB,
	"image/gif":                           ContentTypeGIF,
	"text/html":                           ContentTypeHTML,
	"image/vnd.microsoft.ico":             ContentTypeICO,
	"text/calendar":                       ContentTypeICS,
	"application/java-archiv":             ContentTypeJAR,
	"image/jpeg":                          ContentTypeJPEG,
	"text/javascript":                     ContentTypeJS,
	"application/javascript":              ContentTypeJS,
	"application/json":                    ContentTypeJSON,
	"application/ld+json":                 ContentTypeJSONLD,
	"audio/midi audio/x-midi":             ContentTypeMIDI,
	"audio/mpeg":                          ContentTypeMP3,
	"video/mpeg":                          ContentTypeMPEG,
	"application/vnd.apple.installer+xml": ContentTypeMPKG,
	"application/vnd.oasis.opendocument.presentation": ContentTypeMPEG,
	"application/vnd.oasis.opendocument.spreadsheet":  ContentTypeODP,
	"application/vnd.oasis.opendocument.text":         ContentTypeODS,
	"audio/ogg":                     ContentTypeODT,
	"video/ogg":                     ContentTypeOGA,
	"application/ogg":               ContentTypeOGX,
	"font/otf":                      ContentTypeOTF,
	"image/png":                     ContentTypePNG,
	"application/pdf":               ContentTypePDF,
	"application/vnd.ms-powerpoint": ContentTypePPT,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": ContentTypePPTX,

	"application/x-rar-compressed":  ContentTypeRAR,
	"application/rtf":               ContentTypeRTF,
	"application/x-sh":              ContentTypeSH,
	"image/svg+xml":                 ContentTypeSVG,
	"application/x-shockwave-flash": ContentTypeSWF,
	"application/x-tar":             ContentTypeTAR,
	"image/tiff":                    ContentTypeTIFF,
	"font/ttf":                      ContentTypeTTF,
	"text/plain":                    ContentTypeTXT,
	"application/vnd.visio":         ContentTypeVSD,
	"audio/wav":                     ContentTypeWAV,
	"audio/webm":                    ContentTypeWEBA,
	"video/webm":                    ContentTypeWEBM,
	"image/webp":                    ContentTypeWEBP,
	"font/woff":                     ContentTypeWOFF,
	"font/woff2":                    ContentTypeWOFF2,
	"application/xhtml+xml":         ContentTypeXHTML,
	"application/vnd.ms-excel":      ContentTypeXLS,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ContentTypeXLSX,

	"application/xml":                 ContentTypeXML,
	"application/vnd.mozilla.xul+xml": ContentTypeXUL,
	"application/zip":                 ContentTypeZIP,
	"video/3gpp":                      ContentType3GP,
	"audio/3gpp":                      ContentType3GP,
	"video/3gpp2":                     ContentType3G2,
	"audio/3gpp2":                     ContentType3G2,
	"application/x-7z-compressed":     ContentType7Z,
}
