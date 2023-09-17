#!/bin/bash

# Initialize sparse-checkout and pull only the api directory
echo "Cloning the repository..."
git clone --no-checkout https://github.com/BryceWayne/blogger.git ./tempdir
echo "Cloning complete!"
cd ./tempdir
git sparse-checkout init --cone
git sparse-checkout set api
git pull origin master

# Move the api directory and remove tempdir
mv ./api /app/api
cd ..
rm -rf tempdir

# Start your Go application
./main
