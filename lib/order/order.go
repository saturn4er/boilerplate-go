package order

type Direction byte

const (
	Asc Direction = iota + 1
	Desc
)

type Order[FieldType ~byte] struct {
	Field     FieldType
	Direction Direction
}
