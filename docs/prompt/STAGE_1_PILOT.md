# Implementar uma rotina de "sanitização/limpeza/segurança" nos meus repos pelo meu MCP

## Ideia inicial

Implementar uma rotina de "sanitização/limpeza/segurança", primeiro nos meus repos pelo meu MCP.. mas a ideia é abstração pra ser possivelmente global como ferramenta de "apoio".

Acho que vai ser a primeira coisa realmente produtiva que vou integrar no MCP GoBE.. Poderia tipo:

Algo como:

*- Rotar chaves SSH no github automaticamente;
*- Limpar runs de workflows inexistentes (porque foram substituidos ou descontinuados) e runs acima de X tempo;
*- Me enviar no discord (talvez whatsapp, porque já tá integrado mas a API é um saco.. hehhee) o resumo das atividades que ele executou nos repos com um histórico anexado pra eu salvar localmente, já que estaria removendo do github..SEMPRE bom "manter as coisas salvas após deletá-las".. hahahahha;
*- Me reportar tb de PRs, issues e "inatividade" relevantemente longa em determinado repo que tenha potencial de adesão/produtivo, pra não deixar ideias morrerem. hehehe

---

Refletindo sobre a modularização do projeto GoBE/MCP, que já existe: Eu optei (por hora) por manter todo conteúdo gigante, agrupado no mesmo repo/módulo Go (GoBE) que hoje tem uns 155 arquivos, sendo uns 14 shell e o resto Go.. hehehehe.

Os packages estão SUPER bem separados e quando possível vou adequando blocos dele pra "minha arquitetura meio aloprada.. kkkk".

PS: No Go né:
1: go.mod é exclusivo de cada projeto e representa um Módulo.

2: blocos agrupados na árvore de arquivos/pastas do módulo são packages..

---

Estava tentando no projeto que comecei pro que descrevi acima, implementar uma estrutura de arquivos que facilitasse a organização e a modularização do código. A ideia é tentar usar um design "interface-first",onde TODA estrutura que fosse ser exposta deve ser, primeiro considerada e analisada enquanto interface e SOMENTE se for realmente necessário, exposta como struct. Quando falo sobre EXPOR, digo 'SAIR DO PACKAGE INTERNAL". hehe

O que já existe hoje fora do internal iria continuar existindo, são lógicas que realmente precisam ser expostas, como entrypoints para comandos CLI e API's para outros porojetos acessarem a lófgica do projeto, porém, para lógicas que habitam o internal,iriam acessar através de interfaces expostas por labeled types que teriam só o alias e o constructor exposto.

Enfim: Já escrevi demais, acho que seria ótimo se você desse mma contextualizada inicial agora e já ponderasse coisas em torno do que falei acima, pra eu ir amadurecendo a ideia e quem sabe já ir implementando.

MAS TEM UMA COISA QUE QUERO TIPO: DEIXAR O MAIS CLARO POSSÍVEL E TE PEDIR PARA FAZER O POSSÍVEL E IMPOSSÍVEL PARA NÃO ESQUECER E NÃO DESVIAR DISSO POR NADA. Se realmente for preciso ou muuuito relevantemente positivo, NÓS DECIDIREMOS JUNTOS! A DIRETRIZ É:,

NÃO USAR MOCK! SE POSSÍVEL,EM ABSOLUTAMENTE EM NENHUM LUGAR DESSE PROJETO!

Topa? Bora?

---

Só pra frisar e lembrar um pouco mais (hahahaha):

NÃO USAR MOCK! SE POSSÍVEL,EM ABSOLUTAMENTE EM NENHUM LUGAR DESSE PROJETO!

NÃO USAR MOCK! SE POSSÍVEL,EM ABSOLUTAMENTE EM NENHUM LUGAR DESSE PROJETO!

NÃO USAR MOCK! SE POSSÍVEL,EM ABSOLUTAMENTE EM NENHUM LUGAR DESSE PROJETO!
