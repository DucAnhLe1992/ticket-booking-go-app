import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_URL = process.env.API_URL || 'http://localhost:8080';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const response = await fetch(`${API_URL}/api/users/signin`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    const data = await response.json();

    if (!response.ok) {
      return NextResponse.json(data, { status: response.status });
    }

    // Extract Set-Cookie header from Go backend
    const setCookieHeader = response.headers.get('set-cookie');
    if (setCookieHeader) {
      const cookieStore = await cookies();
      const jwtMatch = setCookieHeader.match(/jwt=([^;]+)/);
      if (jwtMatch) {
        cookieStore.set('jwt', jwtMatch[1], {
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite: 'lax',
          path: '/',
        });
      }
    }

    return NextResponse.json(data);
  } catch (error) {
    console.error('Signin error:', error);
    return NextResponse.json(
      { errors: [{ message: 'Failed to sign in' }] },
      { status: 500 }
    );
  }
}
