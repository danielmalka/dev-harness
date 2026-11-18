---
name: consolidate-memory
description: Arquive a memória em EPOCHAL.md e depois reduza-a
author: malka
argument-hint: ""
metadata:
  roles: [coordinator]
  skills: [context-handoff]
  writes: somente .harness/MEMORY.md e .harness/EPOCHAL.md
disable-model-invocation: true
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Rode isto você mesmo; não despache nenhum especialista, porque estes registros têm um único escritor. Siga a seção "Consolidação a pedido explícito" abaixo, com a skill do kit `context-handoff` carregada para as regras ao redor.

## Pré-requisitos
Nenhuma pergunta sobre outros escritores: a trava de escrita da memória (`context-handoff`, `dh memory hash` e `dh memory write --expect`) recusa a poda quando `.harness/MEMORY.md` mudou desde que foi lido. `.harness/MEMORY.md` deve existir e ter conteúdo operacional; pule o arquivamento e diga isso quando não tiver. Se uma execução anterior foi interrompida depois do arquivamento, localize e verifique esse lote antes de continuar.

## Saída
Um lote acrescentado a `.harness/EPOCHAL.md`: uma cópia literal, claramente delimitada, de `.harness/MEMORY.md` com o id do lote `<AAAA-MM-DD>-<n>` e um timestamp ISO 8601 com deslocamento de fuso, datas e texto originais preservados, qualquer comentário datado mantido fora do bloco bruto. Releia o lote armazenado e compare-o com o original. Só depois dessa verificação, reduza `.harness/MEMORY.md` ao trabalho atual, decisões, bloqueios, evidência necessária, próximo passo e um breve histórico recente, e registre o id do lote e a hora da consolidação. Reporte o lote e o estado ativo restante.

## Limites
- Em qualquer falha ou divergência, deixe `.harness/MEMORY.md` intacto e reporte a operação como incompleta.
- Nunca acrescente de novo um lote já confirmado como armazenado, e nunca remova informação que não foi preservada.
- Deixe `.harness/RISKS.md` intacto.
- Não afirme atomicidade entre os dois arquivos, e nunca consolide em silêncio só porque a memória cresceu. Quando `MEMORY.md` passa de 40000 bytes, o Coordenador oferece este comando uma vez por sessão; a oferta é o único passo automático, e a consolidação só roda quando o dono pede.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`, e `MEMORY.md` só pela trava de escrita da memória.

## Consolidação a pedido explícito
Rode `consolidate-memory` por `/dh:consolidate-memory` só quando o dono pedir para arquivar.
1. Tome `dh memory hash` e leia `MEMORY.md`; pule o arquivamento quando não houver conteúdo operacional. Escolha um id de lote estável com timestamp e fuso e acrescente a `EPOCHAL.md` uma cópia literal delimitada de `MEMORY.md` (comentários datados fora do bloco bruto).
2. Releia o lote, compare-o com o original e tome `dh memory hash` de novo, igual ao primeiro; em caso de diferença, preserve `MEMORY.md` e reporte. Só então reduza `MEMORY.md` com `dh memory write --expect <o primeiro hash>` (saída 3: `MEMORY.md` intacto, a operação reportada como incompleta com o id do lote armazenado, uma nova execução começa um novo lote), mantendo o trabalho atual, decisões, bloqueios, evidência necessária, próximo passo e um breve histórico recente, e registre o id do lote e a hora.
3. Depois de uma interrupção, verifique esse lote antes de continuar; nunca acrescente duas vezes um lote confirmado nem descarte informação não preservada. Deixe `RISKS.md` intacto, reporte o lote e o estado restante, não afirme atomicidade entre os dois arquivos.

## Próximo
`/dh:resume` para continuar o trabalho autorizado restante a partir da memória podada.
