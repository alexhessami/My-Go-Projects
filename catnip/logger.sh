#!/bin/bash
#used as a test logging file

a=1
while true; do
  echo "Hello world. $a" >> output.txt
  a=$((a+1))
  sleep 15
done
