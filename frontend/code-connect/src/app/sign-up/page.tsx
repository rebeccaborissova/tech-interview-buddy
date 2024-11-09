"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import styles from "./SignUp.module.css"; // Import the CSS module for hidden scrollbar

const SignUp = () => {
  const router = useRouter();
  const [firstName, setFirstName] = useState<string>("");
  const [lastName, setLastName] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [confirmPassword, setConfirmPassword] = useState<string>("");
  const [takenDSA, setTakenDSA] = useState<boolean>(false);
  const [schoolYear, setSchoolYear] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState<string>("");

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!firstName.trim() || !lastName.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
      setError("Please fill out all fields.");
      return;
    }

    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    console.log("Sending data:", { firstName, lastName, email, takenDSA, schoolYear, password });
    try {
      const response = await fetch("http://localhost:8000/register", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          "FirstName": firstName,
          "LastName": lastName,
          "Username": email,
          "TakenDSA": takenDSA,
          "SchoolYear": schoolYear,
          "Authorization": password
        })
      });

      const result = await response.json();
      if (response.ok) {
        setError("");
        setSuccess("Account created successfully. Redirecting...");
        setTimeout(() => router.push("/login"), 2000);
      } else {
        setError(result.Message);
      }
    } catch (err) {
      console.error("Network request failed:", err);
      setError("Network response error.");
    }
  };

  return (
    <div style={inlineStyles.container}>
      <div style={inlineStyles.formContainer}>
        <h1 style={inlineStyles.title}>Sign Up</h1>
        <div className={styles.scrollableForm}>
          <form onSubmit={handleSignUp} style={inlineStyles.form}>
            <div style={inlineStyles.inputContainer}>
              <label htmlFor="firstName" style={inlineStyles.label}>First Name</label>
              <input
                type="text"
                id="firstName"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                style={inlineStyles.input}
                required
              />
            </div>
            <div style={inlineStyles.inputContainer}>
              <label htmlFor="lastName" style={inlineStyles.label}>Last Name</label>
              <input
                type="text"
                id="lastName"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                style={inlineStyles.input}
                required
              />
            </div>
            <div style={inlineStyles.inputContainer}>
              <label htmlFor="email" style={inlineStyles.label}>Email</label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                style={inlineStyles.input}
                required
              />
            </div>
            <div style={inlineStyles.inputContainer}>
              <label htmlFor="takenDSA" style={inlineStyles.label}>Taken Data Structures and Algorithms?</label>
              <input
                type="checkbox"
                id="takenDSA"
                checked={takenDSA}
                onChange={(e) => setTakenDSA(e.target.checked)}
                style={inlineStyles.checkbox}
              />
            </div>
            <div style={inlineStyles.inputContainer}>
              <label htmlFor="schoolYear" style={inlineStyles.label}>School Year</label>
              <select
                id="schoolYear"
                value={schoolYear}
                onChange={(e) => setSchoolYear(e.target.value)}
                style={inlineStyles.input}
                required
              >
                <option value="">Select Year</option>
                <option value="Freshman">Freshman</option>
                <option value="Sophomore">Sophomore</option>
                <option value="Junior">Junior</option>
                <option value="Senior">Senior</option>
              </select>
            </div>
            <div style={inlineStyles.inputContainer}>
              <label htmlFor="password" style={inlineStyles.label}>Password</label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                style={inlineStyles.input}
                required
              />
            </div>
            <div style={inlineStyles.inputContainer}>
              <label htmlFor="confirmPassword" style={inlineStyles.label}>Confirm Password</label>
              <input
                type="password"
                id="confirmPassword"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                style={inlineStyles.input}
                required
              />
            </div>
            <div style={inlineStyles.link}>
              <a href="/login" style={inlineStyles.link}>Already have an account? Login</a>
            </div>
            {error && <p style={inlineStyles.error}>{error}</p>}
            {success && <p style={inlineStyles.success}>{success}</p>}
            <button type="submit" style={inlineStyles.button}>Sign Up</button>
          </form>
        </div>
      </div>
    </div>
  );
};

// Inline styles
const inlineStyles = {
  container: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    height: "100vh",
    background: "linear-gradient(to right, #6a11cb, #2575fc)",
  },
  formContainer: {
    maxWidth: "700px",
    width: "100%",
    padding: "40px",
    textAlign: "center" as "center",
    border: "1px solid #ccc",
    borderRadius: "12px",
    backgroundColor: "#fff",
    boxShadow: "0 4px 15px rgba(0, 0, 0, 0.1)",
  },
  title: {
    margin: "0 0 20px 0",
    fontSize: "32px",
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
  checkbox: {
    width: "20px",
    height: "20px",
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

export default SignUp;
