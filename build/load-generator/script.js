import http from 'k6/http';
import { sleep, check } from 'k6';
import { SharedArray } from 'k6/data';

// Read the configuration file
const configData = JSON.parse(open('/data/config.json'));

// Read the env var
const testName = __ENV.TEST_NAME || "test"

export let options = {
    vus: configData.vus,  // number of virtual users
    duration: configData.duration,  // duration of the test
    tags: {
        name: testName,
    },
};

const urlPrefix = configData.urlPrefix;

// Read the JSON file
const jsonData = new SharedArray('routes', function() {
    return JSON.parse(open('/data/routes.json'));
});

// Calculate the total weight
let totalWeight = 0;
for (let route of jsonData) {
    totalWeight += Number(route.weight);
}

export default function() {
    // Select a route based on the weights
    let targetWeight = Math.random() * totalWeight;
    let selectedRoute;
    for (let route of jsonData) {
        targetWeight -= Number(route.weight);
        if (targetWeight <= 0) {
            selectedRoute = route;
            break;
        }
    }

    // Perform an HTTP request to the selected route
    let res;
    const url = urlPrefix + selectedRoute.route;
    if (selectedRoute.method.toLowerCase() === 'get') {
        res = http.get(url);
    } else if (selectedRoute.method.toLowerCase() === 'post') {
        res = http.post(url);
    } else {
        console.error(`Unknown method: ${selectedRoute.method}`);
        return;
    }

    // Check the status of the response
    check(res, {
        'status was 200': (r) => r.status == 200,
        'transaction time OK': (r) => r.timings.duration < 200,
    });

    sleep(1);
}
