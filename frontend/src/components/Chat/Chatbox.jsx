import React from "react";
import axios from "axios";
import { useState } from "react";
import { supabase } from "../../config/supabase";

export default function Chatbox({ sessionId }) {
  const apiUrl = import.meta.env.VITE_DEVELOPMENT_URL;

  const [message, setMessage] = useState("");
  const [role, setRole] = useState("user");

  const getAccessToken = async () => {
    const { data } = await supabase.auth.getSession()
    return data.session?.access_token
  }

  const handleSubmit = async () => {
    const accessToken = await getAccessToken()
    console.log(accessToken)
    const response = await axios.post(`${apiUrl}/chat/sessions/${sessionId}/messages`, {
      message: message,
      role: role,
    }, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      }
    })
    console.log(response)
  }

  return (
    <div className="flex flex-col h-[600px] w-full max-w-2xl mx-auto bg-slate-950 border border-slate-800 rounded-2xl overflow-hidden shadow-2xl p-4 text-white">
      <h2 className="text-xl font-bold mb-4">Chat Session: {sessionId}</h2>
      <div className="flex-1 overflow-y-auto border-t border-slate-800 pt-4">
        <p className="text-slate-400 italic">Chat interface boilerplate. Implement messaging logic here...</p>
      </div>
      <div className="mt-4 flex gap-2">
        <input 
          type="text" 
          placeholder="Type a message..." 
          className="flex-1 bg-slate-800 border-none rounded-xl px-4 py-2 outline-none"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
        />
        <button 
          className="bg-blue-600 px-4 py-2 rounded-xl opacity-50 cursor-not-allowed"
          disabled={message.length === 0}
          onClick={handleSubmit}
        >
          Send
        </button>
      </div>
    </div>
  );
}
