import { initializeApp } from "firebase/app";
import { getMessaging, getToken, onMessage } from "firebase/messaging";
import { getFunctions, httpsCallable } from "firebase/functions";

// Your web app's Firebase configuration
const firebaseConfig = {
    apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
    authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
    projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
    storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
    appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID
  };

// Initialize Firebase
const app = initializeApp(firebaseConfig, "firebase-app");
export const messaging = getMessaging(app);

export const generateToken = async () => {
    const permission = await Notification.requestPermission();
    console.log("notifications permission: ", permission)

    let token = undefined;

    //only if notification permission was granted, we can generate the token
    if (permission == "granted") {
        token = await getToken(messaging, {
            vapidKey: process.env.NEXT_PUBLIC_VAPID_KEY
        });
    
        console.log("messaging token: ", token)
    }

    try {
        const response = await fetch("http://localhost:8000/app/setpushtoken", {
          method: "POST",
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include', // Important for sending cookies
          body: JSON.stringify({
            "Token": token
          })
        });
    
        const result = await response.json();
        if (response.ok) {
          console.log("Push token set successfully:", result);
          return true;
        } else {
          console.error("Failed to set push token:", result);
          return false;
        }
    } catch (error) {
        console.error("Error setting push token:", error);
        return false;
    }
};

    
