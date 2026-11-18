# PRD-902 · Cadastro somente por convite

| Campo | Valor |
|---|---|
| Status | em revisão |
| Dono | dono do workspace |
| Criado / atualizado | 2026-09-20 / 2026-09-20 |
| Histórias | ST-901, ST-902 |

## 1. Problema

O cadastro é aberto, então qualquer pessoa que descubra a URL cria uma conta num workspace em que nunca deveria entrar. Hoje o dono do workspace remove essas contas à mão, depois do fato.

## 2. Resultado esperado

- Sinal 1: contas criadas sem convite · hoje: sem limite · meta: 0.
- Sinal 2: convites pendentes que expiram sem ninguém notar · hoje: sem linha de base · meta: medida depois do primeiro mês.

## 3. Usuários e cenários

| Usuário | Cenário | Hoje | Depois |
|---|---|---|---|
| Admin | Traz um colega para o workspace | Envia a ele a URL de cadastro | Cria um convite vinculado ao e-mail da pessoa |
| Pessoa convidada | Cria uma conta | Cadastra-se diretamente | Resgata o código de convite |

## 4. Escopo

**Dentro do escopo**
- RF-01 Uma pessoa cria uma conta somente ao resgatar um código de convite válido; uma tentativa de cadastro sem código é recusada.
- RF-02 Somente um admin do workspace cria um convite, e o convite é vinculado a exatamente um endereço de e-mail.
- RF-03 Um convite expira 7 dias depois de criado e não pode mais ser resgatado.
- RF-04 O e-mail do convite é enviado uma vez na criação, e o admin o reenvia no máximo 3 vezes por convite.

**Fora do escopo**
- Preços e limites de assentos.
- Login único (SSO).
- Recompensas por indicação.

## 5. Critérios de aceite

- AC-01 (RF-01) Quando uma pessoa envia o formulário de cadastro sem código de convite, então a conta não é criada e a resposta informa que um convite é obrigatório.
- AC-02 (RF-02) Quando um membro que não é admin envia um convite, então o convite não é criado e a resposta informa que somente admins podem convidar. Quando um admin envia um convite para um endereço de e-mail, então o convite é vinculado a esse endereço e nenhum outro endereço pode resgatá-lo.
- AC-03 (RF-03) Quando um convite é resgatado mais de 7 dias depois do seu carimbo de criação, então o resgate é recusado e a resposta informa que o convite expirou.
- AC-04 (RF-04) Quando o admin usa o formulário de convite, então o formulário deve ser rápido.

## 6. Restrições e impacto técnico

- Dados: novo registro de convite com e-mail, código, carimbo de criação, estado e contagem de reenvios.
- Contratos: o endpoint de cadastro ganha um campo obrigatório de código de convite.
- Segurança: o código de convite é de uso único e vinculado a um endereço de e-mail.

## 7. Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Senha compartilhada do workspace | Não pode ser vinculada a uma pessoa nem revogada individualmente |

## 8. Riscos e decisões pendentes

- Risco: um admin revoga o convite errado · mitigação: o convite revogado aparece na lista de convites com o seu estado.

## 9. Referências

- discovery-notes.md
