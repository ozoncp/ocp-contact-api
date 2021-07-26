package models

import "fmt"

type Contact struct {
	Id uint64
	UserId uint64
	Type uint64
	Text string
}

func (c *Contact) String() string { return fmt.Sprintf("Contact(id=%v)", c.Id) }


