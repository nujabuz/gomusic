package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nujabuz/music-cloud/internal/domain"
	repo "github.com/nujabuz/music-cloud/internal/repository/postgres"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB sqlmock.Sqlmock
	db     *gorm.DB
	repo   domain.UserRepository
}

// SetupTest runs before each test to initialize a fresh mock DB and repository
func (s *UserRepositoryTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mockDB, err = sqlmock.New()
	s.Require().NoError(err)

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		DriverName:           "postgres",
		WithoutReturning:     false,
		PreferSimpleProtocol: true,
	})

	s.db, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	s.Require().NoError(err)

	s.repo = repo.NewUserRepository(s.db)
}

// TearDownTest ensures all sql expectations were met
func (s *UserRepositoryTestSuite) TearDownTest() {
	s.NoError(s.mockDB.ExpectationsWereMet())
}

// TestCreate_Success tests successful user creation
func (s *UserRepositoryTestSuite) TestCreate_Success() {
	user := &domain.User{
		Username:  "johndoe",
		Password:  "hashedsecret",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := regexp.QuoteMeta(`INSERT INTO "users" ("username","password","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)
	s.mockDB.ExpectQuery(query).
		WithArgs(user.Username, user.Password, user.CreatedAt, user.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	createdUser, err := s.repo.Create(context.Background(), user)

	s.NoError(err)
	s.NotNil(createdUser)
	s.Equal(int64(1), createdUser.ID)
	s.Equal(user.Username, createdUser.Username)
}

// TestCreate_DuplicateUsername tests user creation when username already exists
func (s *UserRepositoryTestSuite) TestCreate_DuplicateUsername() {
	user := &domain.User{
		Username:  "existinguser",
		Password:  "hashedsecret",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := regexp.QuoteMeta(`INSERT INTO "users" ("username","password","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)
	s.mockDB.ExpectQuery(query).
		WithArgs(user.Username, user.Password, user.CreatedAt, user.UpdatedAt).
		WillReturnError(gorm.ErrDuplicatedKey)

	createdUser, err := s.repo.Create(context.Background(), user)

	s.ErrorIs(err, domain.ErrUserAlreadyExists)
	s.Nil(createdUser)
}

// TestFindByUsername_Success tests fetching an existing user by username
func (s *UserRepositoryTestSuite) TestFindByUsername_Success() {
	username := "johndoe"
	expectedTime := time.Now()

	query := regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 LIMIT $2`)
	rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
		AddRow(1, username, "hashedsecret", expectedTime, expectedTime)

	s.mockDB.ExpectQuery(query).
		WithArgs(username, 1).
		WillReturnRows(rows)

	result, err := s.repo.FindByUsername(context.Background(), username)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(int64(1), result.ID)
	s.Equal(username, result.Username)
}

// TestFindByUsername_NotFound tests fetching a non-existing user
func (s *UserRepositoryTestSuite) TestFindByUsername_NotFound() {
	username := "nonexistent"

	query := regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 LIMIT $2`)
	s.mockDB.ExpectQuery(query).
		WithArgs(username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := s.repo.FindByUsername(context.Background(), username)

	s.ErrorIs(err, domain.ErrUserNotFound)
	s.Nil(result)
}

// TestUpdate_Success tests successfully updating a user by username
func (s *UserRepositoryTestSuite) TestUpdate_Success() {
	user := &domain.User{
		Username: "johndoe",
		Password: "newhashedsecret",
	}

	query := regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"updated_at"=$3 WHERE username = $4`)
	s.mockDB.ExpectExec(query).
		WithArgs(user.Username, user.Password, sqlmock.AnyArg(), user.Username).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedUser, err := s.repo.Update(context.Background(), user)

	s.NoError(err)
	s.NotNil(updatedUser)
	s.Equal(user.Username, updatedUser.Username)
}

// TestUpdate_NotFound tests updating a user that does not exist
func (s *UserRepositoryTestSuite) TestUpdate_NotFound() {
	user := &domain.User{
		Username: "nonexistent",
		Password: "secret",
	}

	query := regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"updated_at"=$3 WHERE username = $4`)
	s.mockDB.ExpectExec(query).
		WithArgs(user.Username, user.Password, sqlmock.AnyArg(), user.Username).
		WillReturnResult(sqlmock.NewResult(0, 0))

	updatedUser, err := s.repo.Update(context.Background(), user)

	s.ErrorIs(err, domain.ErrUserNotFound)
	s.Nil(updatedUser)
}

// TestUpdate_DatabaseError tests update failure due to DB errors
func (s *UserRepositoryTestSuite) TestUpdate_DatabaseError() {
	user := &domain.User{
		Username: "johndoe",
		Password: "secret",
	}

	dbErr := errors.New("connection failed")
	query := regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"updated_at"=$3 WHERE username = $4`)
	s.mockDB.ExpectExec(query).
		WithArgs(user.Username, user.Password, sqlmock.AnyArg(), user.Username).
		WillReturnError(dbErr)

	updatedUser, err := s.repo.Update(context.Background(), user)

	s.ErrorIs(err, dbErr)
	s.Nil(updatedUser)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
