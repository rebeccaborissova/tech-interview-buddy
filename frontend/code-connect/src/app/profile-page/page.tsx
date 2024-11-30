"use client";

import { useEffect, useState } from "react";
import styles from "./Profile.module.css";
import { useRouter } from "next/navigation";
import { getToken } from "../utils/token";

interface UserProfile {
  Email: string;
  FirstName: string;
  LastName: string;
  Year: number;
  Description: string;
}

const ProfilePage = () => {
  const [sessionToken, setSessionToken] = useState<string | null>(null);
  const [userProfile, setUserProfile] = useState<UserProfile>({
    Email: "dummy.user@example.com",
    FirstName: "Dummy",
    LastName: "User",
    Year: 2024,
    Description: "This is a dummy profile for testing purposes.",
  });
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const router = useRouter();

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push("/login");
    }
    setSessionToken(token);
  }, []);

  const handleBack = () => router.push("/dashboard");
  const handleSignOut = () => {
    document.cookie = "session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    router.push("/login");
  };
  const handleDeleteProfile = () => alert("Delete Profile button clicked!");

  const handleEditProfile = () => {
    setIsEditing(true);
  };

  const handleSaveChanges = () => {
    console.log("Saving changes:", userProfile);
    setIsEditing(false);
  };

  const handleInputChange = (field: keyof UserProfile, value: string | number) => {
    setUserProfile((prevProfile) => ({
      ...prevProfile,
      [field]: value,
    }));
  };

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
            <label className={styles.label}>Email</label>
            <input
              type="email"
              value={userProfile.Email}
              onChange={(e) => handleInputChange("Email", e.target.value)}
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
