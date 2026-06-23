package devtools

import (
	"strings"
	"testing"
)

func TestBeautifyHTML_Simple(t *testing.T) {
	input := `<html><head><title>Test</title></head><body><p>Hello</p></body></html>`
	result := BeautifyHTML(input)

	if !result.Success {
		t.Errorf("美化应成功: %s", result.Error)
	}

	output := result.Output
	// 检查输出包含缩进
	if !strings.Contains(output, "  <head>") {
		t.Errorf("应包含缩进标签:\n%s", output)
	}
	// 检查输出包含换行
	if !strings.Contains(output, "\n") {
		t.Errorf("应包含换行符:\n%s", output)
	}
}

func TestBeautifyHTML_SelfClosing(t *testing.T) {
	input := `<div><br/><img src="test.png"/><hr/></div>`
	result := BeautifyHTML(input)

	if !result.Success {
		t.Errorf("美化应成功: %s", result.Error)
	}

	// 自闭合标签应正确处理
	if !strings.Contains(result.Output, "<br/>") {
		t.Errorf("应保留自闭合标签:\n%s", result.Output)
	}
}

func TestBeautifyHTML_Empty(t *testing.T) {
	result := BeautifyHTML("")
	if !result.Success {
		t.Error("空输入应成功")
	}
	if result.Output != "" {
		t.Errorf("空输入应返回空输出: got %q", result.Output)
	}
}

func TestBeautifyCSS_Simple(t *testing.T) {
	input := `body{color:red;font-size:14px}`
	result := BeautifyCSS(input)

	if !result.Success {
		t.Errorf("美化应成功: %s", result.Error)
	}

	output := result.Output
	// 检查包含缩进
	if !strings.Contains(output, "  ") {
		t.Errorf("应包含缩进:\n%s", output)
	}
	// 检查包含换行
	if !strings.Contains(output, "\n") {
		t.Errorf("应包含换行:\n%s", output)
	}
}

func TestBeautifyCSS_Empty(t *testing.T) {
	result := BeautifyCSS("")
	if !result.Success {
		t.Error("空输入应成功")
	}
}

func TestBeautifyJS_Simple(t *testing.T) {
	input := `function test(){return{foo:"bar"}}`
	result := BeautifyJS(input)

	if !result.Success {
		t.Errorf("美化应成功: %s", result.Error)
	}

	output := result.Output
	// 检查包含缩进
	if !strings.Contains(output, "  ") {
		t.Errorf("应包含缩进:\n%s", output)
	}
}

func TestBeautifyJS_Empty(t *testing.T) {
	result := BeautifyJS("")
	if !result.Success {
		t.Error("空输入应成功")
	}
}

func TestBeautifySQL_Simple(t *testing.T) {
	input := `SELECT a,b,c FROM users WHERE id=1`
	result := BeautifySQL(input)

	if !result.Success {
		t.Errorf("美化应成功: %s", result.Error)
	}

	// 检查逗号后有换行
	if !strings.Contains(result.Output, ",\n") {
		t.Errorf("逗号后应有换行:\n%s", result.Output)
	}
}
