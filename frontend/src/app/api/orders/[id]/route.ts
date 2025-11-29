import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_URL = process.env.API_URL || 'http://localhost:8080';

export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    const cookieStore = await cookies();
    const jwt = cookieStore.get('jwt');

    if (!jwt) {
      return NextResponse.json(
        { errors: [{ message: 'Not authenticated' }] },
        { status: 401 }
      );
    }

    const response = await fetch(`${API_URL}/api/orders/${id}`, {
      headers: {
        Cookie: `jwt=${jwt.value}`,
      },
    });

    const data = await response.json();
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Get order error:', error);
    return NextResponse.json(
      { errors: [{ message: 'Failed to fetch order' }] },
      { status: 500 }
    );
  }
}

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    const cookieStore = await cookies();
    const jwt = cookieStore.get('jwt');

    if (!jwt) {
      return NextResponse.json(
        { errors: [{ message: 'Not authenticated' }] },
        { status: 401 }
      );
    }

    const response = await fetch(`${API_URL}/api/orders/${id}`, {
      method: 'DELETE',
      headers: {
        Cookie: `jwt=${jwt.value}`,
      },
    });

    const data = await response.json();
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Cancel order error:', error);
    return NextResponse.json(
      { errors: [{ message: 'Failed to cancel order' }] },
      { status: 500 }
    );
  }
}
