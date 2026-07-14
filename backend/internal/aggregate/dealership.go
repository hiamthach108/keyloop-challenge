package aggregate

import "github.com/hiamthach108/keyloop-challenge/backend/internal/model"

type DealershipAggregate struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (a *DealershipAggregate) FromModel(m *model.Dealership) {
	if a == nil || m == nil {
		return
	}
	a.ID, a.Name, a.Location = m.ID, m.Name, m.Location
}
