package models

import (
	"fmt"
)

// Offer - информаци о выданном офере обучающемуся
type Offer struct {
	Id     uint64
	UserId uint64
	Grade  uint64
	TeamId uint64
}

func (o *Offer) String() string {
	return fmt.Sprintf(
		"Id: %d, UserId: %d, Grade: %d, TeamId: %d",
		o.Id, o.UserId, o.Grade, o.TeamId,
	)
}
