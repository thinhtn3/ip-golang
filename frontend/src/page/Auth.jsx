import { useState } from "react";
import { supabase } from "../config/supabase";
import axios from "axios";

export default function Auth() {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      if (isLogin) {
        // Sign in with password
        const { data, error } = await supabase.auth.signInWithPassword({
          email,
          password,
        });

        if (error) throw error;
        await axios.post(
          "http://localhost:8080/user/profile", {},
          {
            headers: {
              Authorization: `Bearer ${data.session.access_token}`,
            },
          }
        );
        console.log("Logged in:", data.user);
      } else {
        // Sign up with password
        const { data, error } = await supabase.auth.signUp({
          email,
          password,
        });

        if (error) throw error;
        alert("Check your email for verification link!");
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h1 style={styles.title}>
          {isLogin ? "Welcome Back" : "Create Account"}
        </h1>

        {error && <div style={styles.error}>{error}</div>}

        <form onSubmit={handleSubmit} style={styles.form}>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={styles.input}
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={styles.input}
            required
            minLength={6}
          />
          <button type="submit" style={styles.button} disabled={loading}>
            {loading ? "Loading..." : isLogin ? "Sign In" : "Sign Up"}
          </button>
        </form>

        <p style={styles.toggle}>
          {isLogin ? "Don't have an account? " : "Already have an account? "}
          <button
            onClick={() => setIsLogin(!isLogin)}
            style={styles.toggleButton}
          >
            {isLogin ? "Sign Up" : "Sign In"}
          </button>
        </p>
      </div>
    </div>
  );
}

const styles = {
  container: {
    minHeight: "100vh",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    background: "linear-gradient(135deg, #1a1a2e 0%, #16213e 100%)",
  },
  card: {
    background: "#0f0f23",
    padding: "3rem",
    borderRadius: "16px",
    boxShadow: "0 25px 50px -12px rgba(0, 0, 0, 0.5)",
    width: "100%",
    maxWidth: "400px",
    border: "1px solid #2a2a4a",
  },
  title: {
    color: "#fff",
    fontSize: "1.75rem",
    fontWeight: "600",
    marginBottom: "2rem",
    textAlign: "center",
  },
  form: {
    display: "flex",
    flexDirection: "column",
    gap: "1rem",
  },
  input: {
    padding: "0.875rem 1rem",
    fontSize: "1rem",
    borderRadius: "8px",
    border: "1px solid #3a3a5a",
    background: "#1a1a3e",
    color: "#fff",
    outline: "none",
    transition: "border-color 0.2s",
  },
  button: {
    padding: "0.875rem",
    fontSize: "1rem",
    fontWeight: "600",
    borderRadius: "8px",
    border: "none",
    background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
    color: "#fff",
    cursor: "pointer",
    marginTop: "0.5rem",
    transition: "transform 0.2s, box-shadow 0.2s",
  },
  error: {
    background: "#ff4757",
    color: "#fff",
    padding: "0.75rem",
    borderRadius: "8px",
    marginBottom: "1rem",
    fontSize: "0.875rem",
  },
  toggle: {
    color: "#888",
    textAlign: "center",
    marginTop: "1.5rem",
    fontSize: "0.875rem",
  },
  toggleButton: {
    background: "none",
    border: "none",
    color: "#667eea",
    cursor: "pointer",
    fontWeight: "600",
    fontSize: "0.875rem",
  },
};
