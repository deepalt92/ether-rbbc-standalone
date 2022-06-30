set terminal png size 1200,800
set output 'threshold_throughput.png'
set datafile separator ","
set key autotitle columnheader left
set xrange [0:1000]
set yrange [0:3000]
set tics font "Helvetica,15"
set xlabel "Threshold (tx)" font "Helvetica,15"
set ylabel "Throughput (tx/s)" font "Helvetica,15"
set title "Threshold vs Throughput" font "Helvetica,15"
plot for [col=2:5] "./throughput/Threshold.csv" using 1:col with linespoints

set terminal png size 1200,800
set output 'threshold_blocksize.png'
set datafile separator ","
set key autotitle columnheader left
set xrange [0:1000]
set yrange [0:1500]
set tics font "Helvetica,15"
set xlabel "Threshold (tx)" font "Helvetica,15"
set ylabel "Superblock Size (tx)" font "Helvetica,15"
set title "Threshold vs Superblock Size" font "Helvetica,15"
plot for [col=2:5] "./blocksize/Threshold.csv" using 1:col with linespoints