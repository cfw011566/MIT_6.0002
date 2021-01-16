package drunk

import (
	"fmt"
	"math/rand"

	"../location"
)

type Drunk struct {
	name        string
	stepChoices []location.Location
}

func (d *Drunk) Name() string        { return d.name }
func (d *Drunk) SetName(name string) { d.name = name }
func (d *Drunk) String() string {
	return fmt.Sprintf("name=%q, steps=%v", d.name, d.stepChoices)
}

func (d *Drunk) TakeStep() (float64, float64) {
	n := rand.Intn(len(d.stepChoices))
	step := d.stepChoices[n]
	return step.X, step.Y
}

func (d *Drunk) SetStepChoices(steps []location.Location) {
	d.stepChoices = steps[:]
}
