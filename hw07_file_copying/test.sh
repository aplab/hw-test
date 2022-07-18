#!/usr/bin/env bash
clear

set -xeuo pipefail

go build -o go-cp
rm -f out.txt

./go-cp -from testdata/input.txt -to out.txt
cmp out.txt testdata/out_offset0_limit0.txt
rm -f out.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10
cmp out.txt testdata/out_offset0_limit10.txt
rm -f out.txt

./go-cp -from testdata/input.txt -to out.txt -limit 1000
cmp out.txt testdata/out_offset0_limit1000.txt
rm -f out.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10000
cmp out.txt testdata/out_offset0_limit10000.txt
rm -f out.txt

./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 1000
cmp out.txt testdata/out_offset100_limit1000.txt
rm -f out.txt

./go-cp -from testdata/input.txt -to out.txt -offset 6000 -limit 1000
cmp out.txt testdata/out_offset6000_limit1000.txt
rm -f out.txt

# need a big file to progressbar demo because copying too fast
#./go-cp -from /home/polyanin/hwtest/video.mkv -to /home/polyanin/hwtest/video-copy.mkv
#cmp /home/polyanin/hwtest/video.mkv /home/polyanin/hwtest/video-copy.mkv
#rm -f /home/polyanin/hwtest/video-copy.mkv

rm -f go-cp out.txt
echo "PASS"
