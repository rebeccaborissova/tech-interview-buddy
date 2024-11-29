"use client";

import { useState } from "react";
import styles from "./Dashboard.module.css";
import { getToken } from "../utils/token";
import { useEffect } from "react";
import { useRouter } from "next/navigation";

const Dashboard = () => {
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<string | null>(null);
  const [isProfileBarExpanded, setIsProfileBarExpanded] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [sessionToken, setSessionToken] = useState<string | null>(null);
  const router = useRouter();

  const activeUsers = ["Rebecca", "Tim", "Sarah", "Isa", "Gabriel", "Anna"]; // Example active users

  const handleSelectUser = (user: string) => {
    setSelectedUser(user);
    setIsPopupOpen(true); // Open popup for the selected user
  };

  const handleClosePopup = () => {
    setIsPopupOpen(false);
    setSelectedUser(null);
  };

  useEffect(() => {
    const token = getToken();
    if(!token) {
      router.push("/login");
    }
    setSessionToken(token);
  }, []);

  useEffect(() => {
    console.log("Session token:", sessionToken);
    if(sessionToken) {
      fetchActiveUsers();
    }
  }, [sessionToken]);

  const fetchActiveUsers = async () => {
    console.log("sessionToken in api call: ", sessionToken)
    try {
      const response = await fetch("http://localhost:8000/app/activeusers", {
        method: "POST",
        headers: {
          'Content-Type':'application/json',
        },
        credentials: 'include',
        mode: 'cors'
      });

      console.log(response);
      const result = await response.json();
      console.log(result);

    } catch (error) {
      console.error("Error fetching active users:", error);
    }
  };

  return (
    <div className={styles.container}>
      {/* Left Panel */}
      <div className={styles.leftPanel}>
        <h2 className={styles.panelTitle}>Active Users</h2>
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
