import { useEffect, useState, useMemo } from 'react';

interface TimeLeft {
  minutes: number;
  seconds: number;
  expired: boolean;
}

function calculateTimeLeft(targetDate: Date): TimeLeft {
  const difference = targetDate.getTime() - Date.now();

  if (difference <= 0) {
    return { minutes: 0, seconds: 0, expired: true };
  }

  return {
    minutes: Math.floor((difference / 1000 / 60) % 60),
    seconds: Math.floor((difference / 1000) % 60),
    expired: false,
  };
}

export function useCountdown(targetDate: Date | string) {
  const date = useMemo(
    () => (typeof targetDate === 'string' ? new Date(targetDate) : targetDate),
    [targetDate]
  );
  const [timeLeft, setTimeLeft] = useState<TimeLeft>(() => calculateTimeLeft(date));

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft(calculateTimeLeft(date));
    }, 1000);

    return () => clearInterval(timer);
  }, [date]);

  return timeLeft;
}
