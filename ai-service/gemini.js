//use langchain to generate response from gemini

import dotenv from "dotenv";
dotenv.config();

import { ChatGoogleGenerativeAI } from "@langchain/google-genai";
import interviewerPrompt from "./instruction.js";

class Gemini {
    constructor() {
        this.gemini = new ChatGoogleGenerativeAI({
            model: "gemini-3-flash-preview",
            apiKey: process.env.GEMINI_API_KEY,
        });
        this.prompt = interviewerPrompt;
        this.chain = this.prompt.pipe(this.gemini);
    }

    async invoke(prompt) {
        const response = await this.chain.invoke({ input: prompt });
        return response;
    }
}

const gemini = new Gemini();
export default gemini;