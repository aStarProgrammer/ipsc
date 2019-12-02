copy ..\..\IPSC.exe .\X86 /Y
copy ..\..\QuickHelp.txt .\X86 /Y
copy ..\..\FullHelp.txt .\X86 /Y
copy ..\..\config.ini .\X86 /Y

mkdir .\X86\Resources
xcopy ..\..\Resources .\X86\Resources /E /Y /I