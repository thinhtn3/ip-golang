import React from "react";


export default function Chatbox({ sessionId, messages, setMessages, handleSubmit, message, setMessage }) {

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
          onClick={handleSubmit}
        >
          Send
        </button>
      </div>
    </div>
  );
}
