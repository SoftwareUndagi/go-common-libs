package coredata

import (
	"github.com/SoftwareUndagi/go-common-libs/coremodel"
)

//LookupRequestResultWrapper wrapper lookup + still up to date state
type LookupRequestResultWrapper struct {
	//LookupID id of lookuo
	LookupID string `json:"loookupId"`
	//flag is lookup data still up to date
	StillUptodate bool `json:"stillUptodate"`
	//LookupData lookup data
	LookupData *coremodel.LookupHeader `json:"lookupData"`
}
