package coredata

//LookupRequestData param for lookup request
type LookupRequestData struct {

	//LovID ID of lookup
	LovID string `json:"lovId"`
	//Version version to search. for check status of lookup
	Version string `json:"Version"`

	//FilteredCodes codes to be filtered from search
	FilteredCodes []string
}
