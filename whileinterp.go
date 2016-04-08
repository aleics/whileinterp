package main

import (
	"fmt"
	"strings"
	"errors"
	"strconv"
)

var possOP = [...]string{"=", ":=", "<", ">", "==", "!="}
var possFunc = [...]string{"zero", "inc", "dec", "val"}

const whileFuncSTRING = "WHILE"
const doSTRING = "DO"
const odSTRING = "OD"
const assignOPSTRING = "="
const declareOPSTRING = ":="
const littleofOPSTRING = "<"
const biggerofOPSTRING = ">"
const isOPSTRING = "=="
const isNotOPSTRING = "!="

const sizeDeclareOP = len(declareOPSTRING)
const sizeAssignOP = len(assignOPSTRING)

const sizeParenthesis = 1

type variable struct {
	name string
	value int	
}

type stmt struct {
	content string
}

type logicExpr struct {
	op string
	firstVar variable
	secondVar variable
}

func (l *logicExpr) parseExpr(exprString string) error {
	if pos := strings.Index(exprString, littleofOPSTRING); pos != -1 {
		l.op = littleofOPSTRING		
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(littleofOPSTRING):])

		return nil
	} else if pos := strings.Index(exprString, biggerofOPSTRING); pos != -1 {
		l.op = biggerofOPSTRING
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(biggerofOPSTRING):])

		return nil
	} else if pos := strings.Index(exprString, isOPSTRING); pos != -1 {
		l.op = isOPSTRING
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(isOPSTRING):])
		
		return nil
	} else if pos := strings.Index(exprString, isNotOPSTRING); pos != -1 {
		l.op = isNotOPSTRING
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(isNotOPSTRING):])
		
		return nil
	} else {
		return errors.New("parseExpr: operation not defined '" + exprString + "'")
	}
}

func (l *logicExpr) evalLogicExpr() bool {
	switch(l.op) {
		case littleofOPSTRING:
			return l.firstVar.value < l.secondVar.value
		case biggerofOPSTRING:
			return l.firstVar.value > l.secondVar.value		
		case isOPSTRING:
			return l.firstVar.value == l.secondVar.value
		case isNotOPSTRING:
			return l.firstVar.value != l.secondVar.value
		default:
			return false
	}
}

type program struct {
	vars []variable
	stmts []stmt
	isDone bool
}

func initProgram() *program {
	p := new(program)
	p.vars = []variable{}
	p.stmts = []stmt{}
	p.isDone = false
	
	return p
}

func (p *program) getStmts(code string) error{
	listStmts := strings.Split(code, ";")
	
	if len(listStmts) == 0 {
		return errors.New("Code doesn't have any delimiter")
	}
	
	for _, v := range listStmts {
		stmt := new(stmt)
		stmt.content = strings.TrimSpace(v)		
		p.stmts = append(p.stmts, *stmt)
	}
	return nil
}

func (p *program) isVarPresent(name string) bool {
	for _, v := range p.vars {
		if v.name == name {
			return true
		}
	}
	return false
}

func (p *program) setVar(newVar variable) error {
	for i, v := range p.vars {
		if v.name == newVar.name {
			p.vars[i] = newVar
			return nil
		}
	}
	return errors.New("setVar: variable not present")
}

func (p *program) getVar(name string) (variable, error) {
	for _, v := range p.vars {
		if v.name == name {
			return v, nil
		}
	}
	return *new(variable), errors.New("getVar: variable not present")
}

func (p *program) parseWhile(whileCode string) (*logicExpr, string, error) {
	whileCode = strings.TrimSpace(whileCode)
	
	exprString := getExprFromWhile(whileCode)
	stmtString := getStmtFromWhile(whileCode)
	
	expr := new(logicExpr)	
	if err := expr.parseExpr(exprString); err != nil {
		return new(logicExpr), "", err
	}
	
	if val, err := p.getVar(expr.firstVar.name); err == nil {
		expr.firstVar = val
	} else {
		return new(logicExpr), "", errors.New("variable '" + expr.firstVar.name + "' not defined")
	}	
		
	if val, err := p.getVar(expr.secondVar.name); err == nil {
		expr.secondVar = val
	} else {
		return new(logicExpr), "", errors.New("variable '" + expr.secondVar.name + "' not defined")
	}
	
	return expr, stmtString, nil
}

func (p *program) parseProgram() error {
	for _, s := range p.stmts {
		if pos := strings.Index(s.content, "WHILE"); pos != - 1{
			expr, stmtWhile, err := p.parseWhile(s.content)
			if err != nil {
				return err
			}
			
			subprogram := initProgram()
			subprogram.getStmts(stmtWhile)
			subprogram.vars = p.vars
			
			for expr.evalLogicExpr() {				
				subprogram.parseProgram()
				
				expr.firstVar, _ = subprogram.getVar(expr.firstVar.name)
				expr.secondVar, _ = subprogram.getVar(expr.secondVar.name)		
			}
			
			p.vars = subprogram.vars
			
		} else {	
			if pos := strings.Index(s.content, ":="); pos != -1 {
				v := new(variable)
				v.name = strings.TrimSpace(s.content[:pos])
				
				if p.isVarPresent(v.name) {
					return errors.New("error using operator ':='. variable '" + v.name + "' already present.")
				}
				
				val, err := strconv.Atoi(s.content[(pos+sizeDeclareOP)+1:])
				if err != nil {
					val, err = p.execFunc(s.content[(pos+sizeDeclareOP)+1:])
					if err != nil {
						return err
					}
				}
				
				v.value = val		
						
				p.vars = append(p.vars, *v)
			} else if pos := strings.Index(s.content, "="); pos != -1 {
				v := new(variable)
				v.name = strings.TrimSpace(s.content[:pos])
				
				if !p.isVarPresent(v.name) {
					return errors.New("error using operator '='. variable '" + v.name + "' is not present.")
				}
				
				val, err := strconv.Atoi(s.content[(pos+sizeAssignOP)+1:])
				if err != nil {
					val, err = p.execFunc(s.content[(pos+sizeAssignOP)+1:])
					if err != nil {
						return err
					}
				}
			
				v.value = val
				p.setVar(*v)
			}
		}
	}

	return nil
}

func (p *program) execFunc(code string) (int, error) {
	if strings.Index(code, "zero()") != -1 {
		return zero(), nil
	} else if pos := strings.Index(code, "val"); pos != -1 {
		valueString := code[strings.Index(code, "(") + sizeParenthesis:strings.Index(code, ")")]
		
		v, err := strconv.Atoi(valueString)
		if err != nil {
			if currVar, err := p.getVar(valueString); err == nil {
				return val(currVar.value), nil
			}
			return 0, err
		}
		return val(v), nil
	} else if pos := strings.Index(code, "inc"); pos != -1 {
		valueString := code[strings.Index(code, "(") + sizeParenthesis:strings.Index(code, ")")]	
		
		v, err := strconv.Atoi(valueString)
		if err != nil {
			if currVar, err := p.getVar(valueString); err == nil {
				return inc(currVar.value), nil
			}
			return 0, err
		}
		return inc(v), nil
	} else if pos := strings.Index(code, "dec"); pos != -1 {
		valueString := code[strings.Index(code, "(") + sizeParenthesis:strings.Index(code, ")")]
		
		v, err := strconv.Atoi(valueString)
		if err != nil {
			if currVar, err := p.getVar(valueString); err == nil {
				return dec(currVar.value), nil
			}
			return 0, err
			
		}
		return dec(v), nil
	} else {		
		return 0, errors.New("Function not detected!")
	}
}

func (p *program) PrintVars() {
	for _, v := range p.vars {
		fmt.Println(v.name + " => " + strconv.Itoa(v.value))	
	}
}

func (p *program) PrintStmts() {
	for i, v := range p.stmts {
		fmt.Println(" - " + strconv.Itoa(i) + " block: " + v.content)	
	}
}

func getExprFromWhile(whileCode string) string {
	result := whileCode[strings.Index(whileCode, whileFuncSTRING) : strings.Index(whileCode, doSTRING)]
	result = result[strings.Index(result, "(") + sizeParenthesis : strings.Index(result, ")")]

	return result
}

func getStmtFromWhile(whileCode string) string {
	return strings.TrimSpace(whileCode[strings.Index(whileCode, doSTRING) + len(doSTRING) : strings.Index(whileCode, odSTRING)])	
}

func execCode(code string) {
	mainProgram := initProgram()

	fmt.Print("Input program: ")
	fmt.Println(code)

	if err := mainProgram.getStmts(code); err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Code blocks: ")
	
	mainProgram.PrintStmts()
	
	fmt.Println("Loading...")
	if err := mainProgram.parseProgram(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Output: ")
	mainProgram.PrintVars()
}

/*Possible Functions*/
func zero() int {
	return 0
}
func val(a int) int {
	return a
}
func inc(a int) int {
	return a + 1
}
func dec(a int) int {
	return a - 1
}

/*Example:
	"xo := 2; x1 := inc(3); x2 := dec(2); WHILE(xo != x1) DO xo = inc(xo) OD;"
*/

func main() {
	code := "xo := 3; x2 := 3; WHILE(xo == x2) DO x2 = inc(x2) OD;"
	execCode(code)
}