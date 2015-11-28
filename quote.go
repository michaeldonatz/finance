package money

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Quote struct {
	ID        string
	Ticker    string
	Exchange  string
	Price     Money
	PriceCur  Money
	S         string
	LTT       parsedTime // last trade
	LT        parsedTime // last trade
	Change    Money
	ChangePct float64
	CCOL      string
	EL        Money      // after hours
	ELCurrent Money      // after hours price
	ELT       parsedTime // after hours time
	EC        Money      // after hours price change
	ECP       float64    // after hours percent change
	ECCOL     string     // chb
	Dividend  Money
	Yield     float64
	Other     string // forecol
}

var qPos = map[string]int{"id": 0, // Quote position
	"t":      1,
	"e":      2,
	"l":      3,
	"l_cur":  4,
	"s":      5,
	"ltt":    6,
	"lt":     7,
	"c":      8,
	"cp":     9,
	"ccol":   10,
	"el":     11,
	"el_cur": 12,
	"elt":    13,
	"ec":     14,
	"ecp":    15,
	"eccol":  16,
	"div":    17,
	"yld":    18,
}

// GetQuote obtains a security quote of type Quote
// q = Quote
// t = the string ticker ex. "GOOG" "VTI"
// e = the string exchange ex. "NASDAQ" "NYSE"
func (q *Quote) GetQuote(t, e string) *Quote {
	return q.parseQuote(getQuotePage(e, t))
}

// worker funcs for GetQuote

func now() (t parsedTime) {
	lt := time.Now().Local()
	t.Year = lt.Year()
	t.Month = lt.Month()
	t.Weekday = lt.Weekday()
	t.Day = lt.Day()
	t.Zone, _ = lt.Zone()
	return t
}

func (q *Quote) quoteFields(bKey, bVal []byte) *Quote {
	if bVal == nil {
		return q
	}
	bStr := string(bKey)
	for fld, v := range qPos {
		if bStr == fld {
			switch v {
			case 0:
				q.ID = string(bVal)
			case 1:
				q.Ticker = string(bVal)
			case 2:
				q.Exchange = string(bVal)
			case 3:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.Price.Setf(f)
			case 4:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.PriceCur.Setf(f)
			case 5:
				q.S = string(bVal)
			case 6:
				t, err := time.Parse(QLTTTIME, string(bVal))
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.LTT = now()
				q.LTT.Hour = t.Hour()
				q.LTT.Minute = t.Minute()
				q.LTT.Zone, _ = t.Zone()
			case 7:
				t, err := time.Parse(QLTTIME, string(bVal))
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.LT = now()
				q.LT.Month = t.Month()
				q.LT.Day = t.Day()
				q.LT.Hour = t.Hour()
				q.LT.Minute = t.Minute()
				q.LT.Zone, _ = t.Zone()
			case 8:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.Change.Setf(f)
			case 9:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.ChangePct = f
			case 10:
				q.CCOL = string(bVal)
			case 11:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.EL.Setf(f)
			case 12:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.ELCurrent.Setf(f)
			case 13:
				t, err := time.Parse(QLTTIME, string(bVal))
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.ELT = now()
				q.ELT.Month = t.Month()
				q.ELT.Day = t.Day()
				q.ELT.Hour = t.Hour()
				q.ELT.Minute = t.Minute()
				q.ELT.Zone, _ = t.Zone()
			case 14:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.EC.Setf(f)
			case 15:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + string(bVal))
				}
				q.ECP = f
			case 16:
				q.ECCOL = string(bVal)
			case 17:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + err.Error() + " " + string(bVal))
				}
				q.Dividend.Setf(f)
			case 18:
				f, err := strconv.ParseFloat(string(bVal), 64)
				if err != nil {
					panic(UNABLE + fld + " " + err.Error() + " " + string(bVal))
				}
				q.Yield = f
			}
			return q
		}
	}
	panic("no valid symbol loaded list for " + bStr)
}

func (s *Quote) parseQuote(b []byte) *Quote {
	var bKey, bVal []byte
	var pos int
	for pos < len(b) {
		if b[pos] == DOUBLEQUOTE {
			if bKey == nil {
				pos++
				if pos == len(b) {
					panic(BADQUOTE)
				}
				for b[pos] != DOUBLEQUOTE {
					bKey = append(bKey, b[pos])
					pos++
					if pos == len(b) {
						panic(BADQUOTE)
					}
				}
			} else if bVal == nil {
				pos++
				if pos == len(b) {
					panic(BADQUOTE)
				}
				for b[pos] != DOUBLEQUOTE {
					bVal = append(bVal, b[pos])
					pos++
					if pos == len(b) {
						panic(BADQUOTE)
					}
				}
				s.quoteFields(bKey, bVal)
				bKey, bVal = nil, nil
			}
		}
		pos++
	}
	return s
}

func getQuotePage(t, e string) []byte {
	var b []byte
	//  r, _, err := http.Get(QUOTEURL + e + COLON + t)  // http.Get with re-direct
	r, err := http.Get(QUOTEURL + e + COLON + t)
	if err != nil {
		panic(QUOTEFAIL + e + " " + t)
	}
	defer r.Body.Close()
	b, _ = ioutil.ReadAll(r.Body)
	return b
}
