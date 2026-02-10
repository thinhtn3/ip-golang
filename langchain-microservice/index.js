import express from 'express';
import dotenv from 'dotenv';
import gemini from './gemini.js';
import { HumanMessage, AIMessage, SystemMessage } from '@langchain/core/messages';
import { summaryPrompt } from './instruction.js';
dotenv.config();

const app = express();

app.use(express.json());

app.post("/generate", async (req, res) => {
    const { summary, messages } = req.body;
    const chain = messages.map((m) => {
        if (m.role === "user") {
            return new HumanMessage(m.message);
        } else if (m.role === "assistant") {
            return new AIMessage(m.message);
        }
    })
    const response = await gemini.invoke([new SystemMessage(summary), ...chain]);
    console.log(response.content);
    res.json({
        content: response.content,
        role: "assistant"
    });
});

app.post("/summarize", async (req, res) => {
    // get summary and messages array then send it to gemini for summarization
    const { summary, messages } = req.body;
    if (!summary) {
        console.log("Received no summary");
    }
    console.log("Received summary: ", summary);
    console.log(summary, "FOO", messages);
    const response = await gemini.invoke(summaryPrompt, { summary: summary, messages: messages });
    console.log("Response: ", response.content);
    res.json({ content: response.content });
})

app.listen(process.env.PORT || 3000, () => {
    console.log(`Server is running on port ${process.env.PORT}`);
});