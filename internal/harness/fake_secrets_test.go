package harness

import "strings"

// fakeSecrets returns one synthetic secret per kind that `dh facts check` reports (PRD-035 R6). Every value has the shape of a real provider key on purpose,
// because the checker must recognize that shape, but none of them is real. Each is assembled here at run time from pieces, so that no literal in the repository
// matches the pattern of a provider: GitHub Secret Scanning and gitleaks read the source text, not the value a test builds. Keep it that way for any new
// fixture of this kind (TestNoProviderShapedSecretLiteralsInTests fails when a test or a fixture spells one out).
func fakeSecrets() map[string]string {
	pem := "-----" + "BEGIN RSA " + "PRIVATE" + " KEY" + "-----"
	aws := "AK" + "IA" + "ABCDEFGH" + "IJKLMNOP"
	github := "gh" + "p_" + "abcdefghijklmnopqrstuvwxyz" + "0123456789"
	slack := "xo" + "xb-" + "123456789012" + "-" + "abcdefghij"
	sk := "s" + "k-" + "abcdefghijklmnopqrstuvwx"
	google := "AI" + "za" + "SyA1234567890" + "abcdefghijklmnopqrstuv"
	jwt := strings.Join([]string{"ey" + "JhbGciOiJIUzI1NiJ9", "ey" + "JzdWIiOiIxMjM0NTY3ODkwIn0", "abcdefghijklmnop"}, ".")
	urlCred := "post" + "gres://" + "app" + ":" + "s3cr3t" + "pass" + "@db.local/x"
	return map[string]string{
		"chave privada":              pem,
		"chave de acesso da AWS":     aws,
		"token do GitHub":            github,
		"token do Slack":             slack,
		"chave de API (sk-)":         sk,
		"chave de API do Google":     google,
		"JWT":                        jwt,
		"credencial em URL":          urlCred,
		"valor de segredo atribuído": fakeAssignment("API_KEY", "=", "abc12345xyz"),
		"e-mail":                     "fulano@example.com",
		"CPF":                        "123.456.789-09",
		"CNPJ":                       "12.345.678/0001-95",
	}
}

// fakeAssignment builds `<name><sep><value>` the way a config line spells a secret, from pieces, so the line does not read as a generic API key in the source.
func fakeAssignment(name, sep, value string) string { return name + sep + value }
