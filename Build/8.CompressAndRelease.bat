cd ..\..\ipsc_release
set ipscVersion=0.2.0.3
mkdir %ipscVersion%

cd Windows
bandizip c ipsc_%ipscVersion%_Windows_X64.zip Windows_X64
bandizip c ipsc_%ipscVersion%_Windows_X86.zip Windows_X86 

::ping 127.0.0.1 -n 15 > nul 

move /Y ipsc_%ipscVersion%_Windows_X64.zip ..\%ipscVersion%
move /Y ipsc_%ipscVersion%_Windows_X86.zip ..\%ipscVersion%

cd ..
cd Linux

bandizip c ipsc_%ipscVersion%_Linux_X64.tgz Linux_X64
bandizip c ipsc_%ipscVersion%_Linux_X86.tgz Linux_X86 

::ping 127.0.0.1 -n 15 > nul 

move /Y ipsc_%ipscVersion%_Linux_X64.tgz ..\%ipscVersion%
move /Y ipsc_%ipscVersion%_Linux_X86.tgz ..\%ipscVersion%

cd ..
bandizip c ipsc_%ipscVersion%_Darwin64.tgz Darwin64

::ping 127.0.0.1 -n 15 > nul 

move /Y ipsc_%ipscVersion%_Darwin64.tgz %ipscVersion%

bandizip c ipsc_%ipscVersion%_Arm6.tgz Arm6

::ping 127.0.0.1 -n 15 > nul 

move /Y ipsc_%ipscVersion%_Arm6.tgz %ipscVersion%

cd ..\ipsc\Build