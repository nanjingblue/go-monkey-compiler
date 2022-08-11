package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

const (
	ILLEGAL = "ILLEGAl" // 表示遇到未知的词法单元或字符
	EOF     = "EOF"     // 表示文件结尾，通知后续语法分析器停机

	// IDENT INT 标识符 + 字面量
	IDENT = "IDENT"
	INT   = "INT"

	// ASSIGN PLUS 运算符
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// COMMA SEMICOLON 分隔符
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// FUNCTION LET 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// STRING 字符串类型
	STRING = "STRING"

	// 数组
	LBRACKET = "["
	RBRACKET = "]"

	// 哈希表中的 ':'
	COLON = ":"
)

// keywords 关键字表
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent 检查关键字表来判断标识符是否为关键字
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
