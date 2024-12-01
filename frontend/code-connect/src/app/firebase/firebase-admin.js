import "server-only"

const admin = require('firebase-admin');
const serviceAccount = require('../../firebase-admin-config.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
}, "firebase-admin");

export const sendNotification = async (registrationToken, title, body) => {
    const message = {
        data: {
        title: title,
        body: body
        },
        token: registrationToken
    };
    admin.messaging().send(message)
        .then((response) => {
        console.log('Notification sent:', response);
        })
        .catch((error) => {
        console.error('Error sending notification:', error);
        }); 
};
