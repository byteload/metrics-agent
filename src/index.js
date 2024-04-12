import 'dotenv/config';
import express from 'express';
import { getSystemData } from "./systemService.js";

// Initialize the Express app
const app = express();

// Set the port
const port = process.env.PORT || 3000;

// Metrics endpoint
// This endpoint returns the system data in JSON format
// The services parameter is a comma-separated list of services to monitor
app.get('/', async (req, res) => {
    try {
        const data = await getSystemData({
            services: process.env.SERVICES || 'nginx,mysql',
        });

        res.send(data);
    } catch (error) {
        console.error("Error getting system data", error);

        res.status(500).send({
            error: true,
            message: "Error getting system data",
        });
    }
})

// Start the Express app
app.listen(port, () => {
    console.log(`Metrics agent running on port ${port}`)
});
