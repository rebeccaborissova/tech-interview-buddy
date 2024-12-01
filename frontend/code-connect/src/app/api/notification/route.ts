import { NextResponse } from 'next/server'
import { sendNotification } from '../../firebase/firebase-admin'

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { pushToken: token, body: jitsiRoom, username: username } = body;
    
    if (!token) {
      throw new Error('No push token provided');
    }

    const response = await sendNotification(token, jitsiRoom, username);
    return NextResponse.json({ success: true, response });

  } catch (error) {
    console.error('Error sending notification:', error);
    return NextResponse.json(
      { error: error },
      { status: 500 }
    );
  }
}