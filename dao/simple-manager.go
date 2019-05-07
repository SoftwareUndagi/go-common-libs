package dao

import (
	"github.com/SoftwareUndagi/go-common-libs/dao/query"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//SimpleDaoManager dao manager with strip down param. db an log entry passes as member variable on struct
type SimpleDaoManager interface {

	//FindByID find data by id
	FindByID(modelName string, IDstring string) (result interface{}, err error)
	//GenerateDaoWithWhere generate where. so reference is ready to use for query ( i.e. for count, list etc)
	GenerateDaoWithWhere(modelName string, q query.Q) (dbWithWhere *gorm.DB, err error)
	//List query for list of data
	List(modelName string, q query.Q, pageSize int32, page int32) (listData interface{}, count int64, err error)
}

type simpleDaoManagerImpl struct {
	//daoManager dao manager actual
	daoManager Manager
	//db gorm datbase reference
	db *gorm.DB
	//logEntry logrus log entry
	logEntry *logrus.Entry
}

//FindByID find data by id
func (p *simpleDaoManagerImpl) FindByID(modelName string, IDstring string) (result interface{}, err error) {
	return (p.daoManager).FindByID(modelName, IDstring, p.db, p.logEntry)
}

//GenerateDaoWithWhere generate where. so reference is ready to use for query ( i.e. for count, list etc)
func (p *simpleDaoManagerImpl) GenerateDaoWithWhere(modelName string, q query.Q) (dbWithWhere *gorm.DB, err error) {
	return (p.daoManager).GenerateDaoWithWhere(modelName, q, p.db, p.logEntry)
}

//List query for list of data
func (p *simpleDaoManagerImpl) List(modelName string, q query.Q, pageSize int32, page int32) (listData interface{}, count int64, err error) {
	return (p.daoManager).List(modelName, q, pageSize, page, p.db, p.logEntry)
}

//NewSimpleManager generate new simple dao manager
func NewSimpleManager(daoManager Manager, DB *gorm.DB, logEntry *logrus.Entry) SimpleDaoManager {
	return &simpleDaoManagerImpl{daoManager: daoManager, db: DB, logEntry: logEntry}
}
