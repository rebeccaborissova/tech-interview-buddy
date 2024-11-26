// default root page: home page

"use client";

import Navbar from "../utils/navbar";

const HomePage = () => {
  return (
    <div>
      <Navbar />
      <div style={styles.container}>
        <div style={styles.leftSection}>
          <h1 style={styles.title}>Welcome to CodeConnect!</h1>
          <p style={styles.subtitle}>
            Connecting Gator developers through learning and practice.
          </p>
        </div>
        <div style={styles.rightSection}>
          <div style={styles.textBox}>
            <h2 style={styles.textBoxTitle}>About CodeConnect</h2>
            <p style={styles.textBoxContent}>
              CodeConnect is the go-to platform for mastering technical coding
              interviews. Sharpen your skills in data structures and algorithms
              with curated LeetCode and Hackerrank-style challenges. Pair up
              with a fellow Gator and get ready to ace your dream software role!
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

const styles = {
  container: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    height: "100vh",
    background: "linear-gradient(to right, #6a11cb, #2575fc)",
    padding: "2rem",
    color: "white",
  },
  leftSection: {
    flex: 1.3, // Slightly wider to extend text towards the right
    display: "flex",
    flexDirection: "column" as "column",
    justifyContent: "center",
    alignItems: "flex-start",
    padding: "2rem",
    textAlign: "left" as "left",
  },
  title: {
    fontSize: "4rem",
    fontWeight: "bold" as "bold",
    marginBottom: "1rem",
    lineHeight: "1.2",
  },
  subtitle: {
    fontSize: "1.8rem",
    marginTop: "0.5rem",
    lineHeight: "1.5",
    maxWidth: "90%", // Extended width for subtitle
  },
  rightSection: {
    flex: 1,
    display: "flex",
    justifyContent: "center" as "center",
    alignItems: "center" as "center",
    padding: "2rem",
  },
  textBox: {
    background: "white",
    borderRadius: "12px",
    padding: "3rem",
    boxShadow: "0 10px 30px rgba(0, 0, 0, 0.6)", // Enhanced shadow for contrast
    color: "#333",
    maxWidth: "500px",
    textAlign: "center" as "center",
  },
  textBoxTitle: {
    fontSize: "2.2rem",
    fontWeight: "bold" as "bold",
    marginBottom: "1.5rem",
    color: "#6a11cb",
  },
  textBoxContent: {
    fontSize: "1.4rem",
    lineHeight: "1.8",
  },
};

export default HomePage;
