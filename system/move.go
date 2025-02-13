package system

type Move struct {
	UnitRepo UnitRepo
}

func (m *Move) Run() error {
	for id, unit := range m.UnitRepo.List() {
		unit.Position = unit.Position.Add(unit.Velocity)

		m.UnitRepo.Upsert(id, unit)
	}

	return nil
}
