package coredao

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/SoftwareUndagi/go-common-libs/coredata"
	"github.com/SoftwareUndagi/go-common-libs/coremodel"
	"github.com/SoftwareUndagi/go-common-libs/dao"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//FindLookupHeaders read lookup headers with lov id
func FindLookupHeaders(lovIDs []string, db *gorm.DB, logEntry *logrus.Entry) (lookupHeaders []coremodel.LookupHeader, err error) {
	var colName string
	logEntry = logEntry.WithField("lovIds", lovIDs).WithField("model", "LookupHeader")
	colName, err = dao.DefaultDaoManager.GetColumnName("LookupHeader", "ID")
	if err != nil {
		logEntry.WithError(err).Errorf("Fail read column for model : LookupHeader , field:  ID, error: %s", err.Error())
		return
	}
	whSmt := fmt.Sprintf("%s in (?) ", colName)
	if dbRslt := db.Where(whSmt, lovIDs).Find(&lookupHeaders); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithError(err).Errorf("Fail to query for lookup. error: %s", err.Error())
		return
	}
	return
}

//assignToLookupHeader add lookup details to lookup header
func assignToLookupHeader(indexedLookup *map[string]*coremodel.LookupHeader, lookupDetails *[]coremodel.LookupDetail) {
	mapVal := *indexedLookup
	for _, l := range *lookupDetails {
		mapVal[l.LovID].AppendLookupDetail(l)
	}
}

//QueryForLookupWithSQLVersion run query to query for driven by sql query
func QueryForLookupWithSQLVersion(lookupHeaders []coremodel.LookupHeader, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	versionFinderStatment := []string{}
	mapByLovID := make(map[string]*coremodel.LookupHeader)

	for _, h := range lookupHeaders {
		versionFinderStatment = append(versionFinderStatment, h.SQLForVersion)
		mapByLovID[h.ID] = &h

	}
	finalSQL := strings.Join(versionFinderStatment, " union all ")
	rows, errRow := db.Raw(finalSQL).Rows()
	if errRow != nil {
		logEntry.WithError(errRow).WithField("lovIds", reflect.ValueOf(mapByLovID).MapKeys()).Errorf("Fail to run query for lookup version , error: %s", errRow.Error())
		return errRow
	}
	defer rows.Close()
	for rows.Next() {
		var lovID string
		var version string
		rows.Scan(&lovID, &version)
		//theVal :=
		mapByLovID[lovID].Version = version
		//theVal.Version = version
	}
	return
}

//FindSQLDrivenLookupDetails query for custom sql driven lookup. lookup detail will also appended to lookup header
// so there is not needed to append detail to header manually
func FindSQLDrivenLookupDetails(lookupHeaders []coremodel.LookupHeader, db *gorm.DB, logEntry *logrus.Entry) (lookupDetails []coremodel.LookupDetail, err error) {
	if len(lookupHeaders) == 0 {
		return
	}
	dataFinderStatment := []string{}
	mapByLovID := make(map[string]*coremodel.LookupHeader)
	for _, h := range lookupHeaders {
		mapByLovID[h.ID] = &h
		if h.SQLForData == nil {
			continue
		}
		dataFinderStatment = append(dataFinderStatment, *h.SQLForData)

	}
	if len(dataFinderStatment) == 0 {
		logEntry.WithField("lovIds", reflect.ValueOf(mapByLovID).MapKeys()).Warnf("No sql for lookup is found in all lookup data, no sql statement executed")
		return
	}
	finalSQL := strings.Join(dataFinderStatment, " union all ")
	rows, errRow := db.Raw(finalSQL).Rows()
	if errRow != nil {
		logEntry.WithError(errRow).WithField("lovIds", reflect.ValueOf(mapByLovID).MapKeys()).Errorf("Fail to run query for lookup data , error: %s", errRow.Error())
		return nil, errRow
	}
	defer rows.Close()
	for rows.Next() {
		var lkpDetail coremodel.LookupDetail
		db.ScanRows(rows, &lkpDetail)
		lookupDetails = append(lookupDetails, lkpDetail)
		mapByLovID[lkpDetail.LovID].AppendLookupDetail(lkpDetail)
	}
	return
}

//FindSQLDrivenLookupDetailsFiltered find lookup driven bys sql , with code is filtered. detail codes filtered using parameter
func FindSQLDrivenLookupDetailsFiltered(lookupCodeFilters []coredata.LookupRequestData, lookupHeaders []coremodel.LookupHeader, db *gorm.DB, logEntry *logrus.Entry) (lookupDetails []coremodel.LookupDetail, err error) {
	if len(lookupHeaders) == 0 {
		return
	}
	dataFinderStatment := []string{}
	mapByLovID := make(map[string]*coremodel.LookupHeader)
	mapByLovIDFilter := make(map[string]*[]string)
	for _, fltD := range lookupCodeFilters {
		if len(fltD.FilteredCodes) > 0 {
			mapByLovIDFilter[fltD.LovID] = &fltD.FilteredCodes
		}
	}
	for _, h := range lookupHeaders {
		mapByLovID[h.ID] = &h
		if h.SQLForData == nil {
			continue
		}
		if x, ok := mapByLovIDFilter[h.ID]; ok {
			var inVal string
			if h.CodeActualDataType == "string" {
				inVal = "'" + strings.Join(*x, "','") + "'"
			} else {
				inVal = strings.Join(*x, ",")
			}
			finalSQL := strings.Join(strings.Split(h.SQLForDataFiltered, "{{codes}}"), inVal)
			dataFinderStatment = append(dataFinderStatment, finalSQL)
		} else {
			dataFinderStatment = append(dataFinderStatment, *h.SQLForData)
		}

	}
	if len(dataFinderStatment) == 0 {
		logEntry.WithField("lovIds", reflect.ValueOf(mapByLovID).MapKeys()).Warnf("No sql for lookup is found in all lookup data, no sql statement executed")
		return
	}
	finalSQL := strings.Join(dataFinderStatment, " union all ")
	rows, errRow := db.Raw(finalSQL).Rows()
	if errRow != nil {
		logEntry.WithError(errRow).WithField("lovIds", reflect.ValueOf(mapByLovID).MapKeys()).Errorf("Fail to run query for lookup data , error: %s", errRow.Error())
		return nil, errRow
	}
	defer rows.Close()
	for rows.Next() {
		var lkpDetail coremodel.LookupDetail
		db.ScanRows(rows, &lkpDetail)
		lookupDetails = append(lookupDetails, lkpDetail)
		mapByLovID[lkpDetail.LovID].AppendLookupDetail(lkpDetail)
	}
	return

}

//FindSimpleLookupDetails query from tables m_lookup_details by lookup ids
func FindSimpleLookupDetails(lovIDs []string, db *gorm.DB, logEntry *logrus.Entry) (lookupDetails []coremodel.LookupDetail, err error) {
	if len(lovIDs) == 0 {
		return
	}
	logEntry = logEntry.WithField("lovIds", lovIDs).WithField("model", "LookupDetail")
	//colName, err = dao.DefaultDaoManager.GetColumnName("LookupDetail", "LovID")
	whereID := "lov_id in ( ? ) "
	if dbRslt := db.Where(whereID, lovIDs).Order("seq_no").Find(&lookupDetails); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithError(err).Errorf("Fail to query for lookup details. error: %s", err.Error())
		return
	}
	return
}

//FindSimpleLookupDetailFiltered query simple lookup, with filter(if any )
func FindSimpleLookupDetailFiltered(lovIDs []coredata.LookupRequestData, db *gorm.DB, logEntry *logrus.Entry) (lookupDetails []coremodel.LookupDetail, err error) {
	if len(lovIDs) == 0 {
		return
	}
	logEntry = logEntry.WithField("lovIds", lovIDs).WithField("model", "LookupDetail")
	for _, param := range lovIDs {
		if len(param.FilteredCodes) > 0 {
			db = db.Or("(lov_id = ? and detail_code in (?))", param.LovID, param.FilteredCodes)
		} else {
			db = db.Or("lov_id = ?", param.LovID)
		}
	}
	if dbRslt := db.Order("seq_no").Find(&lookupDetails); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithError(err).Errorf("Fail to query for lookup details. error: %s", err.Error())
		return
	}
	return
}

//separateAndIndexLookupHeader memisah antara simple lookup dengan lookup driven by sql. include indexing lookup by lookup id
// returns
// - simpleLookupsID : id of simple lookups. content will be fetch from table m_lookup_details
// - simpleLookups : slice of simple lookups
// - withSQLLookupsID : id of lookups with sql source
func separateAndIndexLookupHeader(lookupHeaders *[]coremodel.LookupHeader) (simpleLookupsID []string, simpleLookups []*coremodel.LookupHeader, withSQLLookupsID []string, withSQLLookups []*coremodel.LookupHeader, indexedLookup map[string]*coremodel.LookupHeader) {
	for _, lk := range *lookupHeaders {
		if lk.FlagUseCustomSQL == "Y" {
			simpleLookups = append(simpleLookups, &lk)
			simpleLookupsID = append(simpleLookupsID, lk.ID)
		} else {
			withSQLLookups = append(withSQLLookups, &lk)
			withSQLLookupsID = append(withSQLLookupsID, lk.ID)
		}
		indexedLookup[lk.ID] = &lk
	}
	return
}
