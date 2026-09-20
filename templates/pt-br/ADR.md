# ADR-000 · <decisão em uma frase>

<!-- Local sugerido: .harness/adr/ADR-000.md
Um ADR registra UMA decisão durável e o motivo. Não é design doc: para
propor e discutir, use RFC. Escreva só o suficiente para alguém daqui a
um ano entender por que foi assim e quando deve ser revisto. -->

| Campo | Valor |
|---|---|
| Status | proposto / aceito / substituído por ADR-000 / revogado |
| Data | AAAA-MM-DD |
| Decisor | <quem aprovou> |
| Origem | PRD-000 / ST-000 / RFC-000 / incidente RISKS.md#id |
| Reversibilidade | barata e local / cara ou com muitos dependentes |

## 1. Contexto

<!-- A situação que exige a decisão e as forças que realmente separam as opções (não liste tudo). Fato marcado como fato, hipótese como hipótese. -->

- Pergunta: <o que muda se for pelo outro caminho>
- Forças: <restrição técnica, dado, operação, segurança, prazo>

## 2. Opções consideradas

<!-- Inclua "manter como está" quando for viável. Preencha só as dimensões que discriminam. -->

| Dimensão | A · <nome> | B · <nome> | C · manter como está |
|---|---|---|---|
| Complexidade | | | |
| Custo de construir | | | |
| Custo de operar | | | |
| Impacto em dados | | | |
| Impacto em segurança | | | |
| Familiaridade do time | | | |

## 3. Decisão

- Escolhida: <opção>
- Motivo decisivo: <a razão que desempatou>
- Descartadas: <opção> porque <motivo específico>
- Preferência declarada: <fator que é gosto ou hábito, nomeado como tal> ou nenhuma

## 4. Consequências

- Fica mais fácil: <...>
- Fica mais difícil: <...>
- Dívida assumida: <o que aceitamos conviver e até quando>
- Contratos afetados: <componente A → B; dado que cruza; compatibilidade>

## 5. Validação

- Check que prova que funciona: <teste, métrica ou comando>
- Sinal de que está falhando: <o que observar em operação>

## 6. Revisar quando

<!-- Condição observável, não data vaga. Sem isto o ADR vira permanente por acidente. -->

- <ex.: volume passar de N, lib X lançar versão Y, incidente na área>

## 7. Referências

- <caminho relativo: RFC, PRD, incidente, benchmark>
