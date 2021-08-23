package repo_test

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-contact-api/internal/models"
	"github.com/ozoncp/ocp-contact-api/internal/repo"
)

var _ = Describe("Repo", func() {
	const tableName = "contact"

	var (
		db       *sql.DB
		sqlxDB   *sqlx.DB
		mock     sqlmock.Sqlmock

		ctx      context.Context
		r        repo.Repo
		contacts []models.Contact
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")

		ctx = context.Background()
		r = repo.NewRepo(sqlxDB)

		contacts = []models.Contact {
			0: {1 , 1, 41, "one"},
			1: {2 , 2, 42, "two"},
			2: {3 , 3, 43, "three"},
			3: {4 , 4, 44, "four"},
		}
	})

	AfterEach(func() {
		mock.ExpectClose()
		err := db.Close()
		Expect(err).Should(BeNil())
	})

	Context("AddContacts", func() {
		BeforeEach(func() {
			mock.ExpectExec("INSERT INTO " + tableName).
				WithArgs(
					contacts[0].UserId, contacts[0].Type, contacts[0].Text,
					contacts[1].UserId, contacts[1].Type, contacts[1].Text,
					contacts[2].UserId, contacts[2].Type, contacts[2].Text,
					contacts[3].UserId, contacts[3].Type, contacts[3].Text,
				).
				WillReturnResult(sqlmock.NewResult(4, 4))
		})

		It("add multiple contacts", func() {
			err := r.AddContacts(ctx, contacts)
			Expect(err).Should(BeNil())
		})
	})

	Context("DescribeContact", func() {
		BeforeEach(func() {
			rows := sqlmock.NewRows([]string{"id", "user_id", "type", "text"}).AddRow(
				contacts[3].Id,
				contacts[3].UserId,
				contacts[3].Type,
				contacts[3].Text)
			mock.ExpectQuery(
				"SELECT id, user_id, type, text FROM " + tableName + " WHERE").
				WithArgs(contacts[3].Id).
				WillReturnRows(rows)
		})

		It("return contact", func() {
			contact, err := r.DescribeContact(ctx, contacts[3].Id)
			Expect(err).Should(BeNil())
			Expect(*contact).Should(BeEquivalentTo(contacts[3]))
		})
	})

	Context("CreateContact", func() {
		BeforeEach(func() {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery("INSERT INTO " + tableName).
				WithArgs(
					contacts[0].UserId,
					contacts[0].Type,
					contacts[0].Text,
				).WillReturnRows(rows)
		})

		It("return created contact id", func() {
			contact := models.Contact{Id: 1, UserId: 1, Type: 41, Text: "one"}
			id, err := r.CreateContact(ctx, contact)
			Expect(err).Should(BeNil())
			Expect(id).Should(BeEquivalentTo(1))
		})
	})

	Context("RemoveContact", func() {
		When("id exists in the db", func() {
			BeforeEach(func() {
				query := mock.ExpectExec("DELETE FROM " + tableName + " WHERE")
				query.WithArgs(contacts[3].Id)
				query.WillReturnResult(sqlmock.NewResult(1, 1))
			})
			It("remove contact and return nil", func() {
				err := r.RemoveContact(ctx, contacts[3].Id)
				Expect(err).Should(BeNil())
			})
		})
		When("id not exists in the db", func() {
			BeforeEach(func() {
				query := mock.ExpectExec("DELETE FROM " + tableName + " WHERE")
				query.WithArgs(999)
				query.WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("test error")))
			})
			It("remove contact and return nil", func() {
				err := r.RemoveContact(ctx, 999)
				Expect(err).ShouldNot(BeNil())
			})
		})
	})

	Context("ListContacts", func() {
		const limit uint64 = 4
		const offset uint64 = 0

		BeforeEach(func() {
			rows := sqlmock.NewRows([]string{"id", "user_id", "type", "text"}).
				AddRow(contacts[0].Id, contacts[0].UserId, contacts[0].Type, contacts[0].Text).
				AddRow(contacts[1].Id, contacts[1].UserId, contacts[1].Type, contacts[1].Text).
				AddRow(contacts[2].Id, contacts[2].UserId, contacts[2].Type, contacts[2].Text).
				AddRow(contacts[3].Id, contacts[3].UserId, contacts[3].Type, contacts[3].Text)

			query := fmt.Sprintf("SELECT id, user_id, type, text FROM %s LIMIT %v OFFSET %v",
				tableName, limit, offset)
			mock.ExpectQuery(query).WillReturnRows(rows)
		})

		It("return list contacts without error", func() {
			actualContacts, err := r.ListContacts(ctx, limit, offset)
			Expect(err).Should(BeNil())
			Expect(actualContacts).Should(BeEquivalentTo(contacts))
		})
	})
})
