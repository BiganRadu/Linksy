import * as React from 'react';
import Grid from '@mui/material/Grid2';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import AnalyticsRoundedIcon from '@mui/icons-material/AnalyticsRounded';
import LinkIcon from '@mui/icons-material/Link';
import QrCodeIcon from '@mui/icons-material/QrCode';
import Button from '@mui/material/Button';


export default function MainGrid() {
  return (
    <Box sx={{ width: '100%', maxWidth: '800px', mt: 3}}>
      <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
        Overview
      </Typography>
	  <Box sx={{ display: 'flex'}}>
		<LinkIcon />
		<Box>
			<Typography variant="h6" component="div">
				Links
			</Typography>
			<Button>
				View Links
			</Button>
		</Box>
	  </Box>
      <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
        Details
      </Typography>
    </Box>
  );
}
