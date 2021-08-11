function randomInterval(min, max) {
    if (min > max) {
        const tmp = min;
        min = max;
        max = tmp;
    }
    return Math.floor(min + Math.random() * (max - min + 1));
}

function randStr() {
    return 'xxxxxxxxxxxxxxxxxxxxxxxxxxx'.replace(/[x]/g, function(c) {
        let r = Math.random() * 16 | 0;
        return r.toString(16);
      })
}

export {randomInterval, randStr};