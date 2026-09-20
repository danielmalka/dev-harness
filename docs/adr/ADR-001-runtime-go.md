# ADR-001 · Scripts do kit em Go, num único binário `dh` com builds pré-compilados no pacote

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-09-20 |
| Decisor | Daniel Lemos |
| Origem | Grill de 20/09/2026 sobre a extensão VS Code; plano §4.1 |
| Reversibilidade | cara: toca empacotador, validador, doctor, hooks e tutoriais |

## 1. Contexto

- Pergunta: em que linguagem ficam os scripts que o plugin executa na máquina do consumidor (status line, hooks de estado, doctor) e os de manutenção do kit (validate, build)?
- Forças: a status line e os hooks rodam a cada evento e precisam responder em milissegundos; o consumidor não deve precisar instalar interpretador para o pacote funcionar; o dono mantém o `ck` do Context Kit em Go e tem mais familiaridade com Go do que com Python; o plano proíbe dependência oculta da máquina do autor e download de conteúdo em execução; a extensão VS Code é obrigatoriamente TypeScript.

## 2. Opções consideradas

| Dimensão | A · Go, binário único `dh`, builds commitados em `dist/` | B · manter Python 3 + bash | C · Go via `go run` no consumidor |
|---|---|---|---|
| Complexidade | módulo Go, CI de cross-compile, wrapper por plataforma | nenhuma nova | wrapper + toolchain no consumidor |
| Custo de construir | migrar validate, build, doctor; escrever snapshot | só escrever snapshot | igual a A sem CI |
| Custo de operar | binários de ~3 MB por alvo, cinco alvos, crescem o git a cada release | requer Python 3 no consumidor; startup de ~50 ms por chamada | primeira chamada compila (segundos) e trava a status line; exige Go instalado |
| Impacto em dados | nenhum | nenhum | nenhum |
| Impacto em segurança | binário commitado precisa de build reproduzível e checksum | scripts legíveis | idem A |
| Familiaridade do time | alta (Go) | média | alta |

## 3. Decisão

- Escolhida: A.
- Motivo decisivo: latência da status line e ausência de interpretador no consumidor, com a promessa "clonou, funcionou" preservada; familiaridade do dono pesa na manutenção.
- Descartadas: B, porque exige Python no consumidor e o dono não quer manter Python; C, porque a primeira execução compila e exige Go instalado, contrariando a portabilidade da v1.
- Preferência declarada: familiaridade com Go é preferência do dono, não requisito técnico; o requisito técnico é latência e ausência de dependência.

## 4. Consequências

- Fica mais fácil: distribuir o plugin pelo marketplace sem pré-requisito além do Claude Code; responder na status line sem atraso; um único ponto de manutenção para validar, empacotar, diagnosticar e observar.
- Fica mais difícil: cada release cross-compila cinco alvos e commita binários; o repositório cresce; PR com binário exige build reproduzível na CI e checksum no manifesto.
- Dívida assumida: até a migração pousar, `scripts/validate.py`, `scripts/build-claude-code.sh` e `scripts/doctor.sh` continuam valendo e Python 3 continua requisito declarado da fixture `slice-01` e do stamp HTML.
- Contratos afetados: `harness-manifest.json` ganha inventário de binários com checksum; `dist/` ganha `bin/<os>_<arch>/dh` e `bin/dh` (wrapper); tutoriais trocam `python3 scripts/validate.py` por `dh validate`.

## 5. Validação

- Check que prova que funciona: `dh validate .` e `dh build` reproduzem o pacote byte a byte a partir das fontes em Linux, macOS e WSL2; `dh snapshot` responde em menos de 50 ms num tick da status line.
- Sinal de que está falhando: `doctor` sem binário para a plataforma; status line travando; `dist/` divergindo das fontes após rebuild.

## 6. Revisar quando

- O tamanho de `dist/` passar de 50 MB, ou o Claude Code passar a distribuir binários de plugin por outro mecanismo, ou Windows nativo entrar na matriz e exigir assinatura de binário.

## 7. Referências

- `docs/plano-produto.html` §4.1 (parágrafo de scripts de manutenção)
- `docs/prd/PRD-001-etapa-1.md`
- `CLAUDE.md` do repositório
