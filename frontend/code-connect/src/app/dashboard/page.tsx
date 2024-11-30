"use client";

import { useEffect, useState } from "react";
import styles from "./Dashboard.module.css";
import toast, { Toaster } from 'react-hot-toast';
import { generateToken, messaging } from '../firebase/firebase.js'
//import { sendCallInvite } from '../firebase';
//import { requestNotificationPermission } from '../firebase/firebase.js';
import { onMessage } from "firebase/messaging";
import { getToken } from "../utils/token";
import { useRouter } from "next/navigation";

interface User {
  Email: string;
  FirstName: string;
  LastName: string;
  InvitedBy: string;
  TakenDSA: boolean;
  Year: number;
  Description: string;
}

const Dashboard = () => {
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [isProfileBarExpanded, setIsProfileBarExpanded] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [sessionToken, setSessionToken] = useState<string | null>(null);
  const [activeUsers, setActiveUsers] = useState<User[] | string>([]);
  const [isIncomingCallPopupOpen, setIsIncomingCallPopupOpen] = useState(false);
  const [incomingCallUser, setIncomingCallUser] = useState<string | null>(null);
  const router = useRouter();

  const handleSelectUser = (user: User) => {
    setSelectedUser(user);
    setIsPopupOpen(true); // Open popup for the selected user
  };

  const handleClosePopup = () => {
    setIsPopupOpen(false);
    setSelectedUser(null);
  };

  /*const handleVideoCallRequest = async () => {
    const token = 'RECIPIENT_FCM_TOKEN'; // Replace with the actual recipient's FCM token
    const callerName = 'Your Name'; // Replace with the actual caller's name

    try {
      const result = await sendCallInvite(token, callerName);
      if (result.data.success) {
        console.log('Call invite sent successfully:', result.data.jitsiRoomUrl);
      } else {
        console.error('Failed to send call invite:', result.data.error);
      }
    } catch (error) {
      console.error('Error sending call invite:', error);
    }
  };*/
        
  useEffect(() => {
    setPushToken();
    
    onMessage(messaging, (payload) => {
      console.log(payload);
      if (payload?.notification?.title === 'Incoming Call') {
        setIncomingCallUser(payload?.notification?.body || "Unknown User");
        setIsIncomingCallPopupOpen(true);
      } else {
        toast(payload?.notification?.body || "hi");
      }
    })
        
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

  const setPushToken = async () => {
    const pushToken = await generateToken();
    console.log("Push token:", pushToken);
  };

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

      const result = await response.json();
      const parsedUsers = JSON.parse(result.Message) as User[];
      if (parsedUsers && parsedUsers.length > 0) {
        setActiveUsers(parsedUsers);
      } else {
        setActiveUsers("No active users found.");
      }
      setIsLoading(false);

    } catch (error) {
      console.error("Error fetching active users:", error);
      setActiveUsers("No active users found.");
      setIsLoading(false);
    }
  };

  const handleAcceptCall = () => {
    // Add your call acceptance logic here
    setIsIncomingCallPopupOpen(false);
    setIncomingCallUser(null);
  };

  return (
    <div className={styles.container}>
      <Toaster position="top-right" />
      {/* Left Panel */}
      <div className={styles.leftPanel}>
        <h2 className={styles.panelTitle}>Active Users</h2>
        <ul className={styles.userList}>
          {isLoading ? (
            <li className={styles.messageItem}>Loading...</li>
          ) : typeof activeUsers === 'string' ? (
            <li className={styles.messageItem}>{activeUsers}</li>
          ) : (
            activeUsers.map((user) => (
              <li
                key={user.Email}
                className={styles.userListItem}
                onClick={() => handleSelectUser(user)}
              >
                <div className={styles.greenCircle}></div>
                {user.FirstName} {user.LastName}
              </li>
            ))
          )}
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

      {/* Outgoing Call Request Popup */}
      {isPopupOpen && selectedUser && (
        <div className={styles.popupOverlay}>
          <div className={styles.popup}>
            <h3 className={styles.popupTitle}>
              Video Call Request for {selectedUser.FirstName} {selectedUser.LastName}
            </h3>
            <p className={styles.popupDescription}>
              Year: {selectedUser.Year}<br/>
              DSA Experience: {selectedUser.TakenDSA ? 'Yes' : 'No'}<br/>
              Description: {selectedUser.Description}
            </p>
            <div className={styles.popupButtons}>
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
        </div>
      )}

      {/* Incoming Call Popup */}
      {isIncomingCallPopupOpen && (
        <div className={styles.popupOverlay}>
          <div className={styles.popup}>
            <h3 className={styles.popupTitle}>
              Incoming Call from {incomingCallUser}
            </h3>
            <div className={styles.popupButtons}>
              <button 
                className={styles.videoCallButton}
                onClick={handleAcceptCall}
              >
                Accept Video Call
              </button>
              <button
                onClick={() => setIsIncomingCallPopupOpen(false)}
                className={styles.popupCloseButton}
              >
                Decline
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
