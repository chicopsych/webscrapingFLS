---
title: História do Arduino (MIC571)
url: https://www.newtoncbraga.com.br/microcontroladores/138-atmel/18422-historia-do-arduino-mic571.html
timestamp: 2026-03-10T11:41:49.9410718-03:00
---

Pesquisar Pesquisar
            

        

      


        		
        	
        
        
        
                    
                                    
PrincipalAprendendo OnlineAutomaçãoBancada de ReparaçãoCircuitosDicasLivros TécnicosNovidadesProjetos e ArtigosRevistasmais ...

                                            
            

    
        
        
        
            
            

            
                
    
    
    
        
        
            História do Arduino (MIC571)        
                            
        
        
            

            
                        Detalhes                    

                    
    
                    Escrito por: Marcos de Lima Carlos    
        
        
        
        
        
    
    
    
    
        
                                                
        De acordo com a wikipedia, o projeto do Arduino se iniciou na Itália, mais precisamente na cidade de Ivrea no ano de 2005. A ideia foi juntar numa única placa, um microcontrolador com comunicação serial (as placas que utilizam entradas USB são conversores) sem a necessidade de um gravador e ainda com fácil programação.
 
Para que isso ocorresse existiu alguns projetos que deram um passo para que o ecossistema do Arduino ocorresse. O nome Arduino vem de um bar em Ivrea, onde os fundadores do projeto costumavam se reunir.
 

Roteiro

História do Arduino
Como era antes do Arduino?
Tipos de Arduino
O Arduino UNO

Esse é um artigo patrocinado pela Curto Circuito

 

Mapa de Ivrea

 
Microcontroladores é um item de eletrônica bem antigo. Porém exige determinados tipos de conhecimento técnico que um estudante do ensino médio não tem. A maioria dos microcontroladores eram programados em linguagem de montagem (assembly), possuem ambientes de desenvolvimento próprios (IDEs) e particularidades que tornavam difíceis o uso por qualquer pessoa sem a bagagem técnica necessária. Não é muito longa a curva de aprendizagem, porém é muito técnica e isso restringe muito a divulgação de um produto.
Um dos projetos que colaborou para o surgimento do Arduino foi o processing. Um dos seus objetivos é atuar como uma ferramenta para não programadores iniciados com a programação. Traduzindo: democratizar o uso da programação. Antes do processing existiu outras ferramentas para isso como o Logo da IBM, BASIC e inúmeras outras surgiram. Veremos isso em um outro artigo sobre as formas de programar um Arduino.
Outra tecnologia que auxiliou o Arduino foi o wiring. O wiring é projeto semelhante ao Arduino. Inclusive você não terá problema em programá-lo se você já mexe com Arduino porque é muito parecido. A imagem abaixo é uma versão atual do wiring na sua versão S (atual).
 

Wiring S

 
O Wiring nasceu, em 2003, de uma dissertação de mestrado de Hernando Barragán em Ivrea, no Instituto de Design Interativo de Ivrea. Ele foi o projeto que foi inspiração ao Arduino. A ideia principal do projeto é fazer com que as pessoas pudessem compartilhar ideias sem ter muito conhecimento técnico de eletrônica. Massimo Banzi, o criador do Arduino, foi orientador de Hernando Barragán. Ele ministrava uma aula de programação de microcontroladores no instituto. O Arduino é um fork do projeto do wiring, utilizado como ponto de partida. Por que, então, o wiring não conseguiu visibilidade?
1. O Wiring foi excessivamente acadêmico. Ele foi criado como uma dissertação.
2. A placa inicial de wiring custava cerca de 60 dólares. Não é um custo alto, porém um arduino uno hoje custa aproximadamente 10 dólares. Isso alavanca a popularidade entre estudantes, pessoas geralmente com poder aquisitivo baixo.
3. Ao término de sua dissertação Hernando Barragán voltou para a Colômbia em 2004. Isso dificulta a interação entre os criadores. Provavelmente, tinha outros planos para o Wiring. O wiring continua sendo desenvolvido até hoje.
4. O Arduino é totalmente focado na comunidade e em conteúdo colaborativo.
A dissertação de mestrado de Hernando Barragán você pode baixar aqui.  É um excelente ponto de partida se quiser entender o contexto da época.
 
O Primeiro Arduino
Abaixo a imagem do primeiro arduino:
 

Primeiro Arduino

 
O primeiro Arduino foi desenvolvido por Massimo Banzi, David Curtielles, Tom Igoe, Gianluca Martino e David Mellis. Utilizou um atmega8 com uma conexão serial DB25 como interface entre o PC e a placa. O grande segredo dessas placas é um software chamado BootLoader.
 
O Bootloader
Programar microcontroladores não é uma tarefa muito legal se você está começando em eletrônica. Cada fabricante possui uma série de requisitos, normalmente escrito em um documento chamado DataSheet(tradução literal - Folha de especificações). Junto a isso você precisará de softwares, que são proprietários e eram específicos para cada linha de chips, e um gravador. Desenhar uma placa não é algo que qualquer pessoa faça com facilidade. Lembrando que em 2005 os softwares CAD para eletrônica não eram tão difundidos e soluções opensource não eram muito divulgadas. Houve iniciativas de tornar os microcontroladores populares como o BasicStamp (https://pt.wikipedia.org/wiki/BASIC_Stamp), por exemplo. E se, tudo isso pudesse ser simplificado numa única placa?
Um dos processos que permite que isso ocorra é o bootloader. Sem o bootloader você precisaria enviar o programa através de um gravador.
O microcontrolador possui uma série de configurações iniciais. O bootloader simplifica isso e seu funcionamento é bem parecido com a BIOS do seu computador. Ele é um programa que é executado quando o microcontrolador inicia ou é feito um reset. É um conjunto de instruções que tem por objetivo gravar as instruções no Microcontrolador através da porta serial (Você utiliza um conversor serial para USB no Arduino UNO!).
A primeira tarefa do bootloader é examinar o que causou a reinicialização. Depois disso ele poderá executar a transmissão de um novo programa ou a execução de um código no Arduino. Por cima, e bem por cima, é isso o que um Arduino faz.
 
Tipos de Arduino
Há muitos derivados do Arduino no mercado. Aqui nesse artigo eu vou abordar por cima 3 tipos de Arduino:
1. Arduino Nano
2. Arduino UNO
3. Arduino Mega
O Arduino Nano é a versão compacta do Arduino. Ela e a versão do Arduino UNO são bastante parecidas em termos de especificações. Operam a 5 volts, possuem a mesma capacidade de memória interna, SRAM e EPROM. Uma das diferenças é ser compacta, pesando apenas 7 gramas (nano) em contrapartida de 25 gramas (UNO).
 

Arduino Nano | Clique na imagem para ampliar |

 
 

Arduino Uno

 
 O Arduino Mega é uma versão mais estendida. Possui mais pinos de saída, um microcontrolador mais potente, entre outras coisas. A ideia da versão mega é atender projetos com maior demanda de velocidade, dispositivos de entrada e saída, e portas seriais.
 

Arduino Mega

 
 
Há outros tipos de Arduino que não serão abordados aqui neste artigo.
 
O Arduino UNO
Nessa série de artigos nós vamos utilizar o Arduino UNO. Um dos motivos pelo qual ele foi o escolhido é o fato dele ser de fácil acesso e muito simples de mexer. Há também muitas placas adicionais que você encontra disponível na Curto Circuito que são chamadas de shields. Elas costumam vir num formato padrão onde se encaixa uma placa a outra.
O Arduino Uno possui as seguintes características:
Microcontrolador - ATmega328P
Tensão de operação - 5V
Tensão de entrada (recomendada) - 7 ~ 12V
Tensão de entrada (limite) - 6 ~ 20V
Pinos de Entrada/Saída digital - 14 (dos quais 6 fornecem uma saída PWM)
PWM pinos - 6
Pinos de entrada analógica - 6
Corrente Contínua por pino de I/O - 20mA
Corrente Contínua para pinos em 3V3 - 50mA
Memória Flash - 32KB sendo que 0.5KB é utilizado pelo bootloader
SRAM - 2KB
EEPROM - 1KB
Velocidade de Clock - 16MHz
Peso - 25g
Abaixo uma imagem da pinagem da placa:
 

Arduino Uno Pinout | Clique na imagem para ampliar |

 
 
Arduino IDE
O Arduino IDE é onde iremos executar a parte de programação da placa Arduino UNO nesse artigo. A interface do Arduino IDE é multiplataforma e opensource. Escolha o sistema operacional a ser utilizado e a instale no sistema. O download pode ser efetuado nesse link
O próximo artigo da série abordaremos a ideia do projeto colaborativo.
 
     

                                    

    
            
        
            
                                                            
                                                                    
                    
                        Marcos de Lima Carlos                    
                
                                                                                
                    
                        arduino                    
                
                        
                    
                                        

            
            
        

                    
                
        





            Recentes        
                
    
        
            
                Filtro passa faixa CIR17252S CB9763E (CIR19810)            
        
    
    
        
            
                Medidor de intensidade de campo sintonizado CIR17233S CB9744E (CIR19791)            
        
    
    
        
            
                SMPS DC para DC CIR19131S CB9669E (CIR19716)            
        
    
    
        
            
                Oscilador de Butler de 20 a 100 MHz CIR17094S CB9631E (CIR19679)            
        
    
    
        
            
                Oscilador estável a cristal CIR17118S CB9656E (CIR19703)            
        
    

    


        

    
                                                                                                                                                                                                                                                                            
                            
                        
                                                            
    
                                                                                                                                                                                                                                                                            
                            
                        
                                                            
    
                                                                                                                                                                                                                                                                            
                            
                        
                                                            




            Última Edição        
                

    
    


        

    
Buscador de Datasheets
    N° de Componente  


            
        
                    
                
        

    NO YOUTUBE


 
  


NOSSO PODCAST

 
 
 



            
        
                    
                
        
                

    
                                                                                                                                                                                                                                                                            
                            
                        
                                                            


    


            Banco de Circuitos        
                
                    
                                    Chave de toque Darlington CIR17029S CB9564E (CIR19612)    
    
    
    
    
    
    
    

                                    Conversor corrente em tensão CIR17063S CB9600E (CIR19648)    
    
    
    
    
    
    
    

                                    Oscilador com comutação de cristal CIR17095S CB9632E (CIR19680)    
    
    
    
    
    
    
    

                                    Dobrador de 150 para 300 MHz CIR17284S CB9795E (CIR19841)    
    
    
    
    
    
    
    

                                    Transmissor de temperatura CIR16476S CB9392E (CIR19448)    
    
    
    
    
    
    
    
    
    


            
            

            
            
                

    Instituto Newton C. Braga:
 Mapa do Site - Entre em contato - Como Anunciar - Políticas do Site - Advertise in Brazil
