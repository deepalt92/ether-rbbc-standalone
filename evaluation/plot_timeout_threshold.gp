orange = '#d95319'; blue = '#4dbeee'; red = '#a2142f'; dblue = '#0072bd';

set terminal png size 1200,800
set output 'timeout_threshold_throughput.png'
set datafile separator ","
set yrange [0:2800]
set xrange [0:5]
set style data histogram
set style histogram cluster gap 1
set style fill solid
set boxwidth 0.9
set tics font "Helvetica,15"
set xtics ("4000" 1, "8000" 2, "16000" 3, "32000" 4)
set grid ytics
set ylabel "Throughput (tx/s)" font "Helvetica,15"
set xlabel "Workload (tx)" font "Helvetica,15"
set title "Timeout vs Threshold" font "Helvetica,15"
plot "./throughput/TimeoutvsThreshold.csv" using 2 title "Timeout" linecolor rgb orange,\
'' using 3 title "Threshold" linecolor rgb dblue

set terminal png size 1200,800
set output 'timeout_threshold_blocksize.png'
set datafile separator ","
set yrange [0:250]
set xrange [0:5]
set style data histogram
set style histogram cluster gap 1
set style fill solid
set boxwidth 0.9
set tics font "Helvetica,15"
set xtics ("4000" 1, "8000" 2, "16000" 3, "32000" 4)
set grid ytics
set ylabel "Superblock Size (tx)" font "Helvetica,15"
set xlabel "Workload (tx)" font "Helvetica,15"
set title "Timeout vs Threshold" font "Helvetica,15"
plot "./blocksize/TimeoutvsThreshold.csv" using 2 title "Timeout" linecolor rgb red,\
'' using 3 title "Threshold" linecolor rgb blue