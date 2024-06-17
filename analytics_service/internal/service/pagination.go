package service

import (
	"analytics_service/internal/domain"
	"analytics_service/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name PaginationService
type PaginationService interface {
	GetPaginationInfo(table string, paginationParams domain.PaginationParams) (domain.Pagination, error)
}

type paginationService struct {
	paginationRepo repository.PaginationRepo
}

func NewPaginationService(
	paginationRepo repository.PaginationRepo,
) PaginationService {
	return &paginationService{
		paginationRepo: paginationRepo,
	}
}

func (s *paginationService) GetPaginationInfo(
	table string,
	paginationParams domain.PaginationParams,
) (domain.Pagination, error) {
	recordsCount, err := s.paginationRepo.GetRecordsCount(table)
	if err != nil {
		return domain.Pagination{}, err
	}

	var pagination domain.Pagination

	// Set current/record per page data
	pagination.CurrentPage = paginationParams.Page
	pagination.RecordPerPage = paginationParams.Limit

	totalPages := recordsCount / paginationParams.Limit
	// Calculate Total Page
	remainder := recordsCount % paginationParams.Limit
	if remainder == 0 {
		pagination.TotalPage = totalPages
	} else {
		pagination.TotalPage = totalPages + 1
	}

	// Calculate the Next/Previous Page
	if paginationParams.Page <= 0 {
		pagination.Next = paginationParams.Page + 1
	} else if paginationParams.Page < pagination.TotalPage {
		pagination.Previous = paginationParams.Page - 1
		pagination.Next = paginationParams.Page + 1
	} else if paginationParams.Page == pagination.TotalPage {
		pagination.Previous = paginationParams.Page - 1
		pagination.Next = 0
	}

	return pagination, nil
}
