import React from "react";
import axios from "axios";
import { useState, useEffect } from "react";
import { supabase } from "../../config/supabase";

export default function Chatbox({ sessionId }) {
  const apiUrl = import.meta.env.VITE_DEVELOPMENT_URL;
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const [role, setRole] = useState("user");

  const getAccessToken = async () => {
    const { data } = await supabase.auth.getSession()
    return data.session?.access_token
  }

  useEffect(() => {
    const fetchMessages = async () => {
      const accessToken = await getAccessToken()
      const response = await axios.get(`${apiUrl}/chat/sessions/${sessionId}/messages`, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        }
      })
      console.log(response)
      if (response.status === 200) {
        setMessages(response.data.messages)
      }
      else {
        console.error('Error fetching messages:', response.data.message)
      }
    }
    fetchMessages()
  }, [])

  const handleSubmit = async () => {
    const accessToken = await getAccessToken()
    const response = await axios.post(`${apiUrl}/chat/sessions/${sessionId}/messages`, {
      message: message,
      role: role,
    }, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      }
    })
    if (response.status === 200) {
      setMessages([...messages, response.data.chat])
    }
    else {
      console.error('Error sending message:', response.data.message)
    }
  }

  return (
    <div className="flex flex-col h-[600px] w-full max-w-2xl mx-auto bg-slate-950 border border-slate-800 rounded-2xl overflow-hidden shadow-2xl p-4 text-white">
      <h2 className="text-xl font-bold mb-4">Chat Session: {sessionId}</h2>
      <div className="flex-1 overflow-y-auto border-t border-slate-800 pt-4">
        {messages.map((message) => (
          <div key={message.id} className={`${message.role === "user" ? "text-right" : "text-left"} ${message.role === "user" ? "bg-blue-600" : "bg-slate-800"} p-2 rounded-lg`}>
            <p>{message.message}</p>
            <p>{message.role}</p>
          </div>
        ))}
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
