package drone

type Drone struct {
	X   float64
	Y   float64
	Z   float64
	Vel float64
}

func (d *Drone) Locate(sectorID float64) float64 {
	loc := d.X*sectorID + d.Y*sectorID + d.Z*sectorID + d.Vel
	return loc
}
