package checkers

func BuildLintCheckers(checkList ...ILintChecker) ILintChecker {
	ret := &CompositeLintChecker{}

	for _, c := range checkList {
		ret.AppendChecker(c)
	}

	return ret
}
