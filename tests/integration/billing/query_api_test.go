package billing

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestBillingQueryReturnsAccountStateAndSummary(t *testing.T) {
	handlerPath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "billing", "http", "handler.go")
	content, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("读取账务查询处理器文件失败 %s: %v", handlerPath, err)
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, handlerPath, content, 0)
	if err != nil {
		t.Fatalf("解析账务查询处理器文件失败 %s: %v", handlerPath, err)
	}

	var foundAccountQuotaView bool
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name == nil || typeSpec.Name.Name != "AccountQuotaView" {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("AccountQuotaView 不是结构体")
			}
			foundAccountQuotaView = true
			fields := map[string]bool{}
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					fields[name.Name] = true
				}
			}
			for _, fieldName := range []string{"Status", "QuotaRemaining"} {
				if !fields[fieldName] {
					t.Fatalf("AccountQuotaView 缺少字段 %s", fieldName)
				}
			}
		}
	}

	if !foundAccountQuotaView {
		t.Fatalf("账务查询处理器文件缺少 AccountQuotaView")
	}
}
