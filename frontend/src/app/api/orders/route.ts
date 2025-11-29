import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_URL = process.env.API_URL || 'http://localhost:8080';

export async function GET() {
  try {
    const cookieStore = await cookies();
    const jwt = cookieStore.get('jwt');

    if (!jwt) {
      return NextResponse.json(
        { errors: [{ message: 'Not authenticated' }] },
        { status: 401 }
      );
    }

    const response = await fetch(`${API_URL}/api/orders`, {
      headers: {
        Cookie: `jwt=${jwt.value}`,
      },
    });

    const data = await response.json();
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Orders list error:', error);
    return NextResponse.json(
      { errors: [{ message: 'Failed to fetch orders' }] },
      { status: 500 }
    );
  }
}

export async function POST(request: Request) {
  try {
    const cookieStore = await cookies();
    const jwt = cookieStore.get('jwt');

    if (!jwt) {
      return NextResponse.json(
        { errors: [{ message: 'Not authenticated' }] },
        { status: 401 }
      );
    }

    const body = await request.json();
    const response = await fetch(`${API_URL}/api/orders`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Cookie: `jwt=${jwt.value}`,
      },
      body: JSON.stringify(body),
    });

    const data = await response.json();
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Create order error:', error);
    return NextResponse.json(
      { errors: [{ message: 'Failed to create order' }] },
      { status: 500 }
    );
  }
}
