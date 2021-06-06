package postgresql

import (
	"armani_blog/pkg/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)


type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title, content, user_name string, user_id int ) (int, error) {
	stmt := `INSERT INTO snippets (title,content,created,user_id,user_name)
	VALUES($1, $2 , $3, $4, $5)  RETURNING id;`

	var id int64
	err := m.DB.QueryRow(context.Background(), stmt, title, content, time.Now(), user_id, user_name).Scan(&id)
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (m *SnippetModel) GetById(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, user_id, user_name FROM snippets
	WHERE id = $1`


	row := m.DB.QueryRow(context.Background(),stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.UserId, &s.UserName)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}


func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, user_name FROM snippets ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(context.Background(),stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.Snippet
	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.UserName)
		if err != nil {
			return nil, err }
		snippets = append(snippets, s) }

	if err = rows.Err(); err != nil { return nil, err
	}
	return snippets, nil
}




func (m *SnippetModel) GetByUserId(id int) ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, user_id, user_name FROM snippets
	WHERE user_id = $1`


	rows, err := m.DB.Query(context.Background(),stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.Snippet
	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.UserId, &s.UserName)
		if err != nil {
			return nil, err }
		snippets = append(snippets, s) }

	if err = rows.Err(); err != nil { return nil, err
	}
	return snippets, nil
}


func (m *SnippetModel) UpdateSnippet(id int, user_id int, mySnip *models.Snippet) (*models.Snippet, error) {


	_, err := m.GetById(id)


	if err != nil{
		return nil, err
	}else{
		err = m.deleteSnippet(id)
		if err != nil{
			return nil, err
		}
		stmt := `INSERT INTO snippets (id,title,content,created,user_id,user_name)
		VALUES($1, $2 , $3, $4, $5, $6)`

		_, err := m.DB.Exec(context.Background(), stmt, mySnip.ID , mySnip.Title, mySnip.Content, time.Now(), user_id, mySnip.UserName)
		if err != nil {
			return nil, err
		}
		return mySnip, nil
	}

}


func (m *SnippetModel) deleteSnippet(id int) error {
	stmt := `DELETE FROM snippets WHERE id = $1`

	_, err := m.DB.Exec(context.Background(),stmt, id)

	if err != nil{
		return  err
	}
	return nil
}


