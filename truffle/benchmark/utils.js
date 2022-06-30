let waitForAccountToUnlock = function(address, timeoutInSec = 10) {
    const startTime = new Date();
    const retryInSec = 2;

    return new Promise(async (resolve, reject) => {
        while ( true ) {
            let isUnlock = await isAccountLocked(address);
            if (isUnlock) {
                resolve(txReceipt);
                return;
            }

            const now = new Date();
            if (now.getTime() - startTime.getTime() > timeoutInSec * 1000) {
                reject(`Timeout after ${timeoutInSec} seconds`);
                return;
            }

            await waitFor(retryInSec)
        }
    });
};
