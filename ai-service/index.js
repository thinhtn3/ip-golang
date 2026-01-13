import express from 'express';
import dotenv from 'dotenv';
import gemini from './gemini.js';
dotenv.config();

const app = express();

app.use(express.json());

app.post("/generate", async (req, res) => {
    const { prompt } = req.body;
    const response = await gemini.invoke(prompt);
    res.json({ response });
});

app.listen(process.env.PORT, () => {
    console.log(`Server is running on port ${process.env.PORT}`);
});