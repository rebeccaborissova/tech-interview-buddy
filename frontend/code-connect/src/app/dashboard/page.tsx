"use client";

import { useEffect, useState } from "react";
import styles from "./Dashboard.module.css";
import toast, { Toaster } from 'react-hot-toast';
import { generateToken, messaging } from '../firebase/firebase.js';
import { onMessage } from "firebase/messaging";
import { getToken } from "../utils/token";
import { useRouter } from "next/navigation";
import { generateJitsiRoom } from "../utils/jitsi";

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
  const [isLoading, setIsLoading] = useState(true);
  const [sessionToken, setSessionToken] = useState<string | null>(null);
  const [activeUsers, setActiveUsers] = useState<User[] | string>([]);
  const [isIncomingCallPopupOpen, setIsIncomingCallPopupOpen] = useState(false);
  const [incomingCallUser, setIncomingCallUser] = useState<string | null>(null);
  const [userPushToken, setUserPushToken] = useState<string | null>(null);
  const router = useRouter();

  const handleSelectUser = async (user: User) => {
    const tokeny = await getPushToken(user.Email);
    console.log("Successfully got Push token of user:", tokeny);

    sendPushNotification(tokeny || "");

    setUserPushToken(tokeny);
    setSelectedUser(user);
    setIsPopupOpen(true);
  };

  const handleClosePopup = () => {
    setIsPopupOpen(false);
    setSelectedUser(null);
  };

  useEffect(() => {
    setPushToken();
    
    onMessage(messaging, (payload) => {
      console.log(payload);
      if (payload?.notification?.title === 'Incoming Call') {
        setIncomingCallUser(payload?.notification?.body || "Unknown User");
        setIsIncomingCallPopupOpen(true);
      } else {
        toast(payload?.notification?.body || "Unknown notification");
      }
    })
    const token = getToken();
    if (!token) {
      router.push("/login");
    }
    setSessionToken(token);
  }, []);

  useEffect(() => {
    console.log("Session token:", sessionToken);
    if (sessionToken) {
      fetchActiveUsers();
    }
  }, [sessionToken]);

  const setPushToken = async () => {
    const pushToken = await generateToken();
    console.log("Push token:", pushToken);
  };

  const sendPushNotification = async (pushToken: string) => {
    try {
      console.log('Preparing fetch request to /api/notification');
      const response = await fetch("/api/notification", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          pushToken
        }),
        credentials: 'include',
        mode: 'cors'
      });

      console.log('Response status:', response.status);
      console.log('Response headers:', Object.fromEntries(response.headers));
  
      if (!response.ok) {
        throw new Error(`HTTP error, status: ${response.status}`);
      }
  
      const result = await response.json();
      console.log("Push notification sent:", result);
  
    } catch (error) {
      console.error("Error sending push notification:", error);
    }
  }

  const fetchActiveUsers = async () => {
    console.log("sessionToken in api call: ", sessionToken);
    try {
      const response = await fetch("http://localhost:8000/app/activeusers", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json',
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

  const getPushToken = async (username: string): Promise<string | null> => {
    try {
      const response = await fetch("http://localhost:8000/app/getpushtoken", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          "username": username
        }),
        credentials: 'include',
        mode: 'cors'
      });
  
      if (!response.ok) {
        throw new Error(`HTTP error, status: ${response.status}`);
      }
  
      const result = await response.json();
      if (result.Code === 200) {
        return result.Message;
      }

      console.log(result);
      return result;
  
    } catch (error) {
      console.error("Error fetching push token:", error);
      return null;
    }
  };

  const handleAcceptCall = () => {
    //Navigate to a video call page
    const jitstRoomUrl = generateJitsiRoom();

    setIsIncomingCallPopupOpen(false);
    setIncomingCallUser(null);

    router.push(jitstRoomUrl);
  };

  return (
    <div className={styles.container}>
      <Toaster position="top-right" />
      
      {/* View My Profile Button */}
      <button 
        className={styles.viewProfileButton} 
        onClick={() => router.push('/profile-page')}
      >
        View My Profile
      </button>
      
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


      {/* Outgoing Call Request Popup */}
      {isPopupOpen && selectedUser && (
        <div className={styles.popupOverlay}>
          <div className={styles.popup}>
            <h3 className={styles.popupTitle}>
               Profile for {selectedUser.FirstName} {selectedUser.LastName}
            </h3>
            <p className={styles.popupDescription}>
              Year: {selectedUser.Year}<br/>
              DSA Experience: {selectedUser.TakenDSA ? 'Yes' : 'No'}<br/>
              Description: {selectedUser.Description}
              Joe?: {userPushToken}
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