import http from 'k6/http';

export const options = {
    duration: '10s',
    vus: 50,
    thresholds: {
        http_req_failed: ['rate<0.001'], // http errors should be less than 1%
        http_req_duration: ['p(99)<500'], // 95 percent of response times must be below 500ms
    },
};

export default function () {
    const url = 'http://dns:8090/v1/locate';
    const payload = JSON.stringify({
        "x": "123.12",
        "y": "456.56",
        "z": "789.89",
        "vel": "20.0"
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    http.post(url, payload, params);
}