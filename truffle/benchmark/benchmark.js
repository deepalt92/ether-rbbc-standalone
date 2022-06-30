var Web3 = require('web3');

const host = "http://localhost:8540";
const web3 = new Web3(new Web3.providers.HttpProvider(host));

const senderValue = {
    "version": 3,
    "id": "54fc30ab-d1b7-4397-a39f-1ddd9110102b",
    "address": "c916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e",
    "crypto": {
        "ciphertext": "29beea264a484e72978286a64d3976060cd8dc74e57dd4d8a14ee12db2dfe6f9",
        "cipherparams": {"iv": "6569b60cf13335a41340243b9ab67d9b"},
        "cipher": "aes-128-ctr",
        "kdf": "scrypt",
        "kdfparams": {
            "dklen": 32,
            "salt": "a36f4967cd868fb5dec0380822dc014e748fb7c2bf968f29d6ad11c377189181",
            "n": 8192,
            "r": 8,
            "p": 1
        },
        "mac": "d55bb34a14abff2487ea78fa65d641981e094e576e6e3323e3550b2812081261"
    }
};
const senderPwd = "WelcomeToSirius";

const receiverValue = {
    "address": "bce16ea55bb357b038e612b1722a88879c665a31",
    "crypto": {
        "cipher": "aes-128-ctr",
        "ciphertext": "a295cf7863f43a86be23ba907c39a04fe2054183e99842884e9f06a493344456",
        "cipherparams": {"iv": "f09bc2c84caa20c270cc244d175f5af3"},
        "kdf": "scrypt",
        "kdfparams": {
            "dklen": 32,
            "n": 262144,
            "p": 1,
            "r": 8,
            "salt": "4dbdbfbe57b915535040018e994dfac15ed4354617adf385f3a78878aa80ecb9"
        },
        "mac": "cb5887859994e461f5c4c1a4bb582051c5c613d13981695c7638f1157c4dab77"
    },
    "id": "12bc2f1e-8e26-4e36-93fc-80fb2ecb7e77",
    "version": 3
};
const receiverPwd = "hello";

const sender = web3.eth.accounts.decrypt(senderValue, senderPwd);
const receiver = web3.eth.accounts.decrypt(receiverValue, receiverPwd);

web3.eth.getBalance(receiver.address).then(console.log);

bench = async function() {
    const oldTime = new Date();

    for (var i = 0; i < 19; i++) {
        var rawTx = {
            from: sender.address,
            to: receiver.address,
            value: 1,
            gas:21000,
        };
        sender.signTransaction(rawTx).then(signedTx => web3.eth.sendSignedTransaction(signedTx.rawTransaction));
    }

    await sender.signTransaction(rawTx).then(signedTx => web3.eth.sendSignedTransaction(signedTx.rawTransaction));

    const newTime = new Date();

    console.log(newTime - oldTime);

    web3.eth.getBalance(receiver.address).then(console.log);
};

bench();
