import 'dotenv/config';
import * as http from "node:http";

import { getSystemData } from "./systemService.js";

// Set the port
const port = process.env.PORT || 3000;

// Metrics endpoint
// This endpoint returns the system data in JSON format
// The services parameter is a comma-separated list of services to monitor
const requestListener = async function (req, res) {
    try {
        const data = await getSystemData({
            services: process.env.SERVICES || 'nginx,mysql',
        });

        res.setHeader('Content-Type', 'application/json');
        res.writeHead(200);
        res.end(JSON.stringify(data));
    } catch (error) {
        console.error("Error getting system data", error);

        res.setHeader('Content-Type', 'application/json');
        res.writeHead(500);
        res.end(JSON.stringify({
            error: true,
            message: "Error getting system data",
        }));
    }
}

const server = http.createServer(requestListener);

server.listen(port, null, () => {
    console.log(`Metrics agent running on port ${port}`)
});
