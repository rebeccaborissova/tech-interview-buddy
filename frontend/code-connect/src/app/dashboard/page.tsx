"use client";

import { useEffect, useState } from "react";
import styles from "./Dashboard.module.css";
import toast, { Toaster } from 'react-hot-toast';
import { generateToken, messaging } from '../firebase/firebase.js';
import { onMessage } from "firebase/messaging";
import { getToken, getUsername } from "../utils/cookies";
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
  const [incomingCallMessage, setIncomingCallMessage] = useState<string | null>(null);
  const [userPushToken, setUserPushToken] = useState<string | null>(null);
  const [jitsiRoom, setJitsiRoom] = useState<string | null>(null);
  const router = useRouter();

  const setPushToken = async () => {
    const pushToken = await generateToken();
    console.log("Push token:", pushToken);
  };

  const handleSelectUser = async (user: User) => {
    setSelectedUser(user);
    setIsPopupOpen(true);
  };

  const handleRequestCall = async () => {
    // Generate a Jitsi room
    const requestedJitsiRoom = generateJitsiRoom();

    // Send a push notification to the selected user
    const token = await getPushToken(selectedUser?.Email || "");
    sendPushNotification(token || "", requestedJitsiRoom || "", selectedUser?.Email || "");

    // Set the Jitsi room and user push token
    setJitsiRoom(requestedJitsiRoom);
    setUserPushToken(token);

    // Redirect to the Jitsi room
    router.push(jitsiRoom || "/dashboard");
  }

  const handleClosePopup = () => {
    setIsPopupOpen(false);
    setSelectedUser(null);
  };

  const handleAcceptCall = () => {
    setIsIncomingCallPopupOpen(false);
    //setIncomingCallMessage(null);

    router.push(jitsiRoom || "/dashboard");
  };

  useEffect(() => {
    setPushToken();
    
    // Listen for incoming messages
    onMessage(messaging, (payload) => {
      console.log("Message received. Payload:", payload);
      // If the notification is an incoming call, display a popup
      if (payload?.notification?.title?.startsWith('Incoming Call From')) {
        setJitsiRoom(payload?.notification?.body || "Unknown room");
        setIncomingCallMessage(payload?.notification?.title || "Unknown call");
        setIsIncomingCallPopupOpen(true);
      } else {
        toast(payload?.notification?.body || "Unknown notification");
      }
    })

    // Check if the user is logged in
    const token = getToken();
    if (!token) {
      router.push("/login");
    }

    setSessionToken(token);
    fetchActiveUsers();
  }, []);

  // Send a push notification to the user you want to video call
  const sendPushNotification = async (pushToken: string, jitsiRoom: string, requestedUsername: string) => {
    try {
      const response = await fetch("/api/notification", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          pushToken: pushToken,
          body: jitsiRoom,
          username: requestedUsername
        }),
        credentials: 'include',
        mode: 'cors'
      });
  
      if (!response.ok) {
        throw new Error(`HTTP error, status: ${response.status}`);
      }
  
      const result = await response.json();
      console.log("Push notification sent:", result);
  
    } catch (error) {
      console.error("Error sending push notification:", error);
    }
  }

  // Fetch the list of active users
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

      // Parse the list of active users
      const result = await response.json();
      const parsedUsers = JSON.parse(result.Message) as User[];

      // Exclude the user that matches document.cookie username
      const username = getUsername();
      const filteredUsers = parsedUsers.filter(user => user.Email !== username);

      // Check if the list of active users is valid and set the state with the list of active users
      if (parsedUsers && parsedUsers.length > 0) {
        setActiveUsers(filteredUsers);
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

  // Fetch the push token for a given username
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
      return result;
  
    } catch (error) {
      console.error("Error fetching push token:", error);
      return null;
    }
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
            </p>
            <div className={styles.popupButtons}>
              <button className={styles.videoCallButton} onClick={handleRequestCall}>
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
              {incomingCallMessage}
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