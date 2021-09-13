import { sleep, group } from 'k6';
import { Counter, Rate } from 'k6/metrics';
import tcp from 'k6/x/tcp';
import { C2GVerifyMessage, C2GTestMessage, C2GHeartbeatMessage } from '../test_unit/func_unit.js';
import { randStr, randomInterval } from '../utils/random.js';

const susRate = new Rate('susRate');
const errCounter = new Counter('errCounter');

const strAddr = '127.0.0.1:9000';

let userName = randStr();
let client = new tcp.Client();
let bCreate = false;

export let options = {
    stages: [
        { duration: '10s', target: 1 },
        { duration: '30s', target: 1 },
    ],
}

export function setup() {
    console.log('init test unit array');
    const srcUnits = [
        { weight: 10, name: 'C2GTestMessage' },
        { weight: 10, name: 'C2GHeartbeatMessage' },
    ];
    const dstUnits = [];
    for (const testInfo of srcUnits) {
        for (let i = 0; i < testInfo.weight; i++) {
            dstUnits.push(testInfo.name);
        }
    }
    return dstUnits;
}

export default function (testUnitArr) {
    let curRet = { code: false };

    if (!bCreate) {
        let err = client.connect(strAddr, OnCallBackFromGolang)

        if (err == null) {
            bCreate = true;

            group('C2GVerifyMessage', function () {
                curRet = C2GVerifyMessage(client)
            });

        } else {
            console.info("Failed to create conn.")
        }
    }
    else {
        let randIdx = randomInterval(0, testUnitArr.length - 1);
        let curOnCallFuncName = testUnitArr[randIdx];

        switch (curOnCallFuncName) {
            case 'C2GTestMessage':
                group('C2GTestMessage', function () {
                    curRet = C2GTestMessage(client)
                });
                break;
            case 'C2GHeartbeatMessage':
                group('C2GHeartbeatMessage', function () {
                    curRet = C2GHeartbeatMessage(client)
                });
                break;
            default:
                break;
        }
    }

    susRate.add(curRet.code);
    if (!curRet.code) {
        errCounter.add(1);
    }
    sleep(randomInterval(3, 3));


    //tcp.write(conn, 'some data\n'); // or tcp.writeLn(conn, 'some data')

}

export function OnCallBackFromGolang(msg, sus) {
    console.info(msg, sus)
}
