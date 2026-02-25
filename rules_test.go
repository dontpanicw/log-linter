package loglinter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestCheckLowercaseStart(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		wantError bool
	}{
		{"uppercase start", "Starting server", true},
		{"lowercase start", "starting server", false},
		{"number start", "123 items", false},
		{"empty message", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{}
			fset := token.NewFileSet()
			lit := &ast.BasicLit{
				ValuePos: token.Pos(fset.Base()),
				Kind:     token.STRING,
				Value:    `"` + tt.message + `"`,
			}

			var reported bool
			pass.Report = func(d analysis.Diagnostic) {
				reported = true
			}

			checkLowercaseStart(pass, lit, tt.message)

			if reported != tt.wantError {
				t.Errorf("checkLowercaseStart() reported = %v, want %v", reported, tt.wantError)
			}
		})
	}
}

func TestCheckEnglishOnly(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		wantError bool
	}{
		{"english only", "starting server", false},
		{"cyrillic text", "запуск сервера", true},
		{"mixed text", "server запуск", true},
		{"chinese text", "服务器启动", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{}
			fset := token.NewFileSet()
			lit := &ast.BasicLit{
				ValuePos: token.Pos(fset.Base()),
				Kind:     token.STRING,
				Value:    `"` + tt.message + `"`,
			}

			var reported bool
			pass.Report = func(d analysis.Diagnostic) {
				reported = true
			}

			result := checkEnglishOnly(pass, lit, tt.message)

			if reported != tt.wantError {
				t.Errorf("checkEnglishOnly() reported = %v, want %v", reported, tt.wantError)
			}
			if tt.wantError && result {
				t.Errorf("checkEnglishOnly() returned true, want false when error reported")
			}
			if !tt.wantError && !result {
				t.Errorf("checkEnglishOnly() returned false, want true when no error")
			}
		})
	}
}

func TestCheckSpecialChars(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		wantError bool
	}{
		{"normal text", "server started", false},
		{"emoji", "server started 🚀", true},
		{"triple exclamation", "error!!!", true},
		{"triple dots", "loading...", true},
		{"single punctuation", "server started!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{}
			fset := token.NewFileSet()
			lit := &ast.BasicLit{
				ValuePos: token.Pos(fset.Base()),
				Kind:     token.STRING,
				Value:    `"` + tt.message + `"`,
			}

			var reported bool
			pass.Report = func(d analysis.Diagnostic) {
				reported = true
			}

			checkSpecialChars(pass, lit, tt.message)

			if reported != tt.wantError {
				t.Errorf("checkSpecialChars() reported = %v, want %v", reported, tt.wantError)
			}
		})
	}
}

func TestCheckSensitiveData(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		wantError bool
	}{
		{"contains password", "user password: secret", true},
		{"contains token", "auth token: abc123", true},
		{"contains api_key", "api_key=12345", true},
		{"safe message", "user login successful", false},
		{"uppercase PASSWORD", "PASSWORD reset", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{}
			fset := token.NewFileSet()
			lit := &ast.BasicLit{
				ValuePos: token.Pos(fset.Base()),
				Kind:     token.STRING,
				Value:    `"` + tt.message + `"`,
			}

			var reported bool
			pass.Report = func(d analysis.Diagnostic) {
				reported = true
			}

			checkSensitiveData(pass, lit, tt.message)

			if reported != tt.wantError {
				t.Errorf("checkSensitiveData() reported = %v, want %v", reported, tt.wantError)
			}
		})
	}
}

func TestIsLogCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "slog.Info call",
			code: `package main; import "log/slog"; func f() { slog.Info("test") }`,
			want: true,
		},
		{
			name: "slog.Error call",
			code: `package main; import "log/slog"; func f() { slog.Error("test") }`,
			want: true,
		},
		{
			name: "non-log call",
			code: `package main; func f() { println("test") }`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatal(err)
			}

			var found bool
			ast.Inspect(file, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					if isLogCall(call) {
						found = true
					}
				}
				return true
			})

			if found != tt.want {
				t.Errorf("isLogCall() = %v, want %v", found, tt.want)
			}
		})
	}
}
