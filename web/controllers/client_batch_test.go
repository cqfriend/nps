package controllers

import (
	"testing"
)

func TestParseBatchClientLine(t *testing.T) {
	// 1. 测试空行
	if c, err := ParseBatchClientLine("   "); err != nil || c != nil {
		t.Fatalf("expected nil for empty line, got client: %v, err: %v", c, err)
	}

	// 2. 测试仅有备注的极简行
	c, err := ParseBatchClientLine("my-client-1")
	if err != nil {
		t.Fatalf("unexpected error for single remark: %v", err)
	}
	if c.Remark != "my-client-1" {
		t.Errorf("expected remark 'my-client-1', got '%s'", c.Remark)
	}
	if c.Cnf.U != "" || c.Cnf.P != "" {
		t.Errorf("expected empty basic auth user/pass, got user='%s', pass='%s'", c.Cnf.U, c.Cnf.P)
	}
	if c.VerifyKey != "" {
		t.Errorf("expected empty vkey (to be auto-generated later), got '%s'", c.VerifyKey)
	}
	if !c.ConfigConnAllow {
		t.Errorf("expected ConfigConnAllow default true")
	}

	// 3. 测试 备注 + Basic用户名 + Basic密码 (vkey留空自动生成)
	c, err = ParseBatchClientLine("my-client-2, user02, pass02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Remark != "my-client-2" || c.Cnf.U != "user02" || c.Cnf.P != "pass02" {
		t.Errorf("basic user/pass mismatch: remark=%s, user=%s, pass=%s", c.Remark, c.Cnf.U, c.Cnf.P)
	}
	if c.VerifyKey != "" {
		t.Errorf("expected empty vkey, got '%s'", c.VerifyKey)
	}

	// 4. 测试 备注 + Basic用户名 + Basic密码 + 唯一密钥
	c, err = ParseBatchClientLine("my-client-3, user03, pass03, secret_vkey_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Remark != "my-client-3" || c.Cnf.U != "user03" || c.Cnf.P != "pass03" || c.VerifyKey != "secret_vkey_123" {
		t.Errorf("fields mismatch: %+v", c)
	}

	// 5. 测试完整字段 (包含流量与连接数扩展项)
	fullLine := "office-pc, adminuser, mypass123, pc_vkey_888, 1024, 2048, 50, 10, 0"
	c, err = ParseBatchClientLine(fullLine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Remark != "office-pc" {
		t.Errorf("expected remark 'office-pc', got '%s'", c.Remark)
	}
	if c.Cnf.U != "adminuser" || c.Cnf.P != "mypass123" {
		t.Errorf("expected Basic auth 'adminuser'/'mypass123', got '%s'/'%s'", c.Cnf.U, c.Cnf.P)
	}
	if c.VerifyKey != "pc_vkey_888" {
		t.Errorf("expected vkey 'pc_vkey_888', got '%s'", c.VerifyKey)
	}
	if c.Flow.FlowLimit != 1024 {
		t.Errorf("expected FlowLimit 1024, got %d", c.Flow.FlowLimit)
	}
	if c.RateLimit != 2048 {
		t.Errorf("expected RateLimit 2048, got %d", c.RateLimit)
	}
	if c.MaxConn != 50 {
		t.Errorf("expected MaxConn 50, got %d", c.MaxConn)
	}
	if c.MaxTunnelNum != 10 {
		t.Errorf("expected MaxTunnelNum 10, got %d", c.MaxTunnelNum)
	}
	if c.ConfigConnAllow != false {
		t.Errorf("expected ConfigConnAllow false, got %v", c.ConfigConnAllow)
	}

	// 6. 测试 Tab 分隔
	tabLine := "tab-client\ttab_user\ttab_pass\ttab_vkey"
	c, err = ParseBatchClientLine(tabLine)
	if err != nil {
		t.Fatalf("unexpected error for tab line: %v", err)
	}
	if c.Remark != "tab-client" || c.Cnf.U != "tab_user" || c.Cnf.P != "tab_pass" || c.VerifyKey != "tab_vkey" {
		t.Errorf("tab parsing error: %+v", c)
	}
}

func TestIsBatchHeaderLine(t *testing.T) {
	if !isBatchHeaderLine("备注,Basic认证用户名,Basic认证密码,唯一验证密钥") {
		t.Errorf("expected true for chinese header")
	}
	if !isBatchHeaderLine("# remark, user, pass, vkey") {
		t.Errorf("expected true for english header")
	}
	if isBatchHeaderLine("normal-client-name, user1, pass1, vkey123") {
		t.Errorf("expected false for data line")
	}
}
