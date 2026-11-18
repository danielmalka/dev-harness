# Resumo · slice-01

Rejeitar um nome em branco em `greet`.

## Comportamento atual

`greet(name)` devolve `"Hello, {name}"` para qualquer string, inclusive valores só com espaços em branco.

## Comportamento desejado

Quando `name` é vazio ou só tem espaços em branco, `greet` levanta `ValueError` com uma mensagem que menciona `name`. Nomes que não estão em branco mantêm a saudação atual.

## Aceite

| ID | Quando | Então |
| --- | --- | --- |
| AC-01 | `greet("Ada")` | devolve `Hello, Ada` |
| AC-02 | `greet("  ")` ou `greet("")` | levanta `ValueError` |

## Escopo

- Conjunto de escrita: `greet.py`. Os testes já estão em `test_greet.py`.
- Fora: novas dependências, funções públicas extras, renomear arquivos.

## Checagens

```bash
python3 -m unittest test_greet.py
```

Hoje o AC-01 passa e o AC-02 falha. Depois da fatia, os dois passam. Checagem de tipos não se aplica.
