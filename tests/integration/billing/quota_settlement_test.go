package billing

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestSettleUsagePausesAccountWhenQuotaExhausted(t *testing.T) {
	summaryServicePath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "billing", "summary", "service.go")
	content, err := os.ReadFile(summaryServicePath)
	if err != nil {
		t.Fatalf("读取汇总结算服务文件失败 %s: %v", summaryServicePath, err)
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, summaryServicePath, content, 0)
	if err != nil {
		t.Fatalf("解析汇总结算服务文件失败 %s: %v", summaryServicePath, err)
	}

	var foundUsageDelta bool
	var foundSettle bool
	for _, decl := range file.Decls {
		switch node := decl.(type) {
		case *ast.GenDecl:
			if node.Tok != token.TYPE {
				continue
			}
			for _, spec := range node.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Name == nil || typeSpec.Name.Name != "UsageDelta" {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					t.Fatalf("UsageDelta 不是结构体")
				}
				foundUsageDelta = true
				fields := map[string]bool{}
				for _, field := range structType.Fields.List {
					for _, name := range field.Names {
						fields[name.Name] = true
					}
				}
				for _, fieldName := range []string{"AccountID", "TotalTokens"} {
					if !fields[fieldName] {
						t.Fatalf("UsageDelta 缺少字段 %s", fieldName)
					}
				}
			}
		case *ast.FuncDecl:
			if node.Name == nil || node.Name.Name != "Settle" {
				continue
			}
			if node.Type == nil || node.Type.Params == nil || len(node.Type.Params.List) != 1 {
				t.Fatalf("Settle 签名不正确")
			}
			param := node.Type.Params.List[0]
			if len(param.Names) != 1 || param.Names[0].Name != "delta" {
				t.Fatalf("Settle 参数名不正确")
			}
			ident, ok := param.Type.(*ast.Ident)
			if !ok || ident.Name != "UsageDelta" {
				t.Fatalf("Settle 参数类型不正确")
			}
			if node.Type.Results == nil || len(node.Type.Results.List) != 1 {
				t.Fatalf("Settle 返回值不正确")
			}
			resultIdent, ok := node.Type.Results.List[0].Type.(*ast.Ident)
			if !ok || resultIdent.Name != "error" {
				t.Fatalf("Settle 返回类型不正确")
			}
			foundSettle = true
		}
	}

	if !foundUsageDelta {
		t.Fatalf("汇总结算服务文件缺少 UsageDelta")
	}
	if !foundSettle {
		t.Fatalf("汇总结算服务文件缺少 Settle")
	}
}
