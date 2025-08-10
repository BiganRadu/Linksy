import * as React from 'react';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

import LinkRoundedIcon from '@mui/icons-material/LinkRounded';
import QrCode2RoundedIcon from '@mui/icons-material/QrCode2Rounded';
import BarChartRoundedIcon from '@mui/icons-material/BarChartRounded';
import SecurityRoundedIcon from '@mui/icons-material/SecurityRounded';

import { LinksyIcon } from './CustomIcons';

const items = [
  {
    icon: <LinkRoundedIcon sx={{ color: 'text.secondary' }} />,
    title: 'Fast & reliable link shortening',
    description:
      'Quickly create short, memorable links that are easy to share and never let you down.',
  },
  {
    icon: <QrCode2RoundedIcon sx={{ color: 'text.secondary' }} />,
    title: 'Instant QR code generation',
    description:
      'Generate sleek, scannable QR codes for every link, perfect for print or offline sharing.',
  },
  {
    icon: <BarChartRoundedIcon sx={{ color: 'text.secondary' }} />,
    title: 'Detailed link analytics',
    description:
      'Track clicks, locations, and usage patterns with real-time analytics to measure your reach.',
  },
  {
    icon: <SecurityRoundedIcon sx={{ color: 'text.secondary' }} />,
    title: 'Secure & private',
    description:
      'Your links and analytics are protected with top-notch security to keep your data safe.',
  },
];

export default function Content() {
  return (
    <Stack
      sx={{ flexDirection: 'column', alignSelf: 'center', gap: 4, maxWidth: 450 }}
    >
      <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
        <LinksyIcon />
      </Box>
      {items.map((item, index) => (
        <Stack key={index} direction="row" sx={{ gap: 2 }}>
          {item.icon}
          <div>
            <Typography gutterBottom sx={{ fontWeight: 'medium' }}>
              {item.title}
            </Typography>
            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
              {item.description}
            </Typography>
          </div>
        </Stack>
      ))}
    </Stack>
  );
}
