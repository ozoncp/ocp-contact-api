package models

import "fmt"

type Contact struct {
	Id uint64     `db:"id"`
	UserId uint64 `db:"user_id"`
	Type uint64   `db:"type"`
	Text string   `db:"text"`
}

func (c *Contact) String() string { return fmt.Sprintf("Contact(id=%v)", c.Id) }


