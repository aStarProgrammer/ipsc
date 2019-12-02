copy ..\..\IPSC.exe .\X64 /Y
copy ..\..\QuickHelp.txt .\X64 /Y
copy ..\..\FullHelp.txt .\X64 /Y
copy ..\..\config.ini .\X64 /Y

mkdir .\X64\Resources
xcopy ..\..\Resources .\X64\Resources /E /Y /I