package harness

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// The ceilings of PRD-035 R5. The owner revises them after 30 days of use (D8), so they live in one place.
const (
	FactsMaxLines       = 150
	FactsMaxBytes       = 12 * 1024
	FactsMaxPerCategory = 15
	FactsMaxFactRunes   = 300
)

// FactCategories are the 11 sections of FACTS.md, in the order PRD-035 R4 fixes. Headings match case-insensitively.
var FactCategories = []string{
	"estilo de escrita", "arquitetura", "design system", "dependências", "ferramentas complementares", "particularidades de código",
	"testes", "dados e esquema", "build e CI", "armadilhas", "fronteiras",
}

const (
	factsCandidates = "candidatos"
	factsChanges    = "mudanças"
)

// FactStates are the values of the `Estado:` field.
var FactStates = []string{"candidato", "confirmado", "hipótese"}

// FactProblem is one violation: the 1-based line (0 for the whole file) and what is wrong.
type FactProblem struct {
	Line int
	Msg  string
}

// FactsReport is what CheckFacts found. Facts counts every fact line (categories and `## Candidatos`); Candidates counts the ones in state `candidato`.
type FactsReport struct {
	Facts      int
	Candidates int
	Problems   []FactProblem
}

// A fact is one line (PRD-035 R3): "- F-001 · <frase> · Evidência: <caminho:linha>[, <caminho:linha>...] · Visto: <AAAA-MM-DD>, <papel>, <tarefa> · sha: <sha curto> · Estado: <estado>".
// The separator is " · " (space, middle dot, space). `sha` also accepts the form without the colon.
var (
	factLine   = regexp.MustCompile(`^- (F-[0-9A-Za-z]*)\b`)
	factID     = regexp.MustCompile(`^F-[0-9]{3}$`)
	factSHA    = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
	shaLabel   = regexp.MustCompile(`^sha:?\s+(.*)$`)
	evidenceRe = regexp.MustCompile(`^(.+):([0-9]+)$`)
)

const factSep = " · "

// secretPatterns is the fixed list of secret and personal-data shapes `dh facts check` reports (PRD-035 R6). It is the list the curator redacts with
// `[redigido]`. The match is never printed: the report names the kind and the line.
var secretPatterns = []struct {
	kind  string
	re    *regexp.Regexp
	check func(match []string) bool
}{
	{"chave privada", regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`), nil},
	{"chave de acesso da AWS", regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`), nil},
	{"token do GitHub", regexp.MustCompile(`\b(?:gh[pousr]_[A-Za-z0-9]{20,}|github_pat_[A-Za-z0-9_]{20,})\b`), nil},
	{"token do Slack", regexp.MustCompile(`\bxox[abprs]-[A-Za-z0-9-]{10,}`), nil},
	{"chave de API (sk-)", regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{20,}`), nil},
	{"chave de API do Google", regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}\b`), nil},
	{"JWT", regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`), nil},
	{"credencial em URL", regexp.MustCompile(`[A-Za-z][A-Za-z0-9+.-]*://[^\s/:@]+:[^\s/@]+@`), nil},
	{"valor de segredo atribuído", regexp.MustCompile(`(?i)\b(?:password|passwd|pwd|senha|secret|segredo|token|api[_-]?key|apikey|access[_-]?key|client[_-]?secret|private[_-]?key)\w*\s*[:=]\s*["']?([A-Za-z0-9+/=_.\-]{8,})`), looksLikeValue},
	{"e-mail", regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9-]+(?:\.[A-Za-z0-9-]+)*\.[A-Za-z]{2,}`), nil},
	{"CPF", regexp.MustCompile(`\b[0-9]{3}\.[0-9]{3}\.[0-9]{3}-[0-9]{2}\b`), nil},
	{"CNPJ", regexp.MustCompile(`\b[0-9]{2}\.[0-9]{3}\.[0-9]{3}/[0-9]{4}-[0-9]{2}\b`), nil},
}

// looksLikeValue keeps an assignment like `token: abc12345` and drops a prose word like `senha: bcrypt`: a value needs a letter and a digit.
func looksLikeValue(m []string) bool {
	v := m[1]
	letter, digit := false, false
	for _, r := range v {
		switch {
		case r >= '0' && r <= '9':
			digit = true
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			letter = true
		}
	}
	return letter && digit
}

// CheckFacts checks the content of a FACTS.md against PRD-035 R3 to R6 and, through exists, that every evidence path is in the repository.
// exists receives a slash-separated path relative to the repository root and says whether it is a file there (a folder is not evidence). It writes nothing.
func CheckFacts(content []byte, exists func(rel string) bool) FactsReport {
	var rep FactsReport
	add := func(line int, format string, a ...any) {
		rep.Problems = append(rep.Problems, FactProblem{line, fmt.Sprintf(format, a...)})
	}

	text := string(content)
	if len(content) > FactsMaxBytes {
		add(0, "teto de tamanho excedido: %d bytes, o teto é %d (12 KB)", len(content), FactsMaxBytes)
	}
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if text == "" {
		lines = nil
	}
	if len(lines) > FactsMaxLines {
		add(0, "teto de linhas excedido: %d linhas, o teto é %d", len(lines), FactsMaxLines)
	}

	section := "" // normalized heading of the current ## section; "" before the first one
	perCategory := map[string]int{}
	seenSection := map[string]int{}
	lastCategory := -1
	ids := map[string]int{}
	for i, raw := range lines {
		n := i + 1
		line := strings.TrimRight(raw, "\r")
		for _, p := range secretPatterns {
			if m := p.re.FindStringSubmatch(line); m != nil && (p.check == nil || p.check(m)) {
				add(n, "padrão de segredo ou dado pessoal (%s): troque por [redigido]", p.kind)
				break
			}
		}
		if strings.HasPrefix(line, "## ") {
			section = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, "## ")))
			if first, dup := seenSection[section]; dup {
				add(n, "seção repetida %q (a primeira está na linha %d)", strings.TrimSpace(strings.TrimPrefix(line, "## ")), first)
			}
			seenSection[section] = n
			if section == factsCandidates || section == factsChanges {
				continue
			}
			idx := categoryIndex(section)
			switch {
			case idx < 0:
				add(n, "seção %q não é uma das categorias (%s), nem Candidatos nem Mudanças", strings.TrimSpace(strings.TrimPrefix(line, "## ")), strings.Join(FactCategories, ", "))
			case idx < lastCategory:
				add(n, "categoria %q fora de ordem: as categorias seguem a ordem %s", strings.TrimSpace(strings.TrimPrefix(line, "## ")), strings.Join(FactCategories, ", "))
			default:
				lastCategory = idx
			}
			continue
		}
		inFacts := section == factsCandidates || categoryIndex(section) >= 0
		if !inFacts || strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(strings.TrimSpace(line), "<!--") {
			continue
		}
		if !factLine.MatchString(line) {
			add(n, "linha que não é um fato: esperado `- F-<n> · <frase> · Evidência: … · Visto: … · sha: … · Estado: …`")
			continue
		}
		rep.Facts++
		if section != factsCandidates {
			perCategory[section]++
			if perCategory[section] == FactsMaxPerCategory+1 {
				add(n, "teto de fatos por categoria excedido em %q: o teto é %d", section, FactsMaxPerCategory)
			}
		}
		if utf8.RuneCountInString(line) > FactsMaxFactRunes {
			add(n, "teto de caracteres por fato excedido: %d, o teto é %d", utf8.RuneCountInString(line), FactsMaxFactRunes)
		}
		state := checkFactLine(line, n, ids, exists, add)
		if state == "candidato" {
			rep.Candidates++
		}
	}
	return rep
}

func categoryIndex(section string) int {
	for i, c := range FactCategories {
		if strings.ToLower(c) == section {
			return i
		}
	}
	return -1
}

// checkFactLine checks R3 on one fact line and returns its Estado ("" when missing or invalid).
func checkFactLine(line string, n int, ids map[string]int, exists func(string) bool, add func(int, string, ...any)) string {
	segs := strings.Split(strings.TrimPrefix(line, "- "), factSep)
	id := strings.TrimSpace(segs[0])
	if !factID.MatchString(id) {
		add(n, "campo id: %q não tem o formato F-<n> com 3 dígitos", id)
	} else if first, dup := ids[id]; dup {
		add(n, "campo id: %s repetido (o primeiro está na linha %d)", id, first)
	} else {
		ids[id] = n
	}

	var phrase []string
	var evidence, seen, shaVal, state string
	var haveEvidence, haveSeen, haveSHA, haveState bool
	order := []string{}
	for _, s := range segs[1:] {
		s = strings.TrimSpace(s)
		switch {
		case strings.HasPrefix(s, "Evidência:") && !haveEvidence:
			haveEvidence, evidence = true, strings.TrimSpace(strings.TrimPrefix(s, "Evidência:"))
			order = append(order, "Evidência")
		case strings.HasPrefix(s, "Visto:") && !haveSeen:
			haveSeen, seen = true, strings.TrimSpace(strings.TrimPrefix(s, "Visto:"))
			order = append(order, "Visto")
		case shaLabel.MatchString(s) && !haveSHA:
			haveSHA, shaVal = true, strings.TrimSpace(shaLabel.FindStringSubmatch(s)[1])
			order = append(order, "sha")
		case strings.HasPrefix(s, "Estado:") && !haveState:
			haveState, state = true, strings.TrimSpace(strings.TrimPrefix(s, "Estado:"))
			order = append(order, "Estado")
		case len(order) == 0:
			phrase = append(phrase, s)
		default:
			add(n, "texto fora dos campos depois de %s", order[len(order)-1])
		}
	}
	if strings.TrimSpace(strings.Join(phrase, factSep)) == "" {
		add(n, "campo frase ausente (%s)", id)
	}
	if !haveEvidence {
		add(n, "campo Evidência: ausente (%s)", id)
	}
	if !haveSeen {
		add(n, "campo Visto: ausente (%s)", id)
	}
	if !haveSHA {
		add(n, "campo sha ausente (%s)", id)
	}
	if !haveState {
		add(n, "campo Estado: ausente (%s)", id)
	}
	if want := []string{"Evidência", "Visto", "sha", "Estado"}; haveEvidence && haveSeen && haveSHA && haveState && strings.Join(order, ",") != strings.Join(want, ",") {
		add(n, "campos fora de ordem (%s): a ordem é id, frase, Evidência, Visto, sha, Estado", id)
	}
	if haveEvidence {
		checkEvidence(evidence, n, id, exists, add)
	}
	if haveSeen {
		parts := strings.Split(seen, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		if len(parts) != 3 || parts[1] == "" || parts[2] == "" {
			add(n, "campo Visto: deve ser <AAAA-MM-DD>, <papel>, <tarefa> (%s)", id)
		} else if _, err := time.Parse("2006-01-02", parts[0]); err != nil {
			add(n, "campo Visto: tem a data %q fora do formato AAAA-MM-DD (%s)", parts[0], id)
		}
	}
	if haveSHA && !factSHA.MatchString(shaVal) {
		add(n, "campo sha: %q não é um sha curto (7 a 40 hexadecimais minúsculos) (%s)", shaVal, id)
	}
	if haveState {
		ok := false
		for _, s := range FactStates {
			ok = ok || state == s
		}
		if !ok {
			add(n, "campo Estado: %q deve ser %s (%s)", state, strings.Join(FactStates, ", "), id)
		}
		if ok {
			return state
		}
	}
	return ""
}

func checkEvidence(evidence string, n int, id string, exists func(string) bool, add func(int, string, ...any)) {
	items := strings.Split(evidence, ",")
	if evidence == "" {
		add(n, "campo Evidência: deve ter de 1 a 3 caminhos caminho:linha, não tem nenhum (%s)", id)
		return
	}
	if len(items) > 3 {
		add(n, "campo Evidência: deve ter de 1 a 3 caminhos caminho:linha, tem %d (%s)", len(items), id)
	}
	for _, it := range items {
		it = strings.TrimSpace(it)
		m := evidenceRe.FindStringSubmatch(it)
		if m == nil || m[2] == "0" || strings.TrimLeft(m[2], "0") == "" {
			add(n, "campo Evidência: %q não é caminho:linha com linha maior que zero (%s)", it, id)
			continue
		}
		path := m[1]
		if strings.HasPrefix(path, "/") || strings.Contains(path, `\`) || strings.Contains(path, ":") || hasDotDot(path) {
			add(n, "campo Evidência: o caminho %q deve ser relativo ao repositório, sem .. nem barra invertida (%s)", path, id)
			continue
		}
		if exists != nil && !exists(path) {
			add(n, "campo Evidência: o caminho %q não é um arquivo do repositório (não existe ou é uma pasta) (%s)", path, id)
		}
	}
}

func hasDotDot(p string) bool {
	for _, seg := range strings.Split(p, "/") {
		if seg == ".." {
			return true
		}
	}
	return false
}
