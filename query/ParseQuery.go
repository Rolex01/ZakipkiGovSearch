package parseQuery

import (
	"log"
	"regexp"
	"strings"
)

func dfWrap(s string, strict bool) string {
	if !strict {
		return "+" + s + "*"
	} else {
		return s
	}

}

func dfWrapPl(s string, strict bool) string {
	if !strict {
		return "+" + s + "*"
	} else {
		return s
	}
}

func phWrapPl(s string) string {
	return "" + s + " "
}

func phWrap(s string) string {
	//return "(" + dictHack(s) + ")"
	return "(" + s + ")"
}

func dictHack(s string) string {
	var syns = make(map[string]string)
	syns["А"] = "A"
	syns["В"] = "B"
	syns["Е"] = "E"
	syns["З"] = "3"
	syns["К"] = "K"
	syns["М"] = "M"
	syns["Н"] = "H"
	syns["О"] = "O"
	syns["Р"] = "P"
	syns["С"] = "C"
	syns["Т"] = "T"
	syns["Х"] = "X"
	syns["а"] = "a"
	syns["е"] = "e"
	syns["к"] = "k"
	syns["о"] = "o"
	syns["п"] = "n"
	syns["р"] = "p"
	syns["с"] = "c"
	syns["у"] = "y"
	syns["х"] = "x"
	for k, v := range syns {
		if strings.Contains(s, k) {
			s = "(" + s + ") | (" + strings.Replace(s, k, v, -1) + ")"
		}
		if strings.Contains(s, v) {
			s = "(" + s + ") | (" + strings.Replace(s, v, k, -1) + ")"
		}
	}
	return s
}

func rmSpecialChars(s string) string {
	specialChars, err := regexp.Compile("[+*()\\'\"«»]+")
	if err != nil {
		log.Println("ParseQuery Error: ", err, " Query:", s)
	}
	s = specialChars.ReplaceAllString(s, "")
	return s
}

func PlutoFix(q string, strict bool) string {
	// log.Println("Parsing query:", q)
	// 0. - Trim spaces and commas
	q = strings.TrimSpace(q)
	// 1. OR:
	q1 := strings.Split(q, ",")
	// log.Println("q1:", q1)
	// 2. Trim spaces in q1 and split by spaces
	q2 := *new([][]string)
	for _, s := range q1 {
		if len(s) > 0 {
			s = strings.TrimSpace(s)
			q2 = append(q2, strings.Split(s, " "))
		}
	}
	// log.Println("q2:", q2)
	// 3. Wrapping the words
	q3 := *new([]string)
	for i, phrase := range q2 {
		for j, word := range phrase {
			word = rmSpecialChars(word)
			if len(word) == 0 {
				continue
			}
			if string(word[0]) == "-" {
				if j != 0 {
					if !strict {
						q2[i][j] = word + "*"
					} else {
						q2[i][j] = word
					}
				} else {
					q2[i][j] = dfWrapPl(word[1:], strict)
				}
			} else {
				q2[i][j] = dfWrapPl(word, strict)
			}
		}
		q3 = append(q3, phWrapPl(strings.Join(q2[i], " ")))
	}
	// 4. Join all
	q4 := strings.Join(q3, "|")
	// log.Println("Final:", q4)
	return q4
}

func ParseQuery(q string, strict bool) string {
	// log.Println("Parsing query:", q)
	// 0. - Trim spaces and commas
	q = strings.TrimSpace(q)
	// 1. OR:
	q1 := strings.Split(q, ",")
	// log.Println("q1:", q1)
	// 2. Trim spaces in q1 and split by spaces
	q2 := *new([][]string)
	for _, s := range q1 {
		if len(s) > 0 {
			s = strings.TrimSpace(s)
			q2 = append(q2, strings.Split(s, " "))
		}
	}
	// log.Println("q2:", q2)
	// 3. Wrapping the words
	q3 := *new([]string)
	for i, phrase := range q2 {
		for j, word := range phrase {
			word = rmSpecialChars(word)
			if len(word) == 0 {
				continue
			}
			if string(word[0]) == "-" {
				if j != 0 {
					if !strict {
						q2[i][j] = word + "*"
					} else {
						q2[i][j] = word
					}
				} else {
					q2[i][j] = dfWrap(word[1:], strict)
				}
			} else {
				q2[i][j] = dfWrap(word, strict)
			}
		}
		// dict hack will be inside phWrap
		q3 = append(q3, phWrap(strings.Join(q2[i], " ")))
	}
	// 4. Join all
	q4 := strings.Join(q3, "|")
	// log.Println("Final:", q4)
	return q4
}

/*func main() {
//	test1 := "  a b, c,g -d "
	test2 := "a b -c"
	parseQuery(test2)
}*/
