import http from 'k6/http';
import { check } from 'k6';
import { randomItem, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.0.0/index.js';

const products = [
  '0PUK6V6EV0', '1YMWWN1N4O', '2ZYFJ3GM2N', '66VCHSJNUP',
  '6E92ZMYYFZ', '9SIQT8TOJO', 'L9ECAV7KIM', 'LS4PSXUNUM', 'OLJCESPC7Z'
];

// Retrieve the total arrival rate from an environment variable
const totalArrivalRate = __ENV.TOTAL_ARRIVAL_RATE || 100; // Default to 100 if not provided

// Retrieve the total duration from an environment variable
const duration = __ENV.DURATION || "10m"; // Default to 100 if not provided

const frontendAddr = __ENV.FRONTEND_ADDR || "frontend:8080";
const baseURL = "http://" + frontendAddr

const indexRoute = __ENV.INDEX_ROUTE || "/"

const testName = __ENV.TEST_NAME || "test"

// Define the weights for each route
const weights = {
  index: 1,
  setCurrency: 2,
  browseProduct: 10,
  addToCart: 2,
  viewCart: 3,
  checkout: 1,
};

// Calculate the sum of all weights
const totalWeight = Object.values(weights).reduce((a, b) => a + b, 0);

// Calculate individual route rates
const rates = {}
Object.entries(weights).forEach(([key, value]) => {
  rates[key] = Math.round((parseInt(totalArrivalRate) * value) / totalWeight)
})

const scenarios = {}
Object.entries(rates).forEach(([key, value]) => {
  scenarios[key] = {
    executor: 'constant-arrival-rate',
    rate: value,
    duration: duration,
    preAllocatedVUs: value,
    exec: key,
  }
})

export let options = {
  scenarios: scenarios,
  tags: {
    name: testName,
  },
};

export function index() {
  let res = http.get(baseURL + indexRoute);
  check(res, { 'status was 200': (r) => r.status == 200 });
}

export function setCurrency() {
  const currencies = ['EUR', 'USD', 'JPY', 'CAD'];
  let res = http.post(baseURL + "/setCurrency", { currency_code: randomItem(currencies) });
  check(res, { 'status was 200': (r) => r.status == 200 });
}

export function browseProduct() {
  let res = http.get(`${baseURL}/product/${randomItem(products)}`);
  check(res, { 'status was 200': (r) => r.status == 200 });
}

export function viewCart() {
  let res = http.get(baseURL + "/cart");
  check(res, { 'status was 200': (r) => r.status == 200 });
}

export function addToCart() {
  let product = randomItem(products);
  http.get(`${baseURL}/product/${product}`);
  let res = http.post(baseURL + "/cart", { product_id: product, quantity: randomIntBetween(1, 10) });
  check(res, { 'status was 200': (r) => r.status == 200 });
}

export function checkout() {
  addToCart();
  let res = http.post(baseURL + "/cart/checkout", {
    email: 'someone@example.com',
    street_address: '1600 Amphitheatre Parkway',
    zip_code: '94043',
    city: 'Mountain View',
    state: 'CA',
    country: 'United States',
    credit_card_number: '4432-8015-6152-0454',
    credit_card_expiration_month: '1',
    credit_card_expiration_year: '2039',
    credit_card_cvv: '672',
  });
  check(res, { 'status was 200': (r) => r.status == 200 });
}
