package service

import (
	"errors"
	"testing"

	"analytics_service/internal/domain"
	"analytics_service/internal/repository"
	"analytics_service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetPaginationInfo(t *testing.T) {
	tableName := "table"
	errRowsCnt := errors.New("errors while getting cont")

	testCases := []struct {
		name                string
		buildPaginationRepo func() repository.PaginationRepo
		paginationParams    domain.PaginationParams
		expectedPagination  domain.Pagination
		expectedErr         error
	}{
		{
			name: "100 records",
			buildPaginationRepo: func() repository.PaginationRepo {
				mockRepo := mocks.NewPaginationRepo(t)

				mockRepo.On("GetRecordsCount", tableName).
					Return(100, nil)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 5, Limit: 10},
			expectedPagination: domain.Pagination{
				Next:          6,
				Previous:      4,
				RecordPerPage: 10,
				CurrentPage:   5,
				TotalPage:     10,
			},
			expectedErr: nil,
		},
		{
			name: "100 records. First page",
			buildPaginationRepo: func() repository.PaginationRepo {
				mockRepo := mocks.NewPaginationRepo(t)

				mockRepo.On("GetRecordsCount", tableName).
					Return(100, nil)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 1, Limit: 10},
			expectedPagination: domain.Pagination{
				Next:          2,
				Previous:      0,
				RecordPerPage: 10,
				CurrentPage:   1,
				TotalPage:     10,
			},
			expectedErr: nil,
		},
		{
			name: "100 records. Last page",
			buildPaginationRepo: func() repository.PaginationRepo {
				mockRepo := mocks.NewPaginationRepo(t)

				mockRepo.On("GetRecordsCount", tableName).
					Return(100, nil)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 10, Limit: 10},
			expectedPagination: domain.Pagination{
				Next:          0,
				Previous:      9,
				RecordPerPage: 10,
				CurrentPage:   10,
				TotalPage:     10,
			},
			expectedErr: nil,
		},
		{
			name: "100 records. Param page is greater than total pages",
			buildPaginationRepo: func() repository.PaginationRepo {
				mockRepo := mocks.NewPaginationRepo(t)

				mockRepo.On("GetRecordsCount", tableName).
					Return(100, nil)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 20, Limit: 10},
			expectedPagination: domain.Pagination{
				Next:          0,
				Previous:      0,
				RecordPerPage: 10,
				CurrentPage:   20,
				TotalPage:     10,
			},
			expectedErr: nil,
		},
		{
			name: "101 records.",
			buildPaginationRepo: func() repository.PaginationRepo {
				mockRepo := mocks.NewPaginationRepo(t)

				mockRepo.On("GetRecordsCount", tableName).
					Return(101, nil)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 5, Limit: 10},
			expectedPagination: domain.Pagination{
				Next:          6,
				Previous:      4,
				RecordPerPage: 10,
				CurrentPage:   5,
				TotalPage:     11,
			},
			expectedErr: nil,
		},
		{
			name: "100 records. Page is zero",
			buildPaginationRepo: func() repository.PaginationRepo {
				mockRepo := mocks.NewPaginationRepo(t)

				mockRepo.On("GetRecordsCount", tableName).
					Return(100, nil)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 0, Limit: 10},
			expectedPagination: domain.Pagination{
				Next:          1,
				Previous:      0,
				RecordPerPage: 10,
				CurrentPage:   0,
				TotalPage:     10,
			},
			expectedErr: nil,
		},
		{
			name: "Error while getting rows count from table",
			buildPaginationRepo: func() repository.PaginationRepo {
				mockRepo := mocks.NewPaginationRepo(t)

				mockRepo.On("GetRecordsCount", tableName).
					Return(0, errRowsCnt)

				return mockRepo
			},
			paginationParams:   domain.PaginationParams{Page: 1, Limit: 10},
			expectedPagination: domain.Pagination{},
			expectedErr:        errRowsCnt,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paginationRepo := tc.buildPaginationRepo()
			paginationService := NewPaginationService(paginationRepo)

			pagination, err := paginationService.GetPaginationInfo(tableName, tc.paginationParams)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedPagination, pagination)
		})
	}
}
