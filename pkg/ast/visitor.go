package ast

type Visitor interface {
	VisitBinaryExpr(expr *Binary) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitUnaryExpr(expr *Unary) any
	VisitTernaryExpr(expr *Ternary) any
	VisitVariableExpr(expr *Variable) any
	VisitAssignExpr(expr *Assign) any
	VisitLogicalExpr(expr *Logical) any
	VisitCallExpr(expr *Call) any

	VisitVarStmt(stmt *Var) any
	VisitExpressionStmt(stmt *Expression) any
	VisitIfStmt(stmt *If) any
	VisitPrintStmt(stmt *Print) any
	VisitBlockStmt(stmt *Block) any
	VisitWhileStmt(stmt *While) any
	VisitBreakStatement(stmt *Break) any
	VisitFunctionStmt(stmt *Function) any
	VisitReturnStmt(stmt *Return) any
}
