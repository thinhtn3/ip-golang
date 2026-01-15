import express from 'express';
import dotenv from 'dotenv';
import gemini from './gemini.js';
dotenv.config();

const app = express();

app.use(express.json());

app.post("/generate", async (req, res) => {
    console.log("server says hello");
    res.json({ message: "Hello from AI service" }, 200);
    // const { prompt } = req.body;
    // const response = await gemini.invoke(prompt);
    // res.json({ response });
});

app.listen(process.env.PORT || 3000, () => {
    console.log(`Server is running on port ${process.env.PORT}`);
});