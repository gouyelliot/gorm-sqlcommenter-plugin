package sqlcommenter

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type DummyModel struct {
	ID   uint
	Name string
}

func TestQuery(t *testing.T) {
	tests := []struct {
		name      string
		operation func(db *gorm.DB) *gorm.DB
		want      string
		wantArgs  []driver.Value
	}{
		{
			name: "SimpleTags",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value", "endpoint": "/test/path"})).
					Table("test").Scan(nil)
			},
			want: "SELECT * FROM `test` /*application='value',endpoint='%2Ftest%2Fpath'*/",
		},
		{
			name: "EscapeKey",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"' /application/ '": "value"})).
					Table("test").Scan(nil)
			},
			want: "SELECT * FROM `test` /*%27%20%2Fapplication%2F%20%27='value'*/",
		},
		{
			name: "EscapeValue",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "'  value  '"})).
					Table("test").Scan(nil)
			},
			want: "SELECT * FROM `test` /*application='%27%20%20value%20%20%27'*/",
		},
		{
			name: "MergeClause",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value"})).
					Clauses(NewTags(map[string]string{"hello": "world", "application": "value2"})).
					Table("test").Scan(nil)
			},
			want: "SELECT * FROM `test` /*application='value2',hello='world'*/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			db, _ := gorm.Open(mysql.New(mysql.Config{
				Conn:                      mockDB,
				SkipInitializeWithVersion: true,
			}))
			db.Use(New())

			mock.ExpectQuery(regexp.QuoteMeta(tt.want)).WithArgs(tt.wantArgs...).WillReturnRows(sqlmock.NewRows([]string{}))
			if tt.operation != nil {
				db = tt.operation(db)
			}
			if db.Error != nil {
				t.Errorf(db.Error.Error())
			}
		})
	}
}

func TestCreate(t *testing.T) {

	dummyObject := DummyModel{
		ID:   15,
		Name: "my-name",
	}

	tests := []struct {
		name      string
		operation func(db *gorm.DB) *gorm.DB
		want      string
		wantArgs  []driver.Value
	}{
		{
			name: "SimpleTags",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value", "endpoint": "/test/path"})).
					Table("test").Create(dummyObject)
			},
			want: "INSERT INTO `test` (`name`,`id`) VALUES (?,?) /*application='value',endpoint='%2Ftest%2Fpath'*/",
			wantArgs: []driver.Value{
				"my-name",
				15,
			},
		},
		{
			name: "EscapeKey",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"' /application/ '": "value"})).
					Table("test").Create(dummyObject)
			},
			want: "INSERT INTO `test` (`name`,`id`) VALUES (?,?) /*%27%20%2Fapplication%2F%20%27='value'*/",
			wantArgs: []driver.Value{
				"my-name",
				15,
			},
		},
		{
			name: "EscapeValue",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "'  value  '"})).
					Table("test").Create(dummyObject)
			},
			want: "INSERT INTO `test` (`name`,`id`) VALUES (?,?) /*application='%27%20%20value%20%20%27'*/",
			wantArgs: []driver.Value{
				"my-name",
				15,
			},
		},
		{
			name: "MergeClause",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value"})).
					Clauses(NewTags(map[string]string{"hello": "world", "application": "value2"})).
					Table("test").Create(dummyObject)
			},
			want: "INSERT INTO `test` (`name`,`id`) VALUES (?,?) /*application='value2',hello='world'*/",
			wantArgs: []driver.Value{
				"my-name",
				15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			db, _ := gorm.Open(mysql.New(mysql.Config{
				Conn:                      mockDB,
				SkipInitializeWithVersion: true,
			}))
			db.Use(New())

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(tt.want)).WithArgs(tt.wantArgs...).WillReturnResult(driver.ResultNoRows)
			mock.ExpectCommit()
			if tt.operation != nil {
				db = tt.operation(db)
			}
			if db.Error != nil {
				t.Errorf(db.Error.Error())
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name      string
		operation func(db *gorm.DB) *gorm.DB
		want      string
		wantArgs  []driver.Value
	}{
		{
			name: "SimpleTags",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value", "endpoint": "/test/path"})).
					Table("test").Where("id = ?", 15).Update("name", "test-name")
			},
			want: "UPDATE `test` SET `name`=? WHERE id = ? /*application='value',endpoint='%2Ftest%2Fpath'*/",
			wantArgs: []driver.Value{
				"test-name",
				15,
			},
		},
		{
			name: "EscapeKey",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"' /application/ '": "value"})).
					Table("test").Where("id = ?", 15).Update("name", "test-name")
			},
			want: "UPDATE `test` SET `name`=? WHERE id = ? /*%27%20%2Fapplication%2F%20%27='value'*/",
			wantArgs: []driver.Value{
				"test-name",
				15,
			},
		},
		{
			name: "EscapeValue",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "'  value  '"})).
					Table("test").Where("id = ?", 15).Update("name", "test-name")
			},
			want: "UPDATE `test` SET `name`=? WHERE id = ? /*application='%27%20%20value%20%20%27'*/",
			wantArgs: []driver.Value{
				"test-name",
				15,
			},
		},
		{
			name: "MergeClause",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value"})).
					Clauses(NewTags(map[string]string{"hello": "world", "application": "value2"})).
					Table("test").Where("id = ?", 15).Update("name", "test-name")
			},
			want: "UPDATE `test` SET `name`=? WHERE id = ? /*application='value2',hello='world'*/",
			wantArgs: []driver.Value{
				"test-name",
				15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			db, _ := gorm.Open(mysql.New(mysql.Config{
				Conn:                      mockDB,
				SkipInitializeWithVersion: true,
			}))
			db.Use(New())

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(tt.want)).WithArgs(tt.wantArgs...).WillReturnResult(driver.ResultNoRows)
			mock.ExpectCommit()
			if tt.operation != nil {
				db = tt.operation(db)
			}
			if db.Error != nil {
				t.Errorf(db.Error.Error())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		operation func(db *gorm.DB) *gorm.DB
		want      string
		wantArgs  []driver.Value
	}{
		{
			name: "SimpleTags",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value", "endpoint": "/test/path"})).
					Table("test").Where("id = ?", 15).Delete(nil)
			},
			want: "DELETE FROM `test` WHERE id = ? /*application='value',endpoint='%2Ftest%2Fpath'*/",
			wantArgs: []driver.Value{
				15,
			},
		},
		{
			name: "EscapeKey",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"' /application/ '": "value"})).
					Table("test").Where("id = ?", 15).Delete(nil)
			},
			want: "DELETE FROM `test` WHERE id = ? /*%27%20%2Fapplication%2F%20%27='value'*/",
			wantArgs: []driver.Value{
				15,
			},
		},
		{
			name: "EscapeValue",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "'  value  '"})).
					Table("test").Where("id = ?", 15).Delete(nil)
			},
			want: "DELETE FROM `test` WHERE id = ? /*application='%27%20%20value%20%20%27'*/",
			wantArgs: []driver.Value{
				15,
			},
		},
		{
			name: "MergeClause",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.
					Clauses(NewTags(map[string]string{"application": "value"})).
					Clauses(NewTags(map[string]string{"hello": "world", "application": "value2"})).
					Table("test").Where("id = ?", 15).Delete(nil)
			},
			want: "DELETE FROM `test` WHERE id = ? /*application='value2',hello='world'*/",
			wantArgs: []driver.Value{
				15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			db, _ := gorm.Open(mysql.New(mysql.Config{
				Conn:                      mockDB,
				SkipInitializeWithVersion: true,
			}))
			db.Use(New())

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(tt.want)).WithArgs(tt.wantArgs...).WillReturnResult(driver.ResultNoRows)
			mock.ExpectCommit()
			if tt.operation != nil {
				db = tt.operation(db)
			}
			if db.Error != nil {
				t.Errorf(db.Error.Error())
			}
		})
	}
}
