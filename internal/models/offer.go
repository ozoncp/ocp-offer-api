package models

import (
	"fmt"
)

// Offer - информаци о выданном офере обучающемуся
type Offer struct {
	Id     uint64 `db:"id"`
	UserId uint64 `db:"user_id"`
	TeamId uint64 `db:"team_id"`
	Grade  uint64 `db:"grade"`
}

func (o *Offer) String() string {
	return fmt.Sprintf(
		"Id: %d, UserId: %d, Grade: %d, TeamId: %d",
		o.Id, o.UserId, o.Grade, o.TeamId,
	)
}
