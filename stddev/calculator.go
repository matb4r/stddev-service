package stddev

type StdDev struct {
	Stddev int   `json:"stddev"`
	Data   []int `json:"data"`
}

type Calculator interface {
	CalculateStdDevs(intSlices [][]int) ([]StdDev, error)
}
