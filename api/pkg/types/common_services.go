package types

type CommonService string

const (
	FileUpload         CommonService = "FileUpload"
	FileDownload       CommonService = "FileDownload"
	FirmwareUpdate     CommonService = "FirmwareUpdate"
	Blink              CommonService = "Blink"
	NewDataSetWriterId CommonService = "NewDataSetWriterId"
)
