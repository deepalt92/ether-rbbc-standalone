set terminal png size 1200,800
set output 'cluster_blocksize.png'
set datafile separator ","
set key autotitle columnheader left
set xrange [0:64]
set yrange [0:19000]
set tics font "Helvetica,15"
set xlabel "Cluster Size" font "Helvetica,15"
set ylabel "Superblock Size (tx)" font "Helvetica,15"
set title "Cluster Size vs Superblock Size" font "Helvetica,15"
plot for [col=2:3] "./blocksize/Cluster.csv" using 1:col with linespoints

set terminal png size 1200,800
set output 'cluster_throughput.png'
set datafile separator ","
set key autotitle columnheader right top
set xrange [0:64]
set yrange [0:2800]
set tics font "Helvetica,15"
set xlabel "Cluster Size" font "Helvetica,15"
set ylabel "Throughput (tx/s)" font "Helvetica,15"
set title "Cluster Size vs Throughput" font "Helvetica,15"
set arrow from 0.,879.9 to 64,879.9 nohead front lc rgb "red" lw 1 dashtype "-"
set arrow from 0.,1411.3 to 64,1411.3 nohead front lc rgb "black" lw 1 dashtype "-"
plot for [col=2:3] "./throughput/Cluster.csv" using 1:col with linespoint, 1/0 t "Lightchain Baseline (Single Node)" lc rgb "red" lw 1 dashtype "-", 1/0 t "Ethereum Baseline (Single Node)" lc rgb "black" lw 1 dashtype "-"
