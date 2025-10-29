import { useEffect } from 'react';

export function useEnterKey(onEnter) {
  useEffect(() => {
    const handleEnterKey = (event) => {
      if (event.key === 'Enter') {
        onEnter();
      }
    };

    document.addEventListener('keydown', handleEnterKey);

    return () => {
      document.removeEventListener('keydown', handleEnterKey);
    };
  }, [onEnter]);
}
