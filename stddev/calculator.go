package stddev

type StdDev struct {
	Stddev float64 `json:"stddev"`
	Data   []int   `json:"data"`
}

type Calculator interface {
	CalculateStdDevs(intSlices [][]int) ([]StdDev, error)
}
