package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"unicode"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TokenType представляет тип лексемы
type TokenType int

const (
	TokIdentifier TokenType = iota
	TokConstant             // 0 или 1
	TokAssignment           // :=
	TokOr                   // or
	TokXor                  // xor
	TokAnd                  // and
	TokNot                  // not
	TokLeftParen            // (
	TokRightParen           // )
	TokSemicolon            // ;
	TokComment              // комментарий
	TokError                // ошибка
	TokEOF                  // конец файла
)

// Token представляет лексему
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
	TypeStr string
}

// Lexer - лексический анализатор
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
	tokens       []Token
	errors       []string
}

// NewLexer создает новый лексический анализатор
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar читает следующий символ
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar смотрит следующий символ без продвижения
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace пропускает пробелы и переводы строк
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment пропускает комментарии (/* ... */ и //)
func (l *Lexer) skipComment() bool {
	if l.ch == '/' {
		if l.peekChar() == '/' {
			// Однострочный комментарий
			l.readChar()
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			return true
		} else if l.peekChar() == '*' {
			// Многострочный комментарий
			l.readChar()
			l.readChar()
			for {
				if l.ch == 0 {
					l.errors = append(l.errors, fmt.Sprintf("Незакрытый комментарий (строка %d, столбец %d)", l.line, l.column))
					return true
				}
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar()
					l.readChar()
					break
				}
				l.readChar()
			}
			return true
		}
	}
	return false
}

// readIdentifier читает идентификатор или ключевое слово
func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

// NextToken возвращает следующую лексему
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	// Пропускаем комментарии
	for l.skipComment() {
		l.skipWhitespace()
	}

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TokAssignment, Value: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.makeError("Неожиданный символ ':'")
		}
	case '(':
		tok = Token{Type: TokLeftParen, Value: "(", Line: tok.Line, Column: tok.Column}
	case ')':
		tok = Token{Type: TokRightParen, Value: ")", Line: tok.Line, Column: tok.Column}
	case ';':
		tok = Token{Type: TokSemicolon, Value: ";", Line: tok.Line, Column: tok.Column}
	case '0':
		tok = Token{Type: TokConstant, Value: "0", Line: tok.Line, Column: tok.Column}
	case '1':
		tok = Token{Type: TokConstant, Value: "1", Line: tok.Line, Column: tok.Column}
	case 0:
		tok = Token{Type: TokEOF, Value: "", Line: tok.Line, Column: tok.Column}
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			tokType := lookupKeyword(identifier)
			tok = Token{Type: tokType, Value: identifier, Line: tok.Line, Column: tok.Column}
			tok.TypeStr = getTokenTypeName(tokType)
			return tok
		} else {
			tok = l.makeError(fmt.Sprintf("Неожиданный символ '%c'", l.ch))
		}
	}

	l.readChar()
	tok.TypeStr = getTokenTypeName(tok.Type)
	return tok
}

// makeError создает токен ошибки
func (l *Lexer) makeError(msg string) Token {
	errorMsg := fmt.Sprintf("%s (строка %d, столбец %d)", msg, l.line, l.column)
	l.errors = append(l.errors, errorMsg)
	return Token{
		Type:    TokError,
		Value:   string(l.ch),
		Line:    l.line,
		Column:  l.column,
		TypeStr: "ОШИБКА",
	}
}

// lookupKeyword определяет, является ли идентификатор ключевым словом
func lookupKeyword(ident string) TokenType {
	keywords := map[string]TokenType{
		"or":  TokOr,
		"xor": TokXor,
		"and": TokAnd,
		"not": TokNot,
	}

	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokIdentifier
}

// getTokenTypeName возвращает строковое представление типа токена
func getTokenTypeName(t TokenType) string {
	names := map[TokenType]string{
		TokIdentifier: "Идентификатор",
		TokConstant:   "Константа",
		TokAssignment: "Присваивание",
		TokOr:         "Ключевое слово (OR)",
		TokXor:        "Ключевое слово (XOR)",
		TokAnd:        "Ключевое слово (AND)",
		TokNot:        "Ключевое слово (NOT)",
		TokLeftParen:  "Левая скобка",
		TokRightParen: "Правая скобка",
		TokSemicolon:  "Точка с запятой",
		TokComment:    "Комментарий",
		TokError:      "ОШИБКА",
		TokEOF:        "Конец файла",
	}
	return names[t]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

// Analyze выполняет полный лексический анализ
func (l *Lexer) Analyze() {
	for {
		tok := l.NextToken()
		if tok.Type != TokEOF {
			l.tokens = append(l.tokens, tok)
		}
		if tok.Type == TokEOF {
			break
		}
	}
}

// AnalysisResult содержит результаты анализа
type AnalysisResult struct {
	Tokens       []Token
	Errors       []string
	HasErrors    bool
	InputText    string
	TokenCount   int
	SuccessCount int
	ErrorCount   int
}

// Template renderer
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Template functions
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Настройка шаблонов
	tmpl := template.New("").Funcs(templateFuncs())
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		e.Logger.Fatal(err)
	}

	t := &Template{
		templates: tmpl,
	}
	e.Renderer = t

	// Статические файлы
	e.Static("/static", "static")

	// Роуты
	e.GET("/", indexHandler)
	e.POST("/analyze", analyzeHandler)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func indexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func analyzeHandler(c echo.Context) error {
	inputText := c.FormValue("code")

	if inputText == "" {
		return c.Render(http.StatusOK, "result", AnalysisResult{
			Errors:    []string{"Введите текст для анализа"},
			HasErrors: true,
		})
	}

	lexer := NewLexer(inputText)
	lexer.Analyze()

	successCount := 0
	errorCount := 0
	for _, tok := range lexer.tokens {
		if tok.Type == TokError {
			errorCount++
		} else {
			successCount++
		}
	}

	result := AnalysisResult{
		Tokens:       lexer.tokens,
		Errors:       lexer.errors,
		HasErrors:    len(lexer.errors) > 0,
		InputText:    inputText,
		TokenCount:   len(lexer.tokens),
		SuccessCount: successCount,
		ErrorCount:   errorCount,
	}

	return c.Render(http.StatusOK, "result", result)
}
