package tokenparser

type ruleType uint8

const (
	upTo         ruleType = iota
	skipTo       ruleType = iota
	skip         ruleType = iota
	skipAny      ruleType = iota
	skipMultiple ruleType = iota
)

type Rule struct {
	Type        ruleType
	Symbol      uint8
	Destination *string
}

type RuleList []Rule

func (rl *RuleList) GetIterator() func() *Rule {
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

func (tp *Tokenparser) SkipMultiple(count uint8) {
	tp.rules = append(tp.rules, Rule{skipMultiple, count, nil})
}

func (tp *Tokenparser) ParseString(line string) bool {
	nextRule := tp.rules.GetIterator()

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

		case skipMultiple:
			linePointer += int(currentRule.Symbol)

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
			firstSym := linePointer
			for {
				if linePointer >= len(line) {
					return false
				}
				if line[linePointer] != currentRule.Symbol {
					linePointer++
				} else {
					*currentRule.Destination = line[firstSym:linePointer]
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
