import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_URL = process.env.API_URL || 'http://localhost:8080';

export async function GET() {
  try {
    const cookieStore = await cookies();
    const jwt = cookieStore.get('jwt');

    if (!jwt) {
      return NextResponse.json({ currentUser: null });
    }

    const response = await fetch(`${API_URL}/api/users/currentuser`, {
      headers: {
        Cookie: `jwt=${jwt.value}`,
      },
    });

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Current user error:', error);
    return NextResponse.json({ currentUser: null });
  }
}
