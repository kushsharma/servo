package backup

import "github.com/kushsharma/servo/sshtunnel"

type DBService struct {
	ssh *sshtunnel.Client
}

func (svc *DBService) Prepare() error {

	return nil
}

func (svc *DBService) Migrate() error {

	return nil
}

func NewDBService(client *sshtunnel.Client) *DBService {
	db := new(DBService)
	db.ssh = client
	return db
}
