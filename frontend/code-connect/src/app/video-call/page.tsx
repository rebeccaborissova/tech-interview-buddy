"use client";
import React, { useEffect, useState } from "react";
import { JitsiMeeting } from "@jitsi/react-sdk";

const JitsiMeetComponent = () => {
    const [isClient, setIsClient] = useState(false);

    useEffect(() => {
        setIsClient(true);
    }, []);

    if (!isClient) return null;

    const roomName = "YourRoomNameHere";

    return (
        <div style={{ height: "100vh", display: "flex", flexDirection: "column" }}>
            <JitsiMeeting
                domain="actual-terribly-longhorn.ngrok-free.app"
                roomName={roomName}
                configOverwrite={{
                    startWithAudioMuted: true,
                    disableModeratorIndicator: true,
                    startScreenSharing: true,
                    enableEmailInStats: false,
                }}
                interfaceConfigOverwrite={{
                    DISABLE_JOIN_LEAVE_NOTIFICATIONS: true,
                }}
                userInfo={{
                    displayName: "Hi",
                    email: "hi@example.com",
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
