package models

import (
	"fmt"
)

// Offer - информаци о выданном офере обучающемуся.
type Offer struct {
	ID     uint64 `db:"id"`
	UserID uint64 `db:"user_id"`
	TeamID uint64 `db:"team_id"`
	Grade  uint64 `db:"grade"`
}

func (o *Offer) String() string {
	return fmt.Sprintf(
		"Id: %d, UserId: %d, Grade: %d, TeamId: %d",
		o.ID, o.UserID, o.Grade, o.TeamID,
	)
}
