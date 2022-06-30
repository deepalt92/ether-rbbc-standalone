function bench(size) {
    personal.unlockAccount(eth.accounts[0], "WelcomeToSirius");
    var batch = web3.createBatch();
    start = Date.now();
    for (i = 0; i < size; i++) {
        batch.add(web3.eth.sendTransaction({
            from: eth.accounts[0],
            to: "0xBce16ea55bB357B038e612b1722A88879c665a31",
            value: i
        }))
    }
    batch.execute();
    end = Date.now();
    duration = end - start;
    throughput = size * 1000 /duration;
    console.log("Start: "+start);
    console.log("End: "+end);
    console.log("Transaction Number: "+size);
    console.log("Duration: "+duration);
    console.log("Throughput: "+throughput);
}

function calTps(start, end, size) {
    duration = end - start;
    throughput = size * 1000 /duration;
    console.log("Duration: "+duration);
    console.log("Throughput: "+throughput);
}

// console.log(eth.getBalance("0xBce16ea55bB357B038e612b1722A88879c665a31"));

console.log(calTps(1570972800176, 1570972807330, 1200));