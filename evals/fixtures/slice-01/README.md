# Fixture slice-01

Projeto minúsculo em Python 3 com biblioteca padrão, usado como a primeira entrega guiada. Não há pacotes para instalar.

A árvore inicial é intencionalmente incompleta: `greet("  ")` e `greet("")` deveriam levantar `ValueError` e hoje não levantam. `python3 -m unittest test_greet.py` fica vermelho (3 testes, 2 falhas) até a fatia autorizada ser implementada. O CI confere que ele falha exatamente assim.

## Estrutura

| Arquivo | Papel |
| --- | --- |
| `BRIEF.md` | Pedido autorizado e aceite. |
| `greet.py` | Código sob teste. |
| `test_greet.py` | Testes de aceite. O caminho feliz passa. Os casos de nome em branco e de nome vazio falham. |
| `expected/greet.py` | Resultado conhecido como correto. Não o copie para o lugar durante o exercício. |
| `expected/unittest.txt` | Saída esperada do unittest depois da fatia. |

## Conclusão

1. Carregue o kit contra este diretório (`docs/inicio-rapido.html`).
2. `/dh:setup` e depois `/dh:plan` usando `BRIEF.md`.
3. `/dh:build` da fatia única, depois `/dh:verify` e `/dh:review`.
4. `python3 -m unittest test_greet.py` bate com `expected/unittest.txt`.
5. `/dh:handoff`.

Python 3 é necessário só para esta fixture. O kit em si traz o binário `dh` e não precisa de interpretador; a fixture fica em Python para que uma máquina limpa não precise de Node, Go nem de navegador.
