package techan

import "github.com/sdcoffey/big"

type cumulativeIndicator struct {
	Indicator
	window int
	mult   big.Decimal
}

// NewCumulativeGainsIndicator returns a derivative indicator which returns all gains made in a base indicator for a given
// window.
func NewCumulativeGainsIndicator(indicator Indicator, window int) Indicator {
	return cumulativeIndicator{
		Indicator: indicator,
		window:    window,
		mult:      big.ONE,
	}
}

// NewCumulativeLossesIndicator returns a derivative indicator which returns all losses in a base indicator for a given
// window.
func NewCumulativeLossesIndicator(indicator Indicator, window int) Indicator {
	return cumulativeIndicator{
		Indicator: indicator,
		window:    window,
		mult:      big.ONE.Neg(),
	}
}

func (ci cumulativeIndicator) Calculate(index int) big.Decimal {
	total := big.NewDecimal(0.0)

	for i := Max(1, index-(ci.window-1)); i <= index; i++ {
		diff := ci.Indicator.Calculate(i).Sub(ci.Indicator.Calculate(i - 1))
		if diff.Mul(ci.mult).GT(big.ZERO) {
			total = total.Add(diff.Abs())
		}
	}

	return total
}

type percentChangeIndicator struct {
	Indicator
}

// NewPercentChangeIndicator returns a derivative indicator which returns the percent change (positive or negative)
// made in a base indicator up until the given indicator
func NewPercentChangeIndicator(indicator Indicator) Indicator {
	return percentChangeIndicator{indicator}
}

func (pgi percentChangeIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	cp := pgi.Indicator.Calculate(index)
	cplast := pgi.Indicator.Calculate(index - 1)
	return cp.Div(cplast).Sub(big.ONE)
}
