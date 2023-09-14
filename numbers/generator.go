package numbers

type Generator interface {
	GetIntSlices(numberOfSlices, length int) ([][]int, error)
}
