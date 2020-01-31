package ta

const (
	// ewma 需要接触到这么多数量的样品后，才能算是一个合格的均值
	ewmaWarmUpSamples = 10
)

// EWMA 能够计算 EWMA
type EWMA struct {
	// 衰减因子
	alpha float64
	// 当前值
	value float64
	// 样本数量
	count uint8
}

// NewEWMA 返回 *EWMA
func NewEWMA(N int) *EWMA {
	if N <= 1 {
		panic("N in EWMA should bigger than 1")
	}
	return &EWMA{
		alpha: 2 / (float64(N) + 1),
	}
}

// Update adds a value to the series and updates the moving average.
func (e *EWMA) Update(value float64) {
	switch {
	case e.count < ewmaWarmUpSamples:
		e.count++
		e.value += value
	case e.count == ewmaWarmUpSamples:
		e.count++
		e.value = e.value / ewmaWarmUpSamples
		e.value = e.move(value)
	default:
		e.value = e.move(value)
	}
}

func (e *EWMA) move(value float64) float64 {
	return e.alpha*value + (1-e.alpha)*e.value
}

// IsInited 会返回 true 如果 EWMA 已经预热好了的话。
func (e *EWMA) IsInited() bool {
	return e.count >= ewmaWarmUpSamples
}

// Value 会返回当前的平均值
// 使用前，请检查是否预热完毕，比如
// if e.IsInited {
//   current = e.Value()
// }
func (e *EWMA) Value() float64 {
	return e.value
}
