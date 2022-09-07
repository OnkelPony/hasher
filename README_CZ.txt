Hasher je skript, sloužící ke generování hashů souborů na disku a jejich porovnávání se seznamem hashů.
Pro sestavení ze zdrojových kódů je nutné mít nainstalovaný kompilátor jazyka go a program make.
Binární soubor sestavíme pomocí příkazů
$ make build/windows
případně
$ make build/linux
Funguje jednoduchá nápověda
$ ./hasher --help
Jeho použití je následujicí:
$ ./hasher --hashes=jmeno.soubour --name=jmeno_projektu PATH
kde parametr
--hashes určuje jméno souboru s hashi, se kterými se má porovnávat. Tento soubor je CSV bez mezer, případně lze každý hash umístit na zvláštní řádek.
Pozor na skryté mezery! Příklad souboru s hashi je hashes.txt. Pokud parametr chybí, neprovádí se žádné porovnávání, pouze se generuje soubor s hashi.
--name je jméno "projektu", které bude přítomno v názvech výstupních souborů. Standardních i chybových. Pokud parametr chybí, je použito "hashInfo"
PATH je cesta od které se má hashovat. Na příklad na Linuxu:
$ ./hasher /home/franta
nebo na Windows:
.\hasher.exe C:\Users\franta
Pokud parametr PATH chybí, použije se na windows "c:\" a na linuxu "/"
Výsledkem je CSV soubor, jehož název je ve formátu:
'hashInfo[Aug  9 153853].csv'
s použitím data a času vytvoření CSV souboru.
Vnitřní formát souboru je:
[/home/franta/notes/lw-sockets.txt], [8eb035fe270e5c04e63eca444ce72287], [d800934a25ab9fe91218bc7d70999652e831db44], [7dec79b860d4ee332b2966407769cc6bba218803b4689d2ec92ff50704ee58be]
První položkou je cesta k souboru, druhou MD5 hash, třetí SHA1 hash a čtvrtou SHA256 hash.
Řádky jsou seřazeny podle abecedy, je tedy možné provádět hashování stejného cíle v různých časech a pomocí nástroje diff zjistit rozdíly.
Nalezené shody mezi položkami souboru s hashi a hashovanými soubory se vypisují na standardní výstup.

Poznámka: První nástřel skriptu byl někdy v roce 2021 a původní požadavky nebyly zcela jednoznačné. Proto jsem nekódoval žádné testy, mělo jít téměř o one-liner.
Tuto chybu už neudělám, i na one-liner napíšu nejdřív test. Funkčnost jsem otestoval na lin i win a nenarazil na zásadní problém. Kdyžtak pište ;-)