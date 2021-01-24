#!/bin/bash
go build -o wriggle
sudo cp ./wriggle /usr/local/bin/
echo "built and added to path, have fun :)"
