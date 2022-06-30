set terminal png size 1200,800
set output 'baseline.png'
set datafile separator ","
set yrange [0:2800]
set xrange [0:6]
set style data histogram
set style histogram cluster gap 1
set style fill solid
set boxwidth 0.9
set xtics ("Mulga Chain (4 nodes)" 1, "Mulga Chain (32 nodes)" 2, "Mulga Chain (64 nodes)" 3, "Lightchain (Single node)" 4, "Ethereum (Single node)" 5)
set grid ytics
set ylabel "Throughput (tx/s)"
set xlabel "Blockchain systems"
set title "Mulga Chain Performance Evaluation"
plot "./throughput/baseline.csv" using 1 linecolor rgb "red"