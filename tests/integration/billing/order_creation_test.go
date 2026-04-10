package billing

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateOrderUsesPlanSnapshot(t *testing.T) {
	planServicePath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "billing", "plans", "service.go")
	if _, err := os.Stat(planServicePath); err != nil {
		t.Fatalf("缺少套餐服务文件 %s: %v", planServicePath, err)
	}

	orderServicePath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "billing", "orders", "service.go")
	content, err := os.ReadFile(orderServicePath)
	if err != nil {
		t.Fatalf("读取订单服务文件失败 %s: %v", orderServicePath, err)
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, orderServicePath, content, 0)
	if err != nil {
		t.Fatalf("解析订单服务文件失败 %s: %v", orderServicePath, err)
	}

	var found bool
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name == nil || typeSpec.Name.Name != "CreateOrderInput" {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("CreateOrderInput 不是结构体")
			}
			found = true
			fields := map[string]bool{}
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					fields[name.Name] = true
				}
			}
			for _, fieldName := range []string{"AccountID", "PlanPriceSnapshotID", "Amount", "Currency"} {
				if !fields[fieldName] {
					t.Fatalf("订单输入模型缺少字段 %s", fieldName)
				}
			}
		}
	}

	if !found {
		t.Fatalf("订单服务文件缺少 CreateOrderInput")
	}
}
