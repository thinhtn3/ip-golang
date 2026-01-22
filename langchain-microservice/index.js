import express from 'express';
import dotenv from 'dotenv';
import gemini from './gemini.js';
import { HumanMessage, AIMessage } from '@langchain/core/messages';
dotenv.config();

const app = express();

app.use(express.json());

app.post("/generate", async (req, res) => {
    const { body } = req.body;
    const chain = body.map((m) => {
        if (m.role === "user") {
            return new HumanMessage(m.message);
        } else if (m.role === "assistant") {
            return new AIMessage(m.message);
        }
    })
    const response = await gemini.invoke(chain);
    console.log(response.content);
    res.json({
        content: response.content,
        role: "assistant"
    });
});

app.listen(process.env.PORT || 3000, () => {
    console.log(`Server is running on port ${process.env.PORT}`);
});