package parser

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `def calculate(x, y)
  x + y
end`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{DEF, "def"},
		{IDENT, "calculate"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COMMA, ","},
		{IDENT, "y"},
		{RPAREN, ")"},
		{NEWLINE, "\n"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{NEWLINE, "\n"},
		{END, "end"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestMoreTokens(t *testing.T) {
	input := `if 5 < 10
  return true
else
  return false
end

10 == 10
10 != 9
"hello"
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IF, "if"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{NEWLINE, "\n"},
		{RETURN, "return"},
		{TRUE, "true"},
		{NEWLINE, "\n"},
		{ELSE, "else"},
		{NEWLINE, "\n"},
		{RETURN, "return"},
		{FALSE, "false"},
		{NEWLINE, "\n"},
		{END, "end"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{NEWLINE, "\n"},
		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{NEWLINE, "\n"},
		{STRING, "hello"},
		{NEWLINE, "\n"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
