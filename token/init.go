package token

var operatorList = []string{}

func initOperatorList() {
	buckets := map[int][]string{}
	for t := operatorBegin; t < operatorEnd; t++ {
		s, found := tokenStringMap[t]
		if !found {
			continue
		}

		i := len(s)
		buckets[i] = append(buckets[i], s)
	}

	for i := 10; i >= 1; i-- {
		if ops, found := buckets[i]; found {
			operatorList = append(operatorList, ops...)
		}
	}
}

func init() {
	initOperatorList()
}
