package progressbar

import (
	"fmt"
	"math"
	"strings"
)

const (
	EmptyBlock    rune = '\u2591'
	CompleteBlock rune = '\u2593'
)

type Progressbar interface {
	fmt.Stringer
	SetValue(int64)
	GetLimit() int64
	GetPercentage() int
	Finish()
	Print()
}

func NewProgressbar(limit int64) Progressbar {
	return &progressbar{
		limit: limit,
	}
}

type progressbar struct {
	limit, value int64
	finish       bool
}

func (p *progressbar) String() string {
	prc := p.GetPercentage()
	return strings.Repeat(string(CompleteBlock), prc) +
		strings.Repeat(string(EmptyBlock), 100-prc) +
		fmt.Sprintf(" %d %% ", prc)
}

func (p *progressbar) SetValue(i int64) {
	if i < 0 {
		i = 0
	}
	if i > p.limit {
		i = p.limit
	}
	p.value = i
}

func (p *progressbar) GetLimit() int64 {
	return p.limit
}

func (p *progressbar) GetPercentage() int {
	if p.limit == 0 {
		return 100
	}
	return int(math.Ceil(float64(p.value) * 100 / float64(p.limit)))
}

func (p *progressbar) Finish() {
	p.SetValue(p.limit)
	p.finish = true
}

func (p *progressbar) Print() {
	if p.finish {
		fmt.Print(p.String() + "Ok \r\n")
	} else {
		fmt.Print(p.String(), "\r")
	}
}
