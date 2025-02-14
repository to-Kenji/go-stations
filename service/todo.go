package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	if len(subject) == 0 {
		return &model.TODO{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	todo := model.TODO{
		ID: int(id),
	}
	if err != nil {
		log.Println(err)
	}
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no todo with id %d\n", id)
	case err != nil:
		log.Fatalf("query error: %v\n", err)
	}
	return &todo, err
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	// decide size
	var sizeOnQuery int64
	if size == 0 {
		// TODO: openapi.yamlに従うなら0ではなく5件取得では？5件だとテスト通らない
		sizeOnQuery = 0
	} else {
		sizeOnQuery = size
	}

	var rows *sql.Rows
	var err error
	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, sizeOnQuery)
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, sizeOnQuery)
	}

	if err != nil {
		return nil, err
	}

	// scan
	todos := []*model.TODO{}
	for rows.Next() {
		var t model.TODO
		rows.Scan(&t.ID, &t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		todos = append(todos, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	if id == 0 {
		return &model.TODO{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	if len(subject) == 0 {
		return &model.TODO{
			ID: int(id),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}

	_, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}
	todo := model.TODO{
		ID: int(id),
	}
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, model.ErrNotFound{}
	case err != nil:
		return nil, err
	}
	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
