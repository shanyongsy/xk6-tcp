import { sleep, group } from 'k6';
import { Counter, Rate } from 'k6/metrics';
import tcp from 'k6/x/tcp';
import { OnGetTime, OnNormalFunc } from '../test_unit/func_unit.js';
import { randStr, randomInterval } from '../utils/random.js';

const susRate = new Rate('susRate');
const errCounter = new Counter('errCounter');

const strAddr = '127.0.0.1:8000';

let userName = randStr();
let client = new tcp.Client();
let bCreate = false;

export let options = {
    stages: [
        { duration: '1s', target: 1000 },
        { duration: '20s', target: 1000 },
    ],
}

export function setup() {
    console.log('init test unit array');
    const srcUnits = [
        {weight: 10, name: 'OnGetTime'},
        {weight: 10, name: 'OnNormalFunc'},
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
    let curRet = {code: false};

    if(! bCreate)
    {
        let err = client.connect(strAddr, OnRevMsg)

        if (err == null){
            bCreate = true;

            group('OnNormalFunc', function(){
                curRet = OnNormalFunc(client, userName)
            });

        } else {
            console.info("Failed to create conn.")
            return;
        }
    }
    else
    {
        let randIdx = randomInterval(0, testUnitArr.length - 1);
        let curTestName = testUnitArr[randIdx];

        switch(curTestName){
            case 'OnGetTime':
                group('OnGetTime', function(){
                    curRet = OnGetTime(client)
                });
                break;
            case 'OnNormalFunc':
                group('OnNormalFunc', function(){
                    curRet = OnNormalFunc(client, userName)
                });
                break;
            default:
                break;
        }
    }

    susRate.add(curRet.code);
    if(!curRet.code){
        errCounter.add(1);
    }
    sleep(randomInterval(3, 3));


    //tcp.write(conn, 'some data\n'); // or tcp.writeLn(conn, 'some data')

}

export function OnRevMsg(msg){
    //let strMsg = byteToString(msg)
    //console.info(strMsg);
}