package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"fmt"
	"math"
)

type MathOps struct {
	LeftOp     interface{}
	RightOp    interface{}
	leftFloat  float64
	rightFloat float64
	operations map[string]func() float64
}

func (c *MathOps) Execute(op string) (float64, error) {
	if c.operations == nil {
		c.operations = map[string]func() float64{
			"-":           c.subtract,
			"subtract":    c.subtract,
			"+":           c.add,
			"add":         c.add,
			"*":           c.multiply,
			"multiply":    c.multiply,
			"/":           c.divide,
			"divide":      c.divide,
			"abs":         c.abs,
			"acos":        c.acos,
			"acosh":       c.acosh,
			"asin":        c.asin,
			"atan":        c.atan,
			"atan2":       c.atan2,
			"atanh":       c.atanh,
			"c64":         c.c64,
			"cbrt":        c.cbrt,
			"ceil":        c.ceil,
			"copysign":    c.copysign,
			"cos":         c.cos,
			"cosh":        c.cosh,
			"dim":         c.dim,
			"erf":         c.erf,
			"erfc":        c.erfc,
			"erfcinv":     c.erfcinv,
			"erfinv":      c.erfinv,
			"exp":         c.exp,
			"exp2":        c.exp2,
			"expm1":       c.expm1,
			"floor":       c.floor,
			"hypot":       c.hypot,
			"log":         c.log,
			"log10":       c.log10,
			"log1p":       c.log1p,
			"log2":        c.log2,
			"max":         c.max,
			"min":         c.min,
			"mod":         c.mod,
			"nextafter":   c.nextafter,
			"pow":         c.pow,
			"remainder":   c.remainder,
			"round":       c.round,
			"roundtoeven": c.roundtoeven,
			"sin":         c.sin,
			"sinh":        c.sinh,
			"sqrt":        c.sqrt,
			"tan":         c.tan,
			"tanh":        c.tanh,
			"trunc":       c.trunc,
		}
	}

	_, ok := c.operations[op]
	if !ok {
		return float64(0), fmt.Errorf("compare: invalid operation %s", op)
	}

	l, err := IfaceToFloat(c.LeftOp)
	if err != nil {
		return 0, err
	}
	r, err := IfaceToFloat(c.RightOp)
	if err != nil {
		return 0, err
	}

	c.leftFloat = l
	c.rightFloat = r
	return c.operations[op](), nil
}

func (c *MathOps) add() float64 {
	return c.leftFloat + c.rightFloat
}

func (c *MathOps) subtract() float64 {
	return c.leftFloat - c.rightFloat
}

func (c *MathOps) multiply() float64 {
	return c.leftFloat * c.rightFloat
}

func (c *MathOps) divide() float64 {
	return c.leftFloat / c.rightFloat // dividing by zero is possible.  Don't do that.
}

func (c *MathOps) abs() float64 {
	return math.Abs(c.leftFloat)
}

func (c *MathOps) acos() float64 {
	return math.Acos(c.leftFloat)
}

func (c *MathOps) acosh() float64 {
	return math.Acosh(c.leftFloat)
}

func (c *MathOps) asin() float64 {
	return math.Asin(c.leftFloat)
}

func (c *MathOps) asinh() float64 {
	return math.Asinh(c.leftFloat)
}

func (c *MathOps) atan() float64 {
	return math.Atan(c.leftFloat)
}

func (c *MathOps) atan2() float64 {
	return math.Atan2(c.leftFloat, c.rightFloat)
}

func (c *MathOps) atanh() float64 {
	return math.Atanh(c.leftFloat)
}

func (c *MathOps) c64() float64 {
	return float64(6502)
}

func (c *MathOps) cbrt() float64 {
	return math.Cbrt(c.leftFloat)
}

func (c *MathOps) ceil() float64 {
	return math.Ceil(c.leftFloat)
}

func (c *MathOps) copysign() float64 {
	return math.Copysign(c.leftFloat, c.rightFloat)
}

func (c *MathOps) cos() float64 {
	return math.Cos(c.leftFloat)
}

func (c *MathOps) cosh() float64 {
	return math.Cosh(c.leftFloat)
}

func (c *MathOps) dim() float64 {
	return math.Dim(c.leftFloat, c.rightFloat)
}

func (c *MathOps) erf() float64 {
	return math.Erf(c.leftFloat)
}

func (c *MathOps) erfc() float64 {
	return math.Erfc(c.leftFloat)
}

func (c *MathOps) erfcinv() float64 {
	return math.Erfcinv(c.leftFloat)
}

func (c *MathOps) erfinv() float64 {
	return math.Erfinv(c.leftFloat)
}

func (c *MathOps) exp() float64 {
	return math.Exp(c.leftFloat)
}

func (c *MathOps) exp2() float64 {
	return math.Exp2(c.leftFloat)
}

func (c *MathOps) expm1() float64 {
	return math.Expm1(c.leftFloat)
}

func (c *MathOps) floor() float64 {
	return math.Floor(c.leftFloat)
}

func (c *MathOps) hypot() float64 {
	return math.Hypot(c.leftFloat, c.rightFloat)
}

func (c *MathOps) log() float64 {
	return math.Log(c.leftFloat)
}

func (c *MathOps) log10() float64 {
	return math.Log10(c.leftFloat)
}

func (c *MathOps) log1p() float64 {
	return math.Log1p(c.leftFloat)
}

func (c *MathOps) log2() float64 {
	return math.Log2(c.leftFloat)
}

func (c *MathOps) max() float64 {
	return math.Max(c.leftFloat, c.rightFloat)
}

func (c *MathOps) min() float64 {
	return math.Min(c.leftFloat, c.rightFloat)
}

func (c *MathOps) mod() float64 {
	return math.Mod(c.leftFloat, c.rightFloat)
}

func (c *MathOps) nextafter() float64 {
	return math.Nextafter(c.leftFloat, c.rightFloat)
}

func (c *MathOps) pow() float64 {
	return math.Pow(c.leftFloat, c.rightFloat)
}

func (c *MathOps) remainder() float64 {
	return math.Remainder(c.leftFloat, c.rightFloat)
}

func (c *MathOps) round() float64 {
	return math.Round(c.leftFloat)
}

func (c *MathOps) roundtoeven() float64 {
	return math.RoundToEven(c.leftFloat)
}

func (c *MathOps) sin() float64 {
	return math.Sin(c.leftFloat)
}

func (c *MathOps) sinh() float64 {
	return math.Sinh(c.leftFloat)
}

func (c *MathOps) sqrt() float64 {
	return math.Sqrt(c.leftFloat)
}

func (c *MathOps) tan() float64 {
	return math.Tan(c.leftFloat)
}

func (c *MathOps) tanh() float64 {
	return math.Tanh(c.leftFloat)
}

func (c *MathOps) trunc() float64 {
	return math.Trunc(c.leftFloat)
}
