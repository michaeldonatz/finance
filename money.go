package money

/*

The package contains type Money...

type Money struct {
	M	int64
}


...which usese a fixed-length guard for precision arithmetic: the
int64 variable Guard (and its float64 and int-related variables Guardf
and Guardi.

ROunding is done on float64 to int64 by	the Rnd() function truncating
at values less than (.5 + (1 / Guardf))	or greater than -(.5 + (1 / Guardf))
in the case of negative numbers. The Guard adds four decimal places
of protection to rounding.

DP is the decimal precision, which can be changed in the DecimalPrecision()
function.  DP hold the places after the decimalplace in teh active money struct field M

The following functions are available

Abs Returns the absolute value of Money
	(m *Money) Abs() *Money
Add Adds two Money types
	(m *Money) Add(n *Money) *Money
Div Divides one Money type from another
	(m *Money) Div(n *Money) *Money
Gett gets value of money truncating after DP (see Value() for no truncation)
	(m *Money) Gett() int64
Get gets the float64 value of money (see Value() for int64)
	(m *Money) Get() float64
Mul Multiplies two Money types
	(m *Money) Mul(n *Money) *Money
Mulf Multiplies a Money with a float to return a money-stored type
	(m *Money) Mulf(f float64) *Money
Neg Returns the negative value of Money
	(m *Money) Neg() *Money
Pow is the power of Money
	(m *Money) Pow(r float64) *Money
Set sets the Money field M
	(m *Money) Set(x int64) *Money
Setf sets a float 64 into a Money type for precision calculations
	(m *Money) Setf(f float64) *Money
Sign returns the Sign of Money 1 if positive, -1 if negative
	(m *Money) Sign() int
String for money type representation in basic monetary unit (DOLLARS CENTS)
	(m *Money) String() string
Sub subtracts one Money type from another
	(m *Money) Sub(n *Money) *Money
Value returns in int64 the value of Money (also see Gett, See Get() for float64)
	(m *Money) Value() int64
*/

import (
	"fmt"
	"math"
)

type Money struct {
	M int64 // value of the integer64 Money
}

func (m *Money) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var f float64

	unmarshalErr := unmarshal(&f)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	m.Setf(f)

	return nil
}

// Abs Returns the absolute value of Money
func (m *Money) Abs() *Money {
	if m.M < 0 {
		m.Neg()
	}
	return m
}

// Add Adds two Money types
func (m *Money) Add(n *Money) *Money {
	r := m.M + n.M
	if (r^m.M)&(r^n.M) < 0 {
		panic(OVFL)
	}
	m.M = r
	return m
}

// Div Divides one Money type from another
func (m *Money) Div(n *Money) *Money {
	f := Guardf * DPf * float64(m.M) / float64(n.M) / Guardf
	i := int64(f)
	return m.Set(Rnd(i, f-float64(i)))
}

// Gett gets value of money truncating after DP (see Value() for no truncation)
func (m *Money) Gett() int64 {
	return m.M / DP
}

// Get gets the float64 value of money (see Value() for int64)
func (m *Money) Get() float64 {
	return float64(m.M) / DPf
}

// Mul Multiplies two Money types
func (m *Money) Mul(n *Money) *Money {
	return m.Set(m.M * n.M / DP)
}

// Mulf Multiplies a Money with a float to return a money-stored type
func (m *Money) Mulf(f float64) *Money {
	i := m.M * int64(f*Guardf*DPf)
	r := i / Guard / DP
	return m.Set(Rnd(r, float64(i)/Guardf/DPf-float64(r)))
}

// Neg Returns the negative value of Money
func (m *Money) Neg() *Money {
	if m.M != 0 {
		m.M *= -1
	}
	return m
}

// Pow is the power of Money
func (m *Money) Pow(r float64) *Money {
	return m.Setf(math.Pow(m.Get(), r))
}

// Set sets the Money field M
func (m *Money) Set(x int64) *Money {
	m.M = x
	return m
}

// Setf sets a float64 into a Money type for precision calculations
func (m *Money) Setf(f float64) *Money {
	fDPf := f * DPf
	r := int64(f * DPf)
	return m.Set(Rnd(r, fDPf-float64(r)))
}

// Sign returns the Sign of Money 1 if positive, -1 if negative
func (m *Money) Sign() int {
	if m.M < 0 {
		return -1
	}
	return 1
}

// String for money type representation in basic monetary unit (DOLLARS CENTS)
func (m *Money) String() string {
	return fmt.Sprintf("%d.%02d", m.Value()/DP, m.Abs().Value()%DP)
}

// Sub subtracts one Money type from another
func (m *Money) Sub(n *Money) *Money {
	r := m.M - n.M
	if (r^m.M)&^(r^n.M) < 0 {
		panic(OVFL)
	}
	m.M = r
	return m
}

// Value returns in int64 the value of Money (also see Gett(), See Get() for float64)
func (m *Money) Value() int64 {
	return m.M
}
