package gaga

// A voicing modifier (voiced or semi-voiced sound mark)
type vom rune

const (
	vmNone        vom = 0
	vmVsmNonspace vom = 0x3099 // Combining voiced sound mark
	vmSsmNonspace vom = 0x309A // Combining semi-voiced sound mark
	vmVsmWide     vom = 0x309B // Wide voiced sournd mark
	vmSsmWide     vom = 0x309C // Wide semi-voiced sound mark
	vmVsmNarrow   vom = 0xFF9E // Narrow voiced sound mark
	vmSsmNarrow   vom = 0xFF9F // Narrow semi-voiced sound mark
)

func (m vom) isNone() bool {
	return m == vmNone
}

// Vsm (Voiced sound mark)
func (m vom) isVsm() bool {
	switch m {
	case vmVsmNonspace, vmVsmWide, vmVsmNarrow:
		return true
	default:
		return false
	}
}

// Ssm (Semi-voiced sound mark)
func (m vom) isSsm() bool {
	switch m {
	case vmSsmNonspace, vmSsmWide, vmSsmNarrow:
		return true
	default:
		return false
	}
}

// Vom (Voicing modifier)
func (m vom) isVom() bool {
	return m.isVsm() || m.isSsm()
}
