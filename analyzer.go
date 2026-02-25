package loglinter

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the main analyzer for log linter
var Analyzer = &analysis.Analyzer{
	Name: "loglinter",
	Doc:  "checks log messages for compliance with logging standards",
	Run:  run,
}

// Sensitive keywords that should not appear in logs
var sensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"token", "api_key", "apikey", "api-key",
	"secret", "private_key", "privatekey",
	"credential", "auth",
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if !isLogCall(call) {
				return true
			}

			checkLogMessage(pass, call)
			return true
		})
	}
	return nil, nil
}

// isLogCall checks if the call is a logging function
func isLogCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Check for log/slog methods
	logMethods := []string{"Debug", "Info", "Warn", "Error", "Fatal", "Panic",
		"DebugContext", "InfoContext", "WarnContext", "ErrorContext"}

	methodName := sel.Sel.Name
	for _, method := range logMethods {
		if methodName == method {
			return true
		}
	}

	return false
}

// checkLogMessage validates the log message against all rules
func checkLogMessage(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	// Get the first argument (the message)
	var msgArg ast.Expr
	
	// Handle context-aware methods (first arg is context)
	if isSelectorName(call.Fun, "DebugContext", "InfoContext", "WarnContext", "ErrorContext") {
		if len(call.Args) < 2 {
			return
		}
		msgArg = call.Args[1]
	} else {
		msgArg = call.Args[0]
	}

	// Extract string literal
	lit, ok := msgArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	message := strings.Trim(lit.Value, `"`)
	if message == "" {
		return
	}

	// Rule 1: Check if message starts with lowercase
	checkLowercaseStart(pass, lit, message)

	// Rule 2: Check if message is in English
	checkEnglishOnly(pass, lit, message)

	// Rule 3: Check for special characters and emojis
	checkSpecialChars(pass, lit, message)

	// Rule 4: Check for sensitive data
	checkSensitiveData(pass, lit, message)
}

// isSelectorName checks if the call is one of the specified selector names
func isSelectorName(expr ast.Expr, names ...string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	for _, name := range names {
		if sel.Sel.Name == name {
			return true
		}
	}
	return false
}

// Rule 1: Check if message starts with lowercase
func checkLowercaseStart(pass *analysis.Pass, lit *ast.BasicLit, message string) {
	if len(message) == 0 {
		return
	}

	firstRune := rune(message[0])
	if unicode.IsLetter(firstRune) && unicode.IsUpper(firstRune) {
		pass.Reportf(lit.Pos(), "log message should start with lowercase letter")
	}
}

// Rule 2: Check if message is in English (no Cyrillic or other non-Latin scripts)
func checkEnglishOnly(pass *analysis.Pass, lit *ast.BasicLit, message string) {
	for _, r := range message {
		if unicode.Is(unicode.Cyrillic, r) || unicode.Is(unicode.Han, r) || 
		   unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) {
			pass.Reportf(lit.Pos(), "log message should be in English only")
			return
		}
	}
}

// Rule 3: Check for special characters and emojis
func checkSpecialChars(pass *analysis.Pass, lit *ast.BasicLit, message string) {
	// Check for emojis
	for _, r := range message {
		if r >= 0x1F300 && r <= 0x1F9FF {
			pass.Reportf(lit.Pos(), "log message should not contain emojis")
			return
		}
	}

	// Check for excessive punctuation
	if strings.Contains(message, "!!!") || strings.Contains(message, "...") ||
	   strings.Contains(message, "!🚀") {
		pass.Reportf(lit.Pos(), "log message should not contain excessive punctuation or special characters")
	}
}

// Rule 4: Check for sensitive data keywords
func checkSensitiveData(pass *analysis.Pass, lit *ast.BasicLit, message string) {
	lowerMsg := strings.ToLower(message)
	
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lowerMsg, keyword) {
			pass.Reportf(lit.Pos(), "log message may contain sensitive data (keyword: %s)", keyword)
			return
		}
	}
}
