package internal

type Predator struct {
	*Location
	Emoji *rune
}

func NewPredator(l *Location, e *rune) *Predator {
	return &Predator{
		Location: l,
		Emoji:    e,
	}
}

func (o *Predator) Collides(l Location) bool {
	return o.X == l.X && o.Y == l.Y
}
