"use client";

import { useState, useEffect } from "react";
import styles from "./Dashboard.module.css";

interface User {
  name: string;
}

const Dashboard = () => {
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<string | null>(null);
  const [isProfileBarExpanded, setIsProfileBarExpanded] = useState(false);
  const [activeUsers, setActiveUsers] = useState<string[]>([]); // List of active users
  const [error, setError] = useState<string | null>(null); // Error message

  useEffect(() => {
    const fetchActiveUsers = async () => {
      try {
        const response = await fetch("http://localhost:8000/app/activeusers", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        });

        const result = await response.json(); // Parse JSON response
        console.log(result)
        if (response.ok) {
          setError(""); // Clear any previous errors

          // Assuming the result is an array of user objects with the required fields
          setActiveUsers(result);
        } else {
          // If the response is not OK, show error message from backend
          setError(result.Message || "Failed to fetch active users.");
        }
      } catch (err) {
        console.error("Network request failed:", err); // Log network error
        setError("Unable to load active users. Please check your network.");
      }
    };

    fetchActiveUsers();
  }, []); // Run on component mount

  const handleSelectUser = (user: string) => {
    setSelectedUser(user);
    setIsPopupOpen(true); // Open popup for the selected user
  };

  const handleClosePopup = () => {
    setIsPopupOpen(false);
    setSelectedUser(null);
  };

  return (
    <div className={styles.container}>
      {/* Left Panel */}
      <div className={styles.leftPanel}>
        <h2 className={styles.panelTitle}>Active Users</h2>
        {error ? (
          <p className={styles.errorText}>{error}</p>
        ) : (
          <ul className={styles.userList}>
            {activeUsers.map((user) => (
              <li
                key={user}
                className={styles.userListItem}
                onClick={() => handleSelectUser(user)}
              >
                <div className={styles.greenCircle}></div>
                {user}
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Main Content */}
      <div className={styles.mainContent}>
        <div className={styles.instructionsBackground}>
          <h1 className={styles.instructionsTitle}>CodeConnect Dashboard</h1>
          <p className={styles.instructionsText}>
            Follow the instructions below to practice your technical interviewing skills and land your dream role!
          </p>
          <ul className={styles.instructionsList}>
            <li>1. View the list of active users in the left panel and view their profile by pressing on a user.</li>
            <li>2. Initiate a video call request to a specific user.</li>
            <li>3. While on the video call, communicate your question-type preferences to the interviewer so they can select an appropriate question from LeetCode.</li>
          </ul>
        </div>
      </div>

      {/* Profile Bar */}
      <div
        className={`${styles.profileBar} ${
          isProfileBarExpanded ? styles.profileBarExpanded : ""
        }`}
        onClick={() => setIsProfileBarExpanded(!isProfileBarExpanded)}
      >
        <div className={styles.profileIcon}>👤</div>
        {isProfileBarExpanded && (
          <div className={styles.profileContent}>
            <h3 className={styles.profileContentTitle}>User Profile</h3>
            <p className={styles.profileContentText}>
              Select options and manage your profile here.
            </p>
          </div>
        )}
      </div>

      {/* Popup Modal */}
      {isPopupOpen && (
        <div className={styles.popupOverlay}>
          <div className={styles.popup}>
            <h3 className={styles.popupTitle}>
              Video Call Request for {selectedUser}
            </h3>
            <p className={styles.popupDescription}>
              Click the button below to request a video call.
            </p>
            <button className={styles.videoCallButton}>
              Request Video Call
            </button>
            <button
              onClick={handleClosePopup}
              className={styles.popupCloseButton}
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
