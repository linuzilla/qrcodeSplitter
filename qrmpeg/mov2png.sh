#!/bin/bash

rm -f output_*.png
ffmpeg -i qrcode.mov -r 5 -f image2 output_%04d.png
