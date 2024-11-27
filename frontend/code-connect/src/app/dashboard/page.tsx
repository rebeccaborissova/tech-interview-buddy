"use client";

import { useState } from "react";
import styles from "./Dashboard.module.css";

const Dashboard = () => {
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<string | null>(null);
  const [isProfileBarExpanded, setIsProfileBarExpanded] = useState(false);

  const activeUsers = ["Rebecca", "Tim", "Sarah", "Isa", "Gabriel", "Anna"]; // Example active users

  const handleRequestUser = () => {
    setIsPopupOpen(true);
  };

  const handleSelectUser = (user: string) => {
    setSelectedUser(user);
    // Ensure the popup closes only after state updates
    setTimeout(() => setIsPopupOpen(false), 0);
  };

  return (
    <div className={styles.container}>
      {/* Left Panel */}
      <div className={styles.leftPanel}>
        <h2 className={styles.panelTitle}>Active Users</h2>
        <ul className={styles.userList}>
          {activeUsers.map((user) => (
            <li key={user} className={styles.userListItem}>
              <div className={styles.greenCircle}></div>
              {user}
            </li>
          ))}
        </ul>
      </div>

      {/* Main Content */}
      <div className={styles.mainContent}>
        <h1 className={styles.instructionsTitle}>CodeConnect Dashboard</h1>
        <p className={styles.instructionsText}>
          Follow the instructions below to practice your technical interviewing skills and land your dream role!
        </p>
        <ul className={styles.instructionsList}>
          <li>1. View the list of active users in the left panel and view their profile by pressing on a user.</li>
          <li>2. Press the "Request a User" button to initiate a video call request to a specific user.</li>
          <li>3. After requesting a user, begin the video call by pressing on the video call icon.</li>
          <li>4. While on the video call, communicate your question-type preferences to the interviewer so they can select an appropriate question from LeetCode.</li>
        </ul>
        <button onClick={handleRequestUser} className={styles.button}>
          Request a User
        </button>
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
            <h3 className={styles.popupTitle}>Request a Video Call</h3>
            <ul className={styles.popupUserList}>
              {activeUsers.map((user) => (
                <li
                  key={user}
                  onClick={() => handleSelectUser(user)}
                  className={styles.popupUserListItem}
                >
                  {user}
                </li>
              ))}
            </ul>
            <button
              onClick={() => setIsPopupOpen(false)}
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
