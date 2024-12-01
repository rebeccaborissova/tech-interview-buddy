import "server-only"

const admin = require('firebase-admin');
const serviceAccount = require('../../../firebase-admin-config.json');

try {
  admin.initializeApp({
    credential: admin.credential.cert(serviceAccount),
  })
  console.log('Initialized.')
} catch (error) {
  /*
   * We skip the "already exists" message which is
   * not an actual error when we're hot-reloading.
   */
  if (!/already exists/u.test(error.message)) {
    console.error('Firebase admin initialization error', error.stack)
  }
}

export default admin

export const sendNotification = async (registrationToken, jitsiRoom, username) => {
    const message = {
      notification: {
        title: "Incoming Call From " + username,
        body: jitsiRoom,
      },
      token: registrationToken
    };
    
    return admin.messaging().send(message);
};
