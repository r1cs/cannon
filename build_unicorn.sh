#!/bin/bash
git clone https://github.com/geohot/unicorn.git -b dev unicorn2
#git clone https://github.com/unicorn-engine/unicorn.git -b dev unicorn2
cd unicorn2
#cmake . -DUNICORN_ARCH=mips -DCMAKE_BUILD_TYPE=Debug
cmake . -DUNICORN_ARCH=mips -DCMAKE_BUILD_TYPE=Release
make -j8

