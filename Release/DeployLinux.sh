rm -rf $HOME/ipsc
mkdir $HOME/ipsc
cp IPSC $HOME/ipsc
cp QuickHelp.txt $HOME/ipsc
cp FullHelp.txt $HOME/ipsc
cp config.ini $HOME/ipsc
cp -rf Resources $HOME/ipsc
cd $HOME/ipsc
chmod 777 ./IPSC
echo 'export PATH=$PATH:$HOME/ipsc' >> ~/.bashrc