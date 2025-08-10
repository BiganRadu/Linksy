import * as React from 'react';
import SvgIcon from '@mui/material/SvgIcon';

export function LinksyIcon() {
  return (
    <SvgIcon sx={{ height: 36, width: 150, mr: 2 }} viewBox="0 0 150 36">
      {/* Link icon background */}
      <circle cx="18" cy="18" r="15" fill="#1976d2" />
      <path
        d="M13.2 19.2a5 5 0 0 1 0-7.07l2.65-2.65a5 5 0 0 1 7.07 0l1.33 1.33-1.67 1.67-1.33-1.33a2.5 2.5 0 0 0-3.54 0l-2.65 2.65a2.5 2.5 0 0 0 0 3.54l1.33 1.33-1.67 1.67-1.32-1.33z"
        fill="#fff"
      />
      <path
        d="M22.8 16.8a5 5 0 0 1 0 7.07l-2.65 2.65a5 5 0 0 1-7.07 0l-1.33-1.33 1.67-1.67 1.33 1.33a2.5 2.5 0 0 0 3.54 0l2.65-2.65a2.5 2.5 0 0 0 0-3.54l-1.33-1.33 1.67-1.67 1.32 1.33z"
        fill="#fff"
      />

      {/* Linksy text */}
      <text
        x="42"
        y="23"
        fontFamily="Arial, Helvetica, sans-serif"
        fontSize="20"
        fontWeight="bold"
        fill="#1976d2"
      >
        Linksy
      </text>
    </SvgIcon>
  );
}
