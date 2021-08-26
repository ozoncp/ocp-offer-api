package models

import (
	"fmt"
)

// Offer - информаци о выданном офере обучающемуся.
type Offer struct {
	IsDeleted bool   `db:"is_deleted"`
	ID        uint64 `db:"id"`
	UserID    uint64 `db:"user_id"`
	TeamID    uint64 `db:"team_id"`
	Grade     uint64 `db:"grade"`
}

func (o *Offer) String() string {
	return fmt.Sprintf(
		"ID: %d, UserID: %d, Grade: %d, TeamID: %d, IsDeleted: %v",
		o.ID, o.UserID, o.Grade, o.TeamID, o.IsDeleted,
	)
}
