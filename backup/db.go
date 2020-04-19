package backup

import (
	"github.com/kushsharma/servo/tunnel"
)

type DBService struct {
	tnl tunnel.Executioner
}

func (svc *DBService) Prepare() error {

	return nil
}

func (svc *DBService) Migrate() error {

	return nil
}

func NewDBService(tnl tunnel.Executioner) *DBService {
	db := new(DBService)
	db.tnl = tnl
	return db
}
