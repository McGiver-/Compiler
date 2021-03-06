package Syn

var terminals = map[string]bool{
	"program":  true,
	"class":    true,
	"float":    true,
	"int":      true,
	"id":       true,
	"floatNum": true,
	"intNum":   true,
	":":        true,
	"{":        true,
	"}":        true,
	"(":        true,
	")":        true,
	"[":        true,
	"]":        true,
	",":        true,
	".":        true,
	"*":        true,
	"/":        true,
	"sr":       true,
	";":        true,
	"if":       true,
	"else":     true,
	"then":     true,
	"return":   true,
	"get":      true,
	"put":      true,
	"and":      true,
	"not":      true,
	"or":       true,
	"-":        true,
	"+":        true,
	"for":      true,
	"=":        true,
	"lt":       true,
	"gt":       true,
	"geq":      true,
	"eq":       true,
	"neq":      true,
	"leq":      true,
}

var nonterminals = map[string]bool{
	"Prog":              true,
	"ProgBody":          true,
	"ClassDecl":         true,
	"Inherit":           true,
	"InheritList":       true,
	"ClassVarDecl":      true,
	"Decl":              true,
	"ClassVarDeclTail":  true,
	"FuncHead":          true,
	"FuncHeadTail":      true,
	"FuncBody":          true,
	"FuncDecl":          true,
	"FuncVarDecl":       true,
	"FuncVarDeclTail":   true,
	"ArraySize":         true,
	"ArraySizeTail":     true,
	"Type":              true,
	"FParams":           true,
	"FParamsTail":       true,
	"FuncDefList":       true,
	"FuncDefListTail":   true,
	"FuncDef":           true,
	"AssignStat":        true,
	"AssignStatTail":    true,
	"Statement":         true,
	"StatementTail":     true,
	"StatementNoAssign": true,
	"StatBlock":         true,
	"Expr":              true,
	"RelExpr":           true,
	"RelExprTail":       true,
	"ArithExpr":         true,
	"ArithExprTail":     true,
	"Term":              true,
	"TermTail":          true,
	"Factor":            true,
	"Num":               true,
	"Sign":              true,
	"EvalExprHead":      true,
	"EvalExprTail":      true,
	"EvalIndiceHead":    true,
	"EvalIndiceTail":    true,
	"EvalNestHead":      true,
	"Variable":          true,
	"VariableTail":      true,
	"VarIndiceTail":     true,
	"VarNest":           true,
	"FunctionCallExpr":  true,
	"Indice":            true,
	"AParams":           true,
	"AParamsTail":       true,
	"AssignOp":          true,
	"RelOp":             true,
	"AddOp":             true,
	"MultOp":            true,
}

// This is the production list. There are 96 productions. Number corresponds to what should be pushed onto the stack.
// This push should be done in reverse.
var attributeGrammar = [93]string{
	"@ClassListAndFuncDefList ProgBody",
	"FuncDef Prog",
	"ClassDecl Prog",
	"program FuncBody ;",
	"class id @id Inherit @InheritList { ClassVarDecl @MemberList } ; @ClassMember",
	": id @id @InheritListMember InheritList",
	"EPSILON @EPSILON",
	", id @id @InheritListMember InheritList",
	"EPSILON",
	"Decl ClassVarDeclTail",
	"Type @Type id @id",
	"ArraySizeTail ; @DimList @VarDecl ClassVarDecl",
	"Type @Type id @id FuncHeadTail",
	"@ScopeSpec sr id @id ( FParams @FparamList )",
	"@EmptyScope ( FParams @FparamList )",
	"{ FuncVarDecl @StatBlock }",
	"( FParams @FparamList ) ; @FuncDecl ClassVarDecl",
	"StatementNoAssign StatementTail",
	"id @id ArraySizeTail @DimList @VarDecl ; FuncVarDecl",
	"[ intNum @intNum ]",
	"ArraySize ArraySizeTail",
	"float @float",
	"id @id",
	"int @int",
	"Type @Type id @id ArraySizeTail @DimList @FparamMember FParamsTail",
	", FParams",
	"EPSILON",
	"Decl FuncDefListTail",
	"FuncDef FuncDefList",
	"FuncHead FuncBody @FuncDef ;",
	"Variable AssignStatTail",
	"AssignOp Expr",
	"AssignStat ;",
	"StatementNoAssign",
	"Statement StatementTail",
	"for ( Decl AssignOp Expr ; RelExpr ; AssignStat ) StatBlock ;",
	"{ StatementTail }",
	"Statement",
	"ArithExpr RelExprTail",
	"ArithExpr RelOp ArithExpr",
	"RelOp ArithExpr",
	"EPSILON",
	"Term ArithExprTail",
	"AddOp Term ArithExprTail",
	"Factor TermTail",
	"MultOp Factor TermTail",
	"( ArithExpr )",
	"not Factor",
	"EvalExprHead",
	"floatNum",
	"intNum",
	"+",
	"-",
	"id EvalExprTail",
	"FunctionCallExpr",
	"EvalIndiceTail",
	"Indice EvalIndiceTail",
	"EvalNestHead",
	"EvalIndiceHead",
	". EvalExprHead",
	"id VariableTail",
	"VarNest",
	"VarIndiceTail",
	"Indice VariableTail",
	". Variable",
	"( AParams )",
	"[ ArithExpr ]",
	"Expr AParamsTail",
	"EPSILON",
	", AParams",
	"EPSILON",
	"=",
	"eq",
	"geq",
	"gt",
	"leq",
	"lt",
	"neq",
	"+",
	"-",
	"or",
	"*",
	"/",
	"and",
	"Type @Type FuncVarDeclTail",
	"FuncDecl",
	"VariableTail AssignStatTail ; StatementTail",
	"return ( Expr ) ;",
	"put ( Expr ) ;",
	"get ( Variable ) ;",
	"if ( Expr ) then StatBlock else StatBlock ;",
	"Num",
	"Sign Factor",
}

var semanticAction = map[string]string{
	"something": "something",
}

// This table has as key the Nonterminal. This is a key to another map which has as a key
// the expected terminals and as value the number of the production that should be pushed to the stack.
// These numbers correspond to the productions is the attributeGrammar array.
var predictionTable = map[string]map[string]int{
	"Prog":              map[string]int{"program": 0, "class": 2, "float": 1, "id": 1, "int": 1},
	"ProgBody":          map[string]int{"program": 3},
	"ClassDecl":         map[string]int{"class": 4},
	"Inherit":           map[string]int{":": 5, "{": 6},
	"InheritList":       map[string]int{",": 7, "{": 6},
	"ClassVarDecl":      map[string]int{"float": 9, "id": 9, "int": 9, "}": 6},
	"Decl":              map[string]int{"float": 10, "id": 10, "int": 10},
	"ClassVarDeclTail":  map[string]int{"[": 11, ";": 11, "(": 85},
	"FuncHead":          map[string]int{"float": 12, "id": 12, "int": 12},
	"FuncHeadTail":      map[string]int{"sr": 13, "(": 14, "{": 6},
	"FuncBody":          map[string]int{"{": 15},
	"FuncDecl":          map[string]int{"(": 16},
	"FuncVarDecl":       map[string]int{"float": 84, "id": 84, "int": 84, "for": 17, "if": 17, "get": 17, "put": 17, "return": 17, "}": 6},
	"FuncVarDeclTail":   map[string]int{"id": 18, "[": 86, ".": 86, "=": 86},
	"ArraySize":         map[string]int{"[": 19},
	"ArraySizeTail":     map[string]int{"[": 20, ";": 6, ")": 6, ",": 6},
	"Type":              map[string]int{"float": 21, "id": 22, "int": 23},
	"FParams":           map[string]int{"float": 24, "id": 24, "int": 24, ")": 6},
	"FParamsTail":       map[string]int{",": 25, ")": 6},
	"FuncDefList":       map[string]int{"float": 27, "id": 27, "int": 27, "$": 6},
	"FuncDefListTail":   map[string]int{"float": 28, "id": 28, "int": 28},
	"FuncDef":           map[string]int{"float": 29, "id": 29, "int": 29},
	"AssignStat":        map[string]int{"id": 30},
	"AssignStatTail":    map[string]int{"=": 31},
	"Statement":         map[string]int{"id": 32, "for": 33, "if": 33, "get": 33, "put": 33, "return": 33},
	"StatementTail":     map[string]int{"id": 34, "for": 34, "if": 34, "get": 34, "put": 34, "return": 34, "}": 6},
	"StatementNoAssign": map[string]int{"for": 35, "return": 87, "put": 88, "get": 89, "if": 90},
	"StatBlock":         map[string]int{"{": 36, "id": 37, "for": 37, "if": 37, "get": 37, "put": 37, "return": 37, ";": 6, "else": 6},
	"Expr":              map[string]int{"(": 38, "not": 38, "id": 38, "floatNum": 38, "intNum": 38, "+": 38, "-": 38},
	"RelExpr":           map[string]int{"(": 39, "not": 39, "id": 39, "floatNum": 39, "intNum": 39, "+": 39, "-": 39},
	"RelExprTail":       map[string]int{"eq": 40, "neq": 40, "lt": 40, "gt": 40, "leq": 40, "geq": 40, ",": 6, ";": 6, ")": 6},
	"ArithExpr":         map[string]int{"(": 42, "not": 42, "id": 42, "floatNum": 42, "intNum": 42, "+": 42, "-": 42},
	"ArithExprTail":     map[string]int{"+": 43, "-": 43, "or": 43, ";": 6, ")": 6, ",": 6, "eq": 6, "geq": 6, "gt": 6, "leq": 6, "lt": 6, "neq": 6, "]": 6, "*": 6, "/": 6, "and": 6},
	"Term":              map[string]int{"(": 44, "not": 44, "id": 44, "floatNum": 44, "intNum": 44, "+": 44, "-": 44},
	"TermTail":          map[string]int{"*": 45, "/": 45, "and": 45, ";": 6, ")": 6, ",": 6, "eq": 6, "geq": 6, "gt": 6, "leq": 6, "lt": 6, "neq": 6, "]": 6, "or": 6, "+": 6, "-": 6},
	"Factor":            map[string]int{"(": 46, "not": 47, "id": 48, "floatNum": 91, "intNum": 91, "+": 92, "-": 92},
	"Num":               map[string]int{"intNum": 50, "floatNum": 49},
	"Sign":              map[string]int{"+": 51, "-": 52},
	"EvalExprHead":      map[string]int{"id": 53},
	"EvalExprTail":      map[string]int{"(": 54, ".": 55, "[": 55, ";": 6, ")": 6, ",": 6, "eq": 6, "geq": 6, "gt": 6, "leq": 6, "lt": 6, "neq": 6, "]": 6, "or": 6, "+": 6, "-": 6, "*": 6, "/": 6, "and": 6},
	"EvalIndiceHead":    map[string]int{"[": 56},
	"EvalIndiceTail":    map[string]int{".": 57, "[": 58, ";": 6, ")": 6, ",": 6, "eq": 6, "geq": 6, "gt": 6, "leq": 6, "lt": 6, "neq": 6, "]": 6, "or": 6, "+": 6, "-": 6, "*": 6, "/": 6, "and": 6},
	"EvalNestHead":      map[string]int{".": 59},
	"Variable":          map[string]int{"id": 60},
	"VariableTail":      map[string]int{"[": 62, ".": 61, "=": 6, ")": 6},
	"VarIndiceTail":     map[string]int{"[": 63, "=": 6, ")": 6},
	"VarNest":           map[string]int{".": 64},
	"FunctionCallExpr":  map[string]int{"(": 65},
	"Indice":            map[string]int{"[": 66},
	"AParams":           map[string]int{"(": 67, "not": 67, "id": 67, "floatNum": 67, "intNum": 67, "+": 67, "-": 67, ")": 6},
	"AParamsTail":       map[string]int{",": 69, ")": 6},
	"AssignOp":          map[string]int{"=": 71},
	"RelOp":             map[string]int{"eq": 72, "geq": 73, "gt": 74, "leq": 75, "lt": 76, "neq": 77},
	"AddOp":             map[string]int{"+": 78, "-": 79, "or": 80},
	"MultOp":            map[string]int{"*": 81, "/": 82, "and": 83},
}
