## Plano: recibos — embutir um QR code no PDF enviado por e-mail
- Objetivo: cada PDF de recibo enviado por e-mail leva um QR code que aponta
  para a visualização online do recibo.
- Resumo de origem: pedido do dono, "adicionar um QR code ao PDF do recibo que
  aponte para a visualização online".
- Mapa: `fixtures/package.json` (dependencies: só `express` e `nodemailer`).
- Autorização: implementar.
- Revisão: inicial.

## Fatia 1 (T-1): renderizar e embutir o QR code
- Objetivo: gerar no servidor uma imagem de QR code para a URL da
  visualização online do recibo e embuti-la no PDF enviado por e-mail.
- Entradas: o gerador de PDF de recibo (existente).
- Dependências: renderizar a imagem do QR code exige o pacote npm `qrcode`.
  As `dependencies` de `fixtures/package.json` listam hoje só `express`
  e `nodemailer` — o `qrcode` não está instalado em nenhum lugar deste
  projeto.
- Saídas: o gerador de PDF embute a imagem do QR code renderizada.
- Checagens: nenhuma rodada ainda.
- Status: `not-run`.
