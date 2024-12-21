cd c_mvm
rm -rf build
mkdir build 
cd build
cmake ../
make install

cd ../../linker
rm -rf build
mkdir build 
cd build
cmake ..
make install
 