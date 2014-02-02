package tokenparser

type ruleType int

const (
	upTo    ruleType = iota
	skipTo  ruleType = iota
	skip    ruleType = iota
	skipAny ruleType = iota
)

type Rule struct {
	Type        ruleType
	Symbol      uint8
	Destination *string
}

type RuleList []Rule

func (rl *RuleList) Start() func() *Rule {
	x := 0
	return func() *Rule {
		if len(*rl) > x {
			x++
			return &([]Rule(*rl)[x-1])
		} else {
			return nil
		}
	}
}

type Tokenparser struct {
	rules  RuleList
	Strict bool
}

func New() *Tokenparser {
	tp := new(Tokenparser)
	tp.rules = make(RuleList, 0)
	tp.Strict = false
	return tp
}

func (tp *Tokenparser) UpTo(symbol uint8, destination *string) {
	tp.rules = append(tp.rules, Rule{upTo, symbol, destination})
}

func (tp *Tokenparser) Skip(symbol uint8) {
	tp.rules = append(tp.rules, Rule{skip, symbol, nil})
}

func (tp *Tokenparser) SkipTo(symbol uint8) {
	tp.rules = append(tp.rules, Rule{skipTo, symbol, nil})
}

func (tp *Tokenparser) SkipAny() {
	tp.rules = append(tp.rules, Rule{skipAny, 0, nil})
}

func (tp *Tokenparser) ParseString(line string) bool {
	nextRule := tp.rules.Start()
	currentRule := nextRule()
	linePointer := 0
	for currentRule != nil {
		switch currentRule.Type {
		case skip:
			if linePointer >= len(line) {
				return false
			}
			if line[linePointer] != currentRule.Symbol {
				return false
			} else {
				linePointer++
			}

		case skipAny:
			if linePointer >= len(line) {
				return false
			}
			linePointer++

		case skipTo:
			for {
				if linePointer >= len(line) {
					return false
				}
				if line[linePointer] != currentRule.Symbol {
					linePointer++
				} else {
					break
				}
			}

		case upTo:
			var accumulator string = ""
			for {
				if linePointer >= len(line) {
					return false
				}
				if line[linePointer] != currentRule.Symbol {
					accumulator += string(line[linePointer])
					linePointer++
				} else {
					*currentRule.Destination = accumulator
					break
				}
			}
		}
		currentRule = nextRule()
	}
	if !tp.Strict {
		return true
	}
	return linePointer == len(line)
}
