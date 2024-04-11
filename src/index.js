import express from 'express';
import { getSystemData } from "./systemService.js";

const app = express();
const port = 3000;

app.get('/', async (req, res) => {
    const data = await getSystemData();
    res.send(data);
})

app.listen(port, () => {
    console.log(`Metrics agent running on port ${port}`)
});
