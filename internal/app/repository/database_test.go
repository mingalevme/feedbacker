package repository

import (
	stdDatabaseSql "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabaseAddNoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()
	f1 := model.MakeFeedback()
	data := AddFeedbackData{f1}
	sql := `INSERT INTO feedback \(service, edition, text, extra\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, created_at`
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, f1.CreatedAt)
	mock.ExpectQuery(sql).WithArgs(f1.Service, f1.Edition, f1.Text, &extraColumnHolder{
		Customer: f1.Customer,
		Context:  f1.Context,
	}).WillReturnRows(rows)
	r := NewDatabaseFeedbackRepository(db, nil)
	f2, err := r.Add(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, f2.ID)
}

func TestDatabaseGetNoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()
	f1 := model.MakeFeedback()
	sql := `SELECT service, edition, text, extra, created_at, updated_at FROM feedback WHERE id = \$1`
	rows := sqlmock.NewRows([]string{"service", "edition", "text", "extra", "created_at", "updated_at"})
	rows.AddRow(f1.Service, f1.Edition, f1.Text, &extraColumnHolder{
		Customer: f1.Customer,
		Context:  f1.Context,
	}, f1.CreatedAt, f1.UpdatedAt)
	mock.ExpectQuery(sql).WithArgs(f1.ID).WillReturnRows(rows)
	r := NewDatabaseFeedbackRepository(db, nil)
	f2, err := r.GetById(f1.ID)
	assert.NoError(t, err)
	assert.Equal(t, f1.ID, f2.ID)
	assert.Equal(t, f1.Service, f2.Service)
	assert.Equal(t, f1.Edition, f2.Edition)
	assert.Equal(t, f1.Text, f2.Text)
	assert.Equal(t, f1.Customer, f2.Customer)
	assert.Equal(t, f1.Context, f2.Context)
	assert.Equal(t, f1.CreatedAt.Unix(), f2.CreatedAt.Unix())
	assert.Equal(t, f1.UpdatedAt.Unix(), f2.UpdatedAt.Unix())
}

func TestDatabaseGetErrNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()
	f1 := model.MakeFeedback()
	sql := `SELECT service, edition, text, extra, created_at, updated_at FROM feedback WHERE id = \$1`
	mock.ExpectQuery(sql).WithArgs(f1.ID).WillReturnError(stdDatabaseSql.ErrNoRows)
	r := NewDatabaseFeedbackRepository(db, nil)
	_, err = r.GetById(f1.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}
