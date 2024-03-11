package internal

type Predators []*Predator

func NewPredators(obstacles ...*Predator) Predators {
	return Predators(obstacles)
}

func (o *Predators) ClearAll() {
	(*o) = (*o)[:0]
}

func (o *Predators) Remove() {
	(*o)[0] = nil
	(*o) = (*o)[1:]
}

func (o *Predators) Add(ob *Predator) {
	*o = append(*o, ob)
}

func (o Predators) Index(i int) *Predator {
	if i < 0 || i >= len(o) {
		return nil
	}
	return o[i]
}
