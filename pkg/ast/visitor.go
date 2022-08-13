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

	VisitVarStmt(stmt *Var) any
	VisitExpressionStmt(stmt *Expression) any
	VisitIfStmt(stmt *If) any
	VisitPrintStmt(stmt *Print) any
	VisitBlockStmt(stmt *Block) any
	VisitWhileStmt(stmt *While) any
}
