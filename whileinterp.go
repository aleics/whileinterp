/*
    whileinterp is a while interpreter written in Go.
    
    Rules / Annotations :
        - only natural integers must be declared (N). Any other type may cause problems.
        - assignment of other variables values is possible using the "val(int a)" function.
        - arithmetic operators like: "+", "-", "*", "/", "%" are not defined. Instead use the declared functions.
        - the declaration of variables is used using the operator ":=".
        - setting the variable's value is possible using "=" (the variable must to be already declared).
        - comparator operators are: "<", ">", "==", "!=".
        - for the moment, a ";" has to be used to divide the different parts/blocks of the code.
        - for the moment, just a command can be used inside the WHILE.
    
    Example code:
	   "xo := 2; x1 := inc(3); x2 := dec(2); WHILE(xo != x1) DO xo = inc(xo) OD;"
*/

package whileinterp

import (
	"fmt"
	"strings"
	"errors"
	"strconv"
)

// possOP lists the current operations available
var possOP = [...]string{"=", ":=", "<", ">", "==", "!="}

// possFunc lists the current functions available
var possFunc = [...]string{"zero", "inc", "dec", "val"}

// whileFuncSTRING defines the while syntax in a string
const whileFuncSTRING = "WHILE"

// doSTRING defines the do syntax in a string
const doSTRING = "DO"

// odSTRING defines the od syntax in a string
const odSTRING = "OD"

// assignOPSTRING defines the operator for the assign functionality
const assignOPSTRING = "="

// declareOPSTRING defines the operator for the declare functionality
const declareOPSTRING = ":="

// littleofOPSTRING defines the operator for the "<" functionality
const littleofOPSTRING = "<"

// biggerofOPSTRING defines the operator for the ">" functionality
const biggerofOPSTRING = ">"

// isOPSTRING defines the operator for the comparation "==" functionaliy
const isOPSTRING = "=="

// isNotOPSTRING defines the operator for the comparation "!=" functionality
const isNotOPSTRING = "!="

// sizeDeclareOP defines the size of the declare operator
const sizeDeclareOP = len(declareOPSTRING)

// sizeAssignOP defines the size of the assign operator
const sizeAssignOP = len(assignOPSTRING)

// sizeParenthesis defines the size of a parenthesis
const sizeParenthesis = 1

// variable type is used for defining every variable and his value
type variable struct {
	name string //name of the variable (used as a id for the variable)
	value int	//value of the variable
}

// stmt defines every block of code divided by ";"
type stmt struct {
	content string //code of the stmt
}

//l ogicExpr is any possible logic expression defined (e.g. x1 > 2)
type logicExpr struct {
	op string //operation that is defined 
	firstVar variable //first variable (left)
	secondVar variable //second variable (right)
}

// parseExpr parses an expression and saves it to the current object
// return error
func (l *logicExpr) parseExpr(exprString string) error {
	if pos := strings.Index(exprString, littleofOPSTRING); pos != -1 { //if expression contains <
		l.op = littleofOPSTRING		
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(littleofOPSTRING):]) 

		return nil
	} else if pos := strings.Index(exprString, biggerofOPSTRING); pos != -1 { //if expression contains >
		l.op = biggerofOPSTRING
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(biggerofOPSTRING):])

		return nil
	} else if pos := strings.Index(exprString, isOPSTRING); pos != -1 { //if expression contains ==
		l.op = isOPSTRING
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(isOPSTRING):])
		
		return nil
	} else if pos := strings.Index(exprString, isNotOPSTRING); pos != -1 { //if expression contains !=
		l.op = isNotOPSTRING
		l.firstVar.name = strings.TrimSpace(exprString[:pos]);
		l.secondVar.name = strings.TrimSpace(exprString[pos + len(isNotOPSTRING):])
		
		return nil
	} else { //expression format unknown
		return errors.New("parseExpr: operation not defined '" + exprString + "'")
	}
}

// evalLogicExpr evaluates if the expression is true or false
// return bool
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

// program is the main class that envolves any variable and statement defined
// annotation: program can have different "subprograms"
type program struct {
	vars []variable //slice of the different variables declared on the program
	stmts []stmt //slice of the different statements declared on the program
}

// initProgram initializes the properties of a program
// return program
func initProgram() *program {
	p := new(program)
	p.vars = []variable{}
	p.stmts = []stmt{}
	
	return p
}

// getStmts returns the different statements defined on a code
// return error
func (p *program) getStmts(code string) error {
	listStmts := strings.Split(code, ";")	
	if len(listStmts) == 0 {
		return errors.New("getStmts: code doesn't have any delimiter")
	}
	
	for _, v := range listStmts {
		stmt := new(stmt)
		stmt.content = strings.TrimSpace(v)		
		p.stmts = append(p.stmts, *stmt) //save every statement to the program object
	}
	return nil
}

// isVarPresent checks if a variable has been already declared in a program
// return bool
func (p *program) isVarPresent(name string) bool {
	for _, v := range p.vars {
		if v.name == name {
			return true
		}
	}
	return false
}

// addVar adds a variable in a program
// return error
func (p *program) addVar(newVar *variable) error {
    if !p.isVarPresent(newVar.name) {
        p.vars = append(p.vars, *newVar)
        return nil
    }    
    return errors.New("addVar: variable '" + newVar.name + "' already present")
}

// setVar modifies a variable in a program
// return error
func (p *program) setVar(newVar *variable) error {
	for i, v := range p.vars {
		if v.name == newVar.name {
			p.vars[i] = *newVar
			return nil
		}
	}
	return errors.New("setVar: variable not present")
}

// getVar returns a variable from the program given a name (id)
// return variable, error
func (p *program) getVar(name string) (variable, error) {
	for _, v := range p.vars {
		if v.name == name {
			return v, nil
		}
	}
	return *new(variable), errors.New("getVar: variable not present")
}

// parseWhile parses a while statement given a code and returns the different properties (logic expression, body code, error)
// return *logicExpr, string, error
func (p *program) parseWhile(whileCode string) (*logicExpr, string, error) {
	whileCode = strings.TrimSpace(whileCode)
	
	exprString := getExprFromWhile(whileCode) //get the logic expression
	stmtString := getStmtFromWhile(whileCode) //get the statement to do
	
	expr := new(logicExpr)	
	if err := expr.parseExpr(exprString); err != nil {
		return new(logicExpr), "", err
	}
	
	if val, err := p.getVar(expr.firstVar.name); err == nil { //get the value of the defined variable from the program object
		expr.firstVar = val
	} else {
		return new(logicExpr), "", errors.New("parseWhile: variable '" + expr.firstVar.name + "' not defined")
	}	
		
	if val, err := p.getVar(expr.secondVar.name); err == nil { //get the value of the defined variable from the program object
		expr.secondVar = val
	} else {
		return new(logicExpr), "", errors.New("parseWhile: variable '" + expr.secondVar.name + "' not defined")
	}
	
	return expr, stmtString, nil
}

// getExprFromWhile returns the logic expression from a while block
// return string
func getExprFromWhile(whileCode string) string {
	result := whileCode[strings.Index(whileCode, whileFuncSTRING) : strings.Index(whileCode, doSTRING)]
	result = result[strings.Index(result, "(") + sizeParenthesis : strings.Index(result, ")")]

	return result
}

// getStmtFromWhile returns the statement to do from a while block
// return string
func getStmtFromWhile(whileCode string) string {
	return strings.TrimSpace(whileCode[strings.Index(whileCode, doSTRING) + len(doSTRING) : strings.Index(whileCode, odSTRING)])	
}

// parseProgram parses and executes the whole program
// all the operations made will be saved on the program object
// return error
func (p *program) parseProgram() error {
	for _, s := range p.stmts {
		if pos := strings.Index(s.content, "WHILE"); pos != - 1 { //if is a "WHILE" statement
			expr, stmtWhile, err := p.parseWhile(s.content)
			if err != nil {
				return err
			}
			
			subprogram := initProgram() //a subprogram is the "do" statement from the WHILE block
			subprogram.getStmts(stmtWhile)
			subprogram.vars = p.vars
			
			for expr.evalLogicExpr() { //the subprogram will be executed as much the expr will be false				
				subprogram.parseProgram() //the subprogram has to be parsed
				
				expr.firstVar, _ = subprogram.getVar(expr.firstVar.name) //refresh the changes to the expression variables
				expr.secondVar, _ = subprogram.getVar(expr.secondVar.name)		
			}
			
			p.vars = subprogram.vars //save the changes on the main program, if finished
			
		} else { //normal statement (declarations, assignations)
			if pos := strings.Index(s.content, ":="); pos != -1 { //if a declaration
				v := new(variable)
				v.name = strings.TrimSpace(s.content[:pos])
				
				if p.isVarPresent(v.name) { //if variable already on the program -> error
					return errors.New("parseProgram: error using operator ':='. variable '" + v.name + "' already present.")
				}
				
				val, err := strconv.Atoi(s.content[(pos+sizeDeclareOP)+1:]) //get value of the declaration
				if err != nil {
					val, err = p.execFunc(s.content[(pos+sizeDeclareOP)+1:]) //execute the function for this operator (:=)
					if err != nil {
						return err
					}
				}
				
				v.value = val		
						
				p.addVar(v) //add the variable to the program
			} else if pos := strings.Index(s.content, "="); pos != -1 { //if an assignment
				v := new(variable)
				v.name = strings.TrimSpace(s.content[:pos])
				
				if !p.isVarPresent(v.name) { //if variable is not on the program -> error
					return errors.New("parseProgram: error using operator '='. variable '" + v.name + "' is not present.")
				}
				
				val, err := strconv.Atoi(s.content[(pos+sizeAssignOP)+1:]) //get value of the assignment
				if err != nil {
					val, err = p.execFunc(s.content[(pos+sizeAssignOP)+1:]) //execute the function for this operator (=)
					if err != nil {
						return err
					}
				}
			
				v.value = val 
				p.setVar(v) //set the value of the variable on the program
			}
		}
	}
	return nil
}

// execFunc executes the the different declared functions from a statement
// return int, error
func (p *program) execFunc(code string) (int, error) {
	if strings.Index(code, "zero()") != -1 { //if the zero function was called
		return zero(), nil
	} else if pos := strings.Index(code, "val"); pos != -1 { //if the val function was called
		valueString := code[strings.Index(code, "(") + sizeParenthesis:strings.Index(code, ")")] //extract the parameter's value of the function
		
		v, err := strconv.Atoi(valueString)
		if err != nil { //if value is not a number
			if currVar, err := p.getVar(valueString); err == nil { //check on the program if a variable has this id
				return val(currVar.value), nil
			}
			return 0, err
		}
		return val(v), nil
	} else if pos := strings.Index(code, "inc"); pos != -1 { //if th inc function was called
		valueString := code[strings.Index(code, "(") + sizeParenthesis:strings.Index(code, ")")] //extract the parameter's value of the function	
		
		v, err := strconv.Atoi(valueString)
		if err != nil { //if value is not a number
			if currVar, err := p.getVar(valueString); err == nil { //check on the program if a variable has this id
				return inc(currVar.value), nil
			}
			return 0, err
		}
		return inc(v), nil
	} else if pos := strings.Index(code, "dec"); pos != -1 {
		valueString := code[strings.Index(code, "(") + sizeParenthesis:strings.Index(code, ")")] //extract the parameter's value of the function
		
		v, err := strconv.Atoi(valueString)
		if err != nil { //if value is not a number
			if currVar, err := p.getVar(valueString); err == nil { //check on the program if a variable has this id
				return dec(currVar.value), nil
			}
			return 0, err
			
		}
		return dec(v), nil
	} else {		
		return 0, errors.New("execFunc: function not detected")
	}
}

//printVars prints the different variables of a program
func (p *program) printVars() {
	for _, v := range p.vars {
		fmt.Println(v.name + " => " + strconv.Itoa(v.value))	
	}
}

//printStmts prints the different statements of a program
func (p *program) printStmts() {
	for i, v := range p.stmts {
		fmt.Println(" - " + strconv.Itoa(i) + " block: " + v.content)	
	}
}

// ExecCode executes the code as a parameter (set log to true, to display the progress per console)
// return bool
func ExecCode(code string, log bool) error {
	mainProgram := initProgram() //initialize the program
    
    if log {
	   fmt.Print("Input program: ")
	   fmt.Println(code)
    }

	if err := mainProgram.getStmts(code); err != nil { //get the different statements
		fmt.Println(err)
		return err
	}
	
    if log {
	   fmt.Println("Code blocks: ")	
	   mainProgram.printStmts()
       fmt.Println("Loading...")
    }	
	
	if err := mainProgram.parseProgram(); err != nil { //parse the code and execute it
		fmt.Println(err)
		return err
	}
    
    if log { //if desired, print the output variables
	   fmt.Println("Output: ")
	   mainProgram.printVars()
    }
    
    return nil
}

/*Possible Functions*/
// zero initializes a variable to 0 (analogy to x := 0 is possible too)
// return int
func zero() int {
	return 0
}
// val returns the value of a variable
// return int
func val(a int) int {
	return a
}
// inc increments a variable's value
// return int
func inc(a int) int {
	return a + 1
}
// dec decrease a variable's value
// return int
func dec(a int) int {
	return a - 1
}