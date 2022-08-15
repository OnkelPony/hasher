Hasher je skript, soužící ke generování hashů souborů na disku.
Pro sestavení ze zdrojových kódů je nutné mít nainstalovaný kompilátor jazyka go a program make.
Binární soubor sestavíme pomocí příkazů
$ make build/windows
případně
$ make build/linux
Jeho použití je následujicí:
$ ./hasher PATH
kde parametr PATH je cesta od které se má hashovat. Na příklad na Linuxu:
$ ./hasher /home/franta
nebo na Windows:
.\hasher.exe C:\Users\franta
Výsledkem je CSV soubor, jehož název je ve formátu:
'hashInfo[Aug  9 153853].csv'
s použitím data a času vytvoření CSV souboru.
Vnitřní formát souboru je:
[/home/franta/notes/lw-sockets.txt], [8eb035fe270e5c04e63eca444ce72287], [d800934a25ab9fe91218bc7d70999652e831db44], [7dec79b860d4ee332b2966407769cc6bba218803b4689d2ec92ff50704ee58be]
První položkou je cesta k souboru, druhou MD5 hash, třetí SHA1 hash a čtvrtou SHA256 hash.
