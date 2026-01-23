//use langchain to generate response from gemini

import dotenv from "dotenv";
dotenv.config();

import { ChatGoogleGenerativeAI } from "@langchain/google-genai";

class Gemini {
    constructor() {
        this.gemini = new ChatGoogleGenerativeAI({
            model: "gemini-3-flash-preview",
            apiKey: process.env.GEMINI_API_KEY,
        });
    }

    async invoke(instruction, input = {}) {
        this.prompt = instruction;
        this.chain = this.prompt.pipe(this.gemini);
        const response = await this.chain.invoke(input);
        return response;
    }
}

const gemini = new Gemini();
export default gemini;