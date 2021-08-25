package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ozoncp/ocp-contact-api/internal/models"
)

const (
	contactTable string = "contact"
)

var column = struct {
	Id string
	UserId string
	Type string
	Text string
}{
	Id: "id",
	UserId: "user_id",
	Type: "type",
	Text: "text",
}

type Repo interface {
	AddContacts(ctx context.Context, contacts []models.Contact) error
	ListContacts(ctx context.Context, limit, offset uint64) ([]models.Contact, error)
	DescribeContact(ctx context.Context, contactId uint64) (*models.Contact, error)
	CreateContact(ctx context.Context, contact models.Contact) (uint64, error)
	RemoveContact(ctx context.Context, contactId uint64) error
	UpdateContact(ctx context.Context, contact models.Contact) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) AddContacts(ctx context.Context, contacts []models.Contact) error {
	query := squirrel.
		Insert(contactTable).
		Columns(column.UserId, column.Type, column.Text).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	for _, contact := range contacts {
		query = query.Values(contact.UserId, contact.Type, contact.Text)
	}

	sqlResult, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return errors.New("insert contact failed")
	}

	return nil
}

func (r *repo) ListContacts(ctx context.Context, limit, offset uint64) ([]models.Contact, error) {
	query := squirrel.Select(column.Id, column.UserId, column.Type, column.Text).
		From(contactTable).
		RunWith(r.db).
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(
			&contact.Id,
			&contact.UserId,
			&contact.Type,
			&contact.Text); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (r *repo) DescribeContact(ctx context.Context, contactId uint64) (*models.Contact, error) {
	query := squirrel.Select(column.Id, column.UserId, column.Type, column.Text).
		From(contactTable).
		Where(squirrel.Eq{column.Id: contactId}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	var contact models.Contact

	if err := query.QueryRowContext(ctx).
		Scan(&contact.Id,
			&contact.UserId,
			&contact.Type,
			&contact.Text); err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *repo) CreateContact(ctx context.Context, contact models.Contact) (uint64, error) {
	query := squirrel.Insert(contactTable).
		Columns(column.UserId, column.Type, column.Text).
		Values(contact.UserId, contact.Type, contact.Text).
		Suffix(`RETURNING "id"`).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	err := query.QueryRowContext(ctx).Scan(&contact.Id)
	if err != nil {
		return 0, err
	}

	return contact.Id, nil
}

func (r *repo) RemoveContact(ctx context.Context, contactId uint64) error {
	query := squirrel.Delete(contactTable).
		Where(squirrel.Eq{column.Id: contactId}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("remove contact failed: contact with id %v not found", contactId)
	}

	return nil
}

func (r *repo) UpdateContact(ctx context.Context, contact models.Contact) error {
	query := squirrel.Update(contactTable).
		Set("user_id", contact.UserId).
		Set("type", contact.Type).
		Set("text", contact.Text).
		Where(squirrel.Eq{"id": contact.Id}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	exec, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("update contact failed: contact with id %v not found", contact.Id)
	}

	return nil
}