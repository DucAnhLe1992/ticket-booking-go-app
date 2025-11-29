'use client';

import { useCountdown } from '@/lib/hooks/useCountdown';
import { Alert, AlertDescription } from '@/components/ui/alert';

interface OrderCountdownProps {
  expiresAt: string;
  onExpire?: () => void;
}

export function OrderCountdown({ expiresAt, onExpire }: OrderCountdownProps) {
  const { minutes, seconds, expired } = useCountdown(expiresAt);

  if (expired) {
    onExpire?.();
    return (
      <Alert variant="destructive">
        <AlertDescription className="font-semibold">
          This order has expired
        </AlertDescription>
      </Alert>
    );
  }

  const isUrgent = minutes < 5;

  return (
    <Alert variant={isUrgent ? 'destructive' : 'default'}>
      <AlertDescription className="font-semibold text-center">
        Time remaining: {minutes}:{seconds.toString().padStart(2, '0')}
      </AlertDescription>
    </Alert>
  );
}
