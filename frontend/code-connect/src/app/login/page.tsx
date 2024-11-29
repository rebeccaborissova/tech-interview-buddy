"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { getToken } from "../utils/token";

const Login = () => {
  const router = useRouter();
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState<string>("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email.trim() || !password.trim() ) {
      setError("Please enter your email and password.");
      return;
    }
    console.log("Sending data:", { email, password }); // debugging log to see if correct data sent to POST request
    try {
      const response = await fetch("http://localhost:8000/account/login", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          "Username": email,
          "Authorization": password
        })
      });

      console.log(response);
      const result = await response.json();
      console.log(result);
      const token = result.Session;

      const cookiesToken = getToken();
      if(!cookiesToken) {
        const now = new Date();
        now.setTime(now.getTime() + 10 * 60 * 1000); //expires in 10 minutes
        const expires = now.toUTCString();
        document.cookie = `session_token=${token}; expires=${expires}; path=/;`;
      }

      console.log(document.cookie);

      if (response.ok) {
        setError(" ")
        if (result.Session) {
          setSuccess("Username and password correct. This user exists in the database.");
          setTimeout(() => router.push("/dashboard"));
        } else {
          setError(result.Message)
        }
      } else {
        setError(result.Message)
      }
    } catch (err) {
      console.error("Network request failed:", err);
      setError("Network response error.");
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
    <div style={styles.link}>
      <a href="/sign-up" style={styles.link}>New User? Create Account</a>
    </div>
    {error && <p style={styles.error}>{error}</p>}
    {success && <p style={styles.success}>{success}</p>}
    <button type="submit" style={styles.button}>Login</button>
    </form>
    </div>
    </div>
  );
};

// inline styles for simplicity
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
    color: "#000",
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
  success: {
    color: "green",
    margin: "10px 0",
  },
  link: {
    color: "#0070f3",
    textDecoration: "none",
    fontSize: "14px",
    transition: "color 0.3s",
  },
};

export default Login;

