import 'dotenv/config';
import * as http from "node:http";

import {getDockerData, getSystemData} from "./systemService.js";

// Set the port
const port = process.env.PORT || 3000;

const requestListener = async function (req, res) {
    try {
        // Docker endpoint
        // This endpoint returns the Docker data in JSON format
        if (req.url === '/docker') {
            const data = await getDockerData();
            res.setHeader('Content-Type', 'application/json');
            res.writeHead(200);
            res.end(JSON.stringify(data));
            return;
        }

        // Metrics endpoint
        // This endpoint returns the system data in JSON format
        // The services parameter is a comma-separated list of services to monitor
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
            message: error.message || "Error getting system data",
        }));
    }
}

const server = http.createServer(requestListener);

server.listen(port, () => {
    console.log(`Metrics agent running on port ${port}`)
});
