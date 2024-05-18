package report

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/elastic"
)

type ElasticService struct {
	elastic    elastic.ElasticService
}

func NewReportService(elastic elastic.ElasticService) *ElasticService {
	return &ElasticService{
		elastic:    elastic,
	}
}

