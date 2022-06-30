set terminal png size 1200,800
set output 'timeout_throughput.png'
set datafile separator ","
set key autotitle columnheader left
set xrange [0:1600]
set yrange [0:2200]
set tics font "Helvetica,15"
set xlabel "Timeout (ms)" font "Helvetica,15"
set ylabel "Throughput (tx/s)" font "Helvetica,15"
set title "Timeout vs Throughput" font "Helvetica,15"
plot for [col=2:5] "./throughput/Timeout.csv" using 1:col with linespoints

set terminal png size 1200,800
set output 'timeout_blocksize.png'
set datafile separator ","
set key autotitle columnheader left
set xrange [0:1600]
set yrange [0:1200]
set tics font "Helvetica,15"
set xlabel "Timeout (ms)" font "Helvetica,15"
set ylabel "Superblock Size (tx)" font "Helvetica,15"
set title "Timeout vs Superblock Size" font "Helvetica,15"
plot for [col=2:5] "./blocksize/Timeout.csv" using 1:col with linespoints