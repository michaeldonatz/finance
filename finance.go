package money

import "time"

// parsedTime is the struct representing a parsed time value.
type parsedTime struct {
	Year                 int
	Month                time.Month
	Day                  int
	Hour, Minute, Second int // 15:04:05 is 15, 4, 5.
	Nanosecond           int // Fractional second.
	Weekday              time.Weekday
	ZoneOffset           int    // seconds east of UTC, e.g. -7*60*60 for -0700
	Zone                 string // e.g., "MST"
}

var (
	Guardi int     = 100
	Guard  int64   = int64(Guardi)
	Guardf float64 = float64(Guardi)
	DP     int64   = 100         // for default of 2 decimal places => 10^2 (can be reset)
	DPf    float64 = float64(DP) // for default of 2 decimal places => 10^2 (can be reset)
	Round          = .5
	//  Round  = .5 + (1 / Guardf)
	Roundn        = Round * -1
	call   string = "c"
	put    string = "p"
)

const (
	DBZ     = "Divide by zero"
	DTL     = "Decimal places too large"
	DLZ     = "Decimal places cannot be less than zero"
	INF     = "Calulcations results in infinity"
	INFN    = "Calulcations results in negative infinity"
	NAN     = "Not a Number"
	NOOR    = "Number out of range"
	OVFL    = "Overflow"
	UND     = "Undefined Number: non a number, or infinity"
	STRCONE = "String Conversion error"
	MAXDEC  = 18
)

const ( // for GetQuote
	DOUBLEQUOTE byte   = 34
	COLONBYTE   byte   = 58
	COLON       string = string(COLONBYTE)
	QUOTEURL    string = "http://finance.google.com/finance/info?client=ig&q="
	QLTTTIME    string = "3:04PM MST"
	QLTTIME     string = "Jan 02, 3:04PM MST"
	UNABLE      string = "Unable to convert "
	BADQUOTE    string = "Quote read is bad or invalid"
	QUOTEFAIL   string = "Bad read or quote not found for "
)

// DecimalChange resets the package-wide decimal place (default is 2 decimal places)
func DecimalChange(d int) {
	if d < 0 {
		panic(DLZ)
	}
	if d > MAXDEC {
		panic(DTL)
	}
	var newDecimal int
	if d > 0 {
		newDecimal++
		for i := 0; i < d; i++ {
			newDecimal *= 10
		}
	}
	DPf = float64(newDecimal)
	DP = int64(newDecimal)
	return
}
