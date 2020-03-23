del /S /F /Q ..\..\ipsc_release\Windows\Windows_X64\ 
del /S /F /Q  ..\..\ipsc_release\Windows\Windows_X86\ 
del /S /F /Q  ..\..\ipsc_release\Linux\Linux_X64\ 
del /S /F /Q  ..\..\ipsc_release\Linux\Linux_X86\ 
del /S /F /Q  ..\..\ipsc_release\Darwin64\ 
del /S /F /Q  ..\..\ipsc_release\Arm6\ 

rmdir /S /Q ..\..\ipsc_release\Windows\Windows_X64\Resources 
rmdir /S /Q  ..\..\ipsc_release\Windows\Windows_X86\Resources 
rmdir /S /Q  ..\..\ipsc_release\Linux\Linux_X64\Resources  
rmdir /S /Q  ..\..\ipsc_release\Linux\Linux_X86\Resources  
rmdir /S /Q  ..\..\ipsc_release\Darwin64\Resources  
rmdir /S /Q  ..\..\ipsc_release\Arm6\Resources  