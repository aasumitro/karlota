package sql_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	sqlRepo "karlota.aasumitro.id/internal/repository/sql"
)

type test struct {
	ID   uint   `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (t *test) ToResponse() *test {
	return &test{
		ID:   t.ID,
		Name: t.Name,
	}
}

func (t *test) FromResponse(tm *test) interface{} {
	return &test{
		ID:   tm.ID,
		Name: tm.Name,
	}
}

type sqlRepositoryTestSuite struct {
	suite.Suite
	db      *sql.DB
	gorm    *gorm.DB
	mock    sqlmock.Sqlmock
	qryRepo sqlRepo.ISQLQueryRepository[test, test]
	cmdRepo sqlRepo.ISQLCommandRepository[test, test]
}

func (s *sqlRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))
	require.NoError(s.T(), err)
	s.gorm, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      s.db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(s.T(), err)
	s.qryRepo = sqlRepo.NewSQLQueryRepository[*test, test](s.gorm)
	s.cmdRepo = sqlRepo.NewSQLCommandRepository[*test, test](s.gorm)
}

func (s *sqlRepositoryTestSuite) AfterTest(_, _ string) {
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *sqlRepositoryTestSuite) Test_Insert_ShouldSuccess() {
	data := &test{Name: "lorem"}
	s.mock.ExpectExec("INSERT").
		WithArgs(data.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.NoError(s.cmdRepo.Insert(context.TODO(), data))
}

func (s *sqlRepositoryTestSuite) Test_Insert_ShouldError() {
	data := &test{Name: "lorem"}
	err := errors.New("lorem")
	s.mock.ExpectExec("INSERT").
		WithArgs(data.Name).
		WillReturnError(err)
	s.Errorf(s.cmdRepo.Insert(context.TODO(), data), err.Error())
}

func (s *sqlRepositoryTestSuite) Test_Delete_ShouldSuccess() {
	//data := &test{ID: 1}
	//s.mock.ExpectExec("DELETE").
	//	WithArgs(data.ID).
	//	WillReturnResult(sqlmock.NewResult(1, 1))
	//s.NoError(s.cmdRepo.Delete(context.TODO(), data))
}

func (s *sqlRepositoryTestSuite) Test_Delete_ShouldError() {
	//data := &test{ID: 1}
	//err := errors.New("lorem")
	//s.mock.ExpectExec("DELETE").
	//	WithArgs(data.ID).
	//	WillReturnError(err)
	//s.Errorf(s.cmdRepo.Delete(context.TODO(), data), err.Error())
}

func (s *sqlRepositoryTestSuite) Test_DeleteByID_ShouldSuccess() {
	data := &test{ID: 1}
	s.mock.ExpectExec("DELETE").
		WithArgs(data.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.NoError(s.cmdRepo.DeleteByID(context.TODO(), data.ID))
}

func (s *sqlRepositoryTestSuite) Test_DeleteByID_ShouldError() {
	data := &test{ID: 1}
	err := errors.New("lorem")
	s.mock.ExpectExec("DELETE").
		WithArgs(data.ID).
		WillReturnError(err)
	s.Errorf(s.cmdRepo.DeleteByID(context.TODO(), data.ID), err.Error())
}

func (s *sqlRepositoryTestSuite) Test_Update_ShouldSuccess() {
	data := &test{ID: 1, Name: "ipsum"}
	s.mock.ExpectExec("UPDATE").
		WithArgs(data.Name, data.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.NoError(s.cmdRepo.Update(context.TODO(), data))
}

func (s *sqlRepositoryTestSuite) Test_Update_ShouldError() {
	data := &test{ID: 1, Name: "ipsum"}
	err := errors.New("lorem")
	s.mock.ExpectExec("UPDATE").
		WithArgs(data.Name, data.ID).
		WillReturnError(err)
	s.Errorf(s.cmdRepo.Update(context.TODO(), data), err.Error())
}

func (s *sqlRepositoryTestSuite) Test_Count_ShouldSuccess() {
	data := s.mock.
		NewRows([]string{"count"}).
		AddRow(2)
	s.mock.ExpectQuery("SELECT").
		WillReturnError(nil).
		WillReturnRows(data)
	total, err := s.qryRepo.Count(context.TODO(),
		sqlRepo.And(sqlRepo.GreaterThan("id", 1)),
		sqlRepo.Or(sqlRepo.GreaterOrEqual("id", 1)),
		sqlRepo.Not(sqlRepo.Equal("id", 1)),
		sqlRepo.LessThan("id", 1),
		sqlRepo.LessOrEqual("id", 1),
		sqlRepo.In("id", []int{1, 2}),
		sqlRepo.IsNull("id"),
	)
	s.Equal(int64(2), total)
	s.NoError(err)
}

func (s *sqlRepositoryTestSuite) Test_Count_ShouldError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("lorem"))
	total, err := s.qryRepo.Count(context.TODO())
	s.Equal(int64(0), total)
	s.Errorf(err, "lorem")
}

func (s *sqlRepositoryTestSuite) Test_FindAll_ShouldSuccess() {
	data := s.mock.
		NewRows([]string{"id", "name"}).
		AddRow(1, "test").
		AddRow(2, "test 2")
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(data)
	items, err := s.qryRepo.FindAll(context.TODO())
	s.Equal(2, len(items))
	s.NoError(err)
}

func (s *sqlRepositoryTestSuite) Test_FindAll_ShouldError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("lorem"))
	items, err := s.qryRepo.FindAll(context.TODO())
	s.Equal(0, len(items))
	s.Errorf(err, "lorem")
}

func (s *sqlRepositoryTestSuite) Test_FindWithLimit_ShouldSuccess() {
	data := s.mock.
		NewRows([]string{"id", "name"}).
		AddRow(1, "test").
		AddRow(2, "test 2")
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(data)
	items, err := s.qryRepo.FindWithLimit(context.TODO(), 2, 0)
	s.Equal(2, len(items))
	s.NoError(err)
}

func (s *sqlRepositoryTestSuite) Test_FindWithLimit_ShouldError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("lorem"))
	items, err := s.qryRepo.FindWithLimit(context.TODO(), 2, 0)
	s.Equal(0, len(items))
	s.Errorf(err, "lorem")
}

func (s *sqlRepositoryTestSuite) Test_FindByID_ShouldSuccess() {
	data := s.mock.
		NewRows([]string{"id", "name"}).
		AddRow(1, "test")
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(data)
	item, err := s.qryRepo.FindByID(context.TODO(), 1)
	s.NotNil(item)
	s.Equal(item.ID, uint(1))
	s.NoError(err)
}

func (s *sqlRepositoryTestSuite) Test_FindByID_ShouldError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("lorem"))
	item, err := s.qryRepo.FindByID(context.TODO(), 1)
	s.Equal(item, &test{})
	s.Errorf(err, "lorem")
}

func TestSQLRepository(t *testing.T) {
	suite.Run(t, new(sqlRepositoryTestSuite))
}
