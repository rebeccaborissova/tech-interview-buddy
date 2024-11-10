"use client";
import React from "react";
import { JitsiMeeting } from "@jitsi/react-sdk";

const JitsiMeetComponent = () => {
    const roomName = "YourRoomNameHere"; // Replace with your preferred room name

    return (
        <div style={{ height: "100vh", display: "flex", flexDirection: "column" }}>
            <JitsiMeeting
                domain="meet.jit.si"
                roomName={roomName}
                configOverwrite={{
                    startWithAudioMuted: true,
                    disableModeratorIndicator: true,
                    startScreenSharing: true,
                    enableEmailInStats: false
                }}
                interfaceConfigOverwrite={{
                    DISABLE_JOIN_LEAVE_NOTIFICATIONS: true
                }}
                userInfo={{
                    displayName: 'Hi',
                    email: "hi@example.com"
                }}
                onApiReady={(externalApi) => {
                    // Attach custom event listeners here if needed
                }}
                getIFrameRef={(iframeRef) => {
                    iframeRef.style.height = "100vh";
                    iframeRef.style.width = "100%";
                }}
            />
        </div>
    );
};

export default JitsiMeetComponent;
