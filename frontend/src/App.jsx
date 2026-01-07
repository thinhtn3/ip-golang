import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./context/authProvider";
import Auth from "./page/Auth";
import Home from "./page/Home";
import Chat from "./page/Chat";
import "./App.css";

function App() {
  const { user, loading } = useAuth();

  if (loading) return <div>Loading...</div>;

  return (
    <Routes>
      <Route path="/" element={user ? <Home /> : <Navigate to="/login" />} />
      <Route path="/login" element={!user ? <Auth /> : <Navigate to="/" />} />
      <Route path="/chat/:sessionId" element={<Chat />} />
    </Routes>
  );
}

export default App;
