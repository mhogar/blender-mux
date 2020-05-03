package dependencies

import "github.com/blendermux/server/database/in_memory_db"

type DependencyResolver struct {
	Database
}

func InitDependencyResolver() *DependencyResolver {
	return &DependencyResolver{
		inmemorydb.InitInMemoryDB(),
	}
}
