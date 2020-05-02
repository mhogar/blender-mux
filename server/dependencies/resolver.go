package dependencies

import "github.com/blendermux/server/database"

type DependencyResolver struct {
	Database
	UserCRUD
}

func InitDependencyResolver() *DependencyResolver {
	db := database.InitInMemoryDB()

	return &DependencyResolver{
		db,
		db,
	}
}
