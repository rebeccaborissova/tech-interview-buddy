"use client";

import { useEffect, useState } from "react";
import styles from "./Profile.module.css";
import { useRouter } from "next/navigation";
import { getToken } from "../utils/token";

interface UserProfile {
  Username: string;
  FirstName: string;
  LastName: string;
  Year: number;
  Description: string;
  InvitedBy?: string;
  TakenDSA?: boolean;
}

const ProfilePage = () => {
  const [sessionToken, setSessionToken] = useState<string | null>(null);
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null);
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const router = useRouter();

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push("/login");
      return;
    }
    setSessionToken(token);

    // Fetch user profile
    fetchUserProfile();
  }, []);

  const fetchUserProfile = async () => {
    try {
      const response = await fetch("http://localhost:8000/app/userinfo", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        mode: "cors",
      });

      if (!response.ok) {
        throw new Error("Failed to fetch user profile.");
      }

      const data = await response.json();
      console.log("Fetched User Profile Data:", data);

      // Map response data to UserProfile
      const parsedData: UserProfile = {
        Username: data.Username,
        FirstName: data.FirstName,
        LastName: data.LastName,
        Year: Number(data.Year), 
        Description: data.Description,
        InvitedBy: data.InvitedBy, 
        TakenDSA: Boolean(data.TakenDSA), 
      };

      setUserProfile(parsedData);
    } catch (error) {
      console.error("Error fetching user profile:", error);
      router.push("/login"); // Redirect to login on error
    }
  };

  const handleBack = () => router.push("/dashboard");

   
  const handleSignOut = async () => {
    // requests for user sign out
    try {
      const response = await fetch("http://localhost:8000/account/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to sign out user.");
      }

      const data = await response.json();
      console.log("Signed user out:", data);

    } catch (error) {
      console.error("Error signing out:", error);
    }

    document.cookie = "session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    router.push("/login");
  };

  // requests for user sign out inside of this method
  const handleDeleteProfile = async () => {
    alert("Delete Profile button clicked!");

    try {
      const response = await fetch("http://localhost:8000/app/userdelete", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to delete account.");
      }

      const data = await response.json();
      console.log("Deleted user account:", data)

    } catch (error) {
      console.log("Error deleting account:", error)
    }

    document.cookie = "session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    router.push("/");
  }

  const handleEditProfile = () => {
    setIsEditing(true);
  };

  const handleSaveChanges = async () => {
    if (!userProfile) return;

    try {
      const response = await fetch("http://localhost:8000/app/useredit", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          "FirstName": userProfile.FirstName,
          "LastName": userProfile.LastName,
          "Username": userProfile.Username,
          "TakenDSA": userProfile.TakenDSA,
          "Year": userProfile.Year,
          "Description": userProfile.Description,
        }),
        credentials: "include",
        mode: 'cors',
      });

      if (!response.ok) {
        throw new Error("Failed to update profile.");
      }

      const data = await response.json();
      console.log("Updated Profile Data:", data);

      setIsEditing(false);
      alert("Profile updated successfully!");
    } catch (error) {
      console.error("Error updating profile:", error);
      alert("Failed to update profile. Please try again.");
    }
  };

  const handleInputChange = (field: keyof UserProfile, value: string | number | boolean) => {
    setUserProfile((prevProfile) => {
      if (!prevProfile) return null;

      // Type enforcement for specific fields
      if (field === "Year" && typeof value === "string") {
        value = Number(value); // Ensure Year is always a number
      }
      if (field === "TakenDSA" && typeof value === "string") {
        value = value === "true"; // Convert string to boolean if necessary
      }

      return {
        ...prevProfile,
        [field]: value,
      };
    });
  };

  if (!userProfile) {
    return <div>Loading...</div>; // Show loading while user data is being fetched
  }

  return (
    <div className={styles.container}>
      {/* Back Button */}
      <button className={styles.backButton} onClick={handleBack}>
        Return to Dashboard
      </button>

      {/* Sign Out Button */}
      <button className={styles.signOutButton} onClick={handleSignOut}>
        Sign Out
      </button>

      <div className={styles.profileCard}>
        <h1 className={styles.title}>My Profile</h1>
        <form className={styles.form}>
          <div className={styles.inputContainer}>
            <label className={styles.label}>First Name</label>
            <input
              type="text"
              value={userProfile.FirstName}
              onChange={(e) => handleInputChange("FirstName", e.target.value)}
              className={styles.input}
              disabled={!isEditing}
            />
          </div>
          <div className={styles.inputContainer}>
            <label className={styles.label}>Last Name</label>
            <input
              type="text"
              value={userProfile.LastName}
              onChange={(e) => handleInputChange("LastName", e.target.value)}
              className={styles.input}
              disabled={!isEditing}
            />
          </div>
          <div className={styles.inputContainer}>
            <label className={styles.label}>Year</label>
            <input
              type="number"
              value={userProfile.Year}
              onChange={(e) => handleInputChange("Year", Number(e.target.value))}
              className={styles.input}
              disabled={!isEditing}
            />
          </div>
          <div className={styles.inputContainer}>
            <label className={styles.label}>Description</label>
            <textarea
              value={userProfile.Description}
              onChange={(e) => handleInputChange("Description", e.target.value)}
              className={styles.textarea}
              disabled={!isEditing}
            />
          </div>

          {userProfile.InvitedBy && (
            <div className={styles.inputContainer}>
              <label className={styles.label}>Invited By</label>
              <input
                type="text"
                value={userProfile.InvitedBy}
                className={styles.input}
                disabled
              />
            </div>
          )}

          {userProfile.TakenDSA !== undefined && (
            <div className={styles.inputContainer}>
              <label className={styles.label}>Taken DSA</label>
              <input
                type="checkbox"
                checked={userProfile.TakenDSA}
                onChange={(e) => handleInputChange("TakenDSA", e.target.checked)}
                className={styles.checkbox}
                disabled={!isEditing}
              />
            </div>
          )}
        </form>

        {/* Edit/Save Changes Button */}
        <button
          onClick={isEditing ? handleSaveChanges : handleEditProfile}
          className={styles.editSaveButton}
        >
          {isEditing ? "Save All Changes" : "Edit Profile"}
        </button>
      </div>

      {/* Delete Profile Button */}
      <button className={styles.deleteProfileButton} onClick={handleDeleteProfile}>
        Delete Profile
      </button>
    </div>
  );
};

export default ProfilePage;
