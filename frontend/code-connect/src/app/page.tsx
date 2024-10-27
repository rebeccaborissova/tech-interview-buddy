"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

const Login = () => {
  const router = useRouter();
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email || !password) {
      setError("Please enter your email and password.");
      return;
    }

    try {
      // Placeholder for login logic
      console.log("Logging in with", { email, password });
      router.push("/dashboard");
    } catch (err) {
      setError("Login failed. Please check your credentials.");
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.formContainer}>
        <h1 style={styles.title}>CodeConnect</h1>
        <form onSubmit={handleLogin} style={styles.form}>
          <div style={styles.inputContainer}>
            <label htmlFor="email" style={styles.label}>Email</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              style={styles.input}
              required
            />
          </div>
          <div style={styles.inputContainer}>
            <label htmlFor="password" style={styles.label}>Password</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              style={styles.input}
              required
            />
          </div>
          {error && <p style={styles.error}>{error}</p>}
          <button type="submit" style={styles.button}>Login</button>
        </form>
      </div>
    </div>
  );
};

// Inline styles for simplicity
const styles = {
  container: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    height: "100vh",
    background: "linear-gradient(to right, #6a11cb, #2575fc)",
  },
  formContainer: {
    maxWidth: "400px",
    width: "100%",
    padding: "30px",
    textAlign: "center" as "center",
    border: "1px solid #ccc",
    borderRadius: "12px",
    backgroundColor: "#fff",
    boxShadow: "0 4px 15px rgba(0, 0, 0, 0.1)",
  },
  title: {
    margin: "0 0 20px 0",
    fontSize: "28px",
    fontWeight: "bold" as "bold",
    color: "#333",
  },
  form: {
    display: "flex",
    flexDirection: "column" as "column",
  },
  inputContainer: {
    margin: "15px 0",
  },
  label: {
    display: "block",
    marginBottom: "5px",
    fontSize: "14px",
    fontWeight: "600",
    color: "#444",
    textAlign: "left" as "left",
  },
  input: {
    width: "100%",
    padding: "12px",
    borderRadius: "8px",
    border: "1px solid #ccc",
    fontSize: "16px",
  },
  button: {
    padding: "12px",
    backgroundColor: "#0070f3",
    color: "#fff",
    border: "none",
    borderRadius: "8px",
    cursor: "pointer",
    fontSize: "16px",
    transition: "background 0.3s",
    marginTop: "10px",
  },
  buttonHover: {
    backgroundColor: "#005bb5",
  },
  error: {
    color: "red",
    margin: "10px 0",
  },
  logoContainer: {
    display: "flex",
    justifyContent: "center",
    gap: "10px",
    marginTop: "20px",
  },
  logo: {
    width: "30px",
    height: "30px",
  },
};

export default Login;
