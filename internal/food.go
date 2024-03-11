package internal

type Food struct {
	*Location
	Emoji *rune
}

func NewFood(l *Location, e *rune) *Food {
	return &Food{
		Location: l,
		Emoji:    e,
	}
}

func (o *Food) Collides(l Location) bool {
	return o.X == l.X && o.Y == l.Y
}
