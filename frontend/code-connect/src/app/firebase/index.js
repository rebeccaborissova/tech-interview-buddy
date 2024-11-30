const functions = require('firebase-functions');
const admin = require('firebase-admin');
admin.initializeApp();

const generateJitsiRoom = () => {
  const roomName = `room-${Math.random().toString(36).substring(2, 15)}`;
  const jitsiDomain = "meet.jit.si";
  return `https://${jitsiDomain}/${roomName}`;
};

exports.sendCallInvite = functions.https.onCall((data, context) => {
  const { token, callerName } = data;
  const jitsiRoomUrl = generateJitsiRoom();

  const payload = {
    notification: {
      title: `${callerName} is inviting you to a video call`,
      body: 'Click to join the call',
      click_action: jitsiRoomUrl
    }
  };

  return admin.messaging().sendToDevice(token, payload)
    .then(response => {
      console.log('Successfully sent message:', response);
      return { success: true, jitsiRoomUrl };
    })
    .catch(error => {
      console.log('Error sending message:', error);
      return { success: false, error };
    });
});