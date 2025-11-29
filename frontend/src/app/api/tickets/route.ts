import { NextResponse } from 'next/server';

const API_URL = process.env.API_URL || 'http://localhost:8080';

export async function GET() {
  try {
    const response = await fetch(`${API_URL}/api/tickets`);
    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Tickets list error:', error);
    return NextResponse.json(
      { errors: [{ message: 'Failed to fetch tickets' }] },
      { status: 500 }
    );
  }
}

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const response = await fetch(`${API_URL}/api/tickets`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    const data = await response.json();
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Create ticket error:', error);
    return NextResponse.json(
      { errors: [{ message: 'Failed to create ticket' }] },
      { status: 500 }
    );
  }
}
