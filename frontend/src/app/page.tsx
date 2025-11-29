import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export default function HomePage() {
  return (
    <div className="space-y-12">
      {/* Hero Section */}
      <section className="text-center py-12">
        <h1 className="text-5xl font-bold mb-4">Welcome to TicketHub</h1>
        <p className="text-xl text-muted-foreground mb-8">
          Buy and sell tickets for events, concerts, and more
        </p>
        <div className="flex gap-4 justify-center">
          <Link href="/tickets">
            <Button size="lg">Browse Tickets</Button>
          </Link>
          <Link href="/signup">
            <Button size="lg" variant="outline">Get Started</Button>
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section className="grid md:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>🎫 Easy Buying</CardTitle>
            <CardDescription>
              Browse and purchase tickets in just a few clicks
            </CardDescription>
          </CardHeader>
          <CardContent>
            Simple, secure checkout process with Stripe integration.
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>⚡ Real-time Updates</CardTitle>
            <CardDescription>
              Get instant notifications about your orders
            </CardDescription>
          </CardHeader>
          <CardContent>
            Event-driven architecture ensures you&apos;re always up to date.
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>🔒 Secure Transactions</CardTitle>
            <CardDescription>
              Your payments are safe and encrypted
            </CardDescription>
          </CardHeader>
          <CardContent>
            Built with industry-standard security practices.
          </CardContent>
        </Card>
      </section>

      {/* CTA Section */}
      <section className="text-center py-12 bg-muted rounded-lg">
        <h2 className="text-3xl font-bold mb-4">Ready to get started?</h2>
        <p className="text-muted-foreground mb-6">
          Join thousands of users buying and selling tickets
        </p>
        <Link href="/signup">
          <Button size="lg">Create Account</Button>
        </Link>
      </section>
    </div>
  );
}
