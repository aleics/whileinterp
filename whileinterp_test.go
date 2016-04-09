package whileinterp

import "testing"

const testCode1 = "xo := 2; x1 := inc(3); x2 := dec(2); WHILE(xo != x1) DO xo = inc(xo) OD;"
const testCode2 = "xo := zero(); x1 := 2; x2 := inc(x1); WHILE(x1 < x2) DO x2 = dec(x2) OD;"

const testExprLittleOfSTRING = "xo < x1"
const testExprBiggerOfSTRING = "xo > x1"
const testExprIsSTRING = "xo == x1"
const testExprIsNotSTRING = "xo != x1"


/*********************** TESTING ***********************/
func TestParseExprLittleOf(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprLittleOfSTRING); err != nil {
        t.Error(err)
    }
}

func TestParseExprBiggerOf(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprBiggerOfSTRING); err != nil {
        t.Error(err)
    }
}

func TestParseExprIs(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprIsSTRING); err != nil {
        t.Error(err)
    }
}

func TestParseExprIsNot(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprIsNotSTRING); err != nil {
        t.Error(err)
    }
}

func TestEvalExprLittleOf(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprLittleOfSTRING); err != nil {
        t.Error(err)
    }
    
    le.firstVar.value = 2
    le.secondVar.value = 3
    
    if ret := le.evalLogicExpr(); ret != true {
        t.Error("returned value not valid")
    } 
}

func TestEvalExprBiggerOf(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprBiggerOfSTRING); err != nil {
        t.Error(err)
    }
    
    le.firstVar.value = 3
    le.secondVar.value = 2
    
    if ret := le.evalLogicExpr(); ret != true {
        t.Error("returned value not valid")
    }    
}

func TestEvalExprIs(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprIsSTRING); err != nil {
        t.Error(err)
    }
    
    le.firstVar.value = 3
    le.secondVar.value = 3
    
    if ret := le.evalLogicExpr(); ret != true {
        t.Error("returned value not valid")
    }   
}

func TestEvalExprIsNot(t *testing.T) {
    le := new(logicExpr)
    if err := le.parseExpr(testExprIsNotSTRING); err != nil {
        t.Error(err)
    }
    
    le.firstVar.value = 3
    le.secondVar.value = 2
    
    if ret := le.evalLogicExpr(); ret != true {
        t.Error("returned value not valid")
    }   
}

func TestGetStmts1(t *testing.T) {
    p := initProgram()
    if err := p.getStmts(testCode1); err != nil {
        t.Error(err)
    }
}

func TestIsVarPresent(t *testing.T) {
    p := initProgram()
    v := new(variable)
    v.name = "xo"
    v.value = 2
    
    if p.isVarPresent(v.name) {
        t.Error("returned valued not valid")
    }
}

func TestAddVar(t *testing.T) {
    p := initProgram()
    v := new(variable)
    v.name = "xo"
    v.value = 2
    
    if err := p.addVar(v); err != nil {
        t.Error(err)
    }
}

func TestSetVar(t *testing.T) {
    p := initProgram()
    v := new(variable)
    v.name = "xo"
    v.value = 2
    
    p.vars = append(p.vars, *v)
    v.value = 3
    
    if err := p.setVar(v); err != nil {
        t.Error(err)
    }
}

func TestGetVar(t *testing.T) {
    p := initProgram()
    v := new(variable)
    v.name = "xo"
    v.value = 2
    
    p.vars = append(p.vars, *v)
    retVar, err := p.getVar(v.name)
    if err != nil {
        t.Error(err)
        return
    }
    if retVar != *v {
        t.Error("returned variable not equal")
    }
}

func TestGetExprFromWhile(t *testing.T) {
    whileCode := "WHILE(xo != x1) DO xo = x1 OD"
    expecExpr := "xo != x1"
    
    if retExpr := getExprFromWhile(whileCode); retExpr != expecExpr {
        t.Error("unexpected returned expression:\n returned: ", retExpr, "\n expected: ", expecExpr)
    }
}

func TestGetStmtFromWhile(t *testing.T) {
    whileCode := "WHILE(xo != x1) DO xo = x1 OD"
    expecStmt := "xo = x1"
    
    if retStmt := getStmtFromWhile(whileCode); retStmt != expecStmt {
        t.Error("unexpected returned statement:\n returned: ", retStmt, "\n expected: ", expecStmt)
    }
}

func TestParseWhile1(t *testing.T) {
    p := initProgram()
    xo := new(variable)
    xo.name = "xo"
    xo.value = 2
    
    x1 := new(variable)
    x1.name = "x1"
    x1.value = 3
    
    whileCode := "WHILE(xo != x1) DO xo = x1 OD"
    expectedStmt := "xo = x1"
    
    expectedExpr := new(logicExpr)
    expectedExpr.op = "!="
    expectedExpr.firstVar = *xo
    expectedExpr.secondVar = *x1
    
    doTestParseWhile(whileCode, xo, x1, p, expectedExpr, expectedStmt, t)
}

func TestParseWhile2(t *testing.T) {
    p := initProgram()
    xo := new(variable)
    xo.name = "xo"
    xo.value = 2
    
    x1 := new(variable)
    x1.name = "x1"
    x1.value = 3
    
    whileCode := "WHILE(xo > x1) DO xo = inc(x1) OD"
    expectedStmt := "xo = inc(x1)"
    
    expectedExpr := new(logicExpr)
    expectedExpr.op = ">"
    expectedExpr.firstVar = *xo
    expectedExpr.secondVar = *x1
    
    doTestParseWhile(whileCode, xo, x1, p, expectedExpr, expectedStmt, t)
}

func TestParseWhile3(t *testing.T) {
    p := initProgram()
    xo := new(variable)
    xo.name = "xo"
    xo.value = 2
    
    x1 := new(variable)
    x1.name = "x1"
    x1.value = 3
    
    whileCode := "WHILE(xo == x1) DO xo = dec(x1) OD"
    expectedStmt := "xo = dec(x1)"
    
    expectedExpr := new(logicExpr)
    expectedExpr.op = "=="
    expectedExpr.firstVar = *xo
    expectedExpr.secondVar = *x1
    
    doTestParseWhile(whileCode, xo, x1, p, expectedExpr, expectedStmt, t)
}

func doTestParseWhile(whileCode string, firstVar, secondVar *variable, p *program, expectedExpr *logicExpr, expectedStmt string, t *testing.T) {        
    if err := p.addVar(firstVar); err != nil {
        t.Error(err)
        return
    }    
    if err := p.addVar(secondVar); err != nil {
        t.Error(err)
        return
    }   
    
    retExpr, retStmt, err := p.parseWhile(whileCode)
    if err != nil {
        t.Error(err)
        return
    }    
    if *retExpr != *expectedExpr { //without the pointers not working
        t.Error("unexpected returned expression:\n returned: ", retExpr, "\n expected: ", expectedExpr)
    }    
    if retStmt != expectedStmt {
        t.Error("unexpected returned statement")
    }
}

func TestExecFuncZero(t *testing.T) {
    p := initProgram()    
    expecVal := 0
    
    funcCode := "xo = zero()"
    retVal, err := p.execFunc(funcCode)
    if err != nil {
        t.Error(err)
        return
    }
    
    if retVal != expecVal {
        t.Error("unexpected returned expression:\n returned: ", retVal, "\n expected: ", expecVal)
    }
}

func TestExecFuncValNum(t *testing.T) {
    p := initProgram()    
    expecVal := 2
    
    funcCode := "xo = val(2)"
    retVal, err := p.execFunc(funcCode)
    if err != nil {
        t.Error(err)
        return
    }
    
    if retVal != expecVal {
        t.Error("unexpected returned expression:\n returned: ", retVal, "\n expected: ", expecVal)
    }
}

func TestExecFuncValVar(t *testing.T) {
    p := initProgram()
    x1 := new(variable)
    x1.name = "x1"
    x1.value = 2
    
    if err := p.addVar(x1); err != nil {
        t.Error(err)   
    }
      
    expecVal := 2
    
    funcCode := "xo = val(x1)"
    retVal, err := p.execFunc(funcCode)
    if err != nil {
        t.Error(err)
        return
    }
    
    if retVal != expecVal {
        t.Error("unexpected returned expression:\n returned: ", retVal, "\n expected: ", expecVal)
    }
}

func TestExecFuncIncNum(t *testing.T) {
    p := initProgram()    
    expecVal := 3
    
    funcCode := "xo = inc(2)"
    retVal, err := p.execFunc(funcCode)
    if err != nil {
        t.Error(err)
        return
    }
    
    if retVal != expecVal {
        t.Error("unexpected returned expression:\n returned: ", retVal, "\n expected: ", expecVal)
    }
}

func TestExecFuncIncVar(t *testing.T) {
    p := initProgram()
    x1 := new(variable)
    x1.name = "x1"
    x1.value = 2
    
    if err := p.addVar(x1); err != nil {
        t.Error(err)   
    }
      
    expecVal := 3
    
    funcCode := "xo = inc(x1)"
    retVal, err := p.execFunc(funcCode)
    if err != nil {
        t.Error(err)
        return
    }
    
    if retVal != expecVal {
        t.Error("unexpected returned expression:\n returned: ", retVal, "\n expected: ", expecVal)
    }
}

func TestExecFuncDecNum(t *testing.T) {
    p := initProgram()    
    expecVal := 1
    
    funcCode := "xo = dec(2)"
    retVal, err := p.execFunc(funcCode)
    if err != nil {
        t.Error(err)
        return
    }
    
    if retVal != expecVal {
        t.Error("unexpected returned expression:\n returned: ", retVal, "\n expected: ", expecVal)
    }
}

func TestExecFuncDecVar(t *testing.T) {
    p := initProgram()
    x1 := new(variable)
    x1.name = "x1"
    x1.value = 2
    
    if err := p.addVar(x1); err != nil {
        t.Error(err)   
    }
      
    expecVal := 1
    
    funcCode := "xo = dec(x1)"
    retVal, err := p.execFunc(funcCode)
    if err != nil {
        t.Error(err)
        return
    }
    
    if retVal != expecVal {
        t.Error("unexpected returned expression:\n returned: ", retVal, "\n expected: ", expecVal)
    }
}

func TestParseProgram1(t *testing.T) {
    p := initProgram()
    p.getStmts(testCode1)
    if err := p.parseProgram(); err != nil {
        t.Error(err)
    }
}

func TestExecCode1(t *testing.T) {
    err := ExecCode(testCode1, false)
    if err != nil {
        t.Error(err)
    }
}

/*********************** BENCHMARK TESTING ***********************/
func BenchmarkParseProgram1(b *testing.B) {
    p := initProgram()
    p.getStmts(testCode1)
    
    benchmarkParseProgram(p, b)
}

func BenchmarkParseProgram2(b *testing.B) {
    p := initProgram()
    p.getStmts(testCode2)
    
    benchmarkParseProgram(p, b)
}

func benchmarkParseProgram(p *program, b *testing.B) {
    for i:= 0; i < b.N; i++ {
        if err := p.parseProgram(); err != nil {
            b.Error(err)
            return
        }    
        p.vars = []variable{}
    }
}