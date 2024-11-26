// default root page: home page

"use client";

import Navbar from "../utils/navbar";

const HomePage = () => {
  return (
    <div>
      <Navbar />
      <div style={styles.container}>
        <h1 style={styles.title}>Welcome to CodeConnect!</h1>
        <p style={styles.subtitle}>Connecting developers through learning and practice.</p>
      </div>
    </div>
  );
};

const styles = {
  container: {
    textAlign: "center" as "center",
    padding: "2rem",
    background: "linear-gradient(to right, #6a11cb, #2575fc)",
    height: "100vh",
    color: "white",
  },
  title: {
    fontSize: "3rem",
    fontWeight: "bold" as "bold",
  },
  subtitle: {
    fontSize: "1.5rem",
    marginTop: "1rem",
  },
};

export default HomePage;
