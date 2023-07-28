package round

type round struct {
	ID       int32  `db:"id"`
	Name     string `db:"name"`
	PackID   int32  `db:"pack_id"`
	Position int16  `db:"position"`
}
