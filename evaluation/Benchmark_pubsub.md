# Benchmark Results

## Experiment 1

This experiment benchmarks the efficiency and scalability of Mulga Chain with pub/sub support

### AWS Configuration

- Consensus cluster size: 4, 8, 16, 32
- EVM subscriber size (per consensus node): 1, 2, 4, 8
- Instance type: c5.2xlarge
- Availability Zone: ap-southeast-1c
- AMI: ami-020ee6b0bd3aa697b
- OS Image: ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20191002
- CPU: 8 vCPU
- Memory: 16 GB
- Storage: EBS
- Network: Up to 10 Gigabit

### Consensus Setup

DBFT consensus clusters of size 4, 8, 16, 32 is set up for benchmark.

### Subscriber Setup

EVM subscribers for each DBFT consensus node is set to 1, 2, 4, 8

Transaction pool threshold: 500

Timeout: 5000

### Client Setup

For each node, 5 clients are deployed to concurrently send transactions to the target server. Each client sends 500 transactions.

### Results

| Subscriber\Consensus |   4  |   8  |  16  |  32  |
|:--------------------:|:----:|:----:|:----:|:----:|
|           4          | 2803 |   -  |   -  |   -  |
|           8          | 2466 | 2043 |   -  |   -  |
|          16          | 2928 | 1808 | 1884 |   -  |
|          32          | 2116 | 1768 | 1616 | 1143 |
|          64          | 1033 |  963 | 1000 |  847 |
|          128         |  545 |  526 |  499 |  514 |

## Experiment 2

This experiment benchmarks the efficiency and scalability of Mulga Chain with

* pub/sub support

* Merging proposals

* Timer for closing late proposals

### AWS Configuration

- Consensus cluster size: 4, 8, 16, 32
- EVM subscriber size (per consensus node): 1, 2, 4, 8
- Instance type: c5.2xlarge
- Availability Zone: ap-southeast-1c
- AMI: ami-020ee6b0bd3aa697b
- OS Image: ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20191002
- CPU: 8 vCPU
- Memory: 16 GB
- Storage: EBS
- Network: Up to 10 Gigabit

### Consensus Setup

DBFT consensus clusters of size 4, 8, 16, 32 is set up for benchmark.

### Subscriber Setup

EVM subscribers for each DBFT consensus node is set to 1, 2, 4, 8

### Client Setup

The total number of transactions handled by the system is set to be 9600

For each node, 3 clients are deployed to concurrently send transactions to the target server.

### Results

| Subscriber\Consensus |                     4                    |                     8                    |                    16                   |                    32                   |
|:--------------------:|:----------------------------------------:|:----------------------------------------:|:---------------------------------------:|:---------------------------------------:|
|           4          | TPS: 6328 Threshold: 2400 Tx/client: 800 |                     -                    |                    -                    |                    -                    |
|           8          | TPS: 5731 Threshold: 1200 Tx/client: 400 | TPS: 5553 Threshold: 1200 Tx/client: 400 |                    -                    |                    -                    |
|          16          |  TPS: 5420 Threshold: 600 Tx/client: 200 |  TPS: 5171 Threshold: 600 Tx/client: 200 | TPS: 3969 Threshold: 600 Tx/client: 200 |                    -                    |
|          32          |  TPS: 4799 Threshold: 300 Tx/client: 100 |  TPS: 4824 Threshold: 300 Tx/client: 100 | TPS: 3825 Threshold: 300 Tx/client: 100 | TPS: 2151 Threshold: 150 Tx/client: 100 |
|          64          |  TPS: 3119 Threshold: 150 Tx/client: 50  |  TPS: 3989 Threshold: 150 Tx/client: 50  |  TPS: 2610 Threshold: 150 Tx/client: 50 |  TPS: 2135 Threshold: 150 Tx/client: 50 |
|          128         |   TPS: 1736 Threshold: 75 Tx/client: 25  |   TPS: 1996 Threshold: 75 Tx/client: 25  |  TPS: 1494 Threshold: 75 Tx/client: 25  |   TPS: 950 Threshold: 75 Tx/client: 25  |

## Experiment 3

This experiment benchmarks the efficiency and scalability of Mulga Chain with

* pub/sub support

* Merging proposals

* Timer for closing late proposals

### AWS Configuration

- Consensus cluster size: 4, 8, 16, 32
- EVM subscriber size (per consensus node): 1, 2, 4, 8
- Instance type: c5.2xlarge
- Availability Zone: ap-southeast-1c
- AMI: ami-020ee6b0bd3aa697b
- OS Image: ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20191002
- CPU: 8 vCPU
- Memory: 16 GB
- Storage: EBS
- Network: Up to 10 Gigabit

### Consensus Setup

DBFT consensus clusters of size 4, 8, 16, 32 is set up for benchmark.

### Subscriber Setup

EVM subscribers for each DBFT consensus node is set to 1, 2, 4, 8

Transaction pool threshold: 300

Timeout: 5000

### Client Setup

The total number of transactions handled by the system is set to be 38400

For each node, 3 clients are deployed to concurrently send transactions to the target server.

### Results

#### TPS of standalone EVM

TPS: 1487

TPS of transaction commit: 3685

Time in commit/Total Time: 40.4 %

Experimental Results can be found in [EVM](./exp3/EVM)

#### TPS of Mulga Chain

Each cell shows: TPS, TPS of transaction commitment, ratio of committing time in total time.

The data displayed here is approximation of the exact data.

| Subscriber\Consensus |        4        |        8        |        16       |        32       |
|:--------------------:|:---------------:|:---------------:|:---------------:|:---------------:|
|           4          | 3092, 4384, 71% |        -        |        -        |        -        |
|           8          | 3118, 4615, 68% | 2936, 4667, 63% |        -        |        -        |
|          16          | 2851, 4714, 61% | 2745, 4690, 59% | 1923, 4488, 44% |        -        |
|          32          | 2020, 2647, 79% | 1939, 2687, 76% | 1820, 2615, 75% | 1289, 2647, 54% |
|          64          | 1107, 1302, 86% | 1077, 1295, 84% |  988, 1106, 90% |  907, 1260, 74% |
|          128         |  563, 656, 86%  |  540, 636, 85%  |  535, 629, 86%  |  507, 656, 79%  |

Raw data can be found in folder [exp3](./exp3)

## Results from non-pubsub system

| Cluster size | Threshold=300 | Threshold=500 |
|:------------:|:-------------:|:-------------:|
|       4      |     2475.6    |     2418.1    |
|       8      |     1897.1    |     2131.5    |
|      16      |     1714.1    |     1823.5    |
|      24      |      1674     |     1681.8    |
|      32      |     1613.1    |     1494.4    |
|      40      |     1406.2    |     1369.7    |
|      48      |     1058.8    |     1178.3    |
|      56      |     988.0     |     1188.1    |
|      64      |     880.2     |     1187.2    |

![](./plots/cluster_throughput.png)