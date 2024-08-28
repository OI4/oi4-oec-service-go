cd api || exit
./gomod.sh

cd ..
cd container || exit
./gomod.sh

cd ..
cd dnp || exit
./gomod.sh

cd ..
cd service || exit
./gomod.sh

cd ..
cd demo || exit
./gomod.sh
