# Checklist de fronteira de confiança

Use depois do passo 1 do procedimento, quando as fronteiras e as exposições em escopo já são conhecidas. Cada linha nomeia o que procurar e a evidência necessária para classificar. Vulnerabilidades demonstradas exigem um caminho alcançável sem mitigação; exposições confirmadas exigem evidência de exposição, sem exigir uma cadeia de exploração na aplicação. Mantenha a evidência incompleta qualificada como hipótese. Uma linha que não se aplica à mudança é pulada, não respondida.

## Identidade e acesso

| Área | O que procurar | O que rastrear antes de reportar |
| --- | --- | --- |
| Autenticação | Pontos de entrada novos ou movidos que pulam a checagem de sessão; validação de token que ignora expiração, audiência ou assinatura; credenciais comparadas sem função de tempo constante | A rota ou o handler exato que chega a trabalho sensível sem identidade verificada |
| Autorização | Uma guarda de rota sem checagem em nível de objeto; um papel lido de um valor fornecido pelo cliente; uma permissão calculada depois que o dado já foi carregado e devolvido; uma checagem que não falha fechada diante de identidade ausente ou parcial; chaves de cache ou de índice sem escopo completo de tenant, usuário e objeto | Qual chamador pode nomear qual identificador de objeto, e onde a posse é imposta |
| Sessão | Tokens em URLs ou logs; falta de rotação após mudança de privilégio; cookies sem as flags que o transporte do projeto exige | Como um valor de sessão roubado ou fixado chega a uma ação autenticada |
| Fronteiras de privilégio | Endpoints internos alcançáveis de fora; flags administrativas derivadas de entrada do usuário; personificação sem registro de auditoria | O caminho mais curto de uma posição sem privilégio até a operação privilegiada |

## Entrada e interpretação

| Área | O que procurar | O que rastrear antes de reportar |
| --- | --- | --- |
| Construção de consulta | Concatenação ou interpolação de strings em SQL ou em outra linguagem de consulta; um identificador dinâmico, como tabela, coluna ou campo de ordenação | Se todo caminho até o destino usa parâmetros, e o que o ramo sem parâmetros aceita |
| Execução de comando | Chamada de shell com argumentos interpolados; um shell usado onde uma chamada direta de processo bastaria; argumentos montados a partir de configuração que um usuário pode definir | Quais caracteres sobrevivem até o shell e o que eles podem iniciar |
| Tratamento de caminho | Nomes de arquivo vindos de requisições ou de arquivos compactados juntados a um diretório base; sequências de travessia; links simbólicos seguidos na extração | Se o caminho resolvido é confirmado como contido na raiz pretendida depois da normalização |
| Desserialização | Grafos de objetos reconstruídos a partir de bytes não confiáveis; formatos que podem nomear um tipo; YAML inseguro ou equivalentes de pickle | Quais tipos o payload pode instanciar e o que a construção deles executa |
| Renderização de template | Dado do usuário passado como fonte do template em vez de dado do template; escape desligado por conveniência | Se o escape fica ligado por padrão no caminho que renderiza o valor não confiável |
| Fronteira de modelo e ferramenta | Texto não confiável que chega a um modelo, ou a uma ferramenta com autoridade, é uma fronteira de confiança; política de ferramenta imposta dentro do prompt em vez de fora do modelo; texto recuperado ou saída de ferramenta tratados como instrução em vez de dado | Onde a política de ferramenta é imposta e se o modelo pode alterá-la |
| Codificação de saída | Valores escritos em HTML, em JSON embutido em HTML, em atributos ou em URLs sem a codificação que aquele contexto exige | O contexto de renderização exato e o codificador aplicado ali |

## Dados e segredos

| Área | O que procurar | O que rastrear antes de reportar |
| --- | --- | --- |
| Segredos na árvore | Chaves, tokens, senhas, strings de conexão em código-fonte, fixtures ou configuração; segredos em texto de erro. Varrer o histórico do controle de versão exige um shell; quando indisponível, registre o histórico como não inspecionado | O arquivo e a linha. Reporte o local e o formato; nunca o valor. Um segredo vazado é rotacionado ou revogado fora do git antes de ser descartado ou colocado na linha de base; uma linha de base guarda só impressões digitais exatas e revisadas, nunca uma exclusão ampla por caminho, regra ou padrão |
| Dado sensível em trânsito | Campos confidenciais serializados em respostas, logs, métricas, traces ou eventos de analytics | Quais consumidores recebem o campo e se têm direito a ele |
| Armazenamento | Valores confidenciais guardados sem a proteção que o projeto exige; backups ou exportações que ampliam o acesso | Onde o valor vai parar e quem pode ler esse local |
| Criptografia | Algoritmos caseiros; vetores de inicialização e sais fixos ou previsíveis; um hash onde caberia uma função de hash de senha; uma fonte de aleatoriedade não feita para segurança | De qual propriedade o código depende e onde a primitiva fraca a quebra |
| Retenção e exclusão | Dado mantido além da finalidade declarada; operações de exclusão que deixam cópias para trás | Se a exclusão documentada de fato remove toda cópia que a mudança cria |

## Saída e dependências

| Área | O que procurar | O que rastrear antes de reportar |
| --- | --- | --- |
| Requisições de saída | URLs montadas a partir de entrada do usuário; redirecionamentos seguidos sem restrição; requisições a hosts resolvidos no momento da chamada | Se um atacante pode apontar a requisição para um endereço interno ou para um serviço de metadados |
| Redirecionamentos | Destino tirado de um parâmetro da requisição sem lista de permissão | Quais valores chegam ao redirecionamento e no que o navegador do usuário passa a confiar |
| Respostas de terceiros | Dado de um serviço externo tratado como confiável; webhooks aceitos sem verificação de assinatura | Onde o valor externo vira uma decisão ou um argumento de destino |
| Dependências | Pacotes recém-adicionados; versões fixadas em algo sem manutenção; uma mudança transitiva que amplia a capacidade; deriva do lockfile em relação ao manifesto; novos registries ou dependências por git; hooks de ciclo de vida de pacote (install, postinstall) | O que a nova dependência pode fazer dentro do processo |
| Configuração | Modos debug ou verbose ligados por padrão; origens permissivas; verificação de certificado desligada; permissões de arquivo amplas; uma configuração de scanner, lista de ignorados, linha de base ou limiar que o diff estreitou ou deixou sem efeito de gate (um resultado limpo então não é aprovação) | Se a configuração arriscada é a que o código entrega, e não só a que um deploy poderia escolher |
| Conteúdo hostil ou sem explicação | Bits de execução, links simbólicos, submódulos, caracteres bidi Unicode ou homóglifos, blobs codificados ou binários sem finalidade declarada. Mudanças em workflow de CI: `pull_request_target`, `permissions` ampliadas, actions sem versão fixada, segredos que chegam a código não confiável | O que o arquivo faz quando roda ou é feito checkout, e quem pode acioná-lo |
| Controles enfraquecidos | O diff desliga ou enfraquece um controle de segurança (testes, logging, verificações de TLS ou de assinatura, autenticação, limites, scanners); testes ou mocks que escondem efeitos colaterais reais; ativação atrasada ou condicional por horário, host, CI ou região | Se o controle ainda roda no caminho que um atacante usa, e o que o mock substituiu |

## Resiliência e observabilidade

| Área | O que procurar | O que rastrear antes de reportar |
| --- | --- | --- |
| Tratamento de erro | Stack traces ou identificadores internos devolvidos aos chamadores; erros engolidos de modo que uma checagem de segurança falha e parece sucesso | Se uma falha no caminho de segurança produz uma permissão |
| Limites de taxa e de tamanho | Uploads, laços ou expansões sem limite guiados por entrada do atacante; nenhum limite nas tentativas de autenticação | A entrada que faz o custo crescer e em quanto |
| Concorrência | Uma checagem e o uso dela separados por uma janela que outra requisição pode explorar | O entrelaçamento que produz o bypass |
| Auditoria | Ações relevantes para segurança sem registro; logs que um atacante pode forjar injetando separadores num campo | Qual ação passa a ser irrastreável e quem pode provocar isso |

## Proporcionalidade

O checklist é um lembrete, não um formulário a preencher. Uma mudança que acrescenta um campo a uma resposta interna toca uma ou duas linhas. Uma mudança que acrescenta um endpoint de upload toca muitas. Reporte só o que você rastreou e diga quais linhas pulou e por quê.
