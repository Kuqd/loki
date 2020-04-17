package regexp

import "regexp/syntax"

// FindEqualMatches finds a set of strict equal matches within a RE2.
// The expression is meant to be used inside anchors. Like Prometheus label matchers.
func FindEqualMatches(re string) []string {
	reg, err := syntax.Parse(re, syntax.Perl)
	if err != nil {
		return nil
	}
	reg = reg.Simplify()

	return findMatches(reg, nil)
}

func findMatches(reg *syntax.Regexp, set []string) []string {
	switch reg.Op {
	case syntax.OpAlternate:
		return findAlternateMatches(reg, set)
	case syntax.OpConcat:
		return findConcatMatches(reg, nil, set)
	case syntax.OpCapture:
		clearCapture(reg)
		return findMatches(reg, set)
	case syntax.OpLiteral:
		if !isCaseInsensitive(reg) {
			set = append(set, string((reg.Rune)))
		}
	case syntax.OpStar:
		return nil
	case syntax.OpEmptyMatch:
		return nil
	}
	return set
}

// findAlternateMatches find matches inside alternates regex ast, when possible, alternate regexp expressions such as:
// (foo|bar) or (foo|(bar|buzz)).
func findAlternateMatches(reg *syntax.Regexp, set []string) []string {
	clearCapture(reg.Sub...)
	for _, sub := range reg.Sub {
		leg := findMatches(sub, set)
		if len(leg) == 0 {
			return nil
		}
		set = leg
	}
	return set
}

func findConcatMatches(reg *syntax.Regexp, baseLiteral []byte, set []string) []string {
	clearCapture(reg.Sub...)

	// if the first and last sub are ^or$ we can ignore them we expect anchored regex.
	var subs []*syntax.Regexp
	for _, s := range reg.Sub {
		if s.Op == syntax.OpBeginText || s.Op == syntax.OpEndText {
			continue
		}
		subs = append(subs, s)
	}
	if len(subs) == 1 && subs[0].Op == syntax.OpAlternate {
		return findAlternateMatches(subs[0], set)
	}
	// We support only concat with 2 legs. such as foo(bar|buzz|f) = foobar|foobuzz|foof
	if len(subs) > 2 {
		return nil
	}

	var ok bool
	literals := 0
	for _, sub := range subs {
		if sub.Op == syntax.OpLiteral {
			// only one literal is allowed.
			if literals != 0 {
				return nil
			}
			literals++
			baseLiteral = append(baseLiteral, []byte(string(sub.Rune))...)
			continue
		}
		// if we have an alternate we must also have a base literal to apply the concatenation with.
		if sub.Op == syntax.OpAlternate && baseLiteral != nil {
			set = findMatchesConcatAlternate(sub, baseLiteral, set)
			ok = true
			continue
		}
		if sub.Op == syntax.OpEmptyMatch {
			ok = true
			set = append(set, string(baseLiteral))
			continue
		}
		return nil
	}
	// if we have a findMatchesConcatAlternate.
	if ok {
		return set
	}
	return nil
}

func findMatchesConcatAlternate(reg *syntax.Regexp, literal []byte, set []string) []string {
	for _, alt := range reg.Sub {
		if isCaseInsensitive(alt) {
			return nil
		}
		switch alt.Op {
		case syntax.OpEmptyMatch:
			set = append(set, string(literal))
		case syntax.OpLiteral:
			// concat the root literal with the alternate one.
			altBytes := []byte(string(alt.Rune))
			altLiteral := make([]byte, 0, len(literal)+len(altBytes))
			altLiteral = append(altLiteral, literal...)
			altLiteral = append(altLiteral, altBytes...)
			set = append(set, string(altLiteral))
		case syntax.OpConcat:
			return findMatchesConcatAlternate(alt, literal, set)
		default:
			return nil
		}
	}
	return set
}
