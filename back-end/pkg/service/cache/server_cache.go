package cache

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
)

type ServerCache interface {
	Set(key string, value *model.Server)
	Get(key string) *model.Server
	Delete(key string) error
	GetMultiRequest(key string) []model.Server
	SetMultiRequest(key string, value []model.Server)  
	ConstructCacheKey(perPage, offset int, status, field, order string) string
}
