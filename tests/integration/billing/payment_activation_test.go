package billing

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestPaymentActivationSkeletonExists(t *testing.T) {
	paymentServicePath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "billing", "payments", "service.go")
	content, err := os.ReadFile(paymentServicePath)
	if err != nil {
		t.Fatalf("读取支付服务文件失败 %s: %v", paymentServicePath, err)
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, paymentServicePath, content, 0)
	if err != nil {
		t.Fatalf("解析支付服务文件失败 %s: %v", paymentServicePath, err)
	}

	var foundPaymentEvent bool
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name == nil || typeSpec.Name.Name != "PaymentEvent" {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("PaymentEvent 不是结构体")
			}
			foundPaymentEvent = true
			fields := map[string]bool{}
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					fields[name.Name] = true
				}
			}
			for _, fieldName := range []string{"ProviderEventID", "OrderID", "EventType"} {
				if !fields[fieldName] {
					t.Fatalf("支付事件模型缺少字段 %s", fieldName)
				}
			}
		}
	}

	if !foundPaymentEvent {
		t.Fatalf("支付服务文件缺少 PaymentEvent")
	}

	subscriptionServicePath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "billing", "subscriptions", "service.go")
	content, err = os.ReadFile(subscriptionServicePath)
	if err != nil {
		t.Fatalf("读取订阅服务文件失败 %s: %v", subscriptionServicePath, err)
	}

	file, err = parser.ParseFile(fileSet, subscriptionServicePath, content, 0)
	if err != nil {
		t.Fatalf("解析订阅服务文件失败 %s: %v", subscriptionServicePath, err)
	}

	var foundActivate bool
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name == nil || funcDecl.Name.Name != "ActivateFromPaidOrder" {
			continue
		}
		if funcDecl.Type == nil || funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) != 1 {
			t.Fatalf("ActivateFromPaidOrder 签名不正确")
		}
		param := funcDecl.Type.Params.List[0]
		if len(param.Names) != 1 || param.Names[0].Name != "orderID" {
			t.Fatalf("ActivateFromPaidOrder 参数名不正确")
		}
		ident, ok := param.Type.(*ast.Ident)
		if !ok || ident.Name != "string" {
			t.Fatalf("ActivateFromPaidOrder 参数类型不正确")
		}
		if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) != 1 {
			t.Fatalf("ActivateFromPaidOrder 返回值不正确")
		}
		resultIdent, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident)
		if !ok || resultIdent.Name != "error" {
			t.Fatalf("ActivateFromPaidOrder 返回类型不正确")
		}
		foundActivate = true
	}

	if !foundActivate {
		t.Fatalf("订阅服务文件缺少 ActivateFromPaidOrder")
	}
}
